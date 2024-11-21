package rabbit_mq

func RabbitMq() {
	producer()
	//consumer()
}
func consumer() {
	go deadLetterConsumer()
	go handlerConsumer()
	for {
		select {}
	}

}
