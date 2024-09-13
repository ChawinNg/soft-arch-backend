package sections

type Instructor struct {
	InstructorID int     `json:"instructor_id"`
	FullName     string  `json:"full_name"`
	Faculty      string  `json:"faculty"`
	DisplayName  string  `json:"display_name"`
	Email        *string `json:"email"`
	PhoneNumber  *string `json:"phone_number"`
}

type Section struct {
	CourseID    string       `json:"courseId"`
	SectionID   int          `json:"id"`
	Section     int          `json:"section"`
	Capacity    int          `json:"capacity"`
	Room        *string      `json:"room"`
	Timeslots   [][]string   `json:"timeslots"`
	Instructors []Instructor `json:"instructors"`
}
