package main

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func init() {
	LogConfigInit()
}

func LogConfigInit() {
	_ = logs.SetLogger(logs.AdapterConsole)
	_ = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/my.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":365,"color":true,"separate":["error", "warning", "info", "debug"]}`)
	//输出文件名和行号
	logs.EnableFuncCallDepth(true)
	//异步输出log
	//logs.Async()
}

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
	//=====================================
	//HttpService()
	LlmChat()
}

func LlmChat() {
	cmd := exec.Command("./app", "如何改进这些代码")
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}

func HttpService() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		fmt.Println(request.RequestURI)
		fmt.Println(request.URL)
		fmt.Println(request.RemoteAddr)

		_, err := writer.Write([]byte("hello world"))
		if err != nil {
			log.Panicln(err)
		}
	})

	logs.Info("it's work~!")

	err := http.ListenAndServe(":9009", nil)
	if err != nil {
		log.Panicln(err)
	}

}
