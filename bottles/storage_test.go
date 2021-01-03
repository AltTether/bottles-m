package bottles

import (
	"fmt"
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)


func TestMessageStorage(t *testing.T) {
	storage := NewMessageStorage()

	text := "This is a Test Message"
	message := &Message{
		Text: &text,
	}
	
	_ = storage.Add(message)
	messageFromStorage, _ := storage.Get()

	assert.Equal(t, *messageFromStorage.Text, text)
}

func TestMessageStoragePostNilMessageText(t *testing.T) {
	storage := NewMessageStorage()

	message := &Message{
		Text: nil,
	}

	err := storage.Add(message)
	assert.Equal(t, fmt.Errorf("Message Text is Nil"), err)
}

func TestGetMessageFromEmptyStorage(t *testing.T) {
	storage := NewMessageStorage()

	message, err := storage.Get()
	assert.Equal(t, err, fmt.Errorf("No Messages"))
	assert.Nil(t, message)
}

func TestMessageStorageAddAndGetInGoRoutine(t *testing.T) {
	storage := NewMessageStorage()

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
					storage.Add(message)
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
					if _, err := storage.Get(); err == nil {
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

func TestTokenStorage(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewTokenStorage(cfg.TokenExpiration)

	tokenStr := "TesT"
	token := &Token{
		Str: &tokenStr,
	}

	_ = storage.Add(token)
	err := storage.Use(token)
	assert.Nil(t, err)
}

func TestTokenStorageAddAndUseInGoRoutine(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewTokenStorage(cfg.TokenExpiration)

	gen := NewRandomStringGenerator(cfg.TokenSize, cfg.Seed)

	n := 10
	for i := 0; i < n; i++ {
		tokenStr := gen.Generate()
		go func() {
			token := &Token{
				Str: &tokenStr,
			}
			storage.Add(token)
		}()

		go func() {
			token := &Token{
				Str: &tokenStr,
			}
			storage.Use(token)
		}()
	}
}

func TestTokenStorageInvalidToken(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewTokenStorage(cfg.TokenExpiration)

	tokenStr := "TesT"
	token := &Token{
		Str: &tokenStr,
	}

	err := storage.Use(token)
	assert.Equal(t, fmt.Errorf("Token is Invalid"), err)
}

func TestTokenStorageSameToken(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewTokenStorage(cfg.TokenExpiration)

	tokenStr1 := "TesT"
	tokenStr2 := "TesT"
	token1 := &Token{
		Str: &tokenStr1,
	}
	token2 := &Token{
		Str: &tokenStr2,
	}

	_ = storage.Add(token1)
	err := storage.Add(token2)
	assert.Equal(t, err, fmt.Errorf("Storage has Same Token"))
}

func TestTokenStorageTokenExpiration(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewTokenStorage(cfg.TokenExpiration)

	tokenStr := "TesT"
	token := &Token{
		Str: &tokenStr,
	}

	_ = storage.Add(token)
	time.Sleep(50 * time.Millisecond)
	err := storage.Use(token)

	assert.Equal(t, fmt.Errorf("Token is Invalid"), err)
}

func TestTokenStorageAddNilToken(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewTokenStorage(cfg.TokenExpiration)

	token := &Token{
		Str: nil,
	}

	err := storage.Add(token)

	assert.Equal(t, fmt.Errorf("Token is Nil"), err)
}
