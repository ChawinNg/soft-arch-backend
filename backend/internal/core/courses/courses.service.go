package courses

import (
    "backend/internal/database"
    "context"
    "database/sql"
    "log"
)

type CourseService struct {
    db *sql.DB
}

func NewCourseService(db *sql.DB) *CourseService {
    return &CourseService{db: db}
}

func (s *CourseService) GetAllCourses() ([]Course, error) {
    rows, err := s.db.Query("SELECT courseId, description, courseType, courseGroupId FROM courses")
    if err != nil {
        log.Println("Error fetching courses:", err)
        return nil, err
    }
    defer rows.Close()

    var courses []Course
    for rows.Next() {
        var course Course
        if err := rows.Scan(&course.CourseID, &course.Description, &course.CourseType, &course.CourseGroupID); err != nil {
            log.Println("Error scanning course:", err)
            return nil, err
        }
        courses = append(courses, course)
    }
    return courses, nil
}

func (s *CourseService) CreateCourse(course Course) error {
    _, err := s.db.Exec("INSERT INTO courses(courseId, description, courseType, courseGroupId) VALUES (?, ?, ?, ?)",
        course.CourseID, course.Description, course.CourseType, course.CourseGroupID)
    if err != nil {
        log.Println("Error creating course:", err)
        return err
    }
    return nil
}

func (s *CourseService) GetCourseByID(id int) (Course, error) {
    var course Course
    err := s.db.QueryRow("SELECT courseId, description, courseType, courseGroupId FROM courses WHERE courseId = ?", id).
        Scan(&course.CourseID, &course.Description, &course.CourseType, &course.CourseGroupID)
    if err == sql.ErrNoRows {
        return course, err
    }
    return course, err
}

func (s *CourseService) UpdateCourse(course Course) error {
    _, err := s.db.Exec("UPDATE courses SET description = ?, courseType = ?, courseGroupId = ? WHERE courseId = ?",
        course.Description, course.CourseType, course.CourseGroupID, course.CourseID)
    if err != nil {
        log.Println("Error updating course:", err)
        return err
    }
    return nil
}

func (s *CourseService) DeleteCourse(id int) error {
    _, err := s.db.Exec("DELETE FROM courses WHERE courseId = ?", id)
    if err != nil {
        log.Println("Error deleting course:", err)
        return err
    }
    return nil
}
