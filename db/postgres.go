package mydb

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect(db *sql.DB) (*sql.DB, error) {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// set database connection parameters
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr) // Convert to int
	if err != nil {
		return nil, err
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// create a new database instance to work with
	mydb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// ping the server to test the connection
	err = mydb.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Successfully connected to Postgres DB: %s! ", dbname)
	return mydb, nil
}
