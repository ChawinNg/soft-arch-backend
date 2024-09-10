package courses

type Course struct {
    CourseID       int    `json:"courseId"`
    Description    string `json:"description"`
    CourseType     string `json:"courseType"`
    CourseGroupID  int    `json:"courseGroupId"`
}
