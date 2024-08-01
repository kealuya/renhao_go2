package rabbit_mq

import (
	"fmt"
	"github.com/gohouse/t"
	"github.com/streadway/amqp"
	"log"
)

func producer() {
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
	//	"test_queue", //队列名
	//	false,        //是否持久化
	//	false,        //是否自动删除
	//	false,        //是否排他
	//	false,        //是否阻塞
	//	nil,          //额外属性
	//)
	//if err != nil {
	//	log.Fatalf("Failed to declare a queue: %v", err)
	//}
	//fmt.Println(q)
	for i := 0; i < 100; i++ {

		body := fmt.Sprintf("my sent message %s", t.New(i).String())
		err = ch.Publish(
			"myexchange", // 交换器
			"a.a",
			//q.Name, // 路由键（队列名）
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Fatalf("Failed to publish a message: %v", err)
		}
		log.Printf(" [x]  Sent %s", body)
	}

}
