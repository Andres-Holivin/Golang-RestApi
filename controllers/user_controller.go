package controllers

import (
	"context"
	"fmt"
	"golang-restapi/configs"
	"golang-restapi/models"
	"golang-restapi/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}
	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.UserResponses{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}
func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()
	fmt.Println(userId)
	objId, _ := primitive.ObjectIDFromHex(userId)
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponses{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}
func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}
	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	var updateUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updateUser)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponses{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updateUser}})
}
func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if result.DeletedCount < 1 {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "user not found"}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponses{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "user deleted"}})
}
func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user []models.User
	defer cancel()
	result, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer result.Close(ctx)
	for result.Next(ctx) {
		var singleUser models.User
		if err = result.Decode(&singleUser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.UserResponses{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		user = append(user, singleUser)
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponses{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}
