package service

import (
	"danmaku/dao"
	"danmaku/middle/rabbitmq"
	"danmaku/model"
	"danmaku/realtime"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SendServiceImpl struct {
	danmakuDao dao.DanmakuDao
	hub        *realtime.Hub
	sendMQ     *rabbitmq.SendMQ
}

func NewSendServiceImpl(danmakuDao dao.DanmakuDao, hub *realtime.Hub, sendMQ *rabbitmq.SendMQ) SendService {
	return &SendServiceImpl{danmakuDao: danmakuDao, hub: hub, sendMQ: sendMQ}
}

func (s *SendServiceImpl) Send(danmaku *SendDanmakuReq) error {
	if danmaku.UserId <= 0 || danmaku.VideoId <= 0 || danmaku.VideoTime < 0 {
		return errors.New("param error")
	}
	if danmaku.Content == "" {
		return errors.New("content is empty")
	}
	if len([]rune(danmaku.Content)) > 200 {
		return errors.New("content is too large")
	}
	msgId := uuid.NewString()
	mqMsg := &model.SendMessageMQ{
		MsgId:      msgId,
		VideoId:    danmaku.VideoId,
		UserId:     danmaku.UserId,
		Content:    danmaku.Content,
		VideoTime:  danmaku.VideoTime,
		CreateTime: time.Now(),
	}
	//err := s.danmakuDao.CreateDanmaku(danmakuSend)
	s.sendMQ.Producer(mqMsg)
	wsMsg := map[string]any{
		"type":       "danmaku",
		"id":         mqMsg.MsgId,
		"videoId":    danmaku.VideoId,
		"userId":     danmaku.UserId,
		"content":    danmaku.Content,
		"videoTime":  danmaku.VideoTime,
		"createTime": time.Now(),
	}
	b, _ := json.Marshal(wsMsg)
	s.hub.Broadcast <- realtime.BroadcastMsg{
		RoomId: strconv.Itoa(danmaku.VideoId),
		Data:   b,
	}
	return nil
}
