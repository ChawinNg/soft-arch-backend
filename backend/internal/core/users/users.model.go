package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name  string             `json:"name" bson:"name"`
    Surname string           `json:"surname" bson:"surname"`
    Email string             `json:"email" bson:"email"`
    Password string          `json:"password" bson:"password"`
}
