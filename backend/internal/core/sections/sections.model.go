package sections

type Instructor struct {
	InstructorID int
	FullName     string
	Faculty      string
	DisplayName  string
	Email        *string
	PhoneNumber  *string
}

type Section struct {
	CourseID    string      `json:"courseId"`
	SectionID   int         `json:"id"`
	Section     int         `json:"section"`
	Capacity    int         `json:"capacity"`
	Room        *string     `json:"room"`
	Timeslots   [][]string `json:"timeslots"`
	Instructors []Instructor `json:"instructors"`
}
