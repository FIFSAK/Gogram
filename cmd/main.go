// @title Gogram API
// @version 1.0
// @description API for Gogram messaging application.
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"fmt"
	_ "github.com/FIFSAK/Gogram/docs"
	"github.com/FIFSAK/Gogram/internal/config"
	"github.com/FIFSAK/Gogram/internal/handlers"
	"github.com/FIFSAK/Gogram/internal/middleware"
	"github.com/FIFSAK/Gogram/internal/models"
	"github.com/FIFSAK/Gogram/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL документации
	))
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})
	r.Get("/users", handler.UserHandler.GetAllUsers)
	r.Post("/register", handler.UserHandler.Register)
	r.Post("/login", handler.UserHandler.Login)

	r.Post("/search", middleware.RequireAuth(handler.UserHandler.Search))
	r.Post("/chat", handler.ChatHandler.CreateChat)
	r.Delete("/chat", handler.ChatHandler.DeleteChat)
	r.Get("/chat", handler.ChatHandler.GetChat)
	r.Get("/chats", handler.ChatHandler.GetUserChats)

	r.Post("/message", handler.MessageHandler.CreateMessage)
	r.Delete("/message", handler.MessageHandler.DeleteMessage)
	r.Get("/messages/search", handler.MessageHandler.SearchMessages)
	r.Get("/messages", handler.MessageHandler.GetMessagesByChat)
	r.Put("/message", handler.MessageHandler.UpdateMessage)

	log.Printf("Starting server on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {

	}

}
