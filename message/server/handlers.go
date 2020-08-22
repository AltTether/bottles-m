package server

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"

	"github.com/message/pool"
	"github.com/message/pool/mysql"
)


type RequestBody struct {
	Text *string `json:"message"`
}

type Handlers struct {
	pool pool.Pool
}

func NewHandlers() (*Handlers) {
	pool := mysql.New()
	return NewHandlersWithPool(pool)
}

func NewHandlersWithPool(pool pool.Pool) (*Handlers) {
	return &Handlers{
		pool: pool,
	}
}

func (h *Handlers) Get(c *gin.Context) {
	message, err := h.pool.Get()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	logMessage := fmt.Sprintf("OUTPUT message=%s", *message.Text())
	c.Set("message", logMessage)

	c.JSON(http.StatusOK, gin.H{
		"message": message.Text(),
	})
}

func (h *Handlers) Post(c *gin.Context) {
	var b RequestBody
	if err := c.BindJSON(&b); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	message := mysql.NewMessage(b.Text)
	if h.pool.Post(message) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	logMessage := fmt.Sprintf("INPUT message=%s", *message.Text())
	c.Set("message", logMessage)

	c.Status(http.StatusOK)
}
