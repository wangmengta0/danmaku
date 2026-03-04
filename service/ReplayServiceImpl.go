package service

import (
	"danmaku/dao"
	"danmaku/model"
	"errors"
)

type ReplayServiceImpl struct {
	danmakuDao dao.DanmakuDao
}

func NewReplayServiceImpl(dao dao.DanmakuDao) ReplayService {
	return &ReplayServiceImpl{danmakuDao: dao}
}
func (r *ReplayServiceImpl) Replay(req ReplayReq) ([]model.Danmaku, error) {
	if req.VideoId < 0 || req.Start < 0 || req.End < 0 || req.Start > req.End {
		return nil, errors.New("param error")
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 500
	}
	if limit > 2000 {
		limit = 2000
	}
	return r.danmakuDao.DanmakuList(req.VideoId, req.Start, req.End, limit)
}
