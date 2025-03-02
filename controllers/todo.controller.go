package todo_controller

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yashchopra25/go-lang-todo-app/database"
	todo_model "github.com/yashchopra25/go-lang-todo-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection *mongo.Collection

//	func InitTodoCollection(c *mongo.Collection) {
//		collection = c
//	}
func GetAllTodos(c *fiber.Ctx) error {
	var allTodos []todo_model.TODO
	cursor, err := database.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Unable to query in the database"})
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var todo todo_model.TODO
		err = cursor.Decode(&todo)
		if err != nil {
			return c.Status(501).JSON(fiber.Map{"success": false, "message": "Something went wrong", "err": err.Error()})
		}
		allTodos = append(allTodos, todo)
	}
	return c.Status(200).JSON(fiber.Map{"success": false, "message": "Todos fetch succesfully", "data": allTodos})
}

func AddTodo(c *fiber.Ctx) error {
	// todo := &TODO{}
	todo := new(todo_model.TODO)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}
	if todo.BODY == "" {
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Body is required"})
	}
	response, err := database.Collection.InsertOne(context.Background(), todo)
	if err != nil {
		fmt.Print(err)
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Failed to insert into the database", "err": err.Error()})
	}
	todo.Id = response.InsertedID.(primitive.ObjectID)
	fmt.Print(response)

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Data inserted into the database", "data": todo})
}

func GetOneTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid todo id"})
	}
	fetchTodo := new(todo_model.TODO)
	response := database.Collection.FindOne(context.Background(), bson.M{"_id": objectId})
	fmt.Print(response.Raw())
	err = response.Decode(fetchTodo)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Todo not found", "err": err.Error()})
	}

	return c.Status(401).JSON(fiber.Map{"success": true, "message": "todo fetch", "data": fetchTodo})
}

func ChangeStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid todo id"})
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}
	response, err := database.Collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Todo not found", "err": err.Error()})
	}

	return c.Status(401).JSON(fiber.Map{"success": true, "message": "todo status updated", "data": response})
}
