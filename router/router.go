package router

import (
	"danmaku/controller"
	"danmaku/dao"
	"danmaku/realtime"
	"danmaku/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	sendSvc := service.NewSendServiceImpl(dao.DanmakuDao{})
	sendCtl := controller.NewSendController(sendSvc)

	replaySvc := service.NewReplayServiceImpl(dao.DanmakuDao{})
	replayCtl := controller.NewReplayController(replaySvc)

	hub := realtime.NewHub()
	go hub.Run()
	wsCtl := controller.NewWSController(hub)
	api := router.Group("/api/v1")
	{
		api.POST("/danmaku", sendCtl.Send)
		api.GET("/danmaku/relay", replayCtl.ReplayDanmaku)
		api.GET("/ws", wsCtl.Subscribe)
	}

	return router

}
