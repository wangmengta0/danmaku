package rabbitmq

import (
	"danmaku/dao"
	"danmaku/model"
	"encoding/json"
	"time"

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
	return sendMQ
}

// 使用工作队列模式
func (sendMQ *SendMQ) Producer(msg *model.SendMessageMQ) {
	_, err := sendMQ.channel.QueueDeclare(
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
	_, err := sendMQ.channel.QueueDeclare(
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
				_ = d.Nack(false, true)
			}
		} else {
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

var SendMQDel *SendMQ

func InitSendMQ() {
	SendMQDel = NewSendMQ("send_del")
	go SendMQDel.Consumer()
}
