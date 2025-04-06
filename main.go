package main

import (
	"database/sql"
	"fmt"
	"log"
	mydb "psychic-sniffle/main.go/db"
	"psychic-sniffle/main.go/server"
)

var db *sql.DB

func main() {
	fmt.Println("Hello, World!")

	db, err := mydb.Connect(db)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = mydb.CreateUsersTable(db)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	defer db.Close()

	err = server.StartServer(db)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
