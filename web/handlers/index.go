package handlers

import (
	"text/template"

	"github.com/gin-gonic/gin"
)

var index = template.Must(template.ParseFiles("./site/index.html"))

func (hm *HandlerManager) IndexHandler(ctx *gin.Context) {
	index.Execute(ctx.Writer, nil)
}
