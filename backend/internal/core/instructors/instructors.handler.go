package instructors

import (
	"net/http"
	"strconv"
	"github.com/gofiber/fiber/v2"

    "backend/internal/model"
)

type InstructorHandler struct {
    service *InstructorService  
}

func NewInstructorHandler(service *InstructorService) *InstructorHandler {
    return &InstructorHandler{service: service}
}


func (h *InstructorHandler) CreateInstructor(c *fiber.Ctx) error {
    var instructor Instructor
    if err := c.BodyParser(&instructor); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }
    if err := h.service.CreateInstructor(instructor); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error creating instructor",
        })
    }
    return c.Status(http.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "Instructor created successfully",
    })
}

func (h *InstructorHandler) UpdateInstructor(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid instructor ID",
        })
    }
	var instructor Instructor
    if err := c.BodyParser(&instructor); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }
    instructor.InstructorID = id
    if err := h.service.UpdateInstructor(instructor); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error updating instructor",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Instructor updated successfully",
		"data": instructor,
    })
}

func (h *InstructorHandler) SendEmail(c *fiber.Ctx) error {
	var email model.Email
    if err := c.BodyParser(&email); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }
    if err := h.service.SendEmail(email); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error sending email",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": email,
    })
}