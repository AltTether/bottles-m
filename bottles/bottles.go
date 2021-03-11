package bottles

import (
	"time"
	"context"
)

const (
	ADD_BOTTLE_MODE = "add_bottle"
	REQUEST_BOTTLE_MODE = "request_bottle"
	GENERATE_BOTTLE_MODE = "generate_bottle"
)

type Message struct {
	Text string
}

type Bottle struct {
	Message *Message
}

type Engine struct {
	config      *Config
	storage     *Storage
	subscribers map[chan *Bottle]struct{}
}

func New(cfg *Config, storage *Storage) *Engine {
	return &Engine{
		config:      cfg,
		storage:     storage,
		subscribers: make(map[chan *Bottle]struct{}),
	}
}

func waitSend(ch chan *Bottle, bottle *Bottle, delay time.Duration) {
	time.Sleep(delay)
	ch <- bottle
}

func (e *Engine) AddBottle(b *Bottle) {
	if err := e.storage.Add(b.Message); err != nil {
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
						Message: m,
					}

					subscriber <- b
				}
			default:
				break
			}
		}
	}()
}
