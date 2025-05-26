package handlers

import (
	"messenger/core/models"

	"golang.org/x/crypto/bcrypt"
)

type HandlerManager struct {
	Sessions     map[string]*models.User
	ModelManager *models.ModelManager
}

func NewHandlerManager(modelManager *models.ModelManager) *HandlerManager {
	return &HandlerManager{
		Sessions:     make(map[string]*models.User),
		ModelManager: modelManager,
	}
}

func (hm *HandlerManager) GetSession(token string) (*models.User, bool) {
	for key, user := range hm.Sessions {
		if err := bcrypt.CompareHashAndPassword([]byte(token), []byte(key)); err == nil {
			return user, true
		}
	}

	return nil, false
}
