package users

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    
    user, err := h.service.GetUserByID(context.Background(), id)
    if err != nil {
        log.Println("Failed to fetch user:", err)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }
    
    return c.JSON(user)
}
