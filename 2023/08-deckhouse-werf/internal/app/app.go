package app

import (
	"habr_app/internal/common"
	"habr_app/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(common.JsonLogger())

	route.LoadHTMLGlob("templates/*")

	route.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	route.GET("/remember", controllers.RememberController)
	route.GET("/say", controllers.SayController)

	err := route.Run()
	if err != nil {
		return
	}
}
