package models

import (
	"database/sql"
	"errors"
)

type User struct {
	ID           int
	Login        string
	PasswordHash string
	Name         string
	Surname      string
}

// CreateUser создает нового пользователя
func (m *ModelManager) CreateUser(login, passwordHash, name, surname string) (*User, error) {
	user := &User{Login: login, PasswordHash: passwordHash, Name: name, Surname: surname}

	err := m.db.QueryRow(
		`INSERT INTO users (login, password_hash, name, surname) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id;`,
		user.Login, user.PasswordHash, user.Name, user.Surname,
	).Scan(&user.ID)

	return user, err
}

// GetUserByID возвращает пользователя по ID
func (m *ModelManager) GetUserByID(id int) (*User, error) {
	user := &User{}

	err := m.db.QueryRow(
		`SELECT id, login, password_hash, name, surname 
		FROM users WHERE id = $1;`,
		id,
	).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Name, &user.Surname)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

// GetUserByLogin возвращает пользователя по логину
func (m *ModelManager) GetUserByLogin(login string) (*User, error) {
	user := &User{}

	err := m.db.QueryRow(
		`SELECT id, login, password_hash, name, surname
		FROM users WHERE login = $1;`,
		login,
	).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Name, &user.Surname)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

// UpdateUser обновляет данные пользователя
func (m *ModelManager) UpdateUser(user *User) error {
	_, err := m.db.Exec(
		`UPDATE users 
		SET login = $1, password_hash = $2, name = $3, surname = $4 
		WHERE id = $5;`,
		user.Login, user.PasswordHash, user.Name, user.Surname, user.ID,
	)
	return err
}

func (m *ModelManager) DeleteUser(id int) error {
	_, err := m.db.Exec(`DELETE FROM users WHERE id = $1;`, id)
	return err
}
