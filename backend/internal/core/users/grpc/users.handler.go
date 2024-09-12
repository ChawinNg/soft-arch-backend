package users

import (
	"context"
	"errors"

	"backend/internal/genproto/users"
	"backend/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
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

func (h *Handler) GetUser(c context.Context, user_id *users.GetUserRequest) (*users.GetUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(user_id.Id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err2 := h.db.FindOne(c, bson.M{"_id": id}).Decode(&user)
	if err2 == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}

	grpcUser, err := model.ConvertMongoToGrpc(user)
	if err != nil {
		return nil, err
	}

	return &users.GetUserResponse{
		User: grpcUser,
	}, nil
}

func (h *Handler) GetAllUser(c context.Context, _ *users.GetAllUserRequest) (*users.GetAllUserResponse, error) {
	cursor, err := h.db.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var users_set []*users.User
	for cursor.Next(c) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}

		grpcUser, err := model.ConvertMongoToGrpc(user)

		if err != nil {
			return nil, err
		}
		users_set = append(users_set, grpcUser)
	}

	return &users.GetAllUserResponse{User: users_set}, nil
}

func (h *Handler) CreateUser(c context.Context, u *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	_, err := h.db.InsertOne(c, u)
	return nil, err
}

func (h *Handler) UpdateUser(c context.Context, u *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(u.User.Id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"name":    u.User.Name,
			"surname": u.User.Surname,
			"email":   u.User.Email,
		},
	}
	result, err := h.db.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil || result.MatchedCount == 0 {
		return nil, err
	}

	return &users.UpdateUserResponse{User: u.User}, err
}

func (h *Handler) DeleteUser(c context.Context, user_id *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(user_id.Id)
	if err != nil {
		return nil, err
	}

	deleteResult, err := h.db.DeleteOne(c, bson.M{"_id": id})
	if err != nil || deleteResult.DeletedCount == 0 {
		return nil, err
	}

	return &users.DeleteUserResponse{}, err
}
