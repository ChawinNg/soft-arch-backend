package enrollments

type Enrollment struct {
	EnrollmentID string `json:"id"`
	UserID       string `json:"user_id"`
	CourseID     string `json:"course_id"`
	SectionID    int    `json:"section_id"`
	Section      int    `json:"section"`
	Points       int64  `json:"points"`
	Round        string `json:"round"`
}
