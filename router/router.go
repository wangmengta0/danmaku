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
	ctl := controller.NewSendController(sendSvc)

	api := router.Group("/api/v1")
	{
		api.POST("/danmaku", ctl.Send)
	}

	return router

}
