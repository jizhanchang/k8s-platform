package routes

import (
	"k8s-platfrom/controllers"
	"k8s-platfrom/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //设置为上线模式
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
		return
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
		return
	})
	v1 := r.Group("/api/v1")

	//v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/pods", controllers.GetPodsHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})
	return r
}
