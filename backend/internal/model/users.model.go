package model

import (
	"backend/internal/genproto/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in MongoDB
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Surname  string             `bson:"surname"`
	Email    string             `bson:"email"`
	Password int                `bson:"password"`
}

func FromServiceModel(u *users.User) User {
	return User{
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u User) ToServiceModel() *users.User {
	return &users.User{
		Name:  u.Name,
		Email: u.Email,
	}
}
