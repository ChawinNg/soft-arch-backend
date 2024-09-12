package sections

type Section struct {
	CourseID  int `json:"courseId"`
	SectionID int `json:"sectionId"`
	Section   int `json:"section"`
	Capacity  int `json:"capacity"`
}
