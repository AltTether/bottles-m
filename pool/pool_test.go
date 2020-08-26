package pool

import (
	"fmt"
	"time"
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

func TestTokenPool(t *testing.T) {
	expiration := 2 * time.Minute
	pool := NewTokenPool(expiration)

	tokenStr := "TesT"
	token := &engine.Token{
		Str: &tokenStr,
	}

	_ = pool.Add(token)
	err := pool.Use(token)
	assert.Nil(t, err)
}

func TestTokenPoolInvalidToken(t *testing.T) {
	expiration := 2 * time.Minute
	pool := NewTokenPool(expiration)

	tokenStr := "TesT"
	token := &engine.Token{
		Str: &tokenStr,
	}

	err := pool.Use(token)
	assert.Equal(t, fmt.Errorf("Token is Invalid"), err)
}

func TestTokenPoolSameToken(t *testing.T) {
	expiration := 2 * time.Minute
	pool := NewTokenPool(expiration)

	tokenStr1 := "TesT"
	tokenStr2 := "TesT"
	token1 := &engine.Token{
		Str: &tokenStr1,
	}
	token2 := &engine.Token{
		Str: &tokenStr2,
	}

	_ = pool.Add(token1)
	err := pool.Add(token2)
	assert.Equal(t, err, fmt.Errorf("Pool has Same Token"))
}

func TestTokenPoolTokenExpiration(t *testing.T) {
	expiration := 10 * time.Millisecond
	pool := NewTokenPool(expiration)

	tokenStr := "TesT"
	token := &engine.Token{
		Str: &tokenStr,
	}

	_ = pool.Add(token)
	time.Sleep(50 * time.Millisecond)
	err := pool.Use(token)

	assert.Equal(t, fmt.Errorf("Token is Invalid"), err)
}
