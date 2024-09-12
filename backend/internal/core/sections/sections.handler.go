package sections

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SectionHandler struct {
    service *SectionService
}

func NewSectionHandler(service *SectionService) *SectionHandler {
    return &SectionHandler{service: service}
}

func (h *SectionHandler) GetAllSections(c *fiber.Ctx) error {
    sections, err := h.service.GetAllSections()
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error fetching sections",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Sections retrieved successfully",
        "data":    sections,
    })
}

func (h *SectionHandler) GetSectionByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid section ID",
        })
    }
    section, err := h.service.GetSectionByID(id)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error fetching section",
        })
    }
    if (section == Section{}) { 
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "status":  "error",
            "message": "Section not found",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Section retrieved successfully",
        "data":    section,
    })
}

func (h *SectionHandler) CreateSection(c *fiber.Ctx) error {
    var section Section
    if err := c.BodyParser(&section); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }
    if err := h.service.CreateSection(section); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error creating section",
        })
    }
    return c.Status(http.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "Section created successfully",
    })
}

func (h *SectionHandler) UpdateSection(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid section ID",
        })
    }
    var section Section
    if err := c.BodyParser(&section); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }
    section.SectionID = id
    if err := h.service.UpdateSection(section); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error updating section",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Section updated successfully",
    })
}

func (h *SectionHandler) DeleteSection(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid section ID",
        })
    }
    if err := h.service.DeleteSection(id); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error deleting section",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Section deleted successfully",
    })
}
