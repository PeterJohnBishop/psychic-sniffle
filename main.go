package main

import (
	"database/sql"
	"fmt"
	"log"
	mydb "psychic-sniffle/main.go/db"
)

var db *sql.DB

func main() {
	fmt.Println("Hello, World!")

	db, err := mydb.Connect(db)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
}
