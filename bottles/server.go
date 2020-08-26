package bottles

import (
	"time"
	"net/http"
	"math/rand"

	"github.com/gin-gonic/gin"
)

const (
	LETTERS string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
	
	messageGetter := func(b *Bottle) (error) {
		message, err := messagePool.Get()
		if err != nil {
			return err
		}
		b.Message = message

		return nil
	}
	tokenAdder := func(b *Bottle) (error) {
		size := 10
		tokenStr := GenerateRandomString(size)
		token := &Token{
			Str: &tokenStr,
		}
		for tokenPool.Add(token) != nil {
			tokenStr = GenerateRandomString(size)
			token = &Token{
				Str: &tokenStr,
			}
		}
		b.Token = token
		return nil
	}
	getPipeline.AddStage(tokenAdder)
	getPipeline.AddStage(messageGetter)


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

func GenerateRandomString(size int) string {
	seed := 42
	r := rand.New(rand.NewSource(int64(seed)))
	l := []rune(LETTERS)
	b := make([]rune, size)
	for i := range b {
		b[i] = l[r.Intn(len(l))]
	}
	t := string(b)
	return t
}
