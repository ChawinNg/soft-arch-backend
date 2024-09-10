package courses

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

type CourseHandler struct {
    service *CourseService
}

func NewCourseHandler(service *CourseService) *CourseHandler {
    return &CourseHandler{service: service}
}

func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
    courses, err := h.service.GetAllCourses()
    if err != nil {
        http.Error(w, "Error retrieving courses", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courses)
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    var course Course
    if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateCourse(course); err != nil {
        http.Error(w, "Error creating course", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid course ID", http.StatusBadRequest)
        return
    }

    course, err := h.service.GetCourseByID(id)
    if err != nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
    var course Course
    if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.UpdateCourse(course); err != nil {
        http.Error(w, "Error updating course", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid course ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteCourse(id); err != nil {
        http.Error(w, "Error deleting course", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
