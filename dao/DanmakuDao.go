package dao

import (
	"danmaku/model"
	"log"
)

type DanmakuDao struct{}

func (d DanmakuDao) CreateDanmaku(danmaku *model.Danmaku) error {
	result := Db.Create(danmaku).Error
	if result != nil {
		log.Println(result.Error())
		return result
	}
	return nil
}
