package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectPostgres() {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=prestasi_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed connect db:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed ping db:", err)
	}

	DB = db
	fmt.Println("PostgreSQL connected")
}
