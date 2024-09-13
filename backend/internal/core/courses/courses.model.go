package courses

import "time"

type Course struct {
	CourseID       string    `json:"id"`
	Description    string    `json:"description"`
	CourseName     string    `json:"course_name"`
	CourseFullName string    `json:"course_full_name"`
	CourseType     string    `json:"course_type"`
	GradingType    string    `json:"grading_type"`
	Faculty        string    `json:"faculty"`
	MidtermExam    *time.Time `json:"midterm_exam_date"`
	FinalExam      *time.Time `json:"final_exam_date"`
	Credit         int       `json:"credit"`
	CourseGroupID  *int       `json:"course_group_id"`
}
