package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	users "backend/internal/core/users/grpc"
	userService "backend/internal/genproto/users"
)

func main() {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	userHandler := users.NewHandler()

	userService.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", "9000"))
	if err != nil {
		panic(err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
