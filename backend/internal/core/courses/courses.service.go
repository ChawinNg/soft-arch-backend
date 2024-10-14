package courses

import (
	"database/sql"
	"log"
	"strings"
)

type CourseService struct {
	db *sql.DB
}

func NewCourseService(db *sql.DB) *CourseService {
	return &CourseService{db: db}
}

func (s *CourseService) GetAllCourses() ([]Course, error) {
	rows, err := s.db.Query("SELECT id, description, course_name, course_full_name, course_type, grading_type, faculty, midterm_exam_date, final_exam_date, credit, course_group_id FROM courses")
	if err != nil {
		log.Println("Error fetching courses:", err)
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.CourseID, &course.Description, &course.CourseName, &course.CourseFullName, &course.CourseType, &course.GradingType, &course.Faculty, &course.MidtermExam, &course.FinalExam, &course.Credit, &course.CourseGroupID); err != nil {
			log.Println("Error scanning course:", err)
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func (s *CourseService) CreateCourse(course Course) error {
	_, err := s.db.Exec("INSERT INTO courses (id, description, course_name, course_full_name, course_type, grading_type, faculty, midterm_exam_date, final_exam_date, credit, course_group_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		course.CourseID, course.Description, course.CourseName, course.CourseFullName, course.CourseType, course.GradingType, course.Faculty, course.MidtermExam, course.FinalExam, course.Credit, course.CourseGroupID)
	if err != nil {
		log.Println("Error creating course:", err)
		return err
	}
	return nil
}

func (s *CourseService) GetCourseByID(id string) (Course, error) {
	var course Course
	err := s.db.QueryRow("SELECT id, description, course_name, course_full_name, course_type, grading_type, faculty, midterm_exam_date, final_exam_date, credit, course_group_id FROM courses WHERE id = ?", id).
		Scan(&course.CourseID, &course.Description, &course.CourseName, &course.CourseFullName, &course.CourseType, &course.GradingType, &course.Faculty, &course.MidtermExam, &course.FinalExam, &course.Credit, &course.CourseGroupID)
	if err == sql.ErrNoRows {
		return course, err
	}
	return course, err
}

func (s *CourseService) UpdateCourse(course Course) error {
	_, err := s.db.Exec("UPDATE courses SET description = ?, course_name = ?, course_full_name = ?, course_type = ?, grading_type = ?, faculty = ?, midterm_exam_date = ?, final_exam_date = ?, credit = ?, course_group_id = ? WHERE id = ?",
		course.Description, course.CourseName, course.CourseFullName, course.CourseType, course.GradingType, course.Faculty, course.MidtermExam, course.FinalExam, course.Credit, course.CourseGroupID, course.CourseID)
	if err != nil {
		log.Println("Error updating course:", err)
		return err
	}
	return nil
}

func (s *CourseService) DeleteCourse(id string) error {
	_, err := s.db.Exec("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting course:", err)
		return err
	}
	return nil
}

func (s *CourseService) GetCoursesPaginated(offset, limit int) ([]Course, int, error) {
    rows, err := s.db.Query(`SELECT id, description, course_name, course_full_name, course_type, grading_type, faculty, midterm_exam_date, final_exam_date, credit, course_group_id 
                              FROM courses LIMIT ? OFFSET ?`, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var courses []Course
    for rows.Next() {
        var course Course
        if err := rows.Scan(&course.CourseID, &course.Description, &course.CourseName, &course.CourseFullName, &course.CourseType, &course.GradingType, &course.Faculty, &course.MidtermExam, &course.FinalExam, &course.Credit, &course.CourseGroupID); err != nil {
            return nil, 0, err
        }
        courses = append(courses, course)
    }

    var totalCourses int
    err = s.db.QueryRow(`SELECT COUNT(*) FROM courses`).Scan(&totalCourses)
    if err != nil {
        return nil, 0, err
    }

    return courses, totalCourses, nil
}

func (s *CourseService) IndexCourses(params map[string]string) ([]Course, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, description, course_name, course_full_name, course_type, grading_type, faculty, midterm_exam_date, final_exam_date, credit, course_group_id FROM courses WHERE 1=1")

	var args []interface{}
	for key, value := range params {
		switch key {
		case "id":
			queryBuilder.WriteString(" AND id LIKE ?")
			args = append(args, "%"+value+"%")
		case "course_name":
			queryBuilder.WriteString(" AND course_name LIKE ?")
			args = append(args, "%"+value+"%")
		case "faculty":
			queryBuilder.WriteString(" AND faculty = ?")
			args = append(args, value)
		case "course_type":
			queryBuilder.WriteString(" AND course_type = ?")
			args = append(args, value)
		case "grading_type":
			queryBuilder.WriteString(" AND grading_type = ?")
			args = append(args, value)
		case "midterm_exam_date":
			queryBuilder.WriteString(" AND midterm_exam_date = ?")
			args = append(args, value)
		case "final_exam_date":
			queryBuilder.WriteString(" AND final_exam_date = ?")
			args = append(args, value)
		}
	}

	rows, err := s.db.Query(queryBuilder.String(), args...)
	if err != nil {
		log.Println("Error indexing courses:", err)
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.CourseID, &course.Description, &course.CourseName, &course.CourseFullName, &course.CourseType, &course.GradingType, &course.Faculty, &course.MidtermExam, &course.FinalExam, &course.Credit, &course.CourseGroupID); err != nil {
			log.Println("Error scanning course:", err)
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}