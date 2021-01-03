package bottles

import (
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

func TestSetHandlersToEngine(t *testing.T) {
	cfg := NewTestConfig()
	engine := New(cfg)

	bottleAddHandlerFunc := func(ctx context.Context, b *Bottle) {}
	bottleGetHandlerFunc := func(ctx context.Context, b *Bottle) {}

	engine.SetBottleAddHandler(bottleAddHandlerFunc)
	engine.SetBottleGetHandler(bottleGetHandlerFunc)
}

func TestRunEngineWithNoHandler(t *testing.T) {
	cfg := NewTestConfig()
	engine := New(cfg)

	engine.Run()
	engine.Stop()
}

func TestAddFromGateway(t *testing.T) {
	cfg := NewTestConfig()
	engine := New(cfg)
	gateway := engine.Gateway

	engine.Run()

	gateway.Add(&Bottle{})

	engine.Stop()
}

func TestGenerateEmptyBottle(t *testing.T) {
	cfg := NewTestConfig()
	engine := New(cfg)

	messageStorage := createTestMessageStorageWithMessages(make([]*Message, 0))
	tokenStorage := createTestTokenStorageWithTokens(make([]*Token, 0), cfg.TokenExpiration)

	bottleGetHandlerFunc := BottleGetHandler(tokenStorage, messageStorage)
	bottleGenerateHandlerFunc := BottleGenerateHandler(messageStorage)

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

	cfg := NewTestConfig()
	engine.SetConfig(cfg)
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

func CreateTestEngine() *Engine {
	n := 10

	return CreateTestEngineWithData(createTestMessages(n), createTestTokens(n))
}

func CreateTestEngineWithData(messages []*Message, tokens []*Token) *Engine {
	cfg := NewTestConfig()
	engine := New(cfg)

	messageStorage := createTestMessageStorageWithMessages(messages)
	tokenStorage := createTestTokenStorageWithTokens(tokens, cfg.TokenExpiration)

	engine.SetBottleGetHandler(BottleGetHandler(tokenStorage, messageStorage))
	engine.SetBottleAddHandler(BottleAddHandler(tokenStorage, messageStorage))

	return engine
}

func createTestMessageStorageWithMessages(ms []*Message) *MessageStorage{
	s := NewMessageStorage()
	for _, m := range ms {
		_ = s.Add(m)
	}

	return s
}

func createTestTokenStorageWithTokens(ts []*Token, expiration time.Duration) *TokenStorage {
	s := NewTokenStorage(expiration)
	for _, t := range ts {
		_ = s.Add(t)
	}

	return s
}

func createTestMessages(n int) []*Message {
	ms := make([]*Message, n)
	for i := 0; i < n; i++ {
		text := testMessageText
		ms[i] = &Message{
			Text: &text,
		}
	}

	return ms
}

func createTestTokens(n int) []*Token {
	ts := make([]*Token, n)
	for i := 0; i < n; i++ {
		str := testTokenStrFormatter(i)
		ts[i] = &Token{
			Str: &str,
		}
	}

	return ts
}
