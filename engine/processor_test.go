package composer

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

var (
	defaultText = "This is a Message"
	defaultMessage = &Message{
		Text: &defaultText,
	}

	defaultTokenStr = "ToKEN"
	defaultToken = &Token{
		Str: &defaultTokenStr,
	}

	defaultBottle = &Bottle{
		Message: defaultMessage,
		Token:   defaultToken,
	}
)

func TestProcessor(t *testing.T) {
	processor:= New()

	processedBottle, _ := processor.Run(defaultBottle)

	assert.Equal(t, processedBottle, defaultBottle)
}

func TestProcessFunc1(t *testing.T) {
	processor := New()

	replaceToken := "TOken"
	tokenReplacer := func(b *Bottle) (*Bottle, error) {
		b.Token.Str = &replaceToken
		return b, nil
	}
	processor.Use(tokenReplacer)

	processedBottle, _ := processor.Run(defaultBottle)

	assert.Equal(t, *processedBottle.Token.Str, replaceToken)
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

	processedBottle, err := processor.Run(defaultBottle)

	if err != nil {
		assert.NotEqual(t,
			*processedBottle.Message.Text,
			"Func2")
		assert.Equal(t,
			*processedBottle.Message.Text,
			"This is a Message")
	}
}
