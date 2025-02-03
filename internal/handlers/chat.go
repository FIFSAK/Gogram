package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/models"
	"net/http"
	"strconv"
)

type chatHandler struct {
	App *config.Application
}

// CreateChat создает новый чат между двумя пользователями
// @Summary Создать чат
// @Tags chats
// @Accept json
// @Produce json
// @Param input body object{ User1Id int64; User2Id int64 } true "User IDs"
// @Success 201 {string} string "Chat created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /chat [post]
func (h chatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var input struct {
		User1Id int64 `json:"user1_id"`
		User2Id int64 `json:"user2_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chat := models.Chat{
		User1Id: input.User1Id,
		User2Id: input.User2Id,
	}
	err = h.App.Models.Chat.Insert(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// DeleteChat удаляет чат
// @Summary Удалить чат
// @Tags chats
// @Param id query int true "Chat ID"
// @Success 200 {string} string "Chat deleted"
// @Failure 400 {string} string "Invalid Chat ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /chat [delete]
func (h chatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	err = h.App.Models.Chat.Delete(models.Chat{ID: chatID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetChat получает чат по его ID
// @Summary Получить чат
// @Tags chats
// @Produce json
// @Param id query int true "Chat ID"
// @Success 200 {object} models.Chat
// @Failure 400 {string} string "Invalid Chat ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /chat [get]
func (h chatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	chat, err := h.App.Models.Chat.Get(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(chat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding chat: %v", err), http.StatusInternalServerError)
	}
}

// GetUserChats получает все чаты, в которых участвует пользователь
// @Summary Получить все чаты пользователя
// @Tags chats
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} models.Chat
// @Failure 400 {string} string "Invalid User ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /chats [get]
func (h chatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chats, err := h.App.Models.Chat.GetUserChatAll(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding chats: %v", err), http.StatusInternalServerError)
	}
}
