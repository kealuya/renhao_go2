package rabbit_mq

import (
	"fmt"
	"github.com/gohouse/t"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

func handlerConsumer() {
	conn, ch := initConnection()
	defer conn.Close()
	defer ch.Close()
	// 0. 声明队列 （带死信配置）
	// 一旦reject，会自动发送死信队列
	args := amqp.Table{
		"x-dead-letter-exchange":    "hotel_dlx_exchange",
		"x-dead-letter-routing-key": "dlx_key",
	}
	q, err := ch.QueueDeclare(
		"hotel_queue_1", //队列名
		false,           //是否持久化
		false,           //是否自动删除
		false,           //是否排他
		false,           //是否阻塞
		args,            //额外属性
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	// 1. 设置QoS
	err = ch.Qos(
		1,     // prefetch count = 1，每次只处理一条消息、预取数量，为了可以让消费者平分
		0,     // prefetch size = 0，不限制消息大小
		false, // global = false，设置应用于每个消费者
	)
	// 2. 绑定队列到交换机
	err = ch.QueueBind(
		q.Name,           // 队列名
		"hotel_key",      // 路由键
		"hotel_exchange", // 交换机名
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}
	// 3. 开始消费
	msgs, err := ch.Consume(
		q.Name,            // 队列名
		"hotel_consume_2", // 消费者标签
		false,             // 自动确认
		false,             // 排他
		false,             // no-local
		false,             // no-wait
		nil,               // 参数
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			// 获取重试次数
			retryCount := 0
			if d.Headers["retry_count"] != nil {
				retryCount = t.New(d.Headers["retry_count"]).Int()
			}
			// 处理消息（这里模拟处理过程）
			log.Printf("Received a message: %s", d.Body)
			err := processMessage(d.Body)
			if err != nil {
				if retryCount >= 3 {
					// 超过最大重试次数，发送到死信队列
					log.Printf("Message failed after %d retries, sending to DLQ", retryCount)
					d.Reject(false) // 拒绝消息，不重新入队
				} else {
					// 重试逻辑:重试记录+1，重新发回队列等待消费
					headers := amqp.Table{
						"retry_count": retryCount + 1,
					}
					err = ch.Publish(
						q.Name,
						"hotel_key1",
						false,
						false,
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        d.Body,
							Headers:     headers,
						})
					if err != nil {
						log.Fatalf("Failed to publish a message: %v", err)
					}
					d.Ack(false) // 确认原消息
				}
			} else {
				// 处理成功
				d.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 模拟消息处理函数
func processMessage(body []byte) error {
	// 这里实现具体的业务逻辑
	// 返回nil表示处理成功，返回error表示处理失败

	s := string(body)
	if strings.Index(s, "error") >= 0 {
		return fmt.Errorf("error:%s", s)
	}

	return nil
}

// 死信队列消费
func deadLetterConsumer() {
	conn, ch := initConnection()
	defer conn.Close()
	defer ch.Close()
	// 声明死信交换机
	err := ch.ExchangeDeclare(
		"hotel_dlx_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare dead letter exchange: %v", err)
	}

	// 声明死信队列
	dlq, err := ch.QueueDeclare(
		"hotel_dlx_exchange",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare dead letter queue: %v", err)
	}

	// 绑定死信队列到死信交换机
	err = ch.QueueBind(
		dlq.Name,
		"dlx_key",
		"hotel_dlx_exchange",
		false,
		nil,
	)

	msgs, err := ch.Consume(
		"hotel_dlx_exchange",
		"hotel_dlx_consume",
		true, // 自动确认
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register DLQ consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 记录死信消息
			log.Printf("Dead letter message received: %s", string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for dead letter messages. To exit press CTRL+C")
	<-forever
}
