package config

import (
	"database/sql"
	"github.com/FIFSAK/Gogram/internal/models"
)

type Application struct {
	Config *Config
	Models models.Models
}

type Config struct {
	Port string
	DB   *sql.DB
}
