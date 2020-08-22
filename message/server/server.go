package server

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.TestMode)

	r := gin.New()

	r.Use(Logger(), gin.Recovery())

	registerRouters(r)

	r.Run()
}
