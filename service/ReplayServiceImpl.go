package service

import (
	"context"
	"danmaku/dao"
	"danmaku/middle/redis"
	"danmaku/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	goRedis "github.com/redis/go-redis/v9"
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
	ctx := context.Background()
	cacheKey := fmt.Sprintf("danmaku:video:%d", req.VideoId)
	redisOpt := &goRedis.ZRangeBy{
		Min:    strconv.Itoa(req.Start),
		Max:    strconv.Itoa(req.End),
		Offset: 0,
		Count:  int64(limit),
	}
	danmakuCache, err := redis.RdbReplay.ZRangeByScore(ctx, cacheKey, redisOpt).Result()
	if err == nil && len(danmakuCache) > 0 {
		var danmakuList []model.Danmaku
		for _, v := range danmakuCache {
			var d model.Danmaku
			_ = json.Unmarshal([]byte(v), &d)
			danmakuList = append(danmakuList, d)
		}
		return danmakuList, nil
	}
	danmakuList, err := r.danmakuDao.DanmakuList(req.VideoId, req.Start, req.End, limit)
	if err != nil {
		return nil, err
	}
	if len(danmakuList) != 0 {
		go r.writeToRedis(req.VideoId, danmakuList)
	}
	return danmakuList, nil
}
func (r *ReplayServiceImpl) writeToRedis(videoId int, danmakuList []model.Danmaku) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("danmaku:video:%d", videoId)
	var ZSetMembers []goRedis.Z
	for _, d := range danmakuList {
		dJson, _ := json.Marshal(d)
		ZSetMembers = append(ZSetMembers, goRedis.Z{
			Score:  float64(d.VideoTime),
			Member: string(dJson),
		})
	}
	pipe := redis.RdbReplay.Pipeline()
	pipe.ZAdd(ctx, cacheKey, ZSetMembers...)
	pipe.Expire(ctx, cacheKey, 1*time.Hour)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println(err)
	}
}
