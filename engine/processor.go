package engine

import ()


type Message struct {
	Text *string
}

type Bottle struct {
	Message *Message
}

type ProcessFunc func(*Bottle) (*Bottle, error)

type Processor struct {
	processFuncs []ProcessFunc
}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) Run(b *Bottle) (*Bottle, error) {
	var err error
	for _, f := range p.processFuncs {
		b, err = f(b)
		if err != nil {
			return b, err
		}
	}
	return b, nil
}

func (p *Processor) Use(f ProcessFunc) {
	p.processFuncs = append(p.processFuncs, f)
}
