package main

import (
	"danmaku/dao"
	"danmaku/middle/rabbitmq"
	"danmaku/middle/redis"
	"danmaku/router"
)

func main() {
	initDeps()
	r := router.SetupRouter()
	_ = r.Run(":8080")
}
func initDeps() {
	dao.Init()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitSendMQ()

	redis.InitRedis()
}
