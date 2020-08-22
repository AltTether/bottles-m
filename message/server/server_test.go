package server

import (
	"bytes"
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	
	"github.com/stretchr/testify/assert"

	"github.com/clients"
)


func TestPostAndGetRoute(t *testing.T) {
	r := gin.Default()
	registerRouters(r)

	w := httptest.NewRecorder()

	message := &clients.Message{
		Text: "",
	}
	messageStr, _ := json.Marshal(message)
	postReq, _ := http.NewRequest("POST", "/",
		bytes.NewBuffer(messageStr))

	getReq, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, postReq)
	assert.Equal(t, http.StatusOK, w.Code)

	r.ServeHTTP(w, getReq)
	assert.Equal(t, http.StatusOK, w.Code)
}
