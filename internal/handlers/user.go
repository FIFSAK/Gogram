package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FIFSAK/Gogram/internal/config"
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
