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

	fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	return err
}
