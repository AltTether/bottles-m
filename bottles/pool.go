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
	tokens     *sync.Map
}

func NewMessagePool() *MessagePool {
	return &MessagePool{
		mux: &sync.Mutex{},
	}
}

func NewTokenPool(expiration time.Duration) *TokenPool {
	return &TokenPool{
		expiration: expiration,
		tokens:     &sync.Map{},
	}
}

func (p *MessagePool) Get() (*Message, error) {
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
	if _, ok := p.tokens.LoadAndDelete(*t.Str); !ok {
		return fmt.Errorf("Token is Invalid")
	}
	
	return nil
}

func (p *TokenPool) Add(t *Token) (error) {
	if t.Str == nil {
		return fmt.Errorf("Token is Nil")
	}

	if _, ok := p.tokens.LoadOrStore(*t.Str, true); ok {
		return fmt.Errorf("Pool has Same Token")
	}

	go func() {
		time.Sleep(p.expiration)
		p.Use(t)
	}()

	return nil
}
