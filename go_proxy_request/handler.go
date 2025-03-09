package goproxyrequest

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GoProxyRequest() {
	// 创建一个resty客户端
	client := resty.New()

	// 设置代理
	client.SetProxy("http://renhao:renhao666@localhost:8890")

	// 发送请求
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get("http://www.youtube.com")

	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
		return
	}

	// 打印响应
	fmt.Printf("响应状态码: %v\n", resp.StatusCode())
	fmt.Printf("响应内容: %v\n", resp.String())
}
