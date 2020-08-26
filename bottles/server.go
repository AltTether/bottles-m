package bottles

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)


type RequestBody struct {
	Message *string `json:"message" binding:"required"`
	Token   *string `json:"token" binding:"required"`
}

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	getPipeline := NewPipeline()
	postPipeline := NewPipeline()

	messagePool := NewMessagePool()

	tokenPool := NewTokenPool(2 * time.Minute)
	if gin.Mode() == gin.TestMode {
		testTokenStr := "test"
		testToken := &Token{
			Str: &testTokenStr,
		}
		tokenPool.Add(testToken)
	}

	postPipeline.AddStage(ValidateTokenStage(tokenPool))
	postPipeline.AddStage(StoreMessageStage(messagePool))
	
	getPipeline.AddStage(AddTokenStage(tokenPool))
	getPipeline.AddStage(AddMessageStage(messagePool))


	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", func(c *gin.Context) {
			bottle := &Bottle{}
			err := getPipeline.Run(bottle)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": gin.H{
					"text": bottle.Message.Text,
				},
				"token": gin.H{
					"str": bottle.Token.Str,
				},
			})
		})

		v1.POST("/bottle", func(c *gin.Context) {
			var body RequestBody
			if c.BindJSON(&body) != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			bottle := &Bottle{
				Message: &Message{
					Text: body.Message,
				},
				Token:   &Token{
					Str: body.Token,
				},
			}

			err := postPipeline.Run(bottle)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
		})
	}

	return r
}
