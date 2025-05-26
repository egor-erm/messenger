package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var register = template.Must(template.ParseFiles("./site/register.html"))

func (hm *HandlerManager) RegisterHandler(ctx *gin.Context) {
	register.Execute(ctx.Writer, nil)
}

type RegisterData struct {
	Login    string
	Password string
	Name     string
	Surname  string
}

func (hm *HandlerManager) RegisterDataHandler(ctx *gin.Context) {
	var registerData RegisterData
	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Выберите другой пароль!",
			"field": "password",
		})

		return
	}

	user, err := hm.ModelManager.CreateUser(registerData.Login, string(hashedPassword), registerData.Name, registerData.Surname)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Выберите другой логин!",
			"field": "login",
		})

		return
	}

	fmt.Println("Новый пользователь: ", user.Login)
	ctx.JSON(http.StatusOK, gin.H{})
}
