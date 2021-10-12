package main

import (
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/initialize"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	global.VP = initialize.InitViper()
	if global.VP == nil {
		fmt.Println("配置文件初始化失败")
	}

	routers := initialize.Routers()

	s := &http.Server{
		Addr:    ":3337",
		Handler: routers,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Println(`System Serve Start Error`)
	}
}
