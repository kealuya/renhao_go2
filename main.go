package main

import (
	"net/http"
)

func main() {
	//rabbit_mq.RabbitMq()
	//goutil.GoUtil()
	//protobuf.Protobuf()
	//progress_go.ProgressGo()
	//go_house.GoHouse()
	//v8go.V8go()
	//cobra.Cobra()
	//copier.Copier()
	//tray_cmd.TrayCmd()
	//aggregate.FuzzyFind()
	//aggregate.Jwt()
	//badger.Badger()
	//go_config.GoConfig()

	HttpService()
}

func HttpService() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})
	http.ListenAndServe(":9009", nil)

}
