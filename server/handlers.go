package server;

import (
	"io"
	"sync"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bottles/bottles"
)


type RequestBody struct {
	Message *string `json:"message" binding:"required"`
	Token   *string `json:"token" binding:"required"`
}

func GetBottleHandlerFunc(gateway *bottles.Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientGone := c.Writer.CloseNotify()

		bottleOutCh := make(chan *bottles.Bottle)
		gateway.RequestBottle(bottleOutCh)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
		Loop:
			for {
				select {
				case <-clientGone:
					c.Status(http.StatusBadRequest)
					break Loop
				case bottle := <-bottleOutCh:
					c.JSON(http.StatusOK, gin.H{
						"message": gin.H{
							"text": bottle.Message.Text,
						},
						"token": gin.H{
							"str": bottle.Token.Str,
						},
					})
					break Loop
				default:
					break
				}
			}
			wg.Done()
		}()

		wg.Wait()
	}
}

func GetBottleStreamHandlerFunc(gateway *bottles.Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientGone := c.Writer.CloseNotify()
		bottleOutCh := make(chan *bottles.Bottle)
		gateway.RequestBottle(bottleOutCh)

		c.Stream(func(w io.Writer) bool {
			select {
			case <-clientGone:
				return false
			case bottle := <-bottleOutCh:
				c.SSEvent("bottle", gin.H{
					"message": gin.H{
						"text": bottle.Message.Text,
					},
					"token": gin.H{
						"str": bottle.Token.Str,
					},
				})

				gateway.RequestBottle(bottleOutCh)

				return true
			default:
				return true
			}
		})
	}
}

func PostBottleHandlerFunc(gateway *bottles.Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body RequestBody
		if c.BindJSON(&body) != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		bottle := &bottles.Bottle{
			Message: &bottles.Message{
				Text: *body.Message,
			},
			Token:   &bottles.Token{
				Str: *body.Token,
			},
		}

		gateway.AddBottle(bottle)

		c.Status(http.StatusOK)
	}
}
