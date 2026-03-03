package service

import (
	"danmaku/dao"
	"danmaku/model"
	"errors"
	"time"
)

type SendServiceImpl struct {
	danmakuDao dao.DanmakuDao
}

func NewSendServiceImpl(danmakuDao dao.DanmakuDao) SendService {
	return &SendServiceImpl{danmakuDao: danmakuDao}
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
	return nil
}
