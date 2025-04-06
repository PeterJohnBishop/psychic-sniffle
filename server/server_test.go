package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	mydb "psychic-sniffle/main.go/db"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	host := "localhost"
	portStr := os.Getenv("PSQL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid PSQL_PORT:", err)
	}

	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASSWORD")
	dbname := os.Getenv("PSQL_DBNAME")

	// connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db, nil
}

func TestUserEndpoints(t *testing.T) {

	var user mydb.User
	user.Name = fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.Username())
	user.Email = gofakeit.Email()
	user.Password = gofakeit.Password(true, true, true, false, false, 8)

	t.Run("creates a user", func(t *testing.T) {
		fmt.Println("Creating a user...")
		requestBody, _ := json.Marshal(map[string]string{
			"name":     user.Name,
			"email":    user.Email,
			"password": user.Password,
		})

		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		CreateUserHandler(db, response, request)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}
		want := "User Created Successfully"
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}

	})

	t.Run("user logged in", func(t *testing.T) {
		fmt.Println("Logging in a user...")
		requestBody, _ := json.Marshal(map[string]string{
			"email":    user.Email,
			"password": user.Password,
		})

		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		Login(db, response, request)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}

		userData, ok := got["user"].(map[string]interface{})
		if !ok {
			t.Fatal("Failed to parse user data from login response")
		}

		id, ok := userData["id"].(string)
		if !ok {
			t.Fatal("Failed to parse user ID from login response")
		}

		user.ID = id

		pass, ok := userData["password"].(string)
		if !ok {
			t.Fatal("Failed to parse user password from login response")
		}
		user.Password = pass

		want := "Login Success"
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}
	})

	t.Run("finds all users", func(t *testing.T) {
		fmt.Println("Finding all users...")
		request, err := http.NewRequest(http.MethodGet, "/users/", nil)
		if err != nil {
			t.Fatal("Failed to create a new request:", err)
		}
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		GetUsersHandler(db, response, request)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}
		want := "Users found!"
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}
	})

	t.Run("finds a user by id", func(t *testing.T) {
		fmt.Println("Finding a user by ID...")
		request, err := http.NewRequest(http.MethodGet, "/users/id/"+user.ID, nil)
		if err != nil {
			t.Fatal("Failed to create a new request:", err)
		}
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		GetUserByIDHandler(db, response, request, user.ID)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}
		want := "User Found with ID, " + user.ID
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}
	})

	t.Run("finds a user by email", func(t *testing.T) {
		fmt.Println("Finding a user by email...")
		request, err := http.NewRequest(http.MethodGet, "/users/email/"+user.Email, nil)
		if err != nil {
			t.Fatal("Failed to create a new request:", err)
		}
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		GetUserByEmailHandler(db, response, request, user.Email)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}
		want := "User Found with Email, " + user.Email
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}
	})

	t.Run("user updated", func(t *testing.T) {
		fmt.Println("Updating a user...")
		requestBody, _ := json.Marshal(map[string]string{
			"name":     user.Name,
			"email":    gofakeit.Email(),
			"password": user.Password,
			// "updated_at": "2021-09-01T00:00:00Z",
			// "created_at": "2021-09-01T00:00:00Z",
		})

		request, _ := http.NewRequest(http.MethodPut, "/users/update/"+user.ID, bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		UpdateUserHandler(db, response, request, user.ID)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}
		want := "User updated!"
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}
	})

	t.Run("user deleted", func(t *testing.T) {
		fmt.Println("Deleting a user...")
		request, err := http.NewRequest(http.MethodDelete, "/users/delete/"+user.ID, nil)
		if err != nil {
			t.Fatal("Failed to create a new request:", err)
		}
		response := httptest.NewRecorder()

		db, err := ConnectDB()
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		DeleteUserHandler(db, response, request, user.ID)

		var got map[string]interface{}
		err = json.Unmarshal([]byte(response.Body.Bytes()), &got)
		if err != nil {
			t.Fatal("Failed to unmarshal response body:", err)
		}
		want := "User deleted!"
		if msg, ok := got["message"].(string); !ok || msg != want {
			t.Errorf(`Expected message: "%s", but got: "%v"`, want, got["message"])
		} else {
			fmt.Println("PASS")
		}
	})
	fmt.Println("Complete test results for User Endpoints: ")
}
