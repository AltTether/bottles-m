package bottles

import (
	"time"
	"context"
)

const (
	ADD_BOTTLE_MODE = "add_bottle"
	REQUEST_BOTTLE_MODE = "request_bottle"
)

type Message struct {
	Text string
}

type Token struct {
	Str string
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
	In  chan *Query
}

type Query struct {
	Mode string
	Data interface{}
}

type AddBottleQuery struct {
	Bottle *Bottle
}

type RequestBottleQuery struct {
	OutCh chan *Bottle
}

func New(cfg *Config) *Engine {
	ctx, cancelFunc := context.WithCancel(context.Background())
	gateway := &Gateway{
		In:  make(chan *Query),
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
		In:  make(chan *Query),
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
	addHandlerCh := make(chan interface{})
	getHandlerCh := make(chan interface{})

	go func() {
	Loop:
		for {
			select {
			case <- e.Ctx.Done():
				break Loop
			case q := <- e.Gateway.In:
				if (q.Mode == ADD_BOTTLE_MODE) {
					addHandlerCh <- q.Data
				}

				if (q.Mode == REQUEST_BOTTLE_MODE) {
					getHandlerCh <- q.Data
				}
			default:
				break
			}
		}
	}()

	go func() {
	Loop:
		for {
			select {
			case <-e.Ctx.Done():
				break Loop
			case data := <- addHandlerCh:
				b, ok := data.(*Bottle)
				if (!ok) {
					break
				}
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
			case data := <- getHandlerCh:
				ch, ok := data.(chan *Bottle)
				if (!ok) {
					break
				}
				for {
					b := &Bottle{}
					e.BottleGetHandler(e.Ctx, b)
					if b.Token == nil || b.Message == nil {
						continue
					} else {
						go waitSend(ch, b, e.Config.SendBottleDelay)
						break
					}
				}
			default:
				break
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
						Text: text,
					},
				}
				e.BottleGenerateHandler(e.Ctx, b)
			default:
				break
			}
		}
	}()
}

func waitSend(ch chan *Bottle, bottle *Bottle, delay time.Duration) {
	time.Sleep(delay)
	ch <- bottle
}

func (e *Engine) Stop() {
	e.cancelFunc()
}

func (g *Gateway) AddBottle(bottle *Bottle) {
	q := &Query{
		Mode: ADD_BOTTLE_MODE,
		Data: bottle,
	}
	g.In <- q
}

func (g *Gateway) RequestBottle(ch chan *Bottle) {
	q := &Query{
		Mode: REQUEST_BOTTLE_MODE,
		Data: ch,
	}
	g.In <- q
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
			Str: tokenStr,
		}
		for tokenStorage.Add(token) != nil {
			tokenStr = gen.Generate()
			token = &Token{
				Str: tokenStr,
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
