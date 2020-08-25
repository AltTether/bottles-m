package engine

import (
	"fmt"
	"testing"
	
	"github.com/stretchr/testify/assert"
)


func defaultBottle() *Bottle {
	text := "This is a Test Message"
	message := &Message{
		Text: &text,
	}
	bottle := &Bottle{
		Message: message,
	}
	return bottle;
}

func TestPipeline(t *testing.T) {
	pipeline := New()

	bottle := defaultBottle()
	_ = pipeline.Run(bottle)

	assert.Equal(t,
		"This is a Test Message",
		*bottle.Message.Text)
}

func TestStageFunc1(t *testing.T) {
	pipeline := New()

	replaceMessage := "replaced"
	messageReplacer := func(b *Bottle) (error) {
		b.Message.Text = &replaceMessage
		return nil
	}
	pipeline.AddStage(messageReplacer)

	bottle := defaultBottle()
	_ = pipeline.Run(bottle)

	assert.Equal(t, *bottle.Message.Text, replaceMessage)
}

func TestStageFuncError(t *testing.T) {
	pipeline := New()

	stageFunc1 := func(b *Bottle) (error) {
		return fmt.Errorf("Func1 Error")
	}
	stageFunc2 := func(b *Bottle) (error) {
		text := "Func2"
		b.Message.Text = &text
		return nil
	}
	pipeline.AddStage(stageFunc1)
	pipeline.AddStage(stageFunc2)

	bottle := defaultBottle()
	err := pipeline.Run(bottle)

	if err != nil {
		assert.NotEqual(t,
			*bottle.Message.Text,
			"Func2")
		assert.Equal(t,
			*bottle.Message.Text,
			"This is a Test Message")
	}
}
