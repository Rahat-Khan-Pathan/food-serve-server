package Models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Distribution struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StudentID    primitive.ObjectID `json:"studentid"`
	Date         primitive.DateTime `json:"date"`
	Shift        string             `json:"shift"`
	Status       bool               `json:"status"`
	FoodItemList []FoodItem         `json:"fooditemlist"`
}
type DistributionPopulated struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Student      Student            `json:"student"`
	Date         primitive.DateTime `json:"date"`
	Shift        string             `json:"shift"`
	Status       bool               `json:"status"`
	FoodItemList []FoodItem         `json:"fooditemlist"`
}

func (obj Distribution) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Date, validation.Required),
		validation.Field(&obj.Shift, validation.Required),
		validation.Field(&obj.FoodItemList, validation.Length(1, 0)),
	)
}
func (obj *DistributionPopulated) CloneFrom(other Distribution) {
	obj.ID = other.ID
	obj.Date = other.Date
	obj.Shift = other.Shift
	obj.Status = other.Status
	obj.FoodItemList = other.FoodItemList
	obj.Student = Student{}
}
