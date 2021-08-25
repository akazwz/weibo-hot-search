package routers

import (
	"github.com/akazwz/weibo-hot-search/api"
	"github.com/gin-gonic/gin"
)

func InitHotSearch(r *gin.RouterGroup) {
	hotSearchRouter := r.Group("hot-search")
	{
		hotSearchRouter.GET("", api.GetCurrentHotSearchApi)
	}
}
