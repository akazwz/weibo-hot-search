package initialize

import (
	"github.com/akazwz/weibo-hot-search/model/response"
	"github.com/akazwz/weibo-hot-search/routers"
	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	var router = gin.Default()
	router.Static("/public", "./public")

	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))

	router.GET("teapot", response.Teapot)

	routerGroup := router.Group("")
	routers.InitHotSearch(routerGroup)
	return router
}
