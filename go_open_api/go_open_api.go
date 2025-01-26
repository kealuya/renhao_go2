package go_open_api

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io"
)

func GoOpenApi() {
	// 创建客户端
	// 设置自定义基础URL
	config := openai.DefaultConfig("sk-z'x'x")
	config.BaseURL = "https://api.siliconflow.cn/v1"
	client := openai.NewClientWithConfig(config)

	// 创建请求
	req := openai.ChatCompletionRequest{
		Model: "ft:LoRA/Qwen/Qwen2.5-14B-Instruct:cm34107y100gh11v5k8e6fbb9:renhao:orhofnrpwixzpyquwiut",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "极川 笔记本支架电脑散热支架铝合金折叠便携悬空立式增高架电脑桌支架苹果Macbook联想拯救者华为架子",
			},
		},
		MaxTokens: 4096,
		Stream:    true,
	}

	// 创建流式请求
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Printf("流式请求创建失败: %v\n", err)
		return
	}
	defer stream.Close()

	// 接收流式响应
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("流式接收错误: %v\n", err)
			return
		}

		fmt.Print(response.Choices[0].Delta.Content)
	}
}
