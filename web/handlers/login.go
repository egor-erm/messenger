package handlers

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var login = template.Must(template.ParseFiles("./site/login.html"))

func (hm *HandlerManager) LoginHandler(ctx *gin.Context) {
	login.Execute(ctx.Writer, nil)
}

type LoginData struct {
	Login    string
	Password string
}

func (hm *HandlerManager) LoginDataHandler(ctx *gin.Context) {
	var loginData LoginData
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := hm.ModelManager.GetUserByLogin(loginData.Login)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Пользователь не найден!",
			"field": "login",
		})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Пароль не подходит!",
			"field": "password",
		})

		return
	}

	userData := loginData.Login + ":" + strconv.Itoa(int(time.Now().Unix()))
	token, err := bcrypt.GenerateFromPassword([]byte(userData), bcrypt.DefaultCost)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	hm.Sessions[userData] = user
	ctx.JSON(http.StatusOK, gin.H{
		"token": string(token),
	})
}
