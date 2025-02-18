package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct {
	db *sql.DB
}

func (m *UserModel) Insert(user User) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := m.db.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func (m *UserModel) Update(user User) error {
	query := "UPDATE users SET username=$1, password=$2 WHERE username=$1"
	_, err := m.db.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)

	}
	return nil
}

func (m *UserModel) Delete(user User) error {
	query := "DELETE FROM users WHERE username=$1"
	_, err := m.db.Exec(query, user.Username)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

func (m *UserModel) Get(username string) (User, error) {
	query := "SELECT id, username, password FROM users WHERE username=$1"
	row := m.db.QueryRow(query, username)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return User{}, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

func (m *UserModel) GetAll() ([]User, error) {
	query := "SELECT id, username FROM users"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %v", err)
	}
	defer rows.Close()
	if err != nil {
		fmt.Printf("error closing rows: %v", err)
	}
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("error getting users: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *UserModel) FindUser(username string) ([]User, error) {
	query := "SELECT username FROM users WHERE username ILIKE $1"
	rows, err := m.db.Query(query, username+"%")
	if err != nil {
		return nil, fmt.Errorf("error getting users: %v", err)
	}
	defer rows.Close()
	if err != nil {
		fmt.Printf("error closing rows: %v", err)
	}
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username)
		if err != nil {
			return nil, fmt.Errorf("error getting users: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}
