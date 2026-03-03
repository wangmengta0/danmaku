package main

import (
	"danmaku/dao"
	"danmaku/router"
)

func main() {
	dao.Init()
	r := router.SetupRouter()
	_ = r.Run(":8080")
}
