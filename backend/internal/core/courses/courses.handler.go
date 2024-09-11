package courses

import (
    "encoding/json"
    "strconv"

    "github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
    service *CourseService
}

func NewCourseHandler(service *CourseService) *CourseHandler {
    return &CourseHandler{service: service}
}

func (h *CourseHandler) GetCourses(c *fiber.Ctx) error {
    courses, err := h.service.GetAllCourses()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving courses")
    }

    return c.JSON(courses)
}

func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
    var course Course
    if err := c.BodyParser(&course); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.service.CreateCourse(course); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Error creating course")
    }

    return c.JSON(course)
}

func (h *CourseHandler) GetCourse(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid course ID")
    }

    course, err := h.service.GetCourseByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Course not found")
    }

    return c.JSON(course)
}

func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
    var course Course
    if err := c.BodyParser(&course); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.service.UpdateCourse(course); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Error updating course")
    }

    return c.JSON(course)
}

func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid course ID")
    }

    if err := h.service.DeleteCourse(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Error deleting course")
    }

    return c.SendStatus(fiber.StatusNoContent)
}
