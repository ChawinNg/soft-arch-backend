package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"backend/internal/core/enrollments"
	userService "backend/internal/genproto/users"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/streadway/amqp"
)

type EnrollmentAction struct {
	Action       string                 `json:"action"`
	EnrollmentID string                 `json:"id,omitempty"`
	UserID       string                 `json:"user_id,omitempty"`
	CourseID     string                 `json:"course_id,omitempty"`
	Enrollment   enrollments.Enrollment `json:"enrollment,omitempty"`
}

type EnrollmentResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []enrollments.Enrollment `json:"data,omitempty"`
}

type NumEnrollmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    int64  `json:"data,omitempty"`
}

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://root:root@%v/", os.Getenv("RABBITMQ_HOST")))
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

func main() {
	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		"enrollment_queue", // queue name
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	grpc_host := os.Getenv("GRPC_SERVER_HOST")
	grpcConn, err := grpc.NewClient(grpc_host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userConn := userService.NewUserServiceClient(grpcConn)
	
	sqlDSN := os.Getenv("SQL_DB_DSN")
	dbSQL, err := sql.Open("mysql", sqlDSN)

	if err != nil {
		log.Fatal(err)
	}
	defer dbSQL.Close()
	enrollmentService := enrollments.NewEnrollmentService(dbSQL)

	go func() {
		for d := range msgs {
			var action EnrollmentAction
			err := json.Unmarshal(d.Body, &action)

			if err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			if action.Action == "get user enrollment" {
				enrollments, err := enrollmentService.GetUserEnrollment(action.UserID)
				var response EnrollmentResponse
				if err != nil {
					response = EnrollmentResponse{
						Status:  "error",
						Message: "Error fetching enrollments",
						Data:    nil,
					}
				} else {
					response = EnrollmentResponse{
						Status:  "success",
						Message: "Enrollments retrieved successfully",
						Data:    enrollments,
					}
				}
				if err != nil {
					log.Printf("Error fetching enrollments for user %s: %v", action.UserID, err)
				} else if len(enrollments) == 0 {
					log.Printf("No enrollments found for user %s", action.UserID)
				} else {
					log.Printf("Enrollments for user %s: %v", action.UserID, enrollments)
				}
				responseBody, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				err = ch.Publish(
					"",               // exchange
					"response_queue", // response queue
					false,            // mandatory
					false,            // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        responseBody,
					})
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				}
			} else if action.Action == "get course enrollment" {
				enrollments, err := enrollmentService.GetCourseEnrollment(action.CourseID)
				var response EnrollmentResponse
				if err != nil {
					response = EnrollmentResponse{
						Status:  "error",
						Message: "Error fetching enrollments",
						Data:    nil,
					}
				} else {
					response = EnrollmentResponse{
						Status:  "success",
						Message: "Enrollments retrieved successfully",
						Data:    enrollments,
					}
				}
				if err != nil {
					log.Printf("Error fetching enrollments for course %s: %v", action.CourseID, err)
				} else if len(enrollments) == 0 {
					log.Printf("No enrollments found for course %s", action.CourseID)
				} else {
					log.Printf("Enrollments for course %s: %v", action.CourseID, enrollments)
				}
				responseBody, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				err = ch.Publish(
					"",               // exchange
					"response_queue", // response queue
					false,            // mandatory
					false,            // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        responseBody,
					})
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				}
			} else if action.Action == "create" {
				enrollment := action.Enrollment
				id, err := enrollmentService.CreateEnrollment(enrollment)
				if err != nil {
					log.Printf("Error creating enrollment: %v", err)
				} else {
					log.Printf("Enrollment created with ID: %v", id)
				}
				var response NumEnrollmentResponse
				if err != nil {
					response = NumEnrollmentResponse{
						Status:  "error",
						Message: "Error creating enrollments",
					}
				} else {
					response = NumEnrollmentResponse{
						Status:  "success",
						Message: "Enrollments created successfully",
						Data:    id,
					}
				}
				responseBody, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				err = ch.Publish(
					"",               // exchange
					"response_queue", // response queue
					false,            // mandatory
					false,            // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        responseBody,
					})
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				}
			} else if action.Action == "update" {
				enrollment := action.Enrollment
				err := enrollmentService.EditEnrollment(enrollment)
				if err != nil {
					log.Printf("Error updating enrollment with ID %s: %v", enrollment.EnrollmentID, err)
				} else {
					log.Printf("Enrollment with ID %s updated successfully", enrollment.EnrollmentID)
				}
				var response EnrollmentResponse
				if err != nil {
					response = EnrollmentResponse{
						Status:  "error",
						Message: "Error updating enrollments",
					}
				} else {
					response = EnrollmentResponse{
						Status:  "success",
						Message: "Enrollments updated successfully",
					}
				}
				responseBody, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				err = ch.Publish(
					"",               // exchange
					"response_queue", // response queue
					false,            // mandatory
					false,            // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        responseBody,
					})
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				}
			} else if action.Action == "delete" {
				enrollment_id := action.EnrollmentID
				err := enrollmentService.DeleteEnrollment(enrollment_id)
				if err != nil {
					log.Printf("Error deleting enrollment with ID %s: %v", enrollment_id, err)
				} else {
					log.Printf("Enrollment with ID %s deleted successfully", enrollment_id)
				}
				var response EnrollmentResponse
				if err != nil {
					response = EnrollmentResponse{
						Status:  "error",
						Message: "Error deleting enrollments",
					}
				} else {
					response = EnrollmentResponse{
						Status:  "success",
						Message: "Enrollments deleted successfully",
					}
				}
				responseBody, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				err = ch.Publish(
					"",               // exchange
					"response_queue", // response queue
					false,            // mandatory
					false,            // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        responseBody,
					})
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				}
			} else if action.Action == "summarize user enrollment" {
				var points int64
				enrollments, err := enrollmentService.GetUserEnrollment(action.UserID)
				if err != nil {
					log.Printf("Error getting user enrollment with user id %s: %v", action.UserID, err)
				}
				points, err = enrollmentService.SummarizePoints(action.UserID)
				if err != nil {
					log.Printf("Error summarizing points for user %s: %v", action.UserID, err)
				}
				for _, enrollment := range enrollments {
					err := enrollmentService.DeleteEnrollment(enrollment.EnrollmentID)
					if err != nil {
						log.Printf("Error deleting enrollment with ID %s: %v", enrollment.EnrollmentID, err)
					}
				}
				userConn.ReduceUserPoint(context.Background(), &userService.ReduceUserPointRequest{
					Id:          action.UserID,
					ReducePoint: points,
				})
				log.Printf("Summarize user with user ID %s successfully. Reduced %v points", action.UserID, points)
				var response NumEnrollmentResponse
				if err != nil {
					response = NumEnrollmentResponse{
						Status:  "error",
						Message: "Error summarizing user enrollments",
					}
				} else {
					response = NumEnrollmentResponse{
						Status:  "success",
						Message: "User enrollments summarized successfully",
						Data:    points,
					}
				}
				responseBody, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				err = ch.Publish(
					"",               // exchange
					"response_queue", // response queue
					false,            // mandatory
					false,            // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        responseBody,
					})
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				}
			}
		}
	}()

	log.Println("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
