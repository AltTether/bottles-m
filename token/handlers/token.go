package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/token/pool"
	"github.com/token/pool/redis"
)

type TokenHandlers struct {
	pool *redis.Pool
}

func New() (*TokenHandlers) {
	pool_ := redis.New()
	return &TokenHandlers{
		pool: pool_,
	}
}

func (tH *TokenHandlers) GetToken(c *gin.Context) {
	log.Printf("START - GET / is Called")

	token := tH.pool.Generate()
	log.Printf("DATA token=%s", token.Str)

	c.JSON(http.StatusOK, token)
	log.Printf("END - GET / is Called")
}

func (th *TokenHandlers) PostToken(c *gin.Context) {
	log.Printf("START - POST / is Called")

	var token *pool.Token
	if c.BindJSON(&token) != nil {
		log.Printf("INPUT str=%s", token.Str)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf("INPUT str=%s", token.Str)

	err := th.pool.Use(token)
	if err != nil {
		log.Printf(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
	log.Printf("END - POST / is Called")
}
