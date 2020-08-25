package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bottles/pool"
	"github.com/bottles/engine"
)


type RequestBody struct {
	Message *string `json:"message"`
}

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	getPipeline := engine.NewPipeline()
	postPipeline := engine.NewPipeline()

	messagePool := pool.NewMessagePool()

	messageAdder := func(b *engine.Bottle) (error) {
		messagePool.Add(b.Message)
		return nil
	}
	postPipeline.AddStage(messageAdder)
	
	messageGetter := func(b *engine.Bottle) (error) {
		message, err := messagePool.Get()
		if err != nil {
			return err
		}
		b.Message = message

		return nil
	}
	getPipeline.AddStage(messageGetter)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", func(c *gin.Context) {
			bottle := &engine.Bottle{}
			err := getPipeline.Run(bottle)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": gin.H{
					"text": bottle.Message.Text,
				},
			})
		})

		v1.POST("/bottle", func(c *gin.Context) {
			var body RequestBody
			if c.BindJSON(&body) != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			bottle := &engine.Bottle{
				Message: &engine.Message{
					Text: body.Message,
				},
			}

			err := postPipeline.Run(bottle)
			if err != nil {
				c.Status(http.StatusBadRequest)
			}

			c.Status(http.StatusOK)
		})
	}

	return r
}
