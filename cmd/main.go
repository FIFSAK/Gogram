package main

import (
	"fmt"
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/handlers"
	"github.com/FIFSAK/Gogram/internal/middleware"
	"github.com/FIFSAK/Gogram/internal/models"
	"github.com/FIFSAK/Gogram/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	db, err := store.New(dsn)
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
