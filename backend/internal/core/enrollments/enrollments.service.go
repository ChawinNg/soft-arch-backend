package enrollments

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type EnrollmentService struct {
	db *sql.DB
}

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://root:root@%v/",os.Getenv("RABBITMQ_HOST")))

	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

func sendMessage(ch *amqp.Channel, queueName, message string) error {
	_, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return err
}

func NewEnrollmentService(db *sql.DB) *EnrollmentService {
	return &EnrollmentService{
		db: db,
	}
}

func (e *EnrollmentService) GetUserEnrollment(user_id string) ([]Enrollment, error) {
	rows, err := e.db.Query(`
    SELECT e.id, e.user_id, e.course_id, c.course_name, c.credit, e.section_id, e.section, e.points, e.round
    FROM enrollments e
    JOIN courses c ON e.course_id = c.id
    WHERE e.user_id = ?`, user_id)
	if err != nil {
		log.Println("Error fetching Enrollments:", err)
		return nil, err
	}
	defer rows.Close()

	var enrollments []Enrollment

	for rows.Next() {
		var enrollment Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.UserID, &enrollment.CourseID, &enrollment.CourseName, &enrollment.CourseCredit, &enrollment.SectionID, &enrollment.Section, &enrollment.Points, &enrollment.Round); err != nil {
			log.Println("Error scanning Enrollments:", err)
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return enrollments, err
	}
	defer conn.Close()
	defer ch.Close()

	message := fmt.Sprintf("Enrollments for user %s were retrieved", user_id)
	if err := sendMessage(ch, "enrollment", message); err != nil {
		log.Println("Failed to send RabbitMQ message:", err)
	} else {
		log.Printf("[X] Sent %s\n", message)
	}

	return enrollments, nil
}

func (e *EnrollmentService) GetCourseEnrollment(course_id string) ([]Enrollment, error) {
	rows, err := e.db.Query(`
	SELECT e.id, e.user_id, e.course_id, c.course_name, c.credit, e.section_id, e.section, e.points, e.round
	FROM enrollments e
	JOIN courses c ON e.course_id = c.id
	WHERE e.course_id = ?`, course_id)
	if err != nil {
		log.Println("Error fetching Enrollments:", err)
		return nil, err
	}
	defer rows.Close()

	var enrollments []Enrollment

	for rows.Next() {
		var enrollment Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.UserID, &enrollment.CourseID, &enrollment.CourseName, &enrollment.CourseCredit, &enrollment.SectionID, &enrollment.Section, &enrollment.Points, &enrollment.Round); err != nil {
			log.Println("Error scanning Enrollments:", err)
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return enrollments, nil
	}
	defer conn.Close()
	defer ch.Close()

	message := fmt.Sprintf("Course enrollments retrieved for course_id: %s, Enrollments: %v", course_id, enrollments)
	if err := sendMessage(ch, "enrollment", message); err != nil {
		log.Println("Failed to send RabbitMQ message:", err)
	} else {
		log.Printf("[X] Sent %s\n", message)
	}

	return enrollments, nil
}

func (e *EnrollmentService) CreateEnrollment(enrollment Enrollment) (int64, error) {
	var courseName string
	var courseCredit int
	err := e.db.QueryRow("SELECT course_name, credit FROM courses WHERE id = ?", enrollment.CourseID).Scan(&courseName, &courseCredit)
	if err != nil {
		log.Println("Error fetching course info:", err)
		return 0, err
	}

	result, err := e.db.Exec("INSERT INTO enrollments(user_id, course_id, course_name, course_credit, section_id, section, points, round) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		enrollment.UserID, enrollment.CourseID, courseName, courseCredit, enrollment.SectionID, enrollment.Section, enrollment.Points, enrollment.Round)
	if err != nil {
		log.Println("Error creating enrollment:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error fetching last insert ID:", err)
		return 0, err
	}

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return 0, err
	}
	defer conn.Close()
	defer ch.Close()

	message := fmt.Sprintf("New enrollment created: %v", enrollment)
	if err := sendMessage(ch, "enrollment", message); err != nil {
		log.Println("Failed to send RabbitMQ message:", err)
	} else {
		log.Printf("[X] Sent %s\n", message)
	}

	return id, nil
}

func (e *EnrollmentService) EditEnrollment(enrollment Enrollment) error {
	_, err := e.db.Exec("UPDATE enrollments SET user_id = ?, course_id = ?, section_id = ?, section = ?, points = ?, round = ? WHERE id = ?",
		enrollment.UserID, enrollment.CourseID, enrollment.SectionID, enrollment.Section, enrollment.Points, enrollment.Round, enrollment.EnrollmentID)
	if err != nil {
		log.Println("Error updating enrollment:", err)
		return err
	}

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return nil
	}
	defer conn.Close()
	defer ch.Close()

	message := fmt.Sprintf("Enrollment updated: %v", enrollment)
	if err := sendMessage(ch, "enrollment", message); err != nil {
		log.Println("Failed to send RabbitMQ message:", err)
	} else {
		log.Printf("[X] Sent %s\n", message)
	}

	return nil
}

func (e *EnrollmentService) DeleteEnrollment(id string) error {
	_, err := e.db.Exec("DELETE FROM enrollments WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting enrollment:", err)
		return err
	}

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return nil
	}
	defer conn.Close()
	defer ch.Close()

	message := fmt.Sprintf("Enrollment deleted: %v", id)
	if err := sendMessage(ch, "enrollment", message); err != nil {
		log.Println("Failed to send RabbitMQ message:", err)
	} else {
		log.Printf("[X] Sent %s\n", message)
	}

	return nil
}

func (e *EnrollmentService) SummarizePoints(user_id string) (int64, error) {
	var totalPoints int64

	err := e.db.QueryRow(`
        SELECT COALESCE(SUM(points), 0) 
        FROM enrollments 
        WHERE user_id = ?`, user_id).Scan(&totalPoints)
	if err != nil {
		log.Println("Error summarizing points:", err)
		return 0, err
	}

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return 0, err
	}
	defer conn.Close()
	defer ch.Close()

	message := fmt.Sprintf("User %s use total points: %d", user_id, totalPoints)
	if err := sendMessage(ch, "enrollment", message); err != nil {
		log.Println("Failed to send RabbitMQ message:", err)
	} else {
		log.Printf("[X] Sent %s\n", message)
	}

	return totalPoints, nil
}

// func (e *EnrollmentService) SummarizeCourseEnrollmentResult(course_id int) error {
// 	enrollments, err := e.GetCourseEnrollment(course_id)
// 	if err != nil {
// 		log.Println("Error fetching Enrollments:", err)
// 		return err
// 	}

// 	for _, enrollment := range enrollments {

// 	}

// 	return nil
// }
