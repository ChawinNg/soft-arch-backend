package main

import (
	"backend/internal/core/instructors"
	"backend/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app3 := fiber.New()

	dbSQL, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/regdealer")

	if err != nil {
		log.Fatal("mysql connection error : ", err)
	}
	defer dbSQL.Close()
	database.DB = dbSQL
	database.NewSQL()

	instructorService := instructors.NewInstructorService(dbSQL)
	instructorHandler := instructors.NewInstructorHandler(instructorService)

	apiv1 := app3.Group("/api/v1")

	apiv1.Post("/instructors", instructorHandler.CreateInstructor)
	apiv1.Put("/instructors/:id", instructorHandler.UpdateInstructor)
	apiv1.Post("/instructors/contact", instructorHandler.SendEmail)

	log.Fatal(app3.Listen(os.Getenv("http://localhost:8082")))
}

func helloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello"})
}
