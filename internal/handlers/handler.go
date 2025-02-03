package handlers

import "github.com/FIFSAK/Gogram/internal/config"

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
