package bottles

import (
	"time"

	"github.com/gin-gonic/gin"
)


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

func DefaultPipelines() (*Pipeline, *Pipeline) {
	getPipeline := NewPipeline()
	postPipeline := NewPipeline()

	messagePool := NewMessagePool()

	tokenPool := NewTokenPool(2 * time.Minute)
	if gin.Mode() == gin.TestMode {
		testTokenStr := "test"
		testToken := &Token{
			Str: &testTokenStr,
		}
		tokenPool.Add(testToken)
	}

	postPipeline.AddStage(ValidateTokenStage(tokenPool))
	postPipeline.AddStage(StoreMessageStage(messagePool))

	getPipeline.AddStage(AddTokenStage(tokenPool))
	getPipeline.AddStage(AddMessageStage(messagePool))

	return getPipeline, postPipeline
}
