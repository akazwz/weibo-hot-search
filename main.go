package main

import (
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/initialize"
	"github.com/akazwz/weibo-hot-search/utils/influx"
)

func main() {
	global.VP = initialize.InitViper()
	if global.VP == nil {
		fmt.Println("配置文件初始化失败")
	}
	influx.Query()
}
