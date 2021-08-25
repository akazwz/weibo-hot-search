package main

import (
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/initialize"
	"github.com/akazwz/weibo-hot-search/utils/influx"
	"log"
)

func main() {
	global.VP = initialize.InitViper()
	if global.VP == nil {
		fmt.Println("配置文件初始化失败")
	}

	/*routers := initialize.Routers()

	s := &http.Server{
		Addr:    ":8000",
		Handler: routers,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Println(`System Serve Start Error`)
	}*/

	hotSearches, err := influx.GetDurationHotSearch()
	if err != nil {
		log.Println("get duration hot searches error")
	}
	fmt.Println(hotSearches[3])
}
