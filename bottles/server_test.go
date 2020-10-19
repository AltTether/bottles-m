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
	cfg := NewTestConfig()
	engine := New(cfg)
	r := NewServer(engine.Gateway, cfg)

	engine.Run()
	defer engine.Stop()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleStreamRoute(t *testing.T) {
	engine := CreateTestEngine()
	cfg := engine.Config
	r := NewServer(engine.Gateway, cfg)

	engine.Run()
	defer engine.Stop()
	
	w := CreateTestResponseRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle/stream", nil)
	go func () {
		time.Sleep(1000 * time.Millisecond)
		w.closeClient()
	}()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Greater(t,
		strings.Count(w.Body.String(), "{\"text\":\"test\"}"), 1)
}

func TestPostBottleRoute(t *testing.T) {
	engine := CreateTestEngine()
	cfg := engine.Config
	r := NewServer(engine.Gateway, cfg)

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

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestGetBottleRoute(t *testing.T) {
	engine := CreateTestEngine()
	cfg := engine.Config
	r := NewServer(engine.Gateway, cfg)

	engine.Run()
	defer engine.Stop()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
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
