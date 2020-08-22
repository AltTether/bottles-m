package server

import (
	"fmt"
	"bytes"
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/message/pool"
)


func Init() {
	gin.SetMode(gin.TestMode)
}

func TestLogger(t *testing.T) {
	Init()

	formatter := DefaultLogFormatter()
	buffer := new(bytes.Buffer)

	conf := LoggerConfig{
		Formatter: formatter,
		Output:    buffer,
	}
	logger := LoggerWithConfig(conf)

	router := gin.New()
	router.Use(logger)
	router.GET("/example", func(c *gin.Context) {
		var message pool.Message
		c.BindJSON(&message)

		logMessage := fmt.Sprintf("INPUT message=%s", *message.Text)
		c.Set("message", logMessage)
	})

	text := "hoge"
	message := &pool.Message{
		Text: &text,
	}
	body, _ := json.Marshal(message)
	header := Header{
		Key:   "Content-Type",
		Value: "application/json",
	}
	performRequest(router, "GET", "/example?a=100", body, header)

	assert.Contains(t, buffer.String(),
		fmt.Sprintf("\"/example?a=100\" INPUT message=hoge"))
}

type Header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, body []byte, headers ...Header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
