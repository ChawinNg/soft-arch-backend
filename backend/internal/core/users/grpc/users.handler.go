package users

import (
	"context"

	"backend/internal/genproto/users"
)

type Handler struct {
	users.UnimplementedUserServiceServer
	// db *mongo.Database
}

func NewHandler() *Handler {
	return &Handler{
		// db:db,
	}
}

func (h *Handler) GetAllUsers(c context.Context, _ *users.GetAllUserRequest) (*users.GetAllUserResponse, error) {
	res := &users.GetAllUserResponse{
		User: make([]*users.User, 0),
	}
	return res, nil
}

func (h *Handler) CreateUser(c context.Context, u *users.CreateUserRequest) (*users.CreateUserResponse, error) {

	return nil, nil
}
