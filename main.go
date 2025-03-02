package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TODO struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BODY      string             `json:"body" bson:"body"`
	Completed bool               `json:"completed" bson:"completed"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load the env file")
	}
	MONGODB_URI := os.Getenv("MONGODB_URI")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	clientOption := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal("failed to connect to the MONGO")
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("failed to connect to the MONGO")
	}
	fmt.Printf("Database connected to the database.")
	collection = client.Database("todo-go-lang").Collection("todo")

	app := fiber.New()
	app.Get("/api/todos", getAllTodos)
	app.Post("/api/todo", AddTodo)
	app.Get("/api/todo/:id", GetOneTodo)
	app.Put("/api/todo/:id", ChangeStatus)
	log.Fatal(app.Listen(":" + PORT))
}

func getAllTodos(c *fiber.Ctx) error {
	var allTodos []TODO
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Unable to query in the database"})
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var todo TODO
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
	todo := new(TODO)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}
	if todo.BODY == "" {
		return c.Status(501).JSON(fiber.Map{"success": false, "message": "Body is required"})
	}
	response, err := collection.InsertOne(context.Background(), todo)
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
	fetchTodo := new(TODO)
	response := collection.FindOne(context.Background(), bson.M{"_id": objectId})
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
	response, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Todo not found", "err": err.Error()})
	}

	return c.Status(401).JSON(fiber.Map{"success": true, "message": "todo status updated", "data": response})
}
