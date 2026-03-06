package rabbitmq

import (
	"context"
	"danmaku/dao"
	"danmaku/middle/redis"
	"danmaku/model"
	"encoding/json"
	"fmt"
	"time"

	goRedis "github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

type SendMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

func NewSendMQ(queueName string) *SendMQ {
	sendMQ := &SendMQ{
		RabbitMQ:  *RMQ,
		queueName: queueName,
	}
	channel, err := sendMQ.conn.Channel()
	sendMQ.channel = channel
	if err != nil {
		panic(err)
	}
	_, err = sendMQ.channel.QueueDeclare(
		sendMQ.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return sendMQ
}

// 使用工作队列模式
func (sendMQ *SendMQ) Producer(msg *model.SendMessageMQ) {
	body, err := json.Marshal(msg)
	err = sendMQ.channel.Publish(
		"",
		sendMQ.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		panic(err)
	}
}
func (sendMQ *SendMQ) Consumer() {
	_ = sendMQ.channel.Qos(200, 0, false)
	msg, err := sendMQ.channel.Consume(
		sendMQ.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	forever := make(chan bool)
	go sendMQ.consumerSendDanmaku(msg)
	<-forever
}

func (sendMQ *SendMQ) consumerSendDanmaku(msg <-chan amqp.Delivery) {
	const (
		batchSize     = 200
		flushInterval = time.Second
	)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()
	batch := make([]model.Danmaku, 0, batchSize)
	deliveries := make([]amqp.Delivery, 0, batchSize)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := dao.DanmakuDao{}.SaveBatch(batch)
		if err != nil {
			for _, d := range deliveries {
				_ = d.Nack(false, false)
			}
		} else {
			go sendMQ.pushBatchToRedis(batch)
			for _, d := range deliveries {
				_ = d.Ack(false)
			}
		}
		batch = batch[:0]
		deliveries = deliveries[:0]
	}
	for {
		select {
		case d, ok := <-msg:
			if !ok {
				flush()
				return
			}
			var m model.SendMessageMQ
			if err := json.Unmarshal(d.Body, &m); err != nil {
				_ = d.Ack(false)
				continue
			}
			batch = append(batch, model.Danmaku{
				MsgId:      m.MsgId,
				VideoId:    m.VideoId,
				UserId:     m.UserId,
				Content:    m.Content,
				VideoTime:  m.VideoTime,
				CreateTime: m.CreateTime,
			})
			deliveries = append(deliveries, d)

			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (sendMq *SendMQ) pushBatchToRedis(batch []model.Danmaku) {
	if len(batch) == 0 {
		return
	}
	ctx := context.Background()
	videoDanmaku := make(map[int][]model.Danmaku)
	for _, v := range batch {
		videoDanmaku[v.VideoId] = append(videoDanmaku[v.VideoId], v)
	}
	pipe := redis.RdbReplay.Pipeline()
	for videoId, list := range videoDanmaku {
		cacheKey := fmt.Sprintf("danmaku:video:%d", videoId)
		exist, _ := redis.RdbReplay.Exists(ctx, cacheKey).Result()
		if exist > 0 {
			var ZSetMembers []goRedis.Z
			for _, member := range list {
				memberJson, _ := json.Marshal(member)
				ZSetMembers = append(ZSetMembers, goRedis.Z{
					Score:  float64(member.VideoTime),
					Member: memberJson,
				})
			}
			pipe.ZAdd(ctx, cacheKey, ZSetMembers...)
		}
	}
	_, _ = pipe.Exec(ctx)
}

var SendMQDel *SendMQ

func InitSendMQ() {
	SendMQDel = NewSendMQ("send_del")
	go SendMQDel.Consumer()
}
