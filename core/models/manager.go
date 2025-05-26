package models

import "database/sql"

type ModelManager struct {
	db *sql.DB
}

func NewModelManager(db *sql.DB) *ModelManager {
	return &ModelManager{
		db: db,
	}
}
