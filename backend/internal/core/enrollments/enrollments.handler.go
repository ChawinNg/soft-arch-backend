package enrollments

import (
	"net/http"

	userService "backend/internal/genproto/users"

	"github.com/gofiber/fiber/v2"
)

type EnrollmentHandler struct {
	service     *EnrollmentService
	userService userService.UserServiceClient
}

func NewEnrollmentHandler(service *EnrollmentService, userService userService.UserServiceClient) *EnrollmentHandler {
	return &EnrollmentHandler{service: service, userService: userService}
}

func (h *EnrollmentHandler) GetUserEnrollment(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	enrollments, err := h.service.GetUserEnrollment(user_id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching Enrollments",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Enrollments retrieved successfully",
		"data":    enrollments,
	})
}

func (h *EnrollmentHandler) GetCourseEnrollment(c *fiber.Ctx) error {
	course_id := c.Params("course_id")

	enrollments, err := h.service.GetUserEnrollment(course_id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching Enrollments",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Enrollments retrieved successfully",
		"data":    enrollments,
	})
}

func (h *EnrollmentHandler) CreateEnrollment(c *fiber.Ctx) error {
	var enrollment Enrollment
	if err := c.BodyParser(&enrollment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	err := h.service.CreateEnrollment(enrollment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create enrollment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "created",
		"message": "Enrollment created successfully",
	})

}

func (h *EnrollmentHandler) EditEnrollment(c *fiber.Ctx) error {
	id := c.Params("id")

	var enrollment Enrollment
	if err := c.BodyParser(&enrollment); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
		})
	}
	enrollment.EnrollmentID = id
	if err := h.service.EditEnrollment(enrollment); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error updating enrollment",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Enrollment updated successfully",
	})
}

func (h *EnrollmentHandler) DeleteEnrollment(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteEnrollment(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error deleting enrollment",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Enrollment deleted successfully",
	})
}

func (h *EnrollmentHandler) SummarizeUserEnrollmentResult(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	enrollments, err := h.service.GetUserEnrollment(user_id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching Enrollments",
		})
	}

	points, err := h.service.SummarizePoints(user_id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error getting user's points",
		})
	}

	for _, enrollment := range enrollments {
		if err := h.service.DeleteEnrollment(enrollment.EnrollmentID); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error deleting enrollment",
			})
		}
	}

	idStr := c.Params("user_id")
	h.userService.ReduceUserPoint(c.Context(), &userService.ReduceUserPointRequest{Id: idStr, ReducePoint: points})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Enrollment deleted and User's Points reduced successfully",
	})
}

func (h *EnrollmentHandler) SummarizeCourseEnrollmentResult(c *fiber.Ctx) error {
	return nil
}
