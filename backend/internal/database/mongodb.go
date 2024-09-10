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

// func (s *service) Health() map[string]string {
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
// 	defer cancel()

// 	err := s.db.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatalf(fmt.Sprintf("db down: %v", err))
// 	}

// 	return map[string]string{
// 		"message": "Mongo is healthy",
// 	}
// }

// func (s *service) getDatabase() mongo.Database {
// 	return *s.db.Database("reg_dealer")
// }
