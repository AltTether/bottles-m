package pool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	
	"github.com/bottles/engine"
)

func TestMessagePool(t *testing.T) {
	pool := NewMessagePool()

	text := "This is a Test Message"
	message := &engine.Message{
		Text: &text,
	}
	
	_ = pool.Add(message)
	messageFromPool, _ := pool.Get()

	assert.Equal(t, *messageFromPool.Text, text)
}

func TestGetMessageFromEmptyPool(t *testing.T) {
	pool := NewMessagePool()

	message, err := pool.Get()
	assert.Equal(t, err, fmt.Errorf("No Messages"))
	assert.Nil(t, message)
}
