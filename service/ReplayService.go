package service

import "danmaku/model"

type ReplayReq struct {
	VideoId int `form:"videoId"`
	Start   int `form:"start"`
	End     int `form:"end"`
	Limit   int `form:"limit"`
}

type ReplayService interface {
	Replay(req ReplayReq) ([]model.Danmaku, error)
}
