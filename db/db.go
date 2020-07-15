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

func getEmployeeCollection() *mongo.Collection {
	return db.Collection("Employee")
}

func InsertNewUser(data *server_proto.CreateUserRequest) (string, string) {
	userCollection := getUserCollection()
	employeeCollection := getEmployeeCollection()

	insertUserData := models.UserStruct{
		Id:          primitive.NewObjectID(),
		FirstName:   data.GetFirstName(),
		LastName:    data.GetLastName(),
		Email:       data.GetEmail(),
	}
	_, err := userCollection.InsertOne(context.TODO(), insertUserData)
	if err != nil {
		log.Println(err)
		return "", ""
	}

	insertEmpData := models.EmployeeStruct{
		Id: primitive.NewObjectID(),
		UserId: insertUserData.Id.Hex(),
		Designation: data.GetDesignation(),
	}
	_, err = employeeCollection.InsertOne(context.TODO(), insertEmpData)
	if err != nil {
		log.Println(err)
		return "", ""
	}

	return insertUserData.Id.Hex(), insertEmpData.Id.Hex()
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

func GetEmployeeViaUserId(userId string) (models.EmployeeStruct, error) {
	employeeCollection := getEmployeeCollection()

	filter := bson.M{"user_id": userId}
	var resEmp models.EmployeeStruct
	err := employeeCollection.FindOne(context.TODO(), filter).Decode(&resEmp)
	if err != nil {
		log.Println(err)
	}
	return resEmp, err
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
