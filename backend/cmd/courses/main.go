package main

import (
	"backend/internal/core/courses"
	"backend/internal/core/sections"
	"backend/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app2 := fiber.New()
	sqlDSN := os.Getenv("SQL_DB_DSN")
	dbSQL, err := sql.Open("mysql", sqlDSN)

	if err != nil {
		log.Fatal("mysql connection error : ", err)
	}
	defer dbSQL.Close()
	database.DB = dbSQL
	database.NewSQL()

	sectionService := sections.NewSectionService(dbSQL)
	sectionHandler := sections.NewSectionHandler(sectionService)
	courseService := courses.NewCourseService(dbSQL)
	courseHandler := courses.NewCourseHandler(courseService, sectionService)

	apiv1 := app2.Group("/api/v1")

	apiv1.Get("/sections", sectionHandler.GetAllSections)
	apiv1.Get("/sections/courses/:id", sectionHandler.GetSectionsByCourseID)
	apiv1.Get("/sections/:id", sectionHandler.GetSectionByID)
	apiv1.Post("/sections", sectionHandler.CreateSection)
	apiv1.Put("/sections/:id", sectionHandler.UpdateSection)
	apiv1.Delete("/sections/:id", sectionHandler.DeleteSection)

	apiv1.Get("/courses", courseHandler.GetCourses)
	apiv1.Get("/courses/search", courseHandler.IndexCourses)
	apiv1.Get("/courses/paginated", courseHandler.GetCoursesPaginated)
	apiv1.Get("/courses/:id", courseHandler.GetCourse)
	apiv1.Post("/courses", courseHandler.CreateCourse)
	apiv1.Put("/courses/:id", courseHandler.UpdateCourse)
	apiv1.Delete("/courses/:id", courseHandler.DeleteCourse)

	log.Fatal(app2.Listen(os.Getenv("BACKEND_REST")))
}

func helloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello"})
}
