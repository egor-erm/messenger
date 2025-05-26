package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

type ViewData struct {
	ID int
}

var chats = template.Must(template.ParseFiles("./site/chats.html"))

func (hm *HandlerManager) ChatsHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	user, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	chats.Execute(ctx.Writer, ViewData{ID: user.ID})
}

func (hm *HandlerManager) GetChatsHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	user, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	chats, err := hm.ModelManager.GetChatsForUserWithNames(user.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, chats)
}

type UserLoginData struct {
	UserLogin string
}

func (hm *HandlerManager) CreateChatHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	user, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	var userLogin UserLoginData
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	member, err := hm.ModelManager.GetUserByLogin(userLogin.UserLogin)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = hm.ModelManager.CreateChat(user.ID, member.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (hm *HandlerManager) GetChatMessagesHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	_, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	chatIDString := ctx.Params.ByName("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Недопустимый chatID")
		return
	}

	messages, err := hm.ModelManager.GetMessagesForChat(chatID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type MessageData struct {
	Text string
}

func (hm *HandlerManager) SendMessageHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	user, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	chatIDString := ctx.Params.ByName("chatID")
	chatID, err := strconv.Atoi(chatIDString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Недопустимый chatID")
		return
	}

	var sendMessageData MessageData
	if err := ctx.ShouldBindJSON(&sendMessageData); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = hm.ModelManager.CreateMessage(sendMessageData.Text, nil, chatID, user.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (hm *HandlerManager) UpdateMessageHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	user, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	messageIDString := ctx.Params.ByName("messageID")
	messageID, err := strconv.Atoi(messageIDString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Недопустимый messageID")
		return
	}

	var updateMessageData MessageData
	if err := ctx.ShouldBindJSON(&updateMessageData); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	message, err := hm.ModelManager.GetMessageByID(messageID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if message.SenderID != user.ID {
		ctx.String(http.StatusForbidden, "Вы не можете редактировать это сообщение")
		return
	}

	message.Text = updateMessageData.Text
	err = hm.ModelManager.UpdateMessage(message)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (hm *HandlerManager) DeleteMessageHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	_, ok := hm.GetSession(token)
	if !ok {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	messageIDString := ctx.Params.ByName("messageID")
	messageID, err := strconv.Atoi(messageIDString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Недопустимый messageID")
		return
	}

	err = hm.ModelManager.DeleteMessage(messageID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
