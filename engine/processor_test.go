package composer

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

func TestProcessor(t *testing.T) {
	processor := New()

	bottle := defaultBottle()
	processedBottle, _ := processor.Run(bottle)

	assert.Equal(t, processedBottle, bottle)
}

func TestProcessFunc1(t *testing.T) {
	processor := New()

	replaceMessage := "replaced"
	messageReplacer := func(b *Bottle) (*Bottle, error) {
		b.Message.Text = &replaceMessage
		return b, nil
	}
	processor.Use(messageReplacer)

	bottle := defaultBottle()
	processedBottle, _ := processor.Run(bottle)

	assert.Equal(t, *processedBottle.Message.Text, replaceMessage)
}

func TestProcessFuncError(t *testing.T) {
	processor := New()

	processFunc1 := func(b *Bottle) (*Bottle, error) {
		return b, fmt.Errorf("Func1 Error")
	}
	processFunc2 := func(b *Bottle) (*Bottle, error) {
		text := "Func2"
		b.Message.Text = &text
		return b, nil
	}
	processor.Use(processFunc1)
	processor.Use(processFunc2)

	bottle := defaultBottle()
	processedBottle, err := processor.Run(bottle)

	if err != nil {
		assert.NotEqual(t,
			*processedBottle.Message.Text,
			"Func2")
		assert.Equal(t,
			*processedBottle.Message.Text,
			"This is a Test Message")
	}
}
