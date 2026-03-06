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

func (d DanmakuDao) DanmakuList(videoId int, start int, end int, limit int) ([]model.Danmaku, error) {
	var list []model.Danmaku
	err := Db.Where("video_id=? and video_time>=? and video_time<=?", videoId, start, end).
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (d DanmakuDao) SaveBatch(batch []model.Danmaku) error {
	return Db.CreateInBatches(batch, 200).Error
}
