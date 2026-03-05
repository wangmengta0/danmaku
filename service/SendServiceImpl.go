package service

import (
	"danmaku/dao"
	"danmaku/model"
	"danmaku/realtime"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type SendServiceImpl struct {
	danmakuDao dao.DanmakuDao
	hub        *realtime.Hub
}

func NewSendServiceImpl(danmakuDao dao.DanmakuDao, hub *realtime.Hub) SendService {
	return &SendServiceImpl{danmakuDao: danmakuDao, hub: hub}
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
	danmakuSend := &model.Danmaku{
		VideoId:    danmaku.VideoId,
		UserId:     danmaku.UserId,
		Content:    danmaku.Content,
		VideoTime:  danmaku.VideoTime,
		CreateTime: time.Now(),
	}
	err := s.danmakuDao.CreateDanmaku(danmakuSend)
	if err != nil {
		return err
	}
	msg := map[string]any{
		"type":       "danmaku",
		"id":         danmakuSend.Id,
		"videoId":    danmakuSend.VideoId,
		"userId":     danmaku.UserId,
		"content":    danmaku.Content,
		"videoTime":  danmaku.VideoTime,
		"createTime": danmakuSend.CreateTime,
	}
	b, _ := json.Marshal(msg)
	s.hub.Broadcast <- realtime.BroadcastMsg{
		RoomId: strconv.Itoa(danmakuSend.VideoId),
		Data:   b,
	}
	return nil
}
