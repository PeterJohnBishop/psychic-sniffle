package mydb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT UNIQUE NOT NULL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}

func CreateUser(db *sql.DB, user User) (User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := uuid.NewV1()
	if err != nil {
		return User{}, err
	}
	userID := "user_" + id.String()
	hashedPassword, err := HashedPassword(user.Password)
	if err != nil {
		return User{}, err
	}
	query := "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING created_at"
	queryErr := db.QueryRowContext(ctx, query, userID, user.Name, user.Email, hashedPassword).Scan(&user.CreatedAt)
	if queryErr != nil {
		return User{}, queryErr
	}
	return user, nil
}

func GetUsers(db *sql.DB) ([]User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT id, name, email, password, created_at, updated_at FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return []User{}, err
	}
	return users, nil
}

func GetUserByEmail(db *sql.DB, email string) (User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1"
	err := db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByID(db *sql.DB, id string) (User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1"
	err := db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUserByID(db *sql.DB, id string, user User) (User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	UPDATE users 
	SET name = $1, email = $2, password = $3, updated_at = NOW() 
	WHERE id = $4 
	RETURNING id, name, email, password, created_at, updated_at`
	var updatedUser User
	err := db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, id).
		Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Password, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func DeleteUserByID(db *sql.DB, id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM users WHERE id = $1"
	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}
	return nil

}
