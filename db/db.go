package db

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

func ConnectDB(dbName string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		return err
	}
	// Double sure connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	db = client.Database(dbName)
	log.Println("Connected to database successfully !!!")
	return nil
	// TODO: DB connection management
}

func DropDB() {
	db.Drop(context.TODO())
	log.Println("Dropped database successfully !!!")
}

func getUserCollection() *mongo.Collection {
	return db.Collection("User")
}

func InsertNewUser(data *server_proto.CreateUserRequest) (string, string) {
	userCollection := getUserCollection()

	insertData := models.UserStruct{
		Id:          primitive.NewObjectID(),
		EmployeeId:  models.GetEmployeeSeqNumber(),
		FirstName:   data.GetFirstName(),
		LastName:    data.GetLastName(),
		Email:       data.GetEmail(),
		Designation: data.GetDesignation(),
	}
	_, err := userCollection.InsertOne(context.TODO(), insertData)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return insertData.Id.Hex(), insertData.EmployeeId
}

func GetUserViaId(userId string) (models.UserStruct, error) {
	userCollection := getUserCollection()

	idObj, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": idObj}

	var resUser models.UserStruct
	err := userCollection.FindOne(context.TODO(), filter).Decode(&resUser)
	if err != nil {
		log.Println(err)
	}
	return resUser, err
}

func UpdateUser(data *server_proto.UpdateUserRequest) (models.UserStruct, error) {
	userCollection := getUserCollection()

	id, _ := primitive.ObjectIDFromHex(data.GetUserId())

	updateData := models.UserStruct{
		Id:    id,
		Email: data.GetEmail(),
	}
	filter := bson.M{"_id": updateData.Id}
	update := bson.M{
		"$set": bson.M{"email": updateData.Email},
	}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
	return GetUserViaId(data.GetUserId())
}
