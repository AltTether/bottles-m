package bottles

import (
	"time"
	"sync"
	"context"
	"math/rand"
)

const (
	LETTERS string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)


type RandomStringGenerator struct {
	rand *rand.Rand
	size int
	mux  *sync.Mutex
}

func NewRandomStringGenerator(size, seed int) *RandomStringGenerator{
	return &RandomStringGenerator{
		rand: rand.New(rand.NewSource(int64(seed))),
		size: size,
		mux:  &sync.Mutex{},
	}
}

func (g *RandomStringGenerator) Generate() string {
	l := []rune(LETTERS)
	b := make([]rune, g.size)
	g.mux.Lock()
	for i := range b {
		b[i] = l[g.rand.Intn(len(l))]
	}
	g.mux.Unlock()
	t := string(b)
	return t
}

type EmptyMessageAdder struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	messageStorage  *MessageStorage
	intervalTime time.Duration
}

func NewEmptyMessageAdder(messageStorage *MessageStorage, intervalTime time.Duration) *EmptyMessageAdder {
	ctx, cancel := context.WithCancel(context.Background())
	return &EmptyMessageAdder{
		ctx:          ctx,
		cancelFunc:   cancel,
		messageStorage:  messageStorage,
		intervalTime: intervalTime,
	}
}

func (a *EmptyMessageAdder) Run() {
	go func() {
		t := time.NewTicker(a.intervalTime)
		defer t.Stop()
	Loop:
		for {
			select {
			case <-a.ctx.Done():
				break Loop
			case <-t.C:
				text := ""
				m := &Message{
					Text: &text,
				}
				if a.messageStorage.Add(m) != nil {
					break
				}
			default:
				break
			}
		}
	}()
}

func (a *EmptyMessageAdder) Stop() {
	a.cancelFunc()
}
