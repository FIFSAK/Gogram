package main

import (
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/handlers"
	"github.com/FIFSAK/Gogram/internal/models"
	"github.com/FIFSAK/Gogram/internal/store"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	db, err := store.New("postgres://postgres:password@localhost:5432/gogram")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	cfg := &config.Config{DB: db, Port: "8080"}

	app := &config.Application{Config: cfg, Models: models.NewModels(db)}

	r := chi.NewRouter()

	handlers := handlers.New(app)

	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})
	r.Get("/users", handlers.UserHandler.GetAllUsers)

}
