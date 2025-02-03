package models

import "database/sql"

type Message struct {
	ID       int64  `json:"id"`
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id"`
	Text     string `json:"message"`
	SentAt   string `json:"sent_at"`
}

type MessageModel struct {
	db *sql.DB
}

func (m *MessageModel) Insert(message Message) error {
	query := "INSERT INTO message (chat_id, sender_id, text, sent_at) VALUES ($1, $2, $3, $4)"
	_, err := m.db.Exec(query, message.ChatID, message.SenderID, message.Text, message.SentAt)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageModel) Delete(message Message) error {
	query := "DELETE FROM message WHERE id=$1"
	_, err := m.db.Exec(query, message.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageModel) SearchMessage(strMessage string) []Message {
	query := "SELECT * from message where text like '%$1%'"
	rows, err := m.db.Query(query, strMessage)
	if err != nil {
		return nil
	}
	defer rows.Close()
	messages := []Message{}
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Text, &message.SentAt)
		if err != nil {
			return nil
		}
		messages = append(messages, message)
	}
	return messages
}

func (m *MessageModel) GetMessagesByChatID(chatID int64) []Message {
	query := "SELECT * from message where chat_id=$1"
	rows, err := m.db.Query(query, chatID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	messages := []Message{}
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Text, &message.SentAt)
		if err != nil {
			return nil
		}
		messages = append(messages, message)
	}
	return messages
}

func (m *MessageModel) Update(message Message) error {
	query := "UPDATE message SET text=$1 WHERE id=$2"
	_, err := m.db.Exec(query, message.Text, message.ID)
	if err != nil {
		return err
	}
	return nil
}
