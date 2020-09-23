package bottles

import (
	"fmt"
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)


func TestMessagePool(t *testing.T) {
	pool := NewMessagePool()

	text := "This is a Test Message"
	message := &Message{
		Text: &text,
	}
	
	_ = pool.Add(message)
	messageFromPool, _ := pool.Get()

	assert.Equal(t, *messageFromPool.Text, text)
}

func TestMessagePoolPostNilMessageText(t *testing.T) {
	pool := NewMessagePool()

	message := &Message{
		Text: nil,
	}

	err := pool.Add(message)
	assert.Equal(t, fmt.Errorf("Message Text is Nil"), err)
}

func TestGetMessageFromEmptyPool(t *testing.T) {
	pool := NewMessagePool()

	message, err := pool.Get()
	assert.Equal(t, err, fmt.Errorf("No Messages"))
	assert.Nil(t, message)
}

func TestMessagePoolAddAndGetInGoRoutine(t *testing.T) {
	pool := NewMessagePool()

	ctx, cancel := context.WithCancel(context.Background())
	cnt := 0
	n := 10
	for i := 0; i < n; i++ {
		go func() {
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()
		Loop:
			for {
				select {
				case <-ctx.Done():
					break Loop
				case <-ticker.C:
					text := "This is a Test Message"
					message := &Message{
						Text: &text,
					}
					pool.Add(message)
				default:
					break
				}
			}
		}()

		go func() {
			ticker := time.NewTicker(1 * time.Millisecond)
			defer ticker.Stop()
		Loop:
			for {
				select {
				case <-ctx.Done():
					break Loop
				case <-ticker.C:
					if _, err := pool.Get(); err == nil {
						cnt++
					}
				default:
					break
				}
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Greater(t, cnt, 0)
}

func TestTokenPool(t *testing.T) {
	expiration := 2 * time.Minute
	pool := NewTokenPool(expiration)

	tokenStr := "TesT"
	token := &Token{
		Str: &tokenStr,
	}

	_ = pool.Add(token)
	err := pool.Use(token)
	assert.Nil(t, err)
}

func TestTokenPoolAddAndUseInGoRoutine(t *testing.T) {
	expiration := 2 * time.Minute
	pool := NewTokenPool(expiration)

	seed := 42
	size := 10
	gen := NewRandomStringGenerator(size, seed)

	n := 10
	for i := 0; i < n; i++ {
		tokenStr := gen.Generate()
		go func() {
			token := &Token{
				Str: &tokenStr,
			}
			pool.Add(token)
		}()

		go func() {
			token := &Token{
				Str: &tokenStr,
			}
			pool.Use(token)
		}()
	}
}

func TestTokenPoolInvalidToken(t *testing.T) {
	expiration := 2 * time.Minute
	pool := NewTokenPool(expiration)

	tokenStr := "TesT"
	token := &Token{
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
	token1 := &Token{
		Str: &tokenStr1,
	}
	token2 := &Token{
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
	token := &Token{
		Str: &tokenStr,
	}

	_ = pool.Add(token)
	time.Sleep(50 * time.Millisecond)
	err := pool.Use(token)

	assert.Equal(t, fmt.Errorf("Token is Invalid"), err)
}

func TestTokenPoolAddNilToken(t *testing.T) {
	expiration := 10 * time.Second
	pool := NewTokenPool(expiration)

	token := &Token{
		Str: nil,
	}

	err := pool.Add(token)

	assert.Equal(t, fmt.Errorf("Token is Nil"), err)
}
