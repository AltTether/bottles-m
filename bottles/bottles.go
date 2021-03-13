package bottles

import (
	"time"
	"context"
)


type Message struct {
	text string
}

type Bottle struct {
	message *Message
}

type Engine struct {
	config      *Config
	storage     MessageKeeper
	subscribers map[chan *Bottle]struct{}
}

type MessageKeeper interface {
	Get() (*Message, error)
	Add(*Message) error
}

func (m *Message) Text() string {
	return m.text
}

func (b *Bottle) Message() *Message {
	return b.message
}

func New(cfg *Config, storage MessageKeeper) *Engine {
	return &Engine{
		config:      cfg,
		storage:     storage,
		subscribers: make(map[chan *Bottle]struct{}),
	}
}

func (e *Engine) AddBottle(b *Bottle) {
	if err := e.storage.Add(b.Message()); err != nil {
		return
	}
}

func (e *Engine) SubscribeBottle(ch chan *Bottle) {
	e.subscribers[ch] = struct{}{}
}

func (e *Engine) Run(ctx context.Context) {
	go func() {
		t := time.NewTicker(e.config.SendBottlePeriod)
		defer t.Stop()

	Loop:
		for {
			select {
			case <- ctx.Done():
				break Loop
			case <- t.C:
				if (len(e.subscribers) == 0) {
					break
				}

				for subscriber := range e.subscribers {
					m, err := e.storage.Get()
					if err != nil {
						break
					}

					b := &Bottle{
						message: m,
					}

					subscriber <- b
				}
			default:
				break
			}
		}
	}()
}
