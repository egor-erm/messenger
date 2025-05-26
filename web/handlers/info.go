package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

var info = template.Must(template.ParseFiles("./site/info.html"))

func (hm *HandlerManager) InfoHandler(ctx *gin.Context) {
	info.Execute(ctx.Writer, nil)
}

func (hm *HandlerManager) UserInfoHandler(ctx *gin.Context) {
	login := ctx.Query("login")

	if login == "" {
		ctx.String(http.StatusBadRequest, "Логин не указан")
		return
	}

	user, err := hm.ModelManager.GetUserByLogin(login)
	if err != nil {
		ctx.String(http.StatusNotFound, "Пользователь не найден")
		return
	}

	total, err := hm.ModelManager.GetUserMessagesCount(user.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка получения статистики")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_messages": total,
	})
}

func (hm *HandlerManager) AllInfoHandler(ctx *gin.Context) {
	daysStr := ctx.Query("days")

	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		ctx.String(http.StatusBadRequest, "Ошибочные входные данные")
	}

	// Количество сообщений за последние N дней
	total, err := hm.ModelManager.GetRecentMessagesCount(days)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка получения статистики")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_messages": total,
	})
}
