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

func isStudentExisting(roll string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Student
	filter := bson.M{
		"roll": roll,
	}
	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}
func StudentCreateNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Student
	var self Models.Student
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		return err
	}
	_, existing := isStudentExisting(self.Roll)
	if existing != "" {
		return errors.New("Roll already exists to another student")
	}
	_, err = collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Status(200).Send([]byte("Food Item Created Successfully"))
	return nil
}
func StudentGetByID(id primitive.ObjectID) (Models.Student, error) {
	collection := DBManager.SystemCollections.Student
	filter := bson.M{"_id": id}
	var self Models.Student
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
func StudentDelete(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Student
	deleteID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	_, err := StudentGetByID(deleteID)
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
func StudentModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Student
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	var results []Models.Student
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
	var self Models.Student
	c.BodyParser(&self)
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	_, id := isStudentExisting(self.Roll)
	if id != "" && id != objID.Hex() {
		c.Status(500)
		return errors.New("Roll already exists to another student")
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
func StudentGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Student
	page, _ := strconv.Atoi(c.Params("page"))
	rowsPerPage, _ := strconv.Atoi(c.Params("rowsperpage"))
	var results []Models.Student
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
			"results":      results[strt:end],
			"totalstudent": len(results),
		},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}
func StudentSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Student
	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	newValue := true
	if c.Params("new_status") == "inactive" {
		newValue = false
	}
	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying Branch status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}
