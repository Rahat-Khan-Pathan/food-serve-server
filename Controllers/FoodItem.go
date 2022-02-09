package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"example.com/example/DBManager"
	"example.com/example/Models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FoodItemCreateNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.FoodItem
	var self Models.FoodItem
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Status(200).Send([]byte("Food Item Created Successfully"))
	return nil
}
func FoodItemGetByID(id primitive.ObjectID) (Models.FoodItem, error) {
	collection := DBManager.SystemCollections.FoodItem
	filter := bson.M{"_id": id}
	var self Models.FoodItem
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
func FoodItemDelete(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.FoodItem
	deleteID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	_, err := FoodItemGetByID(deleteID)
	if err != nil {
		return errors.New("object not found")
	}
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": deleteID})
	if err != nil {
		c.Status(500)
		return err
	}
	c.Status(200).Send([]byte("deleted successfully"))
	return nil
}
func FoodItemModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.FoodItem
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return errors.New("object not found")
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.FoodItem
	c.BodyParser(&self)
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying Branch Document")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}
func FoodItemGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.FoodItem
	page, _ := strconv.Atoi(c.Params("page"))
	rowsPerPage, _ := strconv.Atoi(c.Params("rowsperpage"))
	var results []Models.FoodItem
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
			"results":       results[strt:end],
			"totalfooditem": len(results),
		},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}
