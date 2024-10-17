package enrollments

import (
	"database/sql"
	"log"
)

type EnrollmentService struct {
	db *sql.DB
}

func NewEnrollmentService(db *sql.DB) *EnrollmentService {
	return &EnrollmentService{
		db: db,
	}
}

func (e *EnrollmentService) GetUserEnrollment(user_id string) ([]Enrollment, error) {
	rows, err := e.db.Query("SELECT id, user_id, course_id, section_id, section, points, round FROM enrollments WHERE user_id = ?", user_id)
	if err != nil {
		log.Println("Error fetching Enrollments:", err)
		return nil, err
	}
	defer rows.Close()

	var enrollments []Enrollment

	for rows.Next() {
		var enrollment Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.UserID, &enrollment.CourseID, &enrollment.SectionID, &enrollment.Section, &enrollment.Points, &enrollment.Round); err != nil {
			log.Println("Error scanning Enrollments:", err)
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

func (e *EnrollmentService) GetCourseEnrollment(course_id string) ([]Enrollment, error) {
	rows, err := e.db.Query("SELECT id, user_id, course_id, section_id, section, points, round FROM enrollments WHERE course_id = ?", course_id)
	if err != nil {
		log.Println("Error fetching Enrollments:", err)
		return nil, err
	}
	defer rows.Close()

	var enrollments []Enrollment

	for rows.Next() {
		var enrollment Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.UserID, &enrollment.CourseID, &enrollment.SectionID, &enrollment.Section, &enrollment.Points, &enrollment.Round); err != nil {
			log.Println("Error scanning Enrollments:", err)
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

func (e *EnrollmentService) CreateEnrollment(enrollment Enrollment) (int64, error) {
	result, err := e.db.Exec("INSERT INTO enrollments(user_id, course_id, section_id, section, points, round) VALUES ( ?, ?, ?, ?, ?, ?)",
		enrollment.UserID, enrollment.CourseID, enrollment.SectionID, enrollment.Section, enrollment.Points, enrollment.Round)
	if err != nil {
		log.Println("Error creating enrollment:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error fetching last insert ID:", err)
		return 0, err
	}

	return id, err
}

func (e *EnrollmentService) EditEnrollment(enrollment Enrollment) error {
	_, err := e.db.Exec("UPDATE enrollments SET user_id = ?, course_id = ?, section_id = ?, section = ?, points = ?, round = ? WHERE id = ?",
		enrollment.UserID, enrollment.CourseID, enrollment.SectionID, enrollment.Section, enrollment.Points, enrollment.Round, enrollment.EnrollmentID)
	if err != nil {
		log.Println("Error updating enrollment:", err)
		return err
	}

	return nil
}

func (e *EnrollmentService) DeleteEnrollment(id string) error {
	_, err := e.db.Exec("DELETE FROM enrollments WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting enrollment:", err)
		return err
	}
	return nil
}

func (e *EnrollmentService) SummarizePoints(user_id string) (int64, error) {
	var totalPoints int64
	totalPoints = 0

	enrollments, err := e.GetUserEnrollment(user_id)
	if err != nil {
		log.Println("Error fetching Enrollments:", err)
		return 0, err
	}

	for _, enrollment := range enrollments {
		totalPoints += enrollment.Points
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
