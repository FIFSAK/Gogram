package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	db *sql.DB
}

func (m *UserModel) Insert(user User) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := m.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("Error inserting user: %v", err)
	}
	return nil
}

func (m *UserModel) Update(user User) error {
	query := "UPDATE users SET name=$1, email=$2, password=$3 WHERE name=$1"
	_, err := m.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("Error updating user: %v", err)

	}
	return nil
}

func (m *UserModel) Delete(user User) error {
	query := "DELETE FROM users WHERE name=$1"
	_, err := m.db.Exec(query, user.Username)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}
	return nil
}

func (m *UserModel) Get(user User) (User, error) {
	query := "SELECT name, email, password FROM users WHERE name=$1"
	row := m.db.QueryRow(query, user.Username)
	err := row.Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		return User{}, fmt.Errorf("Error getting user: %v", err)
	}
	return user, nil
}

func (m *UserModel) GetAll() ([]User, error) {
	query := "SELECT name, email, password FROM users"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error getting users: %v", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("Error getting users: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}
