package main

import (
	"fmt"

	userService "backend/internal/genproto/users"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	user "backend/internal/core/users/rest"
)

func main() {

	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	r := fiber.New()

	userConn := userService.NewUserServiceClient(conn)
	userHandler := user.NewHandler(userConn)
	r.Get("/users", userHandler.GetAllUsers)
	r.Post("/users", userHandler.CreateUser)

	err = r.Listen(fmt.Sprintf(":%v", "8080"))
	if err != nil {
		panic(err)
	}
}
