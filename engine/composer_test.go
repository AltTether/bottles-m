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

func TestConposer(t *testing.T) {
	composer := New()

	composedB := composer.Run(defaultBottle)

	assert.Equal(t, composedB, defaultBottle)
}

func TestComposerRun1(t *testing.T) {
	composer := New()

	replaceToken := "TOken"
	tokenReplacer := func(b *Bottle) *Bottle {
		b.Token.Str = &replaceToken
		return b
	}
	composer.Use(tokenReplacer)

	composedB := composer.Run(defaultBottle)

	assert.Equal(t, *composedB.Token.Str, replaceToken)
}
