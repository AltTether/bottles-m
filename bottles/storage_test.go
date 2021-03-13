package bottles

import (
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)


func TestMessageStorage(t *testing.T) {
	storage := NewStorage()

	text := "This is a Test Message"
	message := &Message{
		text: text,
	}
	
	_ = storage.Add(message)
	messageFromStorage, _ := storage.Get()

	assert.Equal(t, messageFromStorage.Text(), text)
}

func TestGetMessageFromEmptyStorage(t *testing.T) {
	storage := NewStorage()

	message, err := storage.Get()
	assert.EqualError(t, err, "No Messages")
	assert.Nil(t, message)
}

func TestMessageStorageAddAndGetInGoRoutine(t *testing.T) {
	storage := NewStorage()

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
						text: text,
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
