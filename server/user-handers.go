package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	mydb "psychic-sniffle/main.go/db"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user mydb.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	userCreated, err := mydb.CreateUser(db, user)
	if err != nil {
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User Created Successfully",
		"user":    userCreated,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := mydb.GetUserByEmail(db, req.Email)
	if err != nil {
		http.Error(w, `{"error": "Failed to get user by that email"}`, http.StatusInternalServerError)
		return
	}

	if !mydb.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, `{"error": "Password Verification Failed"}`, http.StatusUnauthorized)
		return
	}

	userClaims := mydb.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token, err := mydb.NewAccessToken(userClaims)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate authentication token"}`, http.StatusInternalServerError)
		return
	}

	refreshToken, err := mydb.NewRefreshToken(userClaims.StandardClaims)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate refresh token"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":      "Login Success",
		"token":        token,
		"refreshToken": refreshToken,
		"user":         user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RefreshTokenHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	const userIDKey ContextKey = "userID"

	id, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, `{"error": "ID not found in context"}`, http.StatusInternalServerError)
		return
	}

	user, err := mydb.GetUserByID(db, id)
	if err != nil {
		http.Error(w, `{"error": "Failed to get user by that email"}`, http.StatusInternalServerError)
		return
	}

	userClaims := mydb.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token, err := mydb.NewAccessToken(userClaims)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate authentication token"}`, http.StatusInternalServerError)
		return
	}

	refreshToken, err := mydb.NewRefreshToken(userClaims.StandardClaims)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate refresh token"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":      "Token Refreshed",
		"token":        token,
		"refreshToken": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserByEmailHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, email string) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// pathParts := strings.Split(r.URL.Path, "/")
	// if len(pathParts) < 4 || pathParts[2] != "email" {
	// 	http.Error(w, "Invalid URL", http.StatusBadRequest)
	// 	return
	// }
	// email := pathParts[3]

	var user mydb.User
	foundUser, err := mydb.GetUserByEmail(db, email)
	if err != nil {
		http.Error(w, `{"error": "Failed to find User with that email!"}`, http.StatusInternalServerError)
		return
	}
	user = foundUser

	response := map[string]interface{}{
		"message": "User Found with Email, " + email,
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserByIDHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, id string) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// pathParts := strings.Split(r.URL.Path, "/")
	// if len(pathParts) < 4 || pathParts[2] != "id" {
	// 	http.Error(w, "Invalid URL", http.StatusBadRequest)
	// 	return
	// }
	// id := pathParts[3]

	var user mydb.User
	foundUser, err := mydb.GetUserByID(db, id)
	if err != nil {
		http.Error(w, `{"error": "Failed to find User with that ID!"}`, http.StatusInternalServerError)
		return
	}
	user = foundUser

	response := map[string]interface{}{
		"message": "User Found with ID, " + id,
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUsersHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var users []mydb.User
	allUsers, err := mydb.GetUsers(db)
	if err != nil {
		http.Error(w, `{"error": "Failed to get all Users!"}`, http.StatusInternalServerError)
		return
	}
	users = allUsers

	response := map[string]interface{}{
		"message": "Users found!",
		"users":   users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, id string) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// pathParts := strings.Split(r.URL.Path, "/")
	// if len(pathParts) < 4 || pathParts[2] != "update" {
	// 	http.Error(w, "Invalid URL", http.StatusBadRequest)
	// 	return
	// }
	// id := pathParts[3]

	var user mydb.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedUser, err := mydb.UpdateUserByID(db, id, user)
	if err != nil {
		http.Error(w, `{"error": "Failed to update user!"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User updated!",
		"user":    updatedUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, id string) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// pathParts := strings.Split(r.URL.Path, "/")
	// if len(pathParts) < 4 || pathParts[2] != "delete" {
	// 	http.Error(w, "Invalid URL", http.StatusBadRequest)
	// 	return
	// }
	// id := pathParts[3]

	dbErr := mydb.DeleteUserByID(db, id)
	if dbErr != nil {
		http.Error(w, `{"error": "Failed to delete user!"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User deleted!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
