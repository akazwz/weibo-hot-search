package routers

import (
	"github.com/akazwz/weibo-hot-search/api"
	"github.com/gin-gonic/gin"
)

func InitHotSearch(r *gin.RouterGroup) {
	hotSearchRouter := r.Group("hot-searches")
	{
		hotSearchRouter.GET("/current", api.GetCurrentHotSearchApi)
		hotSearchRouter.GET("", api.GetDurationHotSearchApi)
		hotSearchRouter.GET("/content/:content", api.GetHotSearchesByContentApi)
		hotSearchRouter.GET("/keyword/:keyword", api.GetHotSearchesByKeyWordApi)
	}
}
