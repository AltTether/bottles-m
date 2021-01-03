package bottles

import (
	"time"
	"context"
)


type Message struct {
	Text *string
}

type Token struct {
	Str *string
}

type Bottle struct {
	Message *Message
	Token   *Token
}

type Engine struct {
	Ctx                   context.Context
	cancelFunc            context.CancelFunc
	Config                *Config
	Gateway               *Gateway
	BottleAddHandler      HandlerFunc
	BottleGetHandler      HandlerFunc
	BottleGenerateHandler HandlerFunc
}

type HandlerFunc func(ctx context.Context, bottle *Bottle)

type Gateway struct {
	In  chan *Bottle
	Out chan *Bottle
}

func New(cfg *Config) *Engine {
	ctx, cancelFunc := context.WithCancel(context.Background())
	gateway := &Gateway{
		In:  make(chan *Bottle),
		Out: make(chan *Bottle),
	}

	return &Engine{
		Ctx:              ctx,
		cancelFunc:       cancelFunc,
		Config:           cfg,
		Gateway:          gateway,
		BottleAddHandler: func(ctx context.Context, b *Bottle) {},
		BottleGetHandler: func(ctx context.Context, b *Bottle) {},
		BottleGenerateHandler: func(ctx context.Context, b *Bottle) {},
	}
}

func DefaultEngine() *Engine{
	cfg := NewConfig()
	ctx, cancelFunc := context.WithCancel(context.Background())
	gateway := &Gateway{
		In:  make(chan *Bottle),
		Out: make(chan *Bottle),
	}

	messageStorage := NewMessageStorage()
	tokenStorage := NewTokenStorage(cfg.TokenExpiration)

	return &Engine{
		Ctx:                   ctx,
		cancelFunc:            cancelFunc,
		Config:                cfg,
		Gateway:               gateway,
		BottleAddHandler:      BottleAddHandler(tokenStorage, messageStorage),
		BottleGetHandler:      BottleGetHandler(tokenStorage, messageStorage),
		BottleGenerateHandler: BottleGenerateHandler(messageStorage),
	}
}

func (e *Engine) SetConfig(c *Config) {
	e.Config = c
}

func (e *Engine) SetBottleAddHandler(h HandlerFunc) {
	e.BottleAddHandler = h
}

func (e *Engine) SetBottleGetHandler(h HandlerFunc) {
	e.BottleGetHandler = h
}

func (e *Engine) SetBottleGenerateHandler(h HandlerFunc) {
	e.BottleGenerateHandler = h
}

func (e *Engine) Run() {
	go func() {
	Loop:
		for {
			select {
			case <-e.Ctx.Done():
				break Loop
			case b := <-e.Gateway.In:
				e.BottleAddHandler(e.Ctx, b)
			default:
				break
			}
		}
		return
	}()

	go func() {
	Loop:
		for {
			select {
			case <-e.Ctx.Done():
				break Loop
			default:
				b := &Bottle{}
				e.BottleGetHandler(e.Ctx, b)
				if b.Token == nil || b.Message == nil {
					break
				} else {
					e.Gateway.Out <- b
				}
			}
		}
		return
	}()

	go func() {
		t := time.NewTicker(e.Config.GenerateBottleCoolTime)
		defer t.Stop()
	Loop:
		for {
			select {
			case <- e.Ctx.Done():
				break Loop
			case <-t.C:
				text := ""
				b := &Bottle{
					Message: &Message{
						Text: &text,
					},
				}
				e.BottleGenerateHandler(e.Ctx, b)
			default:
				break
			}
		}
	}()
}

func (e *Engine) Stop() {
	e.cancelFunc()
}

func (g *Gateway) Add(b *Bottle) {
	g.In <- b
}

func (g *Gateway) Get() <-chan *Bottle {
	return g.Out
}

func BottleAddHandler(tokenStorage *TokenStorage, messageStorage *MessageStorage) HandlerFunc {
	return func(ctx context.Context, b *Bottle) {
		if err := tokenStorage.Use(b.Token); err != nil {
			return
		}

		if err := messageStorage.Add(b.Message); err != nil {
			return
		}
	}
}

func BottleGetHandler(tokenStorage *TokenStorage, messageStorage *MessageStorage) HandlerFunc {
	size := 10
	seed := 42
	gen := NewRandomStringGenerator(size, seed)

	return func(ctx context.Context, b *Bottle) {
		message, err := messageStorage.Get()
		if err != nil {
			return
		}
		b.Message = message

		tokenStr := gen.Generate()
		token := &Token{
			Str: &tokenStr,
		}
		for tokenStorage.Add(token) != nil {
			tokenStr = gen.Generate()
			token = &Token{
				Str: &tokenStr,
			}
		}
		b.Token = token
	}
}

func BottleGenerateHandler(messageStorage *MessageStorage) HandlerFunc {
	return func(ctx context.Context, b *Bottle) {
		messageStorage.Add(b.Message)
	}
}
