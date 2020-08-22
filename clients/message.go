package clients

import (
	"os"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
)


type Message struct {
	Text string `json:"message"`
}

type MessageClient struct {
	addr string
}

func NewMessageClient() (*MessageClient) {
	host := os.Getenv("MESSAGE_HOST")
	port := os.Getenv("MESSAGE_PORT")
	addr := fmt.Sprintf(
		"http://%s:%s", host, port)

	return &MessageClient{
		addr:  addr,
	}
}

func (c *MessageClient) Post(m *Message) (error) {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.addr,
		"application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid message")
	}

	return nil
}

func (c *MessageClient) Get() (*Message, error) {
	resp, err := http.Get(c.addr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("")
	}
	defer resp.Body.Close()

	message := &Message{}
	err = json.NewDecoder(resp.Body).Decode(message)
	if err != nil {
		return nil, err
	}

	return message, nil

}
