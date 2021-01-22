package server

import(
	"github.com/gin-gonic/gin"
	
	"github.com/bottles/bottles"
)


func NewServer(gateway *bottles.Gateway, cfg *bottles.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	registerRoute(r, gateway, cfg)

	return r
}

func registerRoute(r *gin.Engine, gateway *bottles.Gateway, cfg *bottles.Config) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", GetBottleHandlerFunc(gateway, cfg))
		v1.POST("/bottle", PostBottleHandlerFunc(gateway))
		v1.GET("/bottle/stream", GetBottleStreamHandlerFunc(gateway, cfg))
	}
}

func Run() {
	e := bottles.DefaultEngine()
	g := e.Gateway
	c := e.Config
	
	e.Run()
	defer e.Stop()

	s := NewServer(g, c)
	s.Run()
}
