package bottles

import (
	"fmt"
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

	bottle := &Bottle{}
	gateway.AddBottle(bottle)

	engine.Stop()
}

func TestGetFromMakenChan(t *testing.T) {
	cfg := NewTestConfig()
	engine := New(cfg)

	messageStorage := createTestMessageStorageWithMessages(make([]*Message, 0))
	tokens := make([]*Token, 1)
	tokens[0] = &Token{ Str: "test_token" }
	tokenStorage := createTestTokenStorageWithTokens(tokens, cfg.TokenExpiration)

	bottleAddHandlerFunc := BottleAddHandler(tokenStorage, messageStorage)
	bottleGetHandlerFunc := BottleGetHandler(tokenStorage, messageStorage)

	engine.SetBottleGetHandler(bottleGetHandlerFunc)
	engine.SetBottleAddHandler(bottleAddHandlerFunc)

	gateway := engine.Gateway
	engine.Run()
	defer engine.Stop()

	addedBottle := &Bottle{
		Message: &Message{
			Text: "test_text",
		},
		Token: &Token{
			Str: "test_token",
		},
	}
	gateway.AddBottle(addedBottle)

	bottleOutCh := make(chan *Bottle)
	gateway.RequestBottle(bottleOutCh)

	gotenBottle := <-bottleOutCh

	assert.Equal(t, gotenBottle.Message.Text, "test_text")
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

	bottleOutCh := make(chan *Bottle)
	gateway.RequestBottle(bottleOutCh)
Loop:
	for {
		select {
		case <-bottleOutCh:
			cnt++
			gateway.RequestBottle(bottleOutCh)
		case <-timeout:
			break Loop
		default:
			break
		}
	}

	assert.Greater(t, cnt, 0)
}

func TestBottleGetDeley(t *testing.T) {
	cfg := NewTestConfig()
	engine := New(cfg)

	messages := make([]*Message, 1)
	messages[0] = &Message{ Text: "test_text" }
	messageStorage := createTestMessageStorageWithMessages(messages)

	tokens := make([]*Token, 1)
	tokens[0] = &Token{ Str: "test_token" }
	tokenStorage := createTestTokenStorageWithTokens(tokens, cfg.TokenExpiration)

	bottleGetHandlerFunc := BottleGetHandler(tokenStorage, messageStorage)
	engine.SetBottleGetHandler(bottleGetHandlerFunc)

	gateway := engine.Gateway

	engine.Run()
	defer engine.Stop()

	bottleOutCh := make(chan *Bottle)
	gateway.RequestBottle(bottleOutCh)

	start := time.Now()
	_ = <-bottleOutCh
	end := time.Now()
	elapsed := end.Sub(start)

	assert.GreaterOrEqual(t, elapsed.Milliseconds(), cfg.SendBottleDelay.Milliseconds())
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
		text := "test"
		ms[i] = &Message{
			Text: text,
		}
	}

	return ms
}

func createTestTokens(n int) []*Token {
	ts := make([]*Token, n)
	for i := 0; i < n; i++ {
		str := fmt.Sprintf("test%d", i)
		ts[i] = &Token{
			Str: str,
		}
	}

	return ts
}
