package pool

import (
	"time"
)

type Token struct {
	Str string `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Pool interface {
	Generate() *Token
	Use(*Token) error
}
