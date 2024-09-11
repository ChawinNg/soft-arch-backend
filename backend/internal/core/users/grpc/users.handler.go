package users

import (
	"context"

	"backend/internal/genproto/users"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"backend/internal/model"
)

type Handler struct {
	users.UnimplementedUserServiceServer
	db *mongo.Collection
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{
		db: db.Collection("users"),
	}
}

func (h *Handler) GetAllUsers(c context.Context, _ *users.GetAllUserRequest) (*users.GetAllUserResponse, error) {
	cursor, err := h.db.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var users []model.User
	for cursor.Next(c) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// res := users.GetAllUserResponse{}
	return nil, nil
}

func (h *Handler) CreateUser(c context.Context, u *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	_, err := h.db.InsertOne(c, u)
	return nil, err
}
