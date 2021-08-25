package api

import (
	"github.com/akazwz/weibo-hot-search/model/response"
	"github.com/akazwz/weibo-hot-search/utils/influx"
	"github.com/gin-gonic/gin"
	"log"
)

func GetCurrentHotSearchApi(c *gin.Context) {
	hotSearch, err := influx.GetCurrentHotSearch()
	if err != nil {
		log.Println("get current hot search error")
		response.CommonFailed(4000, "get current hot search error", c)
	}
	response.CommonSuccess(2000, "success", hotSearch, c)
}
