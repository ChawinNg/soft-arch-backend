package main

import (
	"database/sql"
	"log"
	"os"

	"backend/internal/core/courses"
	user "backend/internal/core/users/rest"
	"backend/internal/database"
	userService "backend/internal/genproto/users"
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize Fiber
	app := fiber.New()

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

	//Define Handlers & Services
	userConn := userService.NewUserServiceClient(conn)
	userHandler := user.NewHandler(userConn)

	courseService := courses.NewCourseService(dbSQL)
	courseHandler := courses.NewCourseHandler(courseService)

	//Define Middleware
	secret := os.Getenv("JWT_SECRET")
	mw := middleware.NewMiddleware(secret)
	apiv1 := app.Group("/api/v1", mw.SessionMiddleware)

	// Define route
	apiv1.Get("/users/me", mw.WithAuthentication(userHandler.GetCurrentUser))
	apiv1.Get("/users/:id", mw.WithAuthentication(userHandler.GetUser))
	apiv1.Get("/users", mw.WithAuthentication(userHandler.GetAllUsers))
	apiv1.Put("/users/:id", mw.WithAuthentication(userHandler.UpdateUser))
	apiv1.Delete("/users/:id", mw.WithAuthentication(userHandler.DeleteUser))

	apiv1.Post("/register", userHandler.RegisterUser)
	apiv1.Post("/login", userHandler.LoginUser)

	app.Get("/courses", courseHandler.GetCourses)
	app.Get("/courses/:id", courseHandler.GetCourse)
	app.Post("/courses", courseHandler.CreateCourse)
	app.Put("/courses/:id", courseHandler.UpdateCourse)
	app.Delete("/courses/:id", courseHandler.DeleteCourse)

	// Start the server
	log.Fatal(app.Listen(":8080"))

}
