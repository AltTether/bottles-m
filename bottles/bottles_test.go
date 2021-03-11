package bottles

import (
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

func TestAddFromGateway(t *testing.T) {
	cfg := NewTestConfig()
	storage := NewStorage()
	engine := New(cfg, storage)

	bottle := &Bottle{}
	engine.AddBottle(bottle)

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()
}

func TestGetFromMakenChan(t *testing.T) {
	cfg := NewTestConfig()
	storage := createTestMessageStorageWithMessages(make([]*Message, 0))
	engine := New(cfg, storage)

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()

	addedBottle := &Bottle{
		Message: &Message{
			Text: "test_text",
		},
	}
	engine.AddBottle(addedBottle)

	bottleOutCh := make(chan *Bottle)
	engine.SubscribeBottle(bottleOutCh)

	gotenBottle := <-bottleOutCh

	assert.Equal(t, gotenBottle.Message.Text, "test_text")
}

func TestBottleGetDeley(t *testing.T) {
	cfg := NewTestConfig()
	messages := make([]*Message, 1)
	messages[0] = &Message{ Text: "test_text" }
	storage := createTestMessageStorageWithMessages(messages)
	engine := New(cfg, storage)

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()

	bottleOutCh := make(chan *Bottle)
	engine.SubscribeBottle(bottleOutCh)

	start := time.Now()
	_ = <-bottleOutCh
	end := time.Now()
	elapsed := end.Sub(start)

	assert.GreaterOrEqual(t, elapsed.Milliseconds(), cfg.SendBottlePeriod.Milliseconds())
}

func CreateTestEngine() *Engine {
	n := 10

	return CreateTestEngineWithData(createTestMessages(n))
}

func CreateTestEngineWithData(messages []*Message) *Engine {
	cfg := NewTestConfig()
	storage := createTestMessageStorageWithMessages(messages)
	engine := New(cfg, storage)

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
