package bottles

import (
	"fmt"
	"time"
	"sync"
)

type MessageStorage struct {
	messages []*Message
	mux      *sync.Mutex
}

type TokenStorage struct {
	expiration time.Duration
	tokens     *sync.Map
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		mux: &sync.Mutex{},
	}
}

func NewTokenStorage(expiration time.Duration) *TokenStorage {
	return &TokenStorage{
		expiration: expiration,
		tokens:     &sync.Map{},
	}
}

func (p *MessageStorage) Get() (*Message, error) {
	p.mux.Lock()
	if len(p.messages) == 0 {
		p.mux.Unlock()
		return nil, fmt.Errorf("No Messages")
	}

	m := p.messages[0]
	p.messages = p.messages[1:]
	p.mux.Unlock()
	return m, nil
}

func (p *MessageStorage) Add(m *Message) error {
	p.mux.Lock()
	p.messages = append(p.messages, m)
	p.mux.Unlock()
	return nil
}

func (p *TokenStorage) Use(t *Token) (error) {
	if _, ok := p.tokens.LoadAndDelete(t.Str); !ok {
		return fmt.Errorf("Token is Invalid")
	}
	
	return nil
}

func (p *TokenStorage) Add(t *Token) (error) {
	if t.Str == "" {
		return fmt.Errorf("Token is Empty")
	}

	if _, ok := p.tokens.LoadOrStore(t.Str, true); ok {
		return fmt.Errorf("Storage has Same Token")
	}

	go func() {
		time.Sleep(p.expiration)
		p.Use(t)
	}()

	return nil
}
