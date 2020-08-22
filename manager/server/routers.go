package server

import (
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	h := NewHandlers()

	r.GET("/api", h.Get)
	r.POST("/api", h.Post)
	r.GET("/api/stream", h.Stream)
}
