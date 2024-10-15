package enrollments

type Enrollment struct {
	EnrollmentID int    `json:"id"`
	UserID       int    `json:"user_id"`
	CourseID     string `json:"course_id"`
	SectionID    int    `json:"section_id"`
	Section      int    `json:"section"`
	Points       int64  `json:"point"`
	Round        string `json:"round"`
}
