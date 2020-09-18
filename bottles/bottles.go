package bottles

import (
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
	Ctx              context.Context
	cancelFunc       context.CancelFunc
	Gateway          *Gateway
	BottleAddHandler HandlerFunc
	BottleGetHandler HandlerFunc
}

type HandlerFunc func(ctx context.Context, bottle *Bottle)

type Gateway struct {
	In  chan *Bottle
	Out chan *Bottle
}

func New() *Engine {
	ctx, cancelFunc := context.WithCancel(context.Background())
	gateway := &Gateway{
		In:  make(chan *Bottle),
		Out: make(chan *Bottle),
	}

	return &Engine{
		Ctx:              ctx,
		cancelFunc:       cancelFunc,
		Gateway:          gateway,
		BottleAddHandler: func(ctx context.Context, b *Bottle) {},
		BottleGetHandler: func(ctx context.Context, b *Bottle) {},
	}
}

func (e *Engine) SetBottleAddHandler(h HandlerFunc) {
	e.BottleAddHandler = h
}

func (e *Engine) SetBottleGetHandler(h HandlerFunc) {
	e.BottleGetHandler = h
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

func BottleAddHandler(tokenPool *TokenPool, messagePool *MessagePool) HandlerFunc {
	return func(ctx context.Context, b *Bottle) {
		if err := tokenPool.Use(b.Token); err != nil {
			return
		}

		if err := messagePool.Add(b.Message); err != nil {
			return
		}
	}
}

func BottleGetHandler(tokenPool *TokenPool, messagePool *MessagePool) HandlerFunc {
	return func(ctx context.Context, b *Bottle) {
		message, err := messagePool.Get()
		if err != nil {
			return
		}
		b.Message = message

		size := 10
		tokenStr := GenerateRandomString(size)
		token := &Token{
			Str: &tokenStr,
		}
		for tokenPool.Add(token) != nil {
			tokenStr = GenerateRandomString(size)
			token = &Token{
				Str: &tokenStr,
			}
		}
		b.Token = token
	}
}