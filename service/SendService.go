package service

type SendDanmakuReq struct {
	VideoId   int    `json:"videoId"`
	UserId    int    `json:"userId"`
	Content   string `json:"content"`
	VideoTime int    `json:"videoTime"`
}

type SendService interface {
	Send(danmaku *SendDanmakuReq) error
}
