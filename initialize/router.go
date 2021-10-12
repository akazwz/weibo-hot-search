package initialize

import (
	"time"

	"github.com/akazwz/weibo-hot-search/middle"
	"github.com/akazwz/weibo-hot-search/model/response"
	"github.com/akazwz/weibo-hot-search/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	var router = gin.Default()
	router.Static("/public", "./public")

	// cors
	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))
	// rate limit
	router.Use(middle.RateLimitMiddleware(time.Millisecond*10, 100))
	// teapot
	router.GET("teapot", response.Teapot)
	routerGroup := router.Group("")
	routers.InitHotSearch(routerGroup)
	return router
}
