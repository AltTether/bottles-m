package bottles

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
	r := Default()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleStreamRoute(t *testing.T) {
	n := 10
	r := CreateTestEngine(n)
	
	w := CreateTestResponseRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle/stream", nil)
	go func () {
		time.Sleep(100 * time.Millisecond)
		w.closeClient()
	}()
	r.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	assert.Greater(t,
		strings.Count(w.Body.String(), "{\"text\":\"test\"}"), 1)
}

func TestPostBottleRoute(t *testing.T) {
	n := 10
	r := CreateTestEngine(n)

	w := httptest.NewRecorder()

	message := testMessageText
	token := testTokenStrFormatter(n - 1)
	testRequestBody := &RequestBody{
		Message: &message,
		Token:   &token,
	}
	body, _ := json.Marshal(testRequestBody)
	
	req, _ := http.NewRequest("POST",
		"/api/v1/bottle",
		bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleRoute(t *testing.T) {
	n := 10
	r := CreateTestEngine(n)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "\"message\":{\"text\":\"test\"}")
}

func TestGenerateToken(t *testing.T) {
	size := 10
	tokenStr := GenerateRandomString(size)
	assert.Equal(t, size, len(tokenStr))
}

func CreateTestEngine(n int) *Engine {
	messagePool := CreateTestMessagePool(n)
	tokenPool := CreateTestTokenPool(n)
	return DefaultWithPools(messagePool, tokenPool)
}

func CreateTestTokenPool(n int) *TokenPool {
	tokenPool := NewTokenPool(2 * time.Minute)
	for i := 0; i < n; i++ {
		str := testTokenStrFormatter(i)
		_ = tokenPool.Add(&Token{
			Str: &str,
		})
	}
	return tokenPool
}

func CreateTestMessagePool(n int) *MessagePool {
	messagePool := NewMessagePool()
	for i := 0; i < n; i++ {
		text := testMessageText
		_ = messagePool.Add(&Message{
			Text: &text,
		})
	}
	return messagePool
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
