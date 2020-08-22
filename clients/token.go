package clients

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
)


type Token struct {
	Str string `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type TokenClient struct {
	addr string
}

func NewTokenClient() (*TokenClient) {
	host := os.Getenv("TOKEN_HOST")
	port := os.Getenv("TOKEN_PORT")
	addr := fmt.Sprintf(
		"http://%s:%s", host, port)

	return &TokenClient{
		addr:  addr,
	}
}

func (c *TokenClient) Post(t *Token) (error) {
	jsonStr, err := json.Marshal(t)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.addr,
		"application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (c *TokenClient) Get() (*Token, error) {
	resp, err := http.Get(c.addr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("")
	}
	defer resp.Body.Close()

	token := &Token{}
	err = json.NewDecoder(resp.Body).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
