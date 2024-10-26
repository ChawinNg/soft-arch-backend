package database

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
}

type service struct {
	db *mongo.Database
}

var (
	// host = os.Getenv("DB_HOST")
	// port = os.Getenv("DB_PORT")
	database = os.Getenv("DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(database))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client.Database("reg_dealer"),
	}
}
