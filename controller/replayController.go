package controller

import (
	"danmaku/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReplayController struct {
	replayService service.ReplayService
}

func NewReplayController(replayService service.ReplayService) *ReplayController {
	return &ReplayController{replayService: replayService}
}

func (r *ReplayController) ReplayDanmaku(c *gin.Context) {
	var req service.ReplayReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "invalid query",
			"error": err.Error(),
		})
		return
	}
	fmt.Println("req:", req)
	list, err := r.replayService.Replay(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": list,
	})
}
