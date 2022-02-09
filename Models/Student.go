package Models

import (
	"reflect"
	"strings"

	"example.com/example/Utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullname"`
	Roll     string             `json:"roll"`
	Age      string             `json:"age"`
	Class    string             `json:"class"`
	Hall     string             `json:"hall"`
	Status   bool               `json:"status"`
}

func (obj Student) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.FullName, validation.Required),
		validation.Field(&obj.Roll, validation.Required),
		validation.Field(&obj.Age, validation.Required),
		validation.Field(&obj.Class, validation.Required),
		validation.Field(&obj.Hall, validation.Required),
	)
}
func (obj Student) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}
