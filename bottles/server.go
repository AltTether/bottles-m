package bottles

import(
	"time"
	
	"github.com/gin-gonic/gin"
)


type RequestBody struct {
	Message *string `json:"message" binding:"required"`
	Token   *string `json:"token" binding:"required"`
}

type Config struct {
	GetPipeline  *Pipeline
	PostPipeline *Pipeline
}

func New(conf Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	registerRoute(r, conf.GetPipeline, conf.PostPipeline)

	return r
}

func Default() *gin.Engine {
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

func DefaultWithPools(messagePool *MessagePool, tokenPool *TokenPool) *gin.Engine {
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

func registerRoute(r *gin.Engine, getPipeline, postPipeline *Pipeline) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", GetBottleHandlerFunc(getPipeline))
		v1.POST("/bottle", PostBottleHandlerFunc(postPipeline))
		v1.GET("/bottle/stream", GetBottleStreamHandlerFunc(getPipeline))
	}
}
