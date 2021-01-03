package bottles

import(
	"github.com/gin-gonic/gin"
)


func NewServer(gateway *Gateway, cfg *Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	registerRoute(r, gateway, cfg)

	return r
}

func registerRoute(r *gin.Engine, gateway *Gateway, cfg *Config) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", GetBottleHandlerFunc(gateway, cfg))
		v1.POST("/bottle", PostBottleHandlerFunc(gateway))
		v1.GET("/bottle/stream", GetBottleStreamHandlerFunc(gateway, cfg))
	}
}

func Run() {
	e := DefaultEngine()
	g := e.Gateway
	c := e.Config
	
	e.Run()
	defer e.Stop()

	s := NewServer(g, c)
	s.Run()
}
