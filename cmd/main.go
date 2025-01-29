package main

import (
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/handlers"
	"github.com/FIFSAK/Gogram/internal/middleware"
	"github.com/FIFSAK/Gogram/internal/models"
	"github.com/FIFSAK/Gogram/internal/store"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	db, err := store.New("postgres://postgres:pass@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("%v", err)
	}

	cfg := &config.Config{DB: db, Port: "8080"}

	app := &config.Application{Config: cfg, Models: models.NewModels(db)}

	r := chi.NewRouter()

	handler := handlers.New(app)

	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})
	r.Get("/users", handler.UserHandler.GetAllUsers)

	r.Post("/register", handler.UserHandler.Register)

	r.Post("/login", handler.UserHandler.Login)

	r.Post("/search", middleware.RequireAuth(handler.UserHandler.Search))

	log.Printf("Starting server on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {

	}

}
