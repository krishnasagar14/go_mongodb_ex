package models

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var EmployeeSeqNo int = 0

type UserStruct struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	EmployeeId  string             `json:"employee_id" bson:"employee_id"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Email       string             `json:"email" bson:"email"`
	Designation string             `json:"designation" bson:"designation"`
}

func GetEmployeeSeqNumber() string {
	EmployeeSeqNo = EmployeeSeqNo + 1
	return strconv.Itoa(EmployeeSeqNo)
}
