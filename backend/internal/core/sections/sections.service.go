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
    rows, err := s.db.Query("SELECT section_id, course_id, section, capacity FROM sections")
    if err != nil {
        log.Println("Error fetching sections:", err)
        return nil, err
    }
    defer rows.Close()

    var sections []Section
    for rows.Next() {
        var section Section
        if err := rows.Scan(&section.SectionID, &section.CourseID, &section.Section, &section.Capacity); err != nil {
            log.Println("Error scanning section:", err)
            return nil, err
        }
        sections = append(sections, section)
    }
    return sections, nil
}

func (s *SectionService) GetSectionByID(id int) (Section, error) {
    var section Section
    err := s.db.QueryRow("SELECT section_id, course_id, section, capacity FROM sections WHERE section_id = ?", id).
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
    _, err := s.db.Exec("INSERT INTO sections(course_id, section, capacity) VALUES (?, ?, ?)",
        section.CourseID, section.Section, section.Capacity)
    if err != nil {
        log.Println("Error creating section:", err)
        return err
    }
    return nil
}

func (s *SectionService) UpdateSection(section Section) error {
    _, err := s.db.Exec("UPDATE sections SET course_id = ?, section = ?, capacity = ? WHERE section_id = ?",
        section.CourseID, section.Section, section.Capacity, section.SectionID)
    if err != nil {
        log.Println("Error updating section:", err)
        return err
    }
    return nil
}

func (s *SectionService) DeleteSection(id int) error {
    _, err := s.db.Exec("DELETE FROM sections WHERE section_id = ?", id)
    if err != nil {
        log.Println("Error deleting section:", err)
        return err
    }
    return nil
}
