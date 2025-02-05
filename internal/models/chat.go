package models

import "database/sql"

type Chat struct {
	ID      int64 `json:"id"`
	User1Id int64 `json:"user1_id"`
	User2Id int64 `json:"user2_id"`
}

type ChatModel struct {
	db *sql.DB
}

func (m *ChatModel) Insert(chat Chat) (int, error) {
	query := "INSERT INTO chat (user1_id, user2_id) VALUES ($1, $2) RETURNING id"

	var id int
	err := m.db.QueryRow(query, chat.User1Id, chat.User2Id).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ChatModel) Delete(chat Chat) error {
	query := "DELETE FROM chat WHERE id=$1"
	_, err := m.db.Exec(query, chat.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *ChatModel) Get(id int64) (Chat, error) {
	query := "SELECT id, user1_id, user2_id FROM chat WHERE id=$1"
	row := m.db.QueryRow(query, id)
	var chat Chat
	err := row.Scan(&chat.ID, &chat.User1Id, &chat.User2Id)
	if err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func (m *ChatModel) CheckExists(user1Id, user2Id int64) (Chat, error) {
	query := "SELECT id, user1_id, user2_id FROM chat WHERE (user1_id=$1 AND user2_id=$2) OR (user1_id=$2 AND user2_id=$1)"
	row := m.db.QueryRow(query, user1Id, user2Id)

	var chat Chat
	err := row.Scan(&chat.ID, &chat.User1Id, &chat.User2Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return Chat{}, nil
		}
		return Chat{}, err
	}
	return chat, nil
}

func (m *ChatModel) GetUserChatAll(userId int64) ([]Chat, error) {
	query := "SELECT id, user1_id, user2_id FROM chat WHERE user1_id=$1 OR user2_id=$1"
	rows, err := m.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := []Chat{}
	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.ID, &chat.User1Id, &chat.User2Id)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
