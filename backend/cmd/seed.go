package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(0.0.0.0:3333)/regdealer")
	// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/regdealer")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	seedDir := "./seeds/"
	seedFiles, err := os.ReadDir(seedDir)

	for _, file := range seedFiles {
		seedFile := seedDir + file.Name()
		content, err := os.ReadFile(seedFile)
		if err != nil {
			log.Fatal(seedFile,err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal(seedFile,err)
		}
	}

	log.Println("Database seeded successfully!")
}
