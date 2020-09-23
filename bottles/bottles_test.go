package bottles

import (
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

func TestSetHandlersToEngine(t *testing.T) {
	engine := New()

	bottleAddHandlerFunc := func(ctx context.Context, b *Bottle) {}
	bottleGetHandlerFunc := func(ctx context.Context, b *Bottle) {}

	engine.SetBottleAddHandler(bottleAddHandlerFunc)
	engine.SetBottleGetHandler(bottleGetHandlerFunc)
}

func TestRunEngineWithNoHandler(t *testing.T) {
	engine := New()

	engine.Run()
	engine.Stop()
}

func TestAddFromGateway(t *testing.T) {
	engine := New()
	gateway := engine.Gateway

	engine.Run()

	gateway.Add(&Bottle{})

	engine.Stop()
}

func TestGenerateEmptyBottle(t *testing.T) {
	engine := New()

	messagePool := CreateTestMessagePool(0)
	tokenPool := CreateTestTokenPool(0)

	bottleGetHandlerFunc := BottleGetHandler(tokenPool, messagePool)
	bottleGenerateHandlerFunc := BottleGenerateHandler(messagePool)

	engine.SetBottleGetHandler(bottleGetHandlerFunc)
	engine.SetBottleGenerateHandler(bottleGenerateHandlerFunc)

	gateway := engine.Gateway
	engine.Run()
	defer engine.Stop()

	cnt := 0
	timeout := time.After(100 * time.Millisecond)
Loop:
	for {
		select {
		case <-gateway.Get():
			cnt++
		case <-timeout:
			break Loop
		default:
			break
		}
	}

	assert.Greater(t, cnt, 0)
}

func TestDefaultEngine(t *testing.T) {
	engine := DefaultEngine()
	gateway := engine.Gateway

	engine.Run()
	defer engine.Stop()

	n := 10
	tokenStrs := make([]*string, n)
	for i := 0; i < n; i++ {
		bottle := <-gateway.Get()
		tokenStrs[i] = bottle.Token.Str
	}

	for i := 0; i < n; i++ {
		tokenStr := tokenStrs[i]
		token := &Token{
			Str: tokenStr,
		}

		messageText := ""
		message := &Message{
			Text: &messageText,
		}

		b := &Bottle{
			Message: message,
			Token:   token,
		}

		gateway.Add(b)
	}
}
