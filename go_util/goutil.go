package go_util

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/encodes"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/goutil/timex"
)

func GoUtil() {

	//时间转换
	now := timex.Now()
	fmt.Println(now.String())

	fmt.Println(now.Format("2006-01-02 15:04:05.000000"))

	tx, err := timex.FromString("2022-04-20 19:40:34")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tx.String())

	//工程路径转换
	fmt.Println(cliutil.Workdir()) // current workdir
	fmt.Println(cliutil.BinDir())  // the program exe file dir
	fmt.Println(cliutil.BinName())
	//cliutil.ReadInput("Your name?")
	//cliutil.ReadPassword("Input password:")
	//ans, _ := cliutil.ReadFirstByte("continue?[y/n] ")
	//fmt.Println(ans)

	//	dumper 处理，格式化输出
	dump.P(
		23,
		[]string{"ab", "cd"},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, // len > 10
		map[string]interface{}{
			"key": "val", "sub": map[string]string{"k": "v"},
		},
		struct {
			ab string
			Cd int
		}{
			"ab", 23,
		},
	)

	fmt.Println(envutil.IsWin())
	fmt.Println(envutil.IsMac())
	fmt.Println(envutil.IsLinux())

	fmt.Println(envutil.EnvMap())
	fmt.Println(envutil.EnvPaths())

	// get ENV value by key, can with default value
	fmt.Println(envutil.Getenv("APP_ENV", "dev"))
	fmt.Println(envutil.GetInt("LOG_LEVEL", 1))
	fmt.Println(envutil.GetBool("APP_DEBUG", true))

	b := encodes.B64Encode(("renhao666"))
	fmt.Println(string(byteutil.Md5("renhao666")))
	fmt.Println(b)
	fmt.Println(encodes.B64Decode(b))

	pj, _ := fmtutil.PrettyJSON(`{"key":"val","sub":{"k":"v"}}`)
	fmt.Println(pj)

	//	加载配置文件
	// 设置选项支持ENV变量解析：当获取的值为string类型时，会尝试解析其中的ENV变量
	config.WithOptions(config.ParseEnv)

	// 添加驱动程序以支持yaml内容解析（除了JSON是默认支持，其他的则是按需使用）
	config.AddDriver(yaml.Driver)

	// 加载配置，可以同时传入多个文件
	err1 := config.LoadFiles("goutil/yaml_other.yml")
	if err1 != nil {
		panic(err1)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())

}
