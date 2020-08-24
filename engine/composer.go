package composer

import ()


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

type ProcessFunc func(*Bottle) *Bottle

type Processor struct {
	processFuncs []ProcessFunc
}

func New() *Processor {
	return &Processor{}
}

func (c *Processor) Run(b *Bottle) *Bottle {
	for _, f := range c.processFuncs {
		b = f(b)
	}
	return b
}

func (c *Processor) Use(f ProcessFunc) {
	c.processFuncs = append(c.processFuncs, f)
}
