package instructors

import (
	"backend/internal/core/sections"
	"backend/internal/model"
	"database/sql"
	"log"
	"net/smtp"
	"os"
)

type Instructor struct {
	sections.Instructor
}

type InstructorService struct {
	db *sql.DB
}

func NewInstructorService(db *sql.DB) *InstructorService {
	return &InstructorService{db: db}
}

func (s *InstructorService) CreateInstructor(instructor Instructor) error {
	_, err := s.db.Exec("INSERT INTO instructors(faculty, full_name, display_name, email, phone_number) VALUES (?, ?, ?, ?, ?)",
		instructor.Faculty, instructor.FullName, instructor.DisplayName, instructor.Email, instructor.PhoneNumber)
	if err != nil {
		log.Println("Error creating instructor:", err)
		return err
	}
	return nil
}

func (s *InstructorService) UpdateInstructor(instructor Instructor) error {
	_, err := s.db.Exec("UPDATE instructors SET faculty = ?, full_name = ?, display_name = ?, email = ?, phone_number = ? WHERE id = ?",
		instructor.Faculty, instructor.FullName, instructor.DisplayName, instructor.Email, instructor.PhoneNumber, instructor.InstructorID)
	if err != nil {
		log.Println("Error updating instructor:", err)
		return err
	}
	return nil
}

func (s *InstructorService) SendEmail(email model.Email) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("EMAIL"),
		os.Getenv("GOOGLE_APP_PASSWORD"),
		"smtp.gmail.com",
	)

	msg := "Subject: From: "+email.FromName+" Subject: "+email.Header+"\n"+email.Body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("EMAIL"),
		[]string{email.ToEmail},
		[]byte(msg),
	)

	if err != nil {
		log.Println("Error sending email :", err)
		return err
	}
	return nil

}


