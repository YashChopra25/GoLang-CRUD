package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	todo_controller "github.com/yashchopra25/go-lang-todo-app/controllers"
	"github.com/yashchopra25/go-lang-todo-app/database" // ✅ Import database package
)

func main() {
	fmt.Println("🚀 Starting server...")

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load the .env file")
	}

	// Connect to Database
	database.ConnectDB() // ✅ Connect to MongoDB

	// Get PORT from .env
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}

	app := fiber.New()

	// Routes
	app.Get("/api/todos", todo_controller.GetAllTodos)
	app.Post("/api/todo", todo_controller.AddTodo)
	app.Get("/api/todo/:id", todo_controller.GetOneTodo)
	app.Put("/api/todo/:id", todo_controller.ChangeStatus)

	// Start server
	log.Fatal(app.Listen(":" + PORT))
}
