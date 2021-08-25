package api

import (
	"github.com/akazwz/weibo-hot-search/model/response"
	"github.com/akazwz/weibo-hot-search/utils/influx"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func GetCurrentHotSearchApi(c *gin.Context) {
	hotSearch, err := influx.GetCurrentHotSearch()
	if err != nil {
		log.Println("get current hot search error")
		response.CommonFailed(4000, "get current hot search error", c)
		return
	}
	response.CommonSuccess(2000, "success", hotSearch, c)
}

func GetDurationHotSearchApi(c *gin.Context) {
	start := c.Query("start")
	stop := c.Query("stop")
	if start == "" || stop == "" {
		GetCurrentHotSearchApi(c)
		return
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println("time load location error")
		response.CommonFailed(4000, "time load location error", c)
		return
	}
	startTime, err := time.ParseInLocation("2006-01-02-15-04", start, location)
	if err != nil {
		log.Println("time parse error")
		response.CommonFailed(4000, "time parse error", c)
		return
	}
	endTime, err := time.ParseInLocation("2006-01-02-15-04", stop, location)
	if err != nil {
		log.Println("time parse error")
		response.CommonFailed(4000, "time parse error", c)
		return
	}
	start = startTime.Format(time.RFC3339)
	stop = endTime.Format(time.RFC3339)
	hotSearches, err := influx.GetDurationHotSearch(start, stop)
	if err != nil {
		log.Println("get current hot search error")
		response.CommonFailed(4000, "get current hot search error", c)
	}
	response.CommonSuccess(2000, "success", hotSearches, c)
}
