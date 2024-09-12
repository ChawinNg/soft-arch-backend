package users

import (
	"backend/internal/genproto/users"
	"backend/internal/model"
	"net/http"

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
	return c.JSON(user)
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUser(c.Context(), nil)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(users)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	u := &model.User{}
	err := c.BodyParser(u)
	if err != nil {
		return c.Status(500).SendString("Invalid input")
	}

	_, err = h.service.CreateUser(c.Context(), &users.CreateUserRequest{
		Sid:      u.Sid,
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(201)
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
	return c.JSON(user)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	user, err := h.service.DeleteUser(c.Context(), &users.DeleteUserRequest{Id: idParam})
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}
	return c.JSON(user)
}
