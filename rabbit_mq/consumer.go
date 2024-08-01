package rabbit_mq

import (
	"github.com/streadway/amqp"
	"log"
)

func consumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
	//
	//q, err := ch.QueueDeclare(
	//	"myqueue", //队列名
	//	true,      //是否持久化
	//	false,     //是否自动删除
	//	false,     //是否排他
	//	false,     //是否阻塞
	//	nil,       //额外属性
	//)
	//if err != nil {
	//	log.Fatalf("Failed to declare a queue: %v", err)
	//}

	msgs, err := ch.Consume(
		"myqueue2", // 队列名
		"",         // 消费者标签
		true,       // 自动确认
		false,      // 排他
		false,      // no-local
		false,      // no-wait
		nil,        // 参数
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
