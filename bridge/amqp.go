package bridge

import (
	"encoding/json"
	"fmt"

	"github.com/haobird/gormq"
)

//Bridge 桥接
type BridgeMQ struct {
	RabbitMQ  *gormq.RabbitMQ
	Publisher *gormq.Publisher
}

//Publish 发布
func (b *BridgeMQ) Publish(e *Element) error {
	fmt.Println("[bridge] ", e)
	// 转换为 二进制
	buf, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return b.Publisher.Pub(buf)
}

//InitBridgeMQ rabbitmq
func InitBridgeMQ() *BridgeMQ {
	addr := "amqp://admin:admin@127.0.0.1:5672"
	// addr = config.AMQP.Addr
	fmt.Println(addr)
	rabbitmq := gormq.New(addr)
	// 队列和交换机 创建
	queueExchange := gormq.QueueExchange{
		QueueName:    "receiver.uniview.v1",
		RoutingKey:   "iot.report.opendoor.insert",
		ExchangeName: "iot.opendoor",
		ExchangeType: "topic",
	}
	// queueExchange = gormq.QueueExchange{
	// 	QueueName:    config.AMQP.PublishChannel.QueueName,
	// 	RoutingKey:   config.AMQP.PublishChannel.RoutingKey,
	// 	ExchangeName: config.AMQP.PublishChannel.ExchangeName,
	// 	ExchangeType: config.AMQP.PublishChannel.ExchangeType,
	// }
	publisher := rabbitmq.NewPublisher(queueExchange)

	return &BridgeMQ{
		RabbitMQ:  rabbitmq,
		Publisher: publisher,
	}

}
