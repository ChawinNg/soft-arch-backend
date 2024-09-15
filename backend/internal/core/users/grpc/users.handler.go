package users

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"backend/internal/genproto/users"
	"backend/internal/model"
	"backend/internal/utils"

	_ "github.com/joho/godotenv/autoload"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

// Hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare hashed password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

func (h *Handler) RegisterUser(c context.Context, u *users.RegisterUserRequest) (*users.RegisterUserResponse, error) {
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	reg_user := bson.M{
		"sid":      u.Sid,
		"name":     u.Name,
		"surname":  u.Surname,
		"email":    u.Email,
		"password": hashedPassword,
	}

	user, err := h.db.InsertOne(c, reg_user)
	if err != nil {
		return nil, err
	}

	return &users.RegisterUserResponse{Id: user.InsertedID.(primitive.ObjectID).Hex()}, err
}

func (h *Handler) UpdateUser(c context.Context, u *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(u.User.Id)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(u.User.Password)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			// "name":     u.User.Name,
			// "surname":  u.User.Surname,
			// "email":    u.User.Email,
			"password": hashedPassword,
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

func (h *Handler) LoginUser(c context.Context, u *users.LoginRequest) (*users.LoginResponse, error) {
	var user model.User
	err := h.db.FindOne(c, bson.M{"sid": u.Sid}).Decode(&user)
	if err != nil {
		return nil, err
	}

	grpcUser, err := model.ConvertMongoToGrpc(user)
	if err != nil {
		return nil, err
	}

	passwordHash := grpcUser.Password
	if !checkPasswordHash(u.Password, passwordHash) {
		return nil, errors.New("invalid email or password")
	}

	session := model.Sessions{
		UserId: grpcUser.Id,
		Email:  grpcUser.Email,
	}

	session_expire, err := strconv.Atoi(os.Getenv("SESSION_EXPIRE"))
	if err != nil {
		log.Fatalf("Error converting SISSION_EXPIRE to int: %v", err)
	}

	jwt_secret := os.Getenv("JWT_SECRET")

	token, err := utils.CreateJwtToken(session, time.Duration(session_expire*int(time.Second)), jwt_secret)
	if err != nil {
		return nil, errors.New("cannot create jwt token")
	}

	return &users.LoginResponse{Token: token}, nil
}
