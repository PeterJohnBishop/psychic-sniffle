package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

var DatabaseStatus string

func StartServer(db *sql.DB) error {

	err := db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	} else {
		DatabaseStatus = "connected"
	}

	mux := http.NewServeMux()
	handler := IdentifyKubernetesPod(mux)

	addTestingRoutes(mux)
	addUserRoutes(db, mux)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	log.Println("Listening on", listener.Addr())
	err = http.ListenAndServe(":8080", handler)
	return err
}

func addTestingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message":         "psychic-sniffle server is running",
			"database_status": DatabaseStatus,
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
