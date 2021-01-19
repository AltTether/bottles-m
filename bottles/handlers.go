package bottles;

import (
	"io"
	"time"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)


type RequestBody struct {
	Message *string `json:"message" binding:"required"`
	Token   *string `json:"token" binding:"required"`
}

func GetBottleHandlerFunc(gateway *Gateway, cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(context.Background())

		bottleOutCh := make(chan *Bottle)
		gateway.RequestBottle(bottleOutCh)

		go func() {
			for {
				select {
				case <-ctx.Done():
					c.Status(http.StatusBadRequest)
					return
				case bottle := <-bottleOutCh:
					c.JSON(http.StatusOK, gin.H{
						"message": gin.H{
							"text": bottle.Message.Text,
						},
						"token": gin.H{
							"str": bottle.Token.Str,
						},
					})
					return
				default:
					break
				}
			}
		}()

		time.Sleep(cfg.SendBottleDelay)
		cancel()
	}
}

func GetBottleStreamHandlerFunc(gateway *Gateway, cfg *Config) gin.HandlerFunc {
	sendDelay := time.Duration(cfg.SendBottleDelay)
	return func(c *gin.Context) {
		clientGone := c.Writer.CloseNotify()
		bottleOutCh := make(chan *Bottle)
		gateway.RequestBottle(bottleOutCh)

		c.Stream(func(w io.Writer) bool {
			select {
			case <-clientGone:
				return false
			case bottle := <-bottleOutCh:
				time.Sleep(sendDelay)
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

func PostBottleHandlerFunc(gateway *Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body RequestBody
		if c.BindJSON(&body) != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		bottle := &Bottle{
			Message: &Message{
				Text: *body.Message,
			},
			Token:   &Token{
				Str: *body.Token,
			},
		}

		gateway.AddBottle(bottle)

		c.Status(http.StatusOK)
	}
}
