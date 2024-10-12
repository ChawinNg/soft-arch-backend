package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	users "backend/internal/core/users/grpc"
	userService "backend/internal/genproto/users"
)

func main() {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	database := os.Getenv("DB_DATABASE")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(database))

	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("reg_dealer")

	userHandler := users.NewHandler(db)

	userService.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", "9000"))
	if err != nil {
		panic(err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
