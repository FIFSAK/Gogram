package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/middleware"
	"github.com/FIFSAK/Gogram/internal/models"
	"net/http"
)

type userHandler struct {
	App *config.Application
}

// GetAllUsers возвращает список всех пользователей
// @Summary Получить всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [get]
func (m userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.App.Models.User.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body object{ Username string; Password string } true "User credentials"
// @Success 201 {string} string "User created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (m userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := middleware.HashPassword(input.Password)
	user := models.User{Username: input.Username, Password: hashedPassword}
	err = m.App.Models.User.Insert(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Login авторизует пользователя
// @Summary Вход в систему
// @Tags users
// @Accept json
// @Produce json
// @Param input body object{ Username string; Password string } true "User credentials"
// @Success 200 {string} string "Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (m userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := m.App.Models.User.Get(input.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if middleware.ComparePasswords(user.Password, input.Password) {
		token, err := middleware.GenerateToken(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// Search ищет пользователя по имени пользователя
// @Summary Поиск пользователя
// @Tags users
// @Produce json
// @Param username query string true "Username"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Internal Server Error"
// @Router /search [get]
func (m userHandler) Search(w http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user, err := m.App.Models.User.Get(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		err = fmt.Errorf("Error encoding user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
