package courses

import (
	"backend/internal/core/sections"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	service        *CourseService
	sectionService *sections.SectionService
}

type LocalSection struct {
	SectionID   int                   `json:"id"`
	CourseID    string                `json:"courseId"`
	Section     int                   `json:"section"`
	Capacity    int                   `json:"capacity"`
	MaxCapacity int                   `json:"max_capacity"`
	Room        *string               `json:"room"`
	Timeslots   []string              `json:"timeslots"`
	Instructors []sections.Instructor `json:"instructors"`
}

func NewCourseHandler(service *CourseService, sectionService *sections.SectionService) *CourseHandler {
	return &CourseHandler{
		service:        service,
		sectionService: sectionService,
	}
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
	id := c.Params("id")
	course, err := h.service.GetCourseByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Course not found")
	}

	return c.JSON(course)
}

func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	var course Course
	if err := c.BodyParser(&course); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	course.CourseID = id

	if err := h.service.UpdateCourse(course); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating course")
	}

	return c.JSON(course)
}

func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteCourse(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting course")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CourseHandler) GetCoursesPaginated(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page size")
	}

	offset := (page - 1) * pageSize
	courses, totalCourses, err := h.service.GetCoursesPaginated(offset, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving courses")
	}

	var coursesWithSections []struct {
		Course   Course         `json:"course"`
		Sections []LocalSection `json:"sections"`
	}

	for _, course := range courses {
		sections, err := h.sectionService.GetSectionsByCourseID(course.CourseID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving sections for course")
		}

		var localSections []LocalSection
		for _, section := range sections {
			localSections = append(localSections, LocalSection{
				SectionID:   section.SectionID,
				CourseID:    section.CourseID,
				Section:     section.Section,
				MaxCapacity: section.MaxCapacity,
				Capacity:    section.Capacity,
				Room:        section.Room,
				Timeslots:   section.Timeslots,
				Instructors: section.Instructors,
			})
		}

		coursesWithSections = append(coursesWithSections, struct {
			Course   Course         `json:"course"`
			Sections []LocalSection `json:"sections"`
		}{
			Course:   course,
			Sections: localSections,
		})
	}

	return c.JSON(fiber.Map{
		"totalCourses": totalCourses,
		"courses":      coursesWithSections,
	})
}
