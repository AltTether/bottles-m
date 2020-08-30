package bottles

import (
	"math/rand"
)

const (
	LETTERS string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	seed = 42
	r = rand.New(rand.NewSource(int64(seed)))
)


func GenerateRandomString(size int) string {
	l := []rune(LETTERS)
	b := make([]rune, size)
	for i := range b {
		b[i] = l[r.Intn(len(l))]
	}
	t := string(b)
	return t
}
