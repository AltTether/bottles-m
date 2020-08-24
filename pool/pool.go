package pool

import (
	"fmt"
	
	"github.com/bottles/engine"
)

type MessagePool struct {
	messages []*engine.Message
}

func NewMessagePool() *MessagePool {
	return &MessagePool{}
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
