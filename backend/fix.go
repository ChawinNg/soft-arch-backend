package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func main() {
	// Replace with your actual connection string
	db, err := sql.Open("mysql", "root:root@tcp(0.0.0.0:3333)/regdealer")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Update capacity to 0
	_, err = db.Exec("UPDATE sections SET capacity = 0")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("UPDATE enrollments SET summarized = FALSE")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("	DELETE FROM enrollment_results")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Capacity updated to 0 for all sections.")
}
