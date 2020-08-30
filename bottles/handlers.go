package bottles;

import (
	"io"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetBottleHandlerFunc(pipeline *Pipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
		bottle := &Bottle{}
		err := pipeline.Run(bottle)
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
	}
}

func GetBottleStreamHandlerFunc(pipeline *Pipeline) gin.HandlerFunc {
	sendDelay := time.Duration(10 * time.Millisecond)
	return func(c *gin.Context) {
		ticker := time.NewTicker(sendDelay)
		defer ticker.Stop()

		clientGone := c.Writer.CloseNotify()
		c.Stream(func(w io.Writer) bool {
			select {
			case <-clientGone:
				return false
			case <-ticker.C:
				bottle := &Bottle{}
				if pipeline.Run(bottle) != nil {
					return true
				}

				c.SSEvent("bottle", gin.H{
					"message": gin.H{
						"text": bottle.Message.Text,
					},
					"token": gin.H{
						"str": bottle.Token.Str,
					},
				})
				return true
			default:
				return true
			}
		})
	}
}

func PostBottleHandlerFunc(pipeline *Pipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		go pipeline.Run(bottle)

		c.Status(http.StatusOK)
	}
}
