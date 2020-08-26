package bottles

import (
	"bytes"
	"testing"
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
	r := New()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/bottle", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestPostBottleRoute(t *testing.T) {
	r := New()

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
	r := New()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(testRequestBody)
	req, _ := http.NewRequest("POST",
		"/api/v1/bottle",
		bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

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
