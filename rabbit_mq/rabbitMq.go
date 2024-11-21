package rabbit_mq

import "renhao_go2/rabbit_mq/batch"

func RabbitMq() {
	//producer()
	//consumer()
	batch.Main()
}
func consumer() {
	go deadLetterConsumer()
	go handlerConsumer()
	for {
		select {}
	}

}
