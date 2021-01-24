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

type Token struct {
	Str string
}

type Bottle struct {
	Message *Message
	Token   *Token
}

type Engine struct {
	Ctx        context.Context
	cancelFunc context.CancelFunc
	Config     *Config
	Gateway    *Gateway
	Handlers   map[string]HandlerFunc
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

	handlers := make(map[string]HandlerFunc)
	handlers[ADD_BOTTLE_MODE] = func(ctx context.Context, b *Bottle) {}
	handlers[REQUEST_BOTTLE_MODE] = func(ctx context.Context, b *Bottle) {}
	handlers[GENERATE_BOTTLE_MODE] = func(ctx context.Context, b *Bottle) {}

	return &Engine{
		Ctx:        ctx,
		cancelFunc: cancelFunc,
		Config:     cfg,
		Gateway:    gateway,
		Handlers:   handlers,
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

	handlers := make(map[string]HandlerFunc)
	handlers[ADD_BOTTLE_MODE] = BottleAddHandler(tokenStorage, messageStorage)
	handlers[REQUEST_BOTTLE_MODE] = BottleGetHandler(tokenStorage, messageStorage)
	handlers[GENERATE_BOTTLE_MODE] = BottleGenerateHandler(messageStorage)

	return &Engine{
		Ctx:        ctx,
		cancelFunc: cancelFunc,
		Config:     cfg,
		Gateway:    gateway,
		Handlers:   handlers,
	}
}

func (e *Engine) SetConfig(c *Config) {
	e.Config = c
}

func (e *Engine) AddHandler(mode string, h HandlerFunc) {
	e.Handlers[mode] = h
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
				e.Handlers[ADD_BOTTLE_MODE](e.Ctx, b)
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
					e.Handlers[REQUEST_BOTTLE_MODE](e.Ctx, b)
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
				e.Handlers[GENERATE_BOTTLE_MODE](e.Ctx, b)
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
