package users

import (
	"backend/internal/genproto/users"
	"backend/internal/model"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service users.UserServiceClient
}

func NewHandler(service users.UserServiceClient) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	// users, err := h.service.GetAllUser(c.Context(), nil)
	// if err != nil {
	// 	return c.Status(500).SendString(err.Error())
	// }

	return c.JSON("hello")
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	u := &model.User{}
	err := c.BodyParser(u)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_, err = h.service.CreateUser(c.Context(), &users.CreateUserRequest{
		User: u.ToServiceModel(),
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(201)
}
