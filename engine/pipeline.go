package engine

import ()


type StageFunc func(*Bottle) (error)

type Pipeline struct {
	stages []StageFunc
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) Run(b *Bottle) (error) {
	var err error
	for _, s := range p.stages {
		err = s(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) AddStage(f StageFunc) {
	p.stages = append(p.stages, f)
}
