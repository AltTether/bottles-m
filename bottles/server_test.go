package bottles

import (
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
	testMessage = "test"
	testToken = "test"
	testRequestBody = &RequestBody{
		Message: &testMessage,
		Token:   &testToken,
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
	messagePool := NewMessagePool()
	tokenPool := NewTokenPool(2 * time.Minute)
	for i := 0; i < 10; i++ {
		text := "test"
		_ = messagePool.Add(&Message{
			Text: &text,
		})
	}

	for i := 0; i < 10; i++ {
		str := GenerateRandomString(10)
		_ = tokenPool.Add(&Token{
			Str: &str,
		})
	}

	r := DefaultWithPools(messagePool, tokenPool)

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
	r := Default()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testRequestBody)
	req, _ := http.NewRequest("POST",
		"/api/v1/bottle",
		bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleRoute(t *testing.T) {
	r := Default()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testRequestBody)
	req, _ := http.NewRequest("POST",
		"/api/v1/bottle",
		bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	// for pipeline, this behavior is expected(?)
	time.Sleep(100 * time.Millisecond)
	req, _ = http.NewRequest("GET", "/api/v1/bottle", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "\"message\":{\"text\":\"test\"}")
}

func TestGenerateToken(t *testing.T) {
	size := 10
	tokenStr := GenerateRandomString(size)
	assert.Equal(t, size, len(tokenStr))
}

func CreateTestTokenPool() *TokenPool {
	tokenPool := NewTokenPool(2 * time.Minute)
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("test%d", i)
		_ = tokenPool.Add(&Token{
			Str: &str,
		})
	}
	return tokenPool
}

func CreateTestMessagePool() *MessagePool {
	messagePool := NewMessagePool()
	for i := 0; i < 10; i++ {
		text := "test"
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
