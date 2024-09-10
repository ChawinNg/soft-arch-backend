package main

import (
	"context"
	"log"
	"os"

	"backend/internal/core/users"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Initialize Fiber
	app := fiber.New()

	//Connect to Database

	database := os.Getenv("DB_DATABASE")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(database))

	if err != nil {
		log.Fatal(err)

	}

	db := client.Database("reg_dealer")

	// Initialize User service and handler
	userService := users.NewUserService(db)
	userHandler := users.NewUserHandler(userService)

	// Define route
	app.Get("/users/:id", userHandler.GetUser)

	// Start the server
	log.Fatal(app.Listen(":8080"))

}
