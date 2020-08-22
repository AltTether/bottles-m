package main

import (
	"os"
	"time"
	"strconv"
	"strings"

	"github.com/clients"
)


type Message struct {
	Text string `json:"message"`
}


func main() {
	messageClients := clients.NewMessageClient()

	defaultText := strings.Replace(
		os.Getenv("EMPTY_BOTTLE_DEFAULT_MESSAGE"),
		`\n`,
		"\n",
		-1)

	coolTimeSec, err := strconv.Atoi(os.Getenv("GENERATOR_COOL_TIME_SEC"))
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(time.Duration(coolTimeSec) * time.Second)
	defer t.Stop()
	for {
		select {
		case <- t.C:
			message := &clients.Message{
				Text: defaultText,
			}

			err := messageClients.Post(message)

			if err != nil {
				continue
			}
		}
	}
}
