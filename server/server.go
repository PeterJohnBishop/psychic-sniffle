package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func StartServer(db *sql.DB) error {

	var databaseStatus string

	err := db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	} else {
		databaseStatus = "connected"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message":         "psychic-sniffle server is running",
			"database_status": databaseStatus,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})

	addUserRoutes(db, mux)

	fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	return err
}

func addUserRoutes(db *sql.DB, mux *http.ServeMux) {
	mux.HandleFunc("/register", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(db, w, r)
	})))

	mux.HandleFunc("/login", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		Login(db, w, r)
	})))

	mux.HandleFunc("/users/", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		GetUsersHandler(db, w, r)
	})))

	mux.HandleFunc("/users/email/{email}", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		email := r.PathValue("email")
		GetUserByEmailHandler(db, w, r, email)
	})))

	mux.HandleFunc("/users/id/{id}", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		GetUserByIDHandler(db, w, r, id)
	})))

	mux.HandleFunc("/users/update/{id}", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		UpdateUserHandler(db, w, r, id)
	})))

	mux.HandleFunc("/users/delete/{id}", LoggerMiddleware(VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		DeleteUserHandler(db, w, r, id)
	})))

	mux.HandleFunc("/users/refresh/", LoggerMiddleware(VerifyRefreshToken(func(w http.ResponseWriter, r *http.Request) {
		RefreshTokenHandler(db, w, r)
	})))
}
