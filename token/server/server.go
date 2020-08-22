package server

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	registerRoutes(r)

	r.Run()
}
