package handlers

import (
	"os"
	"log"
	"fmt"
	"time"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/manager/clients"
)

type Bottle struct {
	Message string `form:"message" json:"message"`
	Token string `form:"token" json:"token"`
}

type Bottle_ struct {
	Message *clients.Message `json:"message"`
	Token *clients.Token `json:"token"`
}

type ManagerHandlers struct {
	messageClient *clients.MessageClient
	tokenClient *clients.TokenClient
}

func NewManagerHandlers() (*ManagerHandlers) {
	tokenClient := clients.NewTokenClient()
	messageClient := clients.NewMessageClient()
	return &ManagerHandlers{
		tokenClient: tokenClient,
		messageClient: messageClient,
	}
}

func (m *ManagerHandlers) Post(c *gin.Context) {
	var bottle Bottle
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

func (m *ManagerHandlers) Get(c *gin.Context) {
	w := c.Writer
	r := c.Request
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cached-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

				bottle := &Bottle_{
					Message: message,
					Token: token,
				}
				bottleJson, err := json.Marshal(bottle)
				if err != nil {
					break
				}

				fmt.Fprintf(w, "%s\n", string(bottleJson))
				flusher.Flush()
			}
		}
	}()
	<-r.Context().Done()
}
