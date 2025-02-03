package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/models"
	"net/http"
	"strconv"
	"time"
)

type messageHandler struct {
	App *config.Application
}

// CreateMessage отправляет сообщение в чат
// @Summary Отправить сообщение
// @Tags messages
// @Accept json
// @Produce json
// @Param input body object{ ChatID int64; SenderID int64; Text string } true "Message data"
// @Success 201 {string} string "Message sent"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /message [post]
func (h messageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatID   int64  `json:"chat_id"`
		SenderID int64  `json:"sender_id"`
		Text     string `json:"text"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message := models.Message{
		ChatID:   input.ChatID,
		SenderID: input.SenderID,
		Text:     input.Text,
		SentAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	err = h.App.Models.Message.Insert(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetMessagesByChat получает сообщения чата
// @Summary Получить сообщения чата
// @Tags messages
// @Produce json
// @Param chat_id query int true "Chat ID"
// @Success 200 {array} models.Message
// @Failure 400 {string} string "Invalid Chat ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /messages [get]
func (h messageHandler) GetMessagesByChat(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.ParseInt(r.URL.Query().Get("chat_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	messages := h.App.Models.Message.GetMessagesByChatID(chatID)
	if messages == nil {
		http.Error(w, "No messages found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// DeleteMessage удаляет сообщение по его ID
// @Summary Удалить сообщение
// @Tags messages
// @Param id query int true "Message ID"
// @Success 200 {string} string "Message deleted"
// @Failure 400 {string} string "Invalid Message ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /message [delete]
func (h messageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = h.App.Models.Message.Delete(models.Message{ID: messageID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SearchMessages ищет сообщения по тексту
// @Summary Искать сообщения по содержимому
// @Tags messages
// @Produce json
// @Param text query string true "Search text"
// @Success 200 {array} models.Message
// @Failure 400 {string} string "Search text is required"
// @Failure 404 {string} string "No messages found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /messages/search [get]
func (h messageHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("text")
	if searchText == "" {
		http.Error(w, "Search text is required", http.StatusBadRequest)
		return
	}

	messages := h.App.Models.Message.SearchMessage(searchText)
	if messages == nil {
		http.Error(w, "No messages found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding messages: %v", err), http.StatusInternalServerError)
	}
}

// UpdateMessage обновляет текст сообщения по его ID
// @Summary Обновить сообщение
// @Tags messages
// @Accept json
// @Produce json
// @Param input body object{ ID int64; Text string } true "Updated message data"
// @Success 200 {string} string "Message updated"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /message [put]
func (h messageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID   int64  `json:"id"`
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.App.Models.Message.Update(models.Message{ID: input.ID, Text: input.Text})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
