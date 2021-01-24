package server

import (
	"fmt"
	"time"
	"bytes"
	"testing"
	"strings"
	"net/http"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/bottles/bottles"
)

var (
	testMessageText = "test"
	testTokenStrFormatter = func(i int) string {
		return fmt.Sprintf("test%d", i)
	}
)


func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetBottleRouteEmpty(t *testing.T) {
	cfg := NewTestConfig()
	engine := bottles.New(cfg)
	r := NewServer(engine.Gateway)

	engine.Run()
	defer engine.Stop()

	w := CreateTestResponseRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	go func () {
		time.Sleep(100 * time.Millisecond)
		w.closeClient()
	}()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleStreamRoute(t *testing.T) {
	engine := CreateTestEngine()
	r := NewServer(engine.Gateway)

	engine.Run()
	defer engine.Stop()
	
	w := CreateTestResponseRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle/stream", nil)
	go func () {
		time.Sleep(100 * time.Millisecond)
		w.closeClient()
	}()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Greater(t,
		strings.Count(w.Body.String(), "{\"text\":\"test\"}"), 1)
}

func TestPostBottleRoute(t *testing.T) {
	engine := CreateTestEngine()
	r := NewServer(engine.Gateway)

	engine.Run()
	defer engine.Stop()

	w := httptest.NewRecorder()

	message := testMessageText
	token := testTokenStrFormatter(0)
	testRequestBody := &RequestBody{
		Message: &message,
		Token:   &token,
	}
	body, _ := json.Marshal(testRequestBody)
	
	req, _ := http.NewRequest("POST",
		"/api/v1/bottle",
		bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleRoute(t *testing.T) {
	engine := CreateTestEngine()
	r := NewServer(engine.Gateway)

	engine.Run()
	defer engine.Stop()

	w := CreateTestResponseRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	go func () {
		time.Sleep(100 * time.Millisecond)
		w.closeClient()
	}()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"message\":{\"text\":\"test\"}")
}

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func NewTestConfig() *bottles.Config {
	return &bottles.Config{
		Seed:                   42,
		TokenSize:              10,
		TokenExpiration:        time.Duration(10 * time.Millisecond),
		SendBottleDelay:        time.Duration(15 * time.Millisecond),
		GenerateBottleCoolTime: time.Duration(20 * time.Millisecond),
	}
}

func CreateTestEngine() *bottles.Engine {
	n := 10

	return CreateTestEngineWithData(createTestMessages(n), createTestTokens(n))
}

func CreateTestEngineWithData(messages []*bottles.Message, tokens []*bottles.Token) *bottles.Engine {
	cfg := NewTestConfig()
	engine := bottles.New(cfg)

	messageStorage := createTestMessageStorageWithMessages(messages)
	tokenStorage := createTestTokenStorageWithTokens(tokens, cfg.TokenExpiration)

	engine.AddHandler(bottles.REQUEST_BOTTLE_MODE, bottles.BottleGetHandler(tokenStorage, messageStorage))
	engine.AddHandler(bottles.ADD_BOTTLE_MODE, bottles.BottleAddHandler(tokenStorage, messageStorage))

	return engine
}

func createTestMessageStorageWithMessages(ms []*bottles.Message) *bottles.MessageStorage{
	s := bottles.NewMessageStorage()
	for _, m := range ms {
		_ = s.Add(m)
	}

	return s
}

func createTestTokenStorageWithTokens(ts []*bottles.Token, expiration time.Duration) *bottles.TokenStorage {
	s := bottles.NewTokenStorage(expiration)
	for _, t := range ts {
		_ = s.Add(t)
	}

	return s
}

func createTestMessages(n int) []*bottles.Message {
	ms := make([]*bottles.Message, n)
	for i := 0; i < n; i++ {
		text := "test"
		ms[i] = &bottles.Message{
			Text: text,
		}
	}

	return ms
}

func createTestTokens(n int) []*bottles.Token {
	ts := make([]*bottles.Token, n)
	for i := 0; i < n; i++ {
		str := fmt.Sprintf("test%d", i)
		ts[i] = &bottles.Token{
			Str: str,
		}
	}

	return ts
}
