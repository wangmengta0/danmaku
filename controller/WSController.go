package controller

import (
	"danmaku/realtime"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSController struct {
	hub *realtime.Hub
}

func NewWSController(hub *realtime.Hub) *WSController {
	return &WSController{hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (w *WSController) Subscribe(c *gin.Context) {
	videoId := c.Query("videoId")
	if videoId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "video id is empty",
		})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &realtime.Client{
		Conn:   conn,
		Send:   make(chan []byte, 256),
		RoomId: videoId,
		Hub:    w.hub,
	}
	w.hub.Register <- client
	client.Send <- []byte(`{"type":"subscribed","videoId":"` + videoId + `"}`)
	go client.WritePump()
	go client.ReadPump()
}
