package app

import (
	"app/internal/common"
	"app/internal/controllers"
	"github.com/gin-gonic/gin"
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
