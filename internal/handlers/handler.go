package handlers

import (
	"fmt"
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Handler struct {
	UserHandler    userHandler
	ChatHandler    chatHandler
	MessageHandler messageHandler
}

func New(app *config.Application) Handler {
	return Handler{
		UserHandler:    userHandler{App: app},
		ChatHandler:    chatHandler{App: app},
		MessageHandler: messageHandler{App: app},
	}
}

func GetUserIDFromContext(r *http.Request) (int64, error) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("user claims not found in context")
	}

	userID, ok := claims["id"].(float64)
	fmt.Println(userID)
	if !ok {
		return 0, fmt.Errorf("invalid user ID format")
	}

	return int64(userID), nil
}
