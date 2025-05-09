package goasynq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
)

// ==================== 任务定义 ====================

// 任务类型常量
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

// EmailDeliveryPayload 表示发送邮件任务的数据
type EmailDeliveryPayload struct {
	UserID     int    `json:"user_id"`
	TemplateID string `json:"template_id"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
}

// ImageResizePayload 表示图片调整大小任务的数据
type ImageResizePayload struct {
	SourceURL string `json:"source_url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// NewEmailDeliveryTask 创建一个新的邮件发送任务
func NewEmailDeliveryTask(userID int, templateID, to, subject string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{
		UserID:     userID,
		TemplateID: templateID,
		To:         to,
		Subject:    subject,
	})
	if err != nil {
		return nil, err
	}
	// 第一个参数是任务类型，第二个参数是任务数据
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

// NewImageResizeTask 创建一个新的图片调整大小任务
func NewImageResizeTask(sourceURL string, width, height int) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{
		SourceURL: sourceURL,
		Width:     width,
		Height:    height,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageResize, payload), nil
}

// HandleEmailDeliveryTask 处理邮件发送任务
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("发送邮件给用户ID: %d", p.UserID)
	log.Printf("模板ID: %s", p.TemplateID)
	log.Printf("收件人: %s", p.To)
	log.Printf("主题: %s", p.Subject)

	// 模拟邮件发送耗时
	time.Sleep(2 * time.Second)

	log.Printf("邮件已成功发送给 %s", p.To)
	return nil
}

// HandleImageResizeTask 处理图片调整大小任务
func HandleImageResizeTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("调整图片大小: %s", p.SourceURL)
	log.Printf("新尺寸: %dx%d", p.Width, p.Height)

	// 模拟图片处理耗时
	time.Sleep(3 * time.Second)

	log.Printf("图片 %s 已成功调整为 %dx%d", p.SourceURL, p.Width, p.Height)
	return nil
}

// ==================== 客户端/生产者 ====================

// Client 是任务生产者
type Client struct {
	client *asynq.Client
}

// NewClient 创建一个新的任务客户端
func NewClient(redisAddr, redisPassword string) *Client {
	// 创建Redis连接
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
	}

	// 创建asynq客户端
	client := asynq.NewClient(redisOpt)

	return &Client{
		client: client,
	}
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	return c.client.Close()
}

// EnqueueEmailTask 将邮件任务加入队列
func (c *Client) EnqueueEmailTask(userID int, templateID, to, subject string) error {
	// 创建邮件任务
	task, err := NewEmailDeliveryTask(userID, templateID, to, subject)
	if err != nil {
		return fmt.Errorf("无法创建邮件任务: %v", err)
	}

	// 设置任务选项
	opts := []asynq.Option{
		asynq.MaxRetry(5),                 // 最多重试5次
		asynq.ProcessIn(10 * time.Second), // 10秒后开始处理
		asynq.Queue("critical"),           // 放入critical队列
	}

	// 加入队列
	info, err := c.client.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("无法将任务加入队列: %v", err)
	}

	log.Printf("已将邮件任务加入队列: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueImageResizeTask 将图片调整大小任务加入队列
func (c *Client) EnqueueImageResizeTask(sourceURL string, width, height int) error {
	// 创建图片调整大小任务
	task, err := NewImageResizeTask(sourceURL, width, height)
	if err != nil {
		return fmt.Errorf("无法创建图片调整大小任务: %v", err)
	}

	// 设置任务选项
	opts := []asynq.Option{
		asynq.MaxRetry(3),                // 最多重试3次
		asynq.ProcessIn(5 * time.Second), // 5秒后开始处理
		asynq.Queue("default"),           // 放入default队列
	}

	// 加入队列
	info, err := c.client.Enqueue(task, opts...)
	// 加入队列
	c.client.Enqueue(task, opts...)
	// 加入队列
	c.client.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("无法将任务加入队列: %v", err)
	}

	log.Printf("已将图片调整大小任务加入队列: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// ==================== 工作者/消费者 ====================

// StartWorker 启动工作者服务器
func StartWorker(redisAddr, redisPassword string) error {
	// 创建Redis连接
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
	}

	// 创建服务器
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// 指定每个队列的并发处理数量
			Queues: map[string]int{
				"critical": 10, // critical队列可以10个并发
				"default":  5,  // default队列可以5个并发
				"low":      2,  // low队列可以2个并发
			},
			// 错误处理器
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("处理任务 %s 时出错: %v", task.Type(), err)
			}),
		},
	)

	// 创建任务多路复用器
	mux := asynq.NewServeMux()

	// 注册任务处理器
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.HandleFunc(TypeImageResize, HandleImageResizeTask)

	log.Println("异步任务工作者已启动...")

	// 启动服务器
	if err := srv.Run(mux); err != nil {
		return err
	}

	return nil
}

// ==================== 主程序 ====================

func Main() {

	// Redis配置
	redisAddr := "124.70.48.30:6379" // Redis地址
	redisPassword := "renhao666"     // Redis密码

	go runWorker(redisAddr, redisPassword)
	runClient(redisAddr, redisPassword)

}

// 运行客户端（任务生产者）
func runClient(redisAddr, redisPassword string) {
	log.Println("启动任务生产者...")

	// 创建客户端
	c := NewClient(redisAddr, redisPassword)
	defer c.Close()

	// 创建一些示例任务
	if err := c.EnqueueEmailTask(
		123,
		"welcome-template",
		"new-user@example.com",
		"欢迎加入我们的平台",
	); err != nil {
		log.Fatalf("无法创建邮件任务: %v", err)
	}

	if err := c.EnqueueImageResizeTask(
		"https://example.com/image.jpg",
		800,
		600,
	); err != nil {
		log.Fatalf("无法创建图片调整大小任务: %v", err)
	}

	log.Println("任务已加入队列，按Ctrl+C退出")

	// 等待中断信号
	waitForInterrupt()
}

// 运行工作者（任务消费者）
func runWorker(redisAddr, redisPassword string) {
	log.Println("启动任务消费者...")

	// 启动工作者
	if err := StartWorker(redisAddr, redisPassword); err != nil {
		log.Fatalf("工作者出错: %v", err)
	}
}

// 等待中断信号
func waitForInterrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
