package rabbit_mq

import (
	"github.com/streadway/amqp"
	"log"
)

// 初始化RabbitMQ连接
func initConnection() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return conn, ch
}
func producer() {
	conn, ch := initConnection()
	defer conn.Close()
	defer ch.Close()
	// 声明交换机
	err := ch.ExchangeDeclare(
		"hotel_exchange",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}
	// 声明重试次数
	headers := amqp.Table{
		"retry_count": 3,
	}

	// 发送消息
	err = ch.Publish(
		"hotel_exchange", // 交换器
		"hotel_key",
		//q.Name, // 路由键（队列名）
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("msg error"),
			Headers:     headers,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

}
