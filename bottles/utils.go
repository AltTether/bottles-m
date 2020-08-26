package bottles

import (
	"math/rand"
)

const (
	LETTERS string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)


func GenerateRandomString(size int) string {
	seed := 42
	r := rand.New(rand.NewSource(int64(seed)))
	l := []rune(LETTERS)
	b := make([]rune, size)
	for i := range b {
		b[i] = l[r.Intn(len(l))]
	}
	t := string(b)
	return t
}
