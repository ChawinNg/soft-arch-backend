package main

import (
	"database/sql"
	"log"
	"os"

	"backend/internal/core/courses"
	"backend/internal/core/enrollments"
	"backend/internal/core/sections"
	user "backend/internal/core/users/rest"
	"backend/internal/database"
	userService "backend/internal/genproto/users"
	"backend/internal/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize Fiber
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",       // Specify the allowed origin(s)
		AllowMethods:     "GET,POST,PUT,DELETE",         // Specify allowed methods
		AllowHeaders:     "Content-Type, Authorization", // Specify allowed headers
		AllowCredentials: true,
	}))

	// gRPC Client
	grpc_host := os.Getenv("GRPC_SERVER_HOST")
	conn, err := grpc.NewClient(grpc_host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	// Connect to MySQL
	sqlDSN := os.Getenv("SQL_DB_DSN")
	dbSQL, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer dbSQL.Close()
	database.DB = dbSQL
	database.NewSQL()

	//Define Handlers & Services
	userConn := userService.NewUserServiceClient(conn)
	userHandler := user.NewHandler(userConn)

	sectionService := sections.NewSectionService(dbSQL)
	sectionHandler := sections.NewSectionHandler(sectionService)

	courseService := courses.NewCourseService(dbSQL)
	courseHandler := courses.NewCourseHandler(courseService, sectionService)

	enrollmentService := enrollments.NewEnrollmentService(dbSQL)
	enrollmentHandler := enrollments.NewEnrollmentHandler(enrollmentService, userConn)

	//Define Middleware
	secret := os.Getenv("JWT_SECRET")
	mw := middleware.NewMiddleware(secret)
	apiv1 := app.Group("/api/v1", mw.SessionMiddleware)

	// Define route
	apiv1.Get("/users/me", mw.WithAuthentication(userHandler.GetCurrentUser))
	apiv1.Post("/users/password", mw.WithAuthentication(userHandler.CheckPassword))
	apiv1.Get("/users/:id", mw.WithAuthentication(userHandler.GetUser))
	apiv1.Get("/users", mw.WithAuthentication(userHandler.GetAllUsers))
	apiv1.Put("/users/:id", mw.WithAuthentication(userHandler.UpdateUser))
	apiv1.Delete("/users/:id", mw.WithAuthentication(userHandler.DeleteUser))

	apiv1.Post("/points/reset", mw.WithAuthentication(userHandler.ResetAllUserPoint))
	apiv1.Get("/points/me", mw.WithAuthentication(userHandler.GetCurrentUserPoint))
	apiv1.Post("/points/:id", mw.WithAuthentication(userHandler.ReduceUserPoint))

	apiv1.Post("/register", userHandler.RegisterUser)
	apiv1.Post("/login", userHandler.LoginUser)

	apiv1.Get("/courses", courseHandler.GetCourses)
	apiv1.Get("/courses/paginated", courseHandler.GetCoursesPaginated)
	apiv1.Get("/courses/:id", courseHandler.GetCourse)
	apiv1.Post("/courses", courseHandler.CreateCourse)
	apiv1.Put("/courses/:id", courseHandler.UpdateCourse)
	apiv1.Delete("/courses/:id", courseHandler.DeleteCourse)

	apiv1.Get("/sections", sectionHandler.GetAllSections)
	apiv1.Get("/sections/courses/:id", sectionHandler.GetSectionsByCourseID)
	apiv1.Get("/sections/:id", sectionHandler.GetSectionByID)
	apiv1.Post("/sections", sectionHandler.CreateSection)
	apiv1.Put("/sections/:id", sectionHandler.UpdateSection)
	apiv1.Delete("/sections/:id", sectionHandler.DeleteSection)

	apiv1.Get("/enrollments/:user_id", enrollmentHandler.GetUserEnrollment)
	apiv1.Get("/enrollments/:course_id", enrollmentHandler.GetCourseEnrollment)
	apiv1.Post("/enrollments", enrollmentHandler.CreateEnrollment)
	apiv1.Put("/enrollments/:id", enrollmentHandler.EditEnrollment)
	apiv1.Delete("/enrollments/:id", enrollmentHandler.DeleteEnrollment)
	apiv1.Delete("/enrollments/summarize/:user_id", enrollmentHandler.SummarizeUserEnrollmentResult)

	// Start the server
	log.Fatal(app.Listen("localhost:8080"))

}
