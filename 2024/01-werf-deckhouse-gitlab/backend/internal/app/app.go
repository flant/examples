package app

import (
	"github.com/gin-gonic/gin"
	"habr_app/internal/common"
	"habr_app/internal/controllers"
)

func Run() {
	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(common.JsonLogger())

	route.GET("/api/remember", controllers.RememberController)
	route.GET("/api/say", controllers.SayController)

	err := route.Run()
	if err != nil {
		return
	}
}
