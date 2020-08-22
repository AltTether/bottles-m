package utils

import (
	"os"
	"strconv"
	"math/rand"
)

const (
	LETTERS string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type TokenGenerator struct {
	size int
	gen *rand.Rand
}

func NewTokenGenerator() (*TokenGenerator) {
	seed, err := strconv.Atoi(os.Getenv("SEED"))
	if err != nil {
		panic(err)
	}
	gen := rand.New(rand.NewSource(int64(seed)))

	size, err := strconv.Atoi(os.Getenv("TOKEN_LENGTH"))
	if err != nil {
		panic(err)
	}
	
	return &TokenGenerator{
		size: size,
		gen: gen,
	}
}

func (tG *TokenGenerator) Generate() (string) {
	rs1Letters := []rune(LETTERS)
	b := make([]rune, tG.size)
	for i := range b {
		b[i] = rs1Letters[tG.gen.Intn(len(rs1Letters))]
	}
	tokenStr := string(b)

	return tokenStr
}
