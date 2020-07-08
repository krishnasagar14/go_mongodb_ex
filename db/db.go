package DB

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go_mongodb_ex/models"
	"go_mongodb_ex/proto"
)

var db *mongo.Database

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Double sure connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("local_db")
	log.Println("Connected to database successfully !!!")
}

func get_user_collection() *mongo.Collection {
	return db.Collection("User")
}

func InsertNewUser(data *server_proto.CreateUserRequest) (string, string) {
	user_collection := get_user_collection()

	insert_data := models.UserStruct{
		Id:          primitive.NewObjectID(),
		EmployeeId:  models.GetEmployeeSeqNumber(),
		FirstName:   data.GetFirstName(),
		LastName:    data.GetLastName(),
		Email:       data.GetEmail(),
		Designation: data.GetDesignation(),
	}
	_, err := user_collection.InsertOne(context.TODO(), insert_data)
	if err != nil {
		log.Println(err)
	}
	return insert_data.Id.Hex(), insert_data.EmployeeId
}

func UpdateUser(data *server_proto.UpdateUserRequest) models.UserStruct {
	user_collection := get_user_collection()

	id, _ := primitive.ObjectIDFromHex(data.GetUserId())

	update_data := models.UserStruct{
		Id:    id,
		Email: data.GetEmail(),
	}
	filter := bson.M{"_id": update_data.Id}
	update := bson.M{
		"$set": bson.M{"email": update_data.Email},
	}

	_, err := user_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
	var res_user models.UserStruct
	err = user_collection.FindOne(context.TODO(), filter).Decode(&res_user)
	return res_user
}
