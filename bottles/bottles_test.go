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

	engine.AddHandler(ADD_BOTTLE_MODE, bottleAddHandlerFunc)
	engine.AddHandler(REQUEST_BOTTLE_MODE, bottleGetHandlerFunc)
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

	bottleAddHandlerFunc := BottleAddHandler(messageStorage)
	bottleGetHandlerFunc := BottleGetHandler(messageStorage)

	engine.AddHandler(REQUEST_BOTTLE_MODE, bottleGetHandlerFunc)
	engine.AddHandler(ADD_BOTTLE_MODE, bottleAddHandlerFunc)

	gateway := engine.Gateway
	engine.Run()
	defer engine.Stop()

	addedBottle := &Bottle{
		Message: &Message{
			Text: "test_text",
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

	bottleGetHandlerFunc := BottleGetHandler(messageStorage)
	bottleGenerateHandlerFunc := BottleGenerateHandler(messageStorage)

	engine.AddHandler(REQUEST_BOTTLE_MODE, bottleGetHandlerFunc)
	engine.AddHandler(GENERATE_BOTTLE_MODE, bottleGenerateHandlerFunc)

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

	bottleGetHandlerFunc := BottleGetHandler(messageStorage)
	engine.AddHandler(REQUEST_BOTTLE_MODE, bottleGetHandlerFunc)

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

	return CreateTestEngineWithData(createTestMessages(n))
}

func CreateTestEngineWithData(messages []*Message) *Engine {
	cfg := NewTestConfig()
	engine := New(cfg)

	messageStorage := createTestMessageStorageWithMessages(messages)

	engine.AddHandler(REQUEST_BOTTLE_MODE, BottleGetHandler(messageStorage))
	engine.AddHandler(ADD_BOTTLE_MODE, BottleAddHandler(messageStorage))

	return engine
}

func createTestMessageStorageWithMessages(ms []*Message) *Storage{
	s := NewStorage()
	for _, m := range ms {
		_ = s.Add(m)
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
