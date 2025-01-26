package main

import (
	"github.com/beego/beego/v2/core/logs"
	"renhao_go2/excel"
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
	//go_open_api.GoOpenApi()
	excel.Go_excel()
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
	//LlmChat()
	//cgo.RunCGo()
	//=====================================
	//HttpService()
	//CoreRun()
	//darwinkit.GoDarwinkit()

}
