package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
    collection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
    return &UserService{
        collection: db.Collection("users"),
    }
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var user User
    filter := bson.M{"_id": objID}
    err = s.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}
