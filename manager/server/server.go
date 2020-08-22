package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
)

func Run() {
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	registerRoutes(r)

	r.Run()
}
