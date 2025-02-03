package models

import (
	"database/sql"
)

type Models struct {
	User    UserModel
	Chat    ChatModel
	Message MessageModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User:    UserModel{db: db},
		Chat:    ChatModel{db: db},
		Message: MessageModel{db: db},
	}
}
