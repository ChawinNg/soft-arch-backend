package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
)

var DB *sql.DB

func NewSQL() {
    var err error
    DB, err = sql.Open("sqlite3", "./mycrudapi.db")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `CREATE TABLE IF NOT EXISTS courses (
        courseid INTEGER PRIMARY KEY,
        description TEXT,
        coursetype TEXT,
		coursegroupid INTEGER,
    );`

    _, err = DB.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
}
