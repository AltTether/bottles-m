package server

import (
	"github.com/gin-gonic/gin"
)


func registerRouters(r *gin.Engine) {
	h := NewHandlers()

	r.GET("/", h.Get)
	r.POST("/", h.Post)
}
