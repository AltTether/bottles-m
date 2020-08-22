package redis

import (
	"os"
	"fmt"
	"log"
	"time"
	"strings"
	"strconv"
	"context"
	
	"github.com/go-redis/redis"

	"github.com/token/pool"
	"github.com/token/utils"
)


type Pool struct {
	client *redis.Client
	expiration time.Duration
	generator *utils.TokenGenerator
}

func New() (*Pool) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	opt := &redis.Options{
		Addr: addr,
		Password: "",
		DB: 0,
	}

	reconnSec, err := strconv.Atoi(
		os.Getenv("DB_RECONNECTION_SEC"))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	for err != nil {
		log.Printf("Waiting for %dsec", time.Second)
		time.Sleep(time.Second * time.Duration(reconnSec))
		_, err = client.Ping(ctx).Result()
	}

	expirationSec, err := strconv.Atoi(
		os.Getenv("TOKEN_EXPIRATION_SEC"))
	if err != nil {
		panic(err)
	}
	expiration := time.Duration(expirationSec) * time.Second

	generator := utils.NewTokenGenerator()

	return &Pool{
		client: client,
		expiration: expiration,
		generator: generator,
	}
}

func (p *Pool) Generate() (*pool.Token) {
	str := p.generator.Generate()
	for p.has(str) {
		str = p.generator.Generate()
	}
	p.register(str)

	createdAt := time.Now().UTC()
	deletedAt := createdAt.Add(p.expiration)
	token := &pool.Token{
		Str: str,
		CreatedAt: createdAt,
		DeletedAt: deletedAt,
	}
	return token
}

func (p *Pool) has(str string) (bool) {
	ctx := context.Background()
	val, err := p.client.Get(ctx, str).Result()
	if err != nil {
		return false
	}

	if strings.Compare(str, val) != 0 {
		return false
	}

	return true
}

func (p *Pool) register(str string) {
	ctx := context.Background()
	_, err := p.client.Set(ctx, str, str, p.expiration).Result()
	if err != nil {
		panic(err)
	}
}

func (p *Pool) Use(token *pool.Token) (error) {
	str := token.Str
	if !p.has(str) {
		return fmt.Errorf("This Token(%s) is Invalid", str)
	}

	ctx := context.Background()
	_, err := p.client.Del(ctx, str).Result()
	if err != nil {
		return err
	}

	return nil
}
