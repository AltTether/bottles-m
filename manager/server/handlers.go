package server

import (
	"os"
	"io"
	"log"
	"fmt"
	"time"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/clients"
)

type RequestBottle struct {
	Message string `form:"message" json:"message"`
	Token string `form:"token" json:"token"`
}

type ResponseBottle struct {
	Message *clients.Message `json:"message"`
	Token *clients.Token `json:"token"`
}

type Handlers struct {
	messageClient *clients.MessageClient
	tokenClient *clients.TokenClient
}

func NewHandlers() (*Handlers) {
	tokenClient := clients.NewTokenClient()
	messageClient := clients.NewMessageClient()
	return &Handlers{
		tokenClient: tokenClient,
		messageClient: messageClient,
	}
}

func (m *Handlers) Post(c *gin.Context) {
	var bottle RequestBottle
	if c.BindJSON(&bottle) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf("Get inputs message=%s token=%s",
		bottle.Message, bottle.Token)

	token := &clients.Token{
		Str: bottle.Token,
	}
	err := m.tokenClient.Post(token)
	if err != nil {
		c.String(http.StatusBadRequest, "token is invalid")
		return
	}

	message := &clients.Message{
		Text: bottle.Message,
	}
	err = m.messageClient.Post(message)
	if err != nil {
		c.String(http.StatusBadRequest, "message is invalid")
		return
	}

	c.Status(http.StatusOK)
}

func (m *Handlers) Get(c *gin.Context) {
	w := c.Writer
	r := c.Request
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cached-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	coolTimeSec, err := strconv.Atoi(os.Getenv("MANAGER_COOL_TIME_SEC"))
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(time.Duration(coolTimeSec) * time.Second)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				token, err := m.tokenClient.Get()
				if err != nil {
					break
				}

				message, err := m.messageClient.Get()
				if err != nil {
					break
				}

				bottle := &ResponseBottle{
					Message: message,
					Token: token,
				}
				bottleJson, err := json.Marshal(bottle)
				if err != nil {
					break
				}

				fmt.Fprintf(w, "event: ping\ndata: %s\n\n", string(bottleJson))
				flusher.Flush()
			}
		}
	}()
	<-r.Context().Done()
}

func (m *Handlers) Stream(c *gin.Context) {
	coolTimeSec, err := strconv.Atoi(os.Getenv("MANAGER_COOL_TIME_SEC"))
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(time.Duration(coolTimeSec) * time.Second)
	defer t.Stop()

	clientGone := c.Writer.CloseNotify()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case <-t.C:
			token, err := m.tokenClient.Get()
			if err != nil {
				break
			}

			message, err := m.messageClient.Get()
			if err != nil {
				break
			}

			bottle := &ResponseBottle{
				Message: message,
				Token: token,
			}
			bottleJson, err := json.Marshal(bottle)
			if err != nil {
				break
			}
			
			c.SSEvent("ping", string(bottleJson))
			return true
		}
		return true
	})
}
