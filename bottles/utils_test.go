package bottles

import (
	"time"
	"sync"
	"testing"
	"math/rand"

	"github.com/stretchr/testify/assert"
)


func TestGenerateRandomString(t *testing.T) {
	size := 10
	seed := 42

	gen := NewRandomStringGenerator(size, seed)
	s := gen.Generate()
	assert.Equal(t, size, len(s))
}

func TestGenerateRandomStrings(t *testing.T) {
	size := 10
	seed := 42
	gen := NewRandomStringGenerator(size, seed)

	n := 5
	randomStrings := make([]string, n)
	for i := 0; i < n; i++ {
		s := gen.Generate()
		assert.NotContains(t, randomStrings, s)
		randomStrings[i] = s
	}
}

func TestGenerateRandomStringsInGoroutine(t *testing.T) {
	size := 10
	seed := 42
	gen := NewRandomStringGenerator(size, seed)

	n := 5
	randomStrings := make([]string, n)
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		idx := i
		go func() {
			s := gen.Generate()
			randomStrings[idx] = s
			wg.Done()
		}()
	}

	wg.Wait()

	copyRandomStrings := make([]string, n)
	for i := 0; i < n; i++ {
		assert.NotContains(t, copyRandomStrings, randomStrings[i])
		copyRandomStrings[i] = randomStrings[i]
	}
}

func TestEmptyMessageAdder(t *testing.T) {
	messageStorage := createTestMessageStorageWithMessages(make([]*Message, 0))
	intervalTime := 50 * time.Millisecond
	gen := NewEmptyMessageAdder(messageStorage, intervalTime)

	gen.Run()
	time.Sleep(100 * time.Millisecond)
	gen.Stop()

	cnt := 0
	_, err := messageStorage.Get()
	for err == nil {
		cnt++;
		_, err = messageStorage.Get()
	}

	assert.Less(t, 0, cnt)
}

const (
	BENCH_STRING_SIZE = 10
)

func BenchmarkRandLetterLock(b *testing.B) {
	mux := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	l := []rune(LETTERS)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < BENCH_STRING_SIZE; i++ {
				mux.Lock()
				_ = l[rand.Intn(len(l))]
				mux.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRandStringLock(b *testing.B) {
	mux := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	l := []rune(LETTERS)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			mux.Lock()
			for i := 0; i < BENCH_STRING_SIZE; i++ {
				_ = l[rand.Intn(len(l))]
			}
			mux.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}
