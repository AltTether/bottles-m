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

type ComponentFunc func(*Bottle) *Bottle

type Composer struct {
	components []ComponentFunc
}

func New() *Composer {
	return &Composer{}
}

func (c *Composer) Run(b *Bottle) *Bottle {
	for _, component := range c.components {
		b = component(b)
	}
	return b
}

func (c *Composer) Use(f ComponentFunc) {
	c.components = append(c.components, f)
}
