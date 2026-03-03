package controller

import (
	"danmaku/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SendController struct {
	sendService service.SendService
}

func NewSendController(sendService service.SendService) *SendController {
	return &SendController{sendService: sendService}
}

func (s *SendController) Send(c *gin.Context) {
	var req service.SendDanmakuReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "invalid json",
			"error": err.Error(),
		})
		return
	}
	if err := s.sendService.Send(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "send failed",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}
