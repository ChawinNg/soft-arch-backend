package enrollments

import (
	"backend/internal/core/sections"
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
	rows, err := e.db.Query(`
    SELECT e.id, e.user_id, e.course_id, c.course_name, c.credit, e.section_id, e.section, e.points, e.round
    FROM enrollments e
    JOIN courses c ON e.course_id = c.id
    WHERE e.user_id = ? AND e.summarized = FALSE`, user_id)
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

	return enrollments, nil
}

func (e *EnrollmentService) GetUserEnrollmentResult(user_id string) ([]EnrollmentSummary, error) {
	rows, err := e.db.Query(`
    SELECT user_id, course_id, course_name, course_credit, 
		 section_id, section, round, points, result
	FROM enrollment_results
    WHERE user_id = ?`, user_id)
	if err != nil {
		log.Println("Error fetching Enrollment results:", err)
		return nil, err
	}
	defer rows.Close()

	var enrollments []EnrollmentSummary
	for rows.Next() {
		var enrollment EnrollmentSummary
		if err := rows.Scan(&enrollment.UserID, &enrollment.CourseID, &enrollment.CourseName, &enrollment.CourseCredit,
			&enrollment.SectionID, &enrollment.Section,
			&enrollment.Round, &enrollment.Points, &enrollment.Result); err != nil {
			log.Println("Error scanning Enrollment results:", err)
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
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

	return enrollments, nil
}

func (e *EnrollmentService) CreateEnrollment(enrollment Enrollment) (int64, error) {
	var courseName string
	var courseCredit int

	var existingID int64
	err := e.db.QueryRow("SELECT id FROM enrollments WHERE user_id = ? AND course_id = ? AND section_id = ? AND summarized = FALSE",
		enrollment.UserID, enrollment.CourseID, enrollment.SectionID).Scan(&existingID)
	if err == nil {
		log.Println("enrollment already exists for this user, course, and section")
		return 0, err
	} else if err != sql.ErrNoRows {
		log.Println("Error checking existing enrollment:", err)
		return 0, err
	}

	err = e.db.QueryRow("SELECT course_name, credit FROM courses WHERE id = ?", enrollment.CourseID).Scan(&courseName, &courseCredit)
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

	return id, nil
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

func (e *EnrollmentService) SummarizePoints(user_id string) (int32, error) {
	var totalPoints int32

	err := e.db.QueryRow(`
        SELECT COALESCE(SUM(points), 0) 
        FROM enrollments 
        WHERE user_id = ?`, user_id).Scan(&totalPoints)
	if err != nil {
		log.Println("Error summarizing points:", err)
		return 0, err
	}

	return totalPoints, nil
}

func InsertEnrollmentResult(db *sql.DB, enrollment EnrollmentSummary) error {
	query := `
		INSERT INTO enrollment_results (
			user_id, course_id, course_name, course_credit, section_id, section, round, points, result
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		enrollment.UserID,
		enrollment.CourseID,
		enrollment.CourseName,
		enrollment.CourseCredit,
		enrollment.SectionID,
		enrollment.Section,
		enrollment.Round,
		enrollment.Points,
		enrollment.Result,
	)

	if err != nil {
		log.Println("Error inserting into enrollment_results:", err)
		return err // Return the error if insert fails
	}

	return nil // Return nil if all inserts succeed
}

func (s *EnrollmentService) SummarizeCourseEnrollmentResult(round string) ([]EnrollmentSummary, []sections.Section, error) {
	query := `
		SELECT  e.user_id, e.course_id, e.course_name,e.course_credit,s.max_capacity,s.id, e.section,e.round, e.points,s.capacity
		FROM enrollments e
		INNER JOIN sections s ON e.section_id = s.id
		WHERE e.round = ? AND e.summarized = FALSE
		ORDER BY e.course_id, e.section, e.points DESC
	`
	rows, err := s.db.Query(query, round)
	if err != nil {
		log.Println("failed to query enrollments: %v", err)
		return nil, nil, err
	}
	defer rows.Close()

	var enrollments []EnrollmentSummary
	var SectionToUpdates []sections.Section
	var prevCourseID string
	var availableCapacity, prevSectionID, prevSection, prevMaxCapa int
	prevSection = 0

	for rows.Next() {
		var enrollment EnrollmentSummary
		if err := rows.Scan(
			&enrollment.UserID, &enrollment.CourseID, &enrollment.CourseName, &enrollment.CourseCredit,
			&enrollment.MaxCapacity, &enrollment.SectionID, &enrollment.Section,
			&enrollment.Round, &enrollment.Points, &enrollment.Capacity); err != nil {
			log.Println("Error scanning enrollments:", err)
			return nil, nil, err
		}
		// the code below assumed initial availableCapacity is ALWAYS > 0
		// if first row
		if prevCourseID == "" && prevSection == 0 && enrollment.MaxCapacity-enrollment.Capacity > 0 {
			prevCourseID = enrollment.CourseID
			prevSection = enrollment.Section
			prevSectionID = enrollment.SectionID
			prevMaxCapa = enrollment.MaxCapacity
			availableCapacity = enrollment.MaxCapacity - enrollment.Capacity - 1

			enrollment.Result = true

			enrollments = append(enrollments, enrollment)
			err := InsertEnrollmentResult(s.db, enrollment)
			if err != nil {
				return nil, nil, err
			}
			continue
		}

		// diff course or diff section
		if enrollment.CourseID != prevCourseID || enrollment.Section != prevSection {
			if availableCapacity < 0 {
				availableCapacity = 0
			}
			//add section to update
			var SectionToUpdate sections.Section
			SectionToUpdate.CourseID = prevCourseID
			SectionToUpdate.Section = prevSection
			SectionToUpdate.MaxCapacity = prevMaxCapa
			SectionToUpdate.Capacity = prevMaxCapa - availableCapacity
			SectionToUpdate.SectionID = prevSectionID
			SectionToUpdates = append(SectionToUpdates, SectionToUpdate)
			//set info of the new one
			//first one is always a success reg user
			prevCourseID = enrollment.CourseID
			prevSection = enrollment.Section
			prevSectionID = enrollment.SectionID
			prevMaxCapa = enrollment.MaxCapacity
			availableCapacity = enrollment.MaxCapacity - enrollment.Capacity - 1

			enrollment.Result = true

			enrollments = append(enrollments, enrollment)
			err := InsertEnrollmentResult(s.db, enrollment)
			if err != nil {
				return nil, nil, err
			}
			continue
		} else {
			//full capacity
			if availableCapacity <= 0 {
				//fail reg user
				enrollment.Result = false

			} else {
				//success reg user
				enrollment.Result = true
			}
		}

		enrollments = append(enrollments, enrollment)
		err := InsertEnrollmentResult(s.db, enrollment)
		if err != nil {
			return nil, nil, err
		}

		availableCapacity--

	}
	if len(enrollments) == 0 {
		return enrollments, nil, nil
	}

	lastEnrollment := enrollments[len(enrollments)-1]
	if len(SectionToUpdates) == 0 {
		if availableCapacity < 0 {
			availableCapacity = 0
		}
		var SectionToUpdate sections.Section
		SectionToUpdate.CourseID = prevCourseID
		SectionToUpdate.Section = prevSection
		SectionToUpdate.MaxCapacity = prevMaxCapa
		SectionToUpdate.Capacity = prevMaxCapa - availableCapacity
		SectionToUpdate.SectionID = prevSectionID
		SectionToUpdates = append(SectionToUpdates, SectionToUpdate)
	} else {
		lastSectionToUpdate := SectionToUpdates[len(SectionToUpdates)-1]
		if (lastEnrollment.CourseID != lastSectionToUpdate.CourseID) || (lastEnrollment.Section != lastSectionToUpdate.Section) {
			var SectionToUpdate sections.Section
			SectionToUpdate.CourseID = prevCourseID
			SectionToUpdate.Section = prevSection
			SectionToUpdate.MaxCapacity = prevMaxCapa
			SectionToUpdate.Capacity = prevMaxCapa - availableCapacity
			SectionToUpdate.SectionID = prevSectionID
			SectionToUpdates = append(SectionToUpdates, SectionToUpdate)
		}
	}

	if err = rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return nil, nil, err
	}

	//set summarized
	_, err2 := s.db.Exec("UPDATE enrollments SET summarized = TRUE WHERE round = ?", round)

	if err2 != nil {
		log.Println("Error updating enrollments:", err2)
		return nil, nil, err2
	}

	return enrollments, SectionToUpdates, nil
}
