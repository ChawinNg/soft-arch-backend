package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func NewSQL() {
    var err error
	sqlDSN := os.Getenv("SQL_DB_DSN")
    DB, err = sql.Open("mysql", sqlDSN)
    if err != nil {
        log.Fatal(err)
    }

    createTable := `CREATE TABLE IF NOT EXISTS courses (
        courseid INTEGER PRIMARY KEY,
        description TEXT,
        coursetype TEXT,
		coursegroupid INTEGER
    );`

    createSectionsTable := `CREATE TABLE IF NOT EXISTS sections (
        section_id INT AUTO_INCREMENT PRIMARY KEY,
        course_id INT,
        section INT,
        capacity INT,
        FOREIGN KEY (course_id) REFERENCES courses(course_id)
    );`
    _, err = DB.Exec(createTable)

    if err != nil {
        log.Fatal(err)
    }
    
    _, err = DB.Exec(createSectionsTable)

    if err != nil {
        log.Fatal(err)
    }
}
