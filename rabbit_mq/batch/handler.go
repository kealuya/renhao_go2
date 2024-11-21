package batch

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"sync"
	"time"
)

// Main
/*
	程序建议就是
	1，初始化：connection，开channel（channel就是用来作业的基础单位），在mq声明交换机(exchange)，声明路由(routing_key)，然后声明队列(queue)，绑定交换机路由和队列。
	2，生产者启动：使用channel，进行消息发送，指定交换机和路由（自然也就存储到指定队列中）。
	3，消费者：使用channel，直接消费对应的队列就可以。（一个goroutine 一个channel）
	如果初始化已经声明了交换机、路由，并且绑定了路由、交换机，队列，那么消费者其实只需要指定对应名称的队列，就可以消费。
	** 每次batch处理，如果中间有错误停掉了程序，记得清空队列，不然计数会发生错误 **
*/
func Main() {

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	// 构造生产者
	p, err := GetMqProductInstance()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 构造消费者
	for i := 0; i < 5; i++ {

		c, err := GetMqConsumeInstance(p, ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		go func() {
			c.Consume(func(bytes []byte) error {
				fmt.Println("handle ::" + string(bytes))
				time.Sleep(1 * time.Second)
				// fixme
				if string(bytes) == "3" {
					return errors.New("error~~ ::" + string(bytes))
				}
				wg.Done()
				return nil
			})
		}()

	}

	// 构造死信队列
	dead, err := GetMqDeadLetterInstance(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		dead.Consume(func(bytes []byte) {
			fmt.Println("dead handle ::" + string(bytes))
			// 记录死信消息
			logs.Warning("Dead letter message received: %s", string(bytes))
			wg.Done()
		})
	}()

	// 发送生产者消息
	for i := 0; i < 100; i++ {
		wg.Add(1)
		err := p.Publish([]byte(strconv.Itoa(i)), 0)
		if err != nil {
			fmt.Println("p.Publishs", err)
			//return
		}

	}
	fmt.Println("发送完成", 100)

	wg.Wait()
	cancel()
	fmt.Println("over")
	time.Sleep(2 * time.Second)
	p.Close()
}
