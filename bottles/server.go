package bottles

import(
	"github.com/gin-gonic/gin"
)


func NewServer(gateway *Gateway) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	registerRoute(r, gateway)

	return r
}

func registerRoute(r *gin.Engine, gateway *Gateway) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/bottle", GetBottleHandlerFunc(gateway))
		v1.POST("/bottle", PostBottleHandlerFunc(gateway))
		v1.GET("/bottle/stream", GetBottleStreamHandlerFunc(gateway))
	}
}

func Run() {
	e := DefaultEngine()
	g := e.Gateway
	e.Run()
	defer e.Stop()

	s := NewServer(g)
	s.Run()
}
