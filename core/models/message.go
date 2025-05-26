package models

import (
	"database/sql"
	"errors"
	"time"
)

type Message struct {
	ID       int
	Text     string
	FileID   *int
	Time     time.Time
	ChatID   int
	SenderID int

	Sender User
	Chat   Chat
}

// CreateMessage создает новое сообщение
func (m *ModelManager) CreateMessage(text string, fileID *int, chatID, senderID int) (*Message, error) {
	message := &Message{
		Text:     text,
		FileID:   fileID,
		Time:     time.Now(),
		ChatID:   chatID,
		SenderID: senderID,
	}

	err := m.db.QueryRow(
		`INSERT INTO messages (text, file_id, time, chat_id, sender_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id;`,
		message.Text, message.FileID, message.Time, message.ChatID, message.SenderID,
	).Scan(&message.ID)

	return message, err
}

// GetMessageByID возвращает сообщение по ID
func (m *ModelManager) GetMessageByID(id int) (*Message, error) {
	message := &Message{}
	var fileID sql.NullInt64

	err := m.db.QueryRow(
		`SELECT id, text, file_id, time, chat_id, sender_id 
		FROM messages WHERE id = $1;`,
		id,
	).Scan(&message.ID, &message.Text, &fileID, &message.Time, &message.ChatID, &message.SenderID)

	if fileID.Valid {
		val := int(fileID.Int64)
		message.FileID = &val
	}

	if err == sql.ErrNoRows {
		return nil, errors.New("message not found")
	}

	return message, err
}

// GetMessagesForChat возвращает все сообщения в чате
func (m *ModelManager) GetMessagesForChat(chatID int) ([]Message, error) {
	rows, err := m.db.Query(
		`SELECT m.id, m.text, m.file_id, m.time, m.chat_id, m.sender_id,
		       u.id, u.login, u.name, u.surname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.chat_id = $1
		ORDER BY m.time ASC;`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		message := Message{}
		var fileID sql.NullInt64

		err := rows.Scan(
			&message.ID, &message.Text, &fileID, &message.Time,
			&message.ChatID, &message.SenderID,
			&message.Sender.ID, &message.Sender.Login,
			&message.Sender.Name, &message.Sender.Surname,
		)
		if err != nil {
			return nil, err
		}

		if fileID.Valid {
			val := int(fileID.Int64)
			message.FileID = &val
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// UpdateMessage обновляет данные сообщения
func (m *ModelManager) UpdateMessage(message *Message) error {
	var fileID interface{}
	if message.FileID != nil {
		fileID = *message.FileID
	} else {
		fileID = nil
	}

	_, err := m.db.Exec(
		`UPDATE messages 
		SET text = $1, file_id = $2, time = $3, chat_id = $4, sender_id = $5 
		WHERE id = $6;`,
		message.Text, fileID, message.Time, message.ChatID, message.SenderID, message.ID,
	)
	return err
}

// DeleteMessage удаляет сообщение по ID
func (m *ModelManager) DeleteMessage(id int) error {
	_, err := m.db.Exec(`DELETE FROM messages WHERE id = $1;`, id)
	return err
}
