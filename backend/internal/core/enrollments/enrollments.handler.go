package enrollments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	userService "backend/internal/genproto/users"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

type EnrollmentHandler struct {
	service     *EnrollmentService
	userService userService.UserServiceClient
	rabbitMQ    *amqp.Channel
}

type EnrollmentAction struct {
	Action       string     `json:"action"`
	EnrollmentID string     `json:"id,omitempty"`
	UserID       string     `json:"user_id,omitempty"`
	CourseID     string     `json:"course_id,omitempty"`
	Enrollment   Enrollment `json:"enrollment,omitempty"`
}

type EnrollmentResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    []Enrollment `json:"data,omitempty"`
}

type NumEnrollmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    int64  `json:"data,omitempty"`
}

func NewEnrollmentHandler(service *EnrollmentService, userService userService.UserServiceClient, rabbitMQ *amqp.Channel) *EnrollmentHandler {
	return &EnrollmentHandler{service: service, userService: userService, rabbitMQ: rabbitMQ}
}

func (h *EnrollmentHandler) PublishMessage(action EnrollmentAction) error {
	body, err := json.Marshal(action)
	if err != nil {
		return err
	}

	// Publish the message directly from the handler
	err = h.rabbitMQ.Publish(
		"",                 // Exchange
		"enrollment_queue", // Routing key (queue name)
		false,              // Mandatory
		false,              // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("[X] Sent %s operation", action.Action)
	return nil
}

func (h *EnrollmentHandler) WaitForResponse(responseQueue string) ([]byte, error) {
	// Create a channel to receive the response
	responseChan := make(chan amqp.Delivery, 1)

	// Consume messages from the response queue
	go func() {
		conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v", err)
			return
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %v", err)
			return
		}
		defer ch.Close()

		msgs, err := ch.Consume(
			responseQueue,
			"",    // consumer
			true,  // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			log.Printf("Failed to register a consumer: %v", err)
			return
		}

		for msg := range msgs {
			responseChan <- msg // Send the message to the response channel
			break               // Break after receiving one message
		}
	}()

	// Wait for a response or timeout
	select {
	case msg := <-responseChan:
		return msg.Body, nil
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

func (h *EnrollmentHandler) GetUserEnrollment(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	action := EnrollmentAction{
		Action: "get user enrollment",
		UserID: user_id,
	}

	err := h.PublishMessage(action)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error sending enrollment request to queue",
		})
	}

	response, err := h.WaitForResponse("response_queue")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get response",
		})
	}

	var enrollmentResponse EnrollmentResponse
	if err := json.Unmarshal(response, &enrollmentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unmarshal response",
		})
	}

	return c.JSON(enrollmentResponse)
}

func (h *EnrollmentHandler) GetCourseEnrollment(c *fiber.Ctx) error {
	course_id := c.Params("course_id")

	action := EnrollmentAction{
		Action:   "get course enrollment",
		CourseID: course_id,
	}

	err := h.PublishMessage(action)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error sending course enrollment request to queue",
		})
	}

	response, err := h.WaitForResponse("response_queue")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get response",
		})
	}

	var enrollmentResponse EnrollmentResponse
	if err := json.Unmarshal(response, &enrollmentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unmarshal response",
		})
	}

	return c.JSON(enrollmentResponse)
}

func (h *EnrollmentHandler) CreateEnrollment(c *fiber.Ctx) error {
	var enrollment Enrollment
	if err := c.BodyParser(&enrollment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	action := EnrollmentAction{
		Action:     "create",
		Enrollment: enrollment,
	}

	err := h.PublishMessage(action)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to publish enrollment creation request",
		})
	}

	response, err := h.WaitForResponse("response_queue")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get response",
		})
	}

	var enrollmentResponse NumEnrollmentResponse
	if err := json.Unmarshal(response, &enrollmentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unmarshal response",
		})
	}

	return c.JSON(enrollmentResponse)
}

func (h *EnrollmentHandler) EditEnrollment(c *fiber.Ctx) error {
	id := c.Params("id")

	var enrollment Enrollment
	if err := c.BodyParser(&enrollment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
		})
	}

	enrollment.EnrollmentID = id

	action := EnrollmentAction{
		Action:     "update",
		Enrollment: enrollment,
	}

	err := h.PublishMessage(action)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to publish enrollment update request",
		})
	}

	response, err := h.WaitForResponse("response_queue")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get response",
		})
	}

	var enrollmentResponse EnrollmentResponse
	if err := json.Unmarshal(response, &enrollmentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unmarshal response",
		})
	}

	return c.JSON(enrollmentResponse)
}

func (h *EnrollmentHandler) DeleteEnrollment(c *fiber.Ctx) error {
	id := c.Params("id")

	action := EnrollmentAction{
		Action:       "delete",
		EnrollmentID: id,
	}

	err := h.PublishMessage(action)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to publish enrollment deletion request",
		})
	}

	response, err := h.WaitForResponse("response_queue")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get response",
		})
	}

	var enrollmentResponse EnrollmentResponse
	if err := json.Unmarshal(response, &enrollmentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unmarshal response",
		})
	}

	return c.JSON(enrollmentResponse)
}

func (h *EnrollmentHandler) SummarizeUserEnrollmentResult(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	action := EnrollmentAction{
		Action: "summarize user enrollment",
		UserID: user_id,
	}

	err := h.PublishMessage(action)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to publish enrollment deletion request",
		})
	}

	response, err := h.WaitForResponse("response_queue")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get response",
		})
	}

	var enrollmentResponse NumEnrollmentResponse
	if err := json.Unmarshal(response, &enrollmentResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unmarshal response",
		})
	}

	return c.JSON(enrollmentResponse)
}

func (h *EnrollmentHandler) SummarizeCourseEnrollmentResult(c *fiber.Ctx) error {
	return nil
}
