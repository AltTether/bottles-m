package bottles

import(
	"os"
	"fmt"
	"time"
	"strings"
	"net/http"
	
	"github.com/gin-gonic/gin"
)


type Engine struct {
	*gin.Engine
}

type Config struct {
	GetPipeline *Pipeline
	PostPipeline *Pipeline
}

func New(conf Config) *Engine {
	r := &Engine{
		gin.New(),
	}
	r.Use(gin.Logger(), gin.Recovery())

	registerRoute(r, conf.GetPipeline, conf.PostPipeline)

	return r
}

func Default() *Engine {
	getPipeline := NewPipeline()
	postPipeline := NewPipeline()

	messagePool := NewMessagePool()
	tokenPool := NewTokenPool(2 * time.Minute)

	postPipeline.AddStage(ValidateTokenStage(tokenPool))
	postPipeline.AddStage(StoreMessageStage(messagePool))

	getPipeline.AddStage(AddTokenStage(tokenPool))
	getPipeline.AddStage(AddMessageStage(messagePool))

	conf := Config{
		GetPipeline:  getPipeline,
		PostPipeline: postPipeline,
	}

	return New(conf)
}

func DefaultWithPools(messagePool *MessagePool, tokenPool *TokenPool) *Engine {
	getPipeline := NewPipeline()
	postPipeline := NewPipeline()

	postPipeline.AddStage(ValidateTokenStage(tokenPool))
	postPipeline.AddStage(StoreMessageStage(messagePool))

	getPipeline.AddStage(AddTokenStage(tokenPool))
	getPipeline.AddStage(AddMessageStage(messagePool))

	conf := Config{
		GetPipeline:  getPipeline,
		PostPipeline: postPipeline,
	}

	return New(conf)
}

func registerRoute(r *Engine, getPipeline, postPipeline *Pipeline) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", GetBottleHandlerFunc(getPipeline))
		v1.POST("/bottle", PostBottleHandlerFunc(postPipeline))
		v1.GET("/bottle/stream", GetBottleStreamHandlerFunc(getPipeline))
	}
}

func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine)
	return
}

func debugPrintError(err error) {
	if err != nil {
		if gin.IsDebugging() {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		}
	}
}

func debugPrint(format string, values ...interface{}) {
	if gin.IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] "+format, values...)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
