package endpoints

import (
	"ggateway/pkg/services"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
)

// SetupRouter registers the endpoint handlers and returns a pointer to the
// server instance.
func SetupRouter() *gin.Engine {
	router := gin.Default()
	/*router.Use(gin_comm.CorsFunc())
	router.Use(gin_comm.AuthMiddleWare())*/
	// api/v1 router group
	v1 := router.Group("/v1", func(c *gin.Context) {
		c.Next()
		c.Header("X-TIME", "time")
	})
	{
		v1.Any("/log", nil)
	}
	ginpprof.Wrap(router)
	return router
}
