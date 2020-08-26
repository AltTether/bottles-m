package bottles

import (
	"github.com/gin-gonic/gin"
)


type RequestBody struct {
	Message *string `json:"message" binding:"required"`
	Token   *string `json:"token" binding:"required"`
}

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	getPipeline, postPipeline := DefaultPipelines()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", GetBottleHandlerFunc(getPipeline))
		v1.POST("/bottle", PostBottleHandlerFunc(postPipeline))
	}

	return r
}
