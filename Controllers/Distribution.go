package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"example.com/example/DBManager"
	"example.com/example/Models"
	"example.com/example/Utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isDistributionExisting(id primitive.ObjectID, shift string, date primitive.DateTime) bool {
	collection := DBManager.SystemCollections.Distribution
	filter := bson.M{
		"studentid": id,
		"shift":     shift,
		"date":      date,
	}
	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}
func DistributionCreateNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Distribution
	var self Models.Distribution
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		return err
	}
	if self.StudentID == primitive.NilObjectID {
		c.Status(500)
		return errors.New("StudentID cannot be empty")
	}
	if isDistributionExisting(self.StudentID, self.Shift, self.Date) {
		c.Status(500)
		return errors.New("Already served")
	}

	_, err = collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Status(200).Send([]byte("Distribution Created Successfully"))
	return nil
}
func DistributionGetByID(id primitive.ObjectID) (Models.Distribution, error) {
	collection := DBManager.SystemCollections.Student
	filter := bson.M{"_id": id}
	var self Models.Distribution
	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return self, errors.New("object not found")
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 0 {
		return self, errors.New("object not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}
func DistributionGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Distribution
	page, _ := strconv.Atoi(c.Params("page"))
	rowsPerPage, _ := strconv.Atoi(c.Params("rowsperpage"))
	var results []Models.Distribution
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	strt := page * rowsPerPage
	end := page*rowsPerPage + rowsPerPage
	if page == -1 {
		strt = 0
	} else if strt > len(results) {
		strt = len(results)
	}
	if end > len(results) || rowsPerPage == -1 {
		end = len(results)
	}
	// Decode
	response, _ := json.Marshal(
		bson.M{
			"results":           results[strt:end],
			"totaldistribution": len(results),
		},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}
func DistributionGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Distribution
	page, _ := strconv.Atoi(c.Params("page"))
	rowsPerPage, _ := strconv.Atoi(c.Params("rowsperpage"))
	var results []Models.Distribution
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)

	populatedResults := []Models.DistributionPopulated{}
	for _, val := range results {
		var resultPopulated Models.DistributionPopulated
		resultPopulated.CloneFrom(val)
		resultPopulated.Student, _ = StudentGetByID(val.StudentID)
		populatedResults = append(populatedResults, resultPopulated)
	}
	strt := page * rowsPerPage
	end := page*rowsPerPage + rowsPerPage
	if page == -1 {
		strt = 0
	} else if strt > len(results) {
		strt = len(results)
	}
	if end > len(results) || rowsPerPage == -1 {
		end = len(results)
	}
	// Decode
	response, _ := json.Marshal(
		bson.M{
			"results":           populatedResults[strt:end],
			"totaldistribution": len(populatedResults),
		},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}
