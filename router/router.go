package router

import (
	"danmaku/controller"
	"danmaku/dao"
	"danmaku/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	sendSvc := service.NewSendServiceImpl(dao.DanmakuDao{})
	sendCtl := controller.NewSendController(sendSvc)

	replaySvc := service.NewReplayServiceImpl(dao.DanmakuDao{})
	replayCtl := controller.NewReplayController(replaySvc)
	api := router.Group("/api/v1")
	{
		api.POST("/danmaku", sendCtl.Send)
		api.GET("/danmaku/relay", replayCtl.ReplayDanmaku)
	}

	return router

}
