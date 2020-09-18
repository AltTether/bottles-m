package bottles

import (
	"sync"
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
