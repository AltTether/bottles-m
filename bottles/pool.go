package bottles

import (
	"fmt"
	"time"
	"sync"
)

type MessagePool struct {
	messages []*Message
	mux      *sync.Mutex
}

type TokenPool struct {
	expiration time.Duration
	tokens     map[string]bool
}

func NewMessagePool() *MessagePool {
	return &MessagePool{
		mux: &sync.Mutex{},
	}
}

func NewTokenPool(expiration time.Duration) *TokenPool {
	return &TokenPool{
		expiration: expiration,
		tokens:     make(map[string]bool),
	}
}

func (p *MessagePool) Get() (*Message, error) {
	p.mux.Lock()
	if len(p.messages) == 0 {
		return nil, fmt.Errorf("No Messages")
	}

	m := p.messages[0]
	p.messages = p.messages[1:]
	p.mux.Unlock()
	return m, nil
}

func (p *MessagePool) Add(m *Message) error {
	if m.Text == nil {
		return fmt.Errorf("Message Text is Nil")
	}

	p.mux.Lock()
	p.messages = append(p.messages, m)
	p.mux.Unlock()
	return nil
}

func (p *TokenPool) Use(t *Token) (error) {
	if _, ok := p.tokens[*t.Str]; !ok {
		return fmt.Errorf("Token is Invalid")
	}

	delete(p.tokens, *t.Str)
	return nil
}

func (p *TokenPool) Add(t *Token) (error) {
	if t.Str == nil {
		return fmt.Errorf("Token is Nil")
	}
	if _, ok := p.tokens[*t.Str]; ok {
		return fmt.Errorf("Pool has Same Token")
	}

	p.tokens[*t.Str] = true
	go func() {
		time.Sleep(p.expiration)
		if _, ok := p.tokens[*t.Str]; !ok {
			return
		}
		delete(p.tokens, *t.Str)
	}()

	return nil
}
