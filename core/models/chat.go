package models

import (
	"database/sql"
	"errors"
	"time"
)

type Chat struct {
	ID        int
	Time      time.Time
	CreatorID int
	MemberID  int

	Creator User
	Member  User
}

type ChatData struct {
	ID        int
	Time      time.Time
	CreatorID int
	MemberID  int

	Name        string
	Surname     string
	LastMessage string
}

// CreateChat создает новый чат
func (m *ModelManager) CreateChat(creatorID, memberID int) (*Chat, error) {
	rows, err := m.db.Query(`SELECT * FROM chats WHERE creator_id = $1 AND member_id = $2;`, creatorID, memberID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return nil, errors.New("чат уже существует")
	}

	rows, err = m.db.Query(`SELECT * FROM chats WHERE creator_id = $1 AND member_id = $2;`, memberID, creatorID)
	if err != nil {
		return nil, err
	}
	
	if rows.Next() {
		return nil, errors.New("чат уже существует")
	}

	chat := &Chat{
		CreatorID: creatorID,
		MemberID:  memberID,
		Time:      time.Now(),
	}

	err = m.db.QueryRow(
		`INSERT INTO chats (time, creator_id, member_id) 
		VALUES ($1, $2, $3) 
		RETURNING id;`,
		chat.Time, chat.CreatorID, chat.MemberID,
	).Scan(&chat.ID)

	return chat, err
}

// GetChatByID возвращает чат по ID
func (m *ModelManager) GetChatByID(id int) (*Chat, error) {
	chat := &Chat{}
	err := m.db.QueryRow(
		`SELECT id, time, creator_id, member_id 
		FROM chats WHERE id = $1;`,
		id,
	).Scan(&chat.ID, &chat.Time, &chat.CreatorID, &chat.MemberID)

	if err == sql.ErrNoRows {
		return nil, errors.New("чат не найден")
	}

	return chat, err
}

// GetChatsForUser возвращает все чаты пользователя
func (m *ModelManager) GetChatsForUser(userID int) ([]Chat, error) {
	rows, err := m.db.Query(
		`SELECT 
            c.id, c.time, c.creator_id, c.member_id,
            creator.id, creator.login, creator.name, creator.surname,
            member.id, member.login, member.name, member.surname
        FROM chats c
        JOIN users creator ON c.creator_id = creator.id
        JOIN users member ON c.member_id = member.id
        WHERE c.creator_id = $1 OR c.member_id = $1
        ORDER BY c.time DESC;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]Chat, 0)
	for rows.Next() {
		chat := Chat{}
		err := rows.Scan(
			&chat.ID, &chat.Time, &chat.CreatorID, &chat.MemberID,
			&chat.Creator.ID, &chat.Creator.Login, &chat.Creator.Name, &chat.Creator.Surname,
			&chat.Member.ID, &chat.Member.Login, &chat.Member.Name, &chat.Member.Surname,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (m *ModelManager) GetChatsForUserWithNames(userID int) ([]ChatData, error) {
	chats, err := m.GetChatsForUser(userID)
	if err != nil {
		return nil, err
	}

	chatsData := make([]ChatData, 0)
	for _, chat := range chats {
		chatData := ChatData{
			ID:        chat.ID,
			Time:      chat.Time,
			CreatorID: chat.CreatorID,
			MemberID:  chat.MemberID,
		}

		if chat.CreatorID == userID {
			chatData.Name = chat.Member.Name
			chatData.Surname = chat.Member.Surname
		} else {
			chatData.Name = chat.Creator.Name
			chatData.Surname = chat.Creator.Surname
		}

		err = m.db.QueryRow("SELECT text FROM messages WHERE chat_id = $1 ORDER BY time DESC LIMIT 1;", chat.ID).Scan(&chatData.LastMessage)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		chatsData = append(chatsData, chatData)
	}

	return chatsData, nil
}

// UpdateChat обновляет данные чата
func (m *ModelManager) UpdateChat(chat *Chat) error {
	_, err := m.db.Exec(
		`UPDATE chats 
		SET time = $1, creator_id = $2, member_id = $3 
		WHERE id = $4;`,
		chat.Time, chat.CreatorID, chat.MemberID, chat.ID,
	)
	return err
}

// DeleteChat удаляет чат по ID
func (m *ModelManager) DeleteChat(id int) error {
	_, err := m.db.Exec(`DELETE FROM chats WHERE id = $1;`, id)
	return err
}
