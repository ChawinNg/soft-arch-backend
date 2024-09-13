package sections

type Section struct {
	CourseID    string      `json:"courseId"`
	SectionID   int         `json:"id"`
	Section     int         `json:"section"`
	Capacity    int         `json:"capacity"`
	Room        *string     `json:"room"`
	Timeslots   [][]string `json:"timeslots"`
	Instructors [][]string `json:"instructors"`
}
