package sections

import (
	"database/sql"
	"log"
)

type SectionService struct {
	db *sql.DB
}

func NewSectionService(db *sql.DB) *SectionService {
	return &SectionService{db: db}
}

func (s *SectionService) GetAllSections() ([]Section, error) {
	rows, err := s.db.Query("SELECT id, course_id, section, capacity, room FROM sections")
	if err != nil {
		log.Println("Error fetching sections:", err)
		return nil, err
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var section Section
		if err := rows.Scan(&section.SectionID, &section.CourseID, &section.Section, &section.Capacity, &section.Room); err != nil {
			log.Println("Error scanning section:", err)
			return nil, err
		}
		sections = append(sections, section)
	}
	return sections, nil
}

func (s *SectionService) GetSectionByID(id int) (Section, error) {
	var section Section
	err := s.db.QueryRow("SELECT id, course_id, section, capacity, room FROM sections WHERE id = ?", id).
		Scan(&section.SectionID, &section.CourseID, &section.Section, &section.Capacity)
	if err == sql.ErrNoRows {
		return section, nil // No result found
	} else if err != nil {
		log.Println("Error fetching section by ID:", err)
		return section, err
	}
	return section, nil
}

func (s *SectionService) CreateSection(section Section) error {
	_, err := s.db.Exec("INSERT INTO sections(course_id, section, capacity, room) VALUES (?, ?, ?, ?)",
		section.CourseID, section.Section, section.Capacity, section.Room)
	if err != nil {
		log.Println("Error creating section:", err)
		return err
	}
	return nil
}

func (s *SectionService) UpdateSection(section Section) error {
	_, err := s.db.Exec("UPDATE sections SET course_id = ?, section = ?, capacity = ? WHERE id = ?",
		section.CourseID, section.Section, section.Capacity, section.SectionID)
	if err != nil {
		log.Println("Error updating section:", err)
		return err
	}
	return nil
}

func (s *SectionService) DeleteSection(id int) error {
	_, err := s.db.Exec("DELETE FROM sections WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting section:", err)
		return err
	}
	return nil
}

func (s *SectionService) GetSectionsByCourseID(id string) ([]Section, error) {
	query := `
    SELECT
        s.id AS section_id,
        s.course_id,
        s.section,
        s.room,
        s.capacity
    FROM
        sections s
    WHERE
        s.course_id = ?`

	rows, err := s.db.Query(query, id)
	if err != nil {
		log.Println("Error querying sections:", err)
		return nil, err
	}
	defer rows.Close()

	var sections []Section

	for rows.Next() {
		var section Section
		if err := rows.Scan(&section.SectionID, &section.CourseID, &section.Section, &section.Room, &section.Capacity); err != nil {
			log.Println("Error scanning section row:", err)
			return nil, err
		}
		section.Timeslots, err = s.getTimeSlotsForSection(section.SectionID)
		section.Instructors, err = s.getInstructorsForSection(section.SectionID)
		sections = append(sections, section)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
		return nil, err
	}

	return sections, nil
}

func (s *SectionService) getInstructorsForSection(sectionID int) ([]Instructor, error) {
	rows, err := s.db.Query(`SELECT instructor_id FROM section_instructors WHERE section_id = ?`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructors []Instructor

	for rows.Next() {
		var instructorID int
		if err := rows.Scan(&instructorID); err != nil {
			return nil, err
		}

		instructorRow := s.db.QueryRow(`SELECT id, full_name, faculty, display_name, email, phone_number FROM instructors WHERE id = ?`, instructorID)

		var instructor Instructor
		if err := instructorRow.Scan(&instructor.InstructorID, &instructor.FullName, &instructor.Faculty, &instructor.DisplayName, &instructor.Email, &instructor.PhoneNumber); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		// email := ""
        // if instructor.Email != nil {
        //     email = *instructor.Email
        // }
        
        // phoneNumber := ""
        // if instructor.PhoneNumber != nil {
        //     phoneNumber = *instructor.PhoneNumber
        // }

		instructors = append(instructors, instructor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return instructors, nil
}

func (s *SectionService) getTimeSlotsForSection(sectionID int) ([][]string, error) {
	query := `
    SELECT
        time
    FROM
        timeslots
    WHERE
        section_id = ?`

	rows, err := s.db.Query(query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeSlots [][]string

	for rows.Next() {
		var timeSlot string
		if err := rows.Scan(&timeSlot); err != nil {
			return nil, err
		}

		timeSlots = append(timeSlots, []string{timeSlot})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timeSlots, nil
}
