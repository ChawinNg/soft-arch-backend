package users

import (
	"backend/internal/genproto/users"
	"backend/internal/model"
	"backend/internal/utils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service users.UserServiceClient
}

func NewHandler(service users.UserServiceClient) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	user, err := h.service.GetUser(c.Context(), &users.GetUserRequest{Id: idParam})
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}
	return c.JSON(fiber.Map{
		"success": true,
		"user":    user.User,
	})
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUser(c.Context(), nil)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users.User,
	})
}

func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	u := &model.User{}
	err := c.BodyParser(u)
	if err != nil {
		return c.Status(500).SendString("Invalid input")
	}

	user, err := h.service.RegisterUser(c.Context(), &users.RegisterUserRequest{
		Sid:      u.Sid,
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user.Id,
	})
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	u := &model.User{}
	err := c.BodyParser(u)
	if err != nil {
		return c.Status(500).SendString("Invalid input")
	}

	user, err := h.service.UpdateUser(c.Context(), &users.UpdateUserRequest{
		User: &users.User{
			Id:       idParam,
			Sid:      u.Sid,
			Name:     u.Name,
			Surname:  u.Surname,
			Email:    u.Email,
			Password: u.Password,
		},
	})
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}
	return c.JSON(fiber.Map{
		"success": true,
		"user":    user.User,
	})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	_, err := h.service.DeleteUser(c.Context(), &users.DeleteUserRequest{Id: idParam})
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    "",
	})
}

func (h *Handler) LoginUser(c *fiber.Ctx) error {
	u := &model.User{}
	err := c.BodyParser(u)
	if err != nil {
		return c.Status(500).SendString("Invalid input")
	}

	token, err := h.service.LoginUser(c.Context(), &users.LoginRequest{
		Sid:      u.Sid,
		Password: u.Password,
	})
	if err != nil {
		return c.Status(500).SendString("Invalid student id or password")
	}

	session_expire, err := strconv.Atoi(os.Getenv("SESSION_EXPIRE"))
	if err != nil {
		log.Fatalf("Error converting SISSION_EXPIRE to int: %v", err)
	}

	c.Cookie(utils.CreateSessionCookie(token.Token, session_expire))
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token.Token,
	})
}

func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	session := c.Locals("session").(model.Sessions)
	user, err := h.service.GetUser(c.Context(), &users.GetUserRequest{
		Id: session.UserId,
	})
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user.User,
	})
}
