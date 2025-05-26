package web

import (
	"messenger/core/models"
	"messenger/web/handlers"

	"github.com/gin-gonic/gin"
)

type Web struct {
	HandlerManager *handlers.HandlerManager
	router         *gin.Engine
}

func NewWeb(modelManager *models.ModelManager) *Web {
	web := &Web{
		HandlerManager: handlers.NewHandlerManager(modelManager),
		router:         gin.Default(),
	}

	web.loadHandlers()

	return web
}

func (w *Web) loadHandlers() {
	w.router.GET("/", w.HandlerManager.IndexHandler)
	w.router.GET("/login", w.HandlerManager.LoginHandler)
	w.router.GET("/register", w.HandlerManager.RegisterHandler)
	w.router.GET("/chats", w.HandlerManager.ChatsHandler)

	w.router.POST("/login", w.HandlerManager.LoginDataHandler)
	w.router.POST("/register", w.HandlerManager.RegisterDataHandler)

	w.router.GET("/api/chats", w.HandlerManager.GetChatsHandler)
	w.router.POST("/api/chats", w.HandlerManager.CreateChatHandler)

	w.router.GET("/api/chats/:chatID/messages", w.HandlerManager.GetChatMessagesHandler)
	w.router.POST("/api/chats/:chatID/messages", w.HandlerManager.SendMessageHandler)
	w.router.PUT("/api/messages/:messageID", w.HandlerManager.UpdateMessageHandler)
	w.router.DELETE("/api/messages/:messageID", w.HandlerManager.DeleteMessageHandler)
}

func (w *Web) Run() {
	if err := w.router.Run(":80"); err != nil {
		panic(err)
	}
}
