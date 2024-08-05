package go_config

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/yaml"
)

func GoConfig() {
	// 设置选项支持ENV变量解析：当获取的值为string类型时，会尝试解析其中的ENV变量
	config.WithOptions(config.ParseEnv)

	// 添加驱动程序以支持yaml内容解析（除了JSON是默认支持，其他的则是按需使用）
	config.AddDriver(yaml.Driver)

	// 加载配置，可以同时传入多个文件
	err := config.LoadFiles("./go_config/config.yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("config data: \n %#v\n", config.Data())
	fmt.Printf("config shell: \n %#v\n", config.String("envKey1"))

	// 添加驱动程序以支持yaml内容解析（除了JSON是默认支持，其他的则是按需使用）
	config.AddDriver(ini.Driver)
	// 加载更多文件
	err = config.LoadFiles("./go_config/config.ini")
	// 也可以一次性加载多个文件
	// err := config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())
	fmt.Printf("config shell: \n %#v\n", config.String("debug666"))

}
