package server

import (
	"os"
	"io"
	"fmt"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogFormatterParams struct {
	Request *http.Request
	TimeStamp time.Time
	StatusCode int
	Latency time.Duration
	ClientIP string
	Method string
	Path string
	ErrorMessage string
	isTerm bool
	BodySize int
	Keys map[string]interface{}
	Message string
}

type LogFormatter func(params LogFormatterParams) string

type LoggerConfig struct {
	Formatter LogFormatter
	Output    io.Writer
	SkipPaths []string
}

func DefaultLogFormatter() LogFormatter {
	return func(params LogFormatterParams) string {
		if params.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			params.Latency = params.Latency - params.Latency%time.Second
		}
		return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %#v %s\n%s",
			params.TimeStamp.Format("2006/01/02 - 15:04:05"),
			params.StatusCode,
			params.Latency,
			params.ClientIP,
			params.Method,
			params.Path,
			params.Message,
			params.ErrorMessage,
		)
	}
}

func DefaultWriter() io.Writer {
	return os.Stdout
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Formatter: DefaultLogFormatter(),
		Output:    DefaultWriter(),
	})
}

func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = DefaultLogFormatter()
	}

	out := conf.Output

	notlogged := conf.SkipPaths

	isTerm := true

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		v, _ := c.Get("message")
		message, ok := v.(string)
		if !ok {
			message = ""
		}

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			params := LogFormatterParams{
				Request: c.Request,
				isTerm:  isTerm,
				Keys:    c.Keys,
			}

			// Stop timer
			params.TimeStamp = time.Now()
			params.Latency = params.TimeStamp.Sub(start)

			params.ClientIP = c.ClientIP()
			params.Method = c.Request.Method
			params.StatusCode = c.Writer.Status()
			params.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			params.BodySize = c.Writer.Size()
			params.Message = (string)(message)

			if raw != "" {
				path = path + "?" + raw
			}

			params.Path = path

			fmt.Fprint(out, formatter(params))
		}
	}
}
