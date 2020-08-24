package composer

import (
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

	processedBottle := processor.Run(defaultBottle)

	assert.Equal(t, processedBottle, defaultBottle)
}

func TestProcessFunc1(t *testing.T) {
	processor := New()

	replaceToken := "TOken"
	tokenReplacer := func(b *Bottle) *Bottle {
		b.Token.Str = &replaceToken
		return b
	}
	processor.Use(tokenReplacer)

	processedBottle := processor.Run(defaultBottle)

	assert.Equal(t, *processedBottle.Token.Str, replaceToken)
}
