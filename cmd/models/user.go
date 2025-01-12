package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	db *sql.DB
}

func (m *UserModel) Insert(user User) error {
	request := "INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3)"
	_, err := m.db.Exec(request, user.Nickname, user.Email, user.Password)
	if err != nil {
		fmt.Errorf("Error inserting user: %v", err)
	}
	return nil
}

func (m *UserModel) Update(user User) error {
	request := "UPDATE users SET nickname=$1, email=$2, password=$3 WHERE nickname=$1"
	_, err := m.db.Exec(request, user.Nickname, user.Email, user.Password)
	if err != nil {
		fmt.Errorf("Error updating user: %v", err)
	}
	return nil
}
