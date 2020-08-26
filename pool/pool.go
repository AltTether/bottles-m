package pool

import (
	"fmt"
	"time"
	
	"github.com/bottles/engine"
)

type MessagePool struct {
	messages []*engine.Message
}

type TokenPool struct {
	expiration time.Duration
	tokens     map[string]bool
}

func NewMessagePool() *MessagePool {
	return &MessagePool{}
}

func NewTokenPool(expiration time.Duration) *TokenPool {
	return &TokenPool{
		expiration: expiration,
		tokens:     make(map[string]bool),
	}
}

func (p *MessagePool) Get() (*engine.Message, error) {
	if len(p.messages) == 0 {
		return nil, fmt.Errorf("No Messages")
	}

	m := p.messages[0]
	p.messages = p.messages[1:]
	return m, nil
}

func (p *MessagePool) Add(m *engine.Message) error {
	p.messages = append(p.messages, m)
	return nil
}

func (p *TokenPool) Use(t *engine.Token) (error) {
	if _, ok := p.tokens[*t.Str]; !ok {
		return fmt.Errorf("Token is Invalid")
	}

	delete(p.tokens, *t.Str)
	return nil
}

func (p *TokenPool) Add(t *engine.Token) (error) {
	if _, ok := p.tokens[*t.Str]; ok {
		return fmt.Errorf("Pool has Same Token")
	}

	p.tokens[*t.Str] = true
	go func() {
		time.Sleep(p.expiration)
		if _, ok := p.tokens[*t.Str]; ok {
			return
		}
		delete(p.tokens, *t.Str)
	}()

	return nil
}
