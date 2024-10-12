package instructors

import (
	"backend/internal/core/sections"
	"database/sql"
	"log"
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
