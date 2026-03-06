package rabbitmq

import "github.com/streadway/amqp"

type RabbitMQ struct {
	conn  *amqp.Connection
	mqUrl string
}

const MqURL = "amqp://guest:guest@127.0.0.1:5672/"

var RMQ *RabbitMQ

func InitRabbitMQ() {
	RMQ = &RabbitMQ{
		mqUrl: MqURL,
	}
	dial, err := amqp.Dial(RMQ.mqUrl)
	if err != nil {
		panic(err)
	}
	RMQ.conn = dial
}
func (r *RabbitMQ) destroy() {
	r.conn.Close()
}
