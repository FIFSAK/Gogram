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

func (m userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.App.Models.User.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		err = fmt.Errorf("Error encoding users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

func (m userHandler) Login(w http.ResponseWriter, request *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(request.Body).Decode(&input)
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

func (m userHandler) Search(w http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("username")
	user, err := m.App.Models.User.Get(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		err = fmt.Errorf("Error encoding user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
