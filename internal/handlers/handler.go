package handlers

import "github.com/FIFSAK/Gogram/internal/config"

type Handler struct {
	UserHandler userHandler
}

func New(app *config.Application) Handler {
	return Handler{
		UserHandler: userHandler{App: app},
	}
}
