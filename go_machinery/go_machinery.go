package gomachinery

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// ==================== 任务定义 ====================

// 任务类型常量
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

// ==================== 任务处理函数 ====================

// HandleEmailDeliveryTask 处理邮件发送任务
func HandleEmailDeliveryTask(userID int, templateID, to, subject string) error {
	log.Printf("发送邮件给用户ID: %d", userID)
	log.Printf("模板ID: %s", templateID)
	log.Printf("收件人: %s", to)
	log.Printf("主题: %s", subject)

	// 模拟邮件发送耗时
	time.Sleep(2 * time.Second)

	log.Printf("邮件已成功发送给 %s", to)
	return nil
}

// HandleImageResizeTask 处理图片调整大小任务
func HandleImageResizeTask(sourceURL string, width, height int) error {
	log.Printf("调整图片大小: %s", sourceURL)
	log.Printf("新尺寸: %dx%d", width, height)

	// 模拟图片处理耗时
	time.Sleep(3 * time.Second)

	log.Printf("图片 %s 已成功调整为 %dx%d", sourceURL, width, height)
	return nil
}

// ==================== 客户端/生产者 ====================

// Client 是任务生产者
type Client struct {
	client *machinery.Server
}

// NewClient 创建一个新的任务客户端
func NewClient(redisAddr, redisPassword string) (*Client, error) {
	// 创建Machinery配置
	cnf := &config.Config{
		Broker:        fmt.Sprintf("redis://%s:%s@%s", "", redisPassword, redisAddr),
		ResultBackend: fmt.Sprintf("redis://%s:%s@%s", "", redisPassword, redisAddr),
		DefaultQueue:  "machinery_tasks",
	}

	// 创建Machinery服务器
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, fmt.Errorf("无法创建Machinery服务器: %v", err)
	}

	// 注册任务
	err = server.RegisterTasks(map[string]interface{}{
		TypeEmailDelivery: HandleEmailDeliveryTask,
		TypeImageResize:   HandleImageResizeTask,
	})
	if err != nil {
		return nil, fmt.Errorf("无法注册任务: %v", err)
	}

	return &Client{
		client: server,
	}, nil
}

// Close 关闭客户端连接 (Machinery不需要显式关闭)
func (c *Client) Close() error {
	// Machinery不需要显式关闭连接
	return nil
}

// EnqueueEmailTask 将邮件任务加入队列
func (c *Client) EnqueueEmailTask(userID int, templateID, to, subject string) error {
	// 创建任务签名
	task := &tasks.Signature{
		Name: TypeEmailDelivery,
		Args: []tasks.Arg{
			{
				Type:  "int",
				Value: userID,
			},
			{
				Type:  "string",
				Value: templateID,
			},
			{
				Type:  "string",
				Value: to,
			},
			{
				Type:  "string",
				Value: subject,
			},
		},
		RetryCount:   5,          // 最多重试5次
		RetryTimeout: 10,         // 10秒后重试
		RoutingKey:   "critical", // 放入critical队列
	}

	// 设置延迟执行（10秒后）
	eta := time.Now().Add(10 * time.Second)
	task.ETA = &eta

	// 加入队列
	asyncResult, err := c.client.SendTask(task)
	if err != nil {
		return fmt.Errorf("无法将任务加入队列: %v", err)
	}

	log.Printf("已将邮件任务加入队列: id=%s", asyncResult.Signature.UUID)
	return nil
}

// EnqueueImageResizeTask 将图片调整大小任务加入队列
func (c *Client) EnqueueImageResizeTask(sourceURL string, width, height int) error {
	// 创建任务签名
	task := &tasks.Signature{
		Name: TypeImageResize,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: sourceURL,
			},
			{
				Type:  "int",
				Value: width,
			},
			{
				Type:  "int",
				Value: height,
			},
		},
		RetryCount:   3,         // 最多重试3次
		RetryTimeout: 5,         // 5秒后重试
		RoutingKey:   "default", // 放入default队列
	}

	// 设置延迟执行（5秒后）
	eta := time.Now().Add(5 * time.Second)
	task.ETA = &eta

	// 加入队列（三次，模拟原代码中的多次调用）
	for i := 0; i < 3; i++ {
		asyncResult, err := c.client.SendTask(task)
		if err != nil {
			return fmt.Errorf("无法将任务加入队列: %v", err)
		}
		log.Printf("已将图片调整大小任务加入队列: id=%s", asyncResult.Signature.UUID)
	}

	return nil
}

// ==================== 工作者/消费者 ====================

// StartWorker 启动工作者服务器
func StartWorker(redisAddr, redisPassword string) error {
	// 创建Machinery配置
	cnf := &config.Config{
		Broker:        fmt.Sprintf("redis://%s:%s@%s", "", redisPassword, redisAddr),
		ResultBackend: fmt.Sprintf("redis://%s:%s@%s", "", redisPassword, redisAddr),
		DefaultQueue:  "machinery_tasks",
	}

	// 创建Machinery服务器
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return fmt.Errorf("无法创建Machinery服务器: %v", err)
	}
	// 注册任务
	err = server.RegisterTasks(map[string]interface{}{
		TypeEmailDelivery: HandleEmailDeliveryTask,
		TypeImageResize:   HandleImageResizeTask,
	})
	if err != nil {
		return fmt.Errorf("无法注册任务: %v", err)
	}

	// 创建Worker
	worker := server.NewWorker("machinery-worker", 10) // 同时处理10个任务

	// 配置不同队列的消费者数量
	// Machinery默认会使用RoutingKey作为队列名称
	// 我们可以在启动worker时为不同队列配置不同数量的worker

	log.Println("异步任务工作者已启动...")

	// 启动Worker
	err = worker.Launch()
	if err != nil {
		return fmt.Errorf("启动工作者失败: %v", err)
	}

	return nil
}

// ==================== 主程序 ====================

func Main() {
	// Redis配置
	redisAddr := "124.70.48.30:6379" // Redis地址
	redisPassword := "renhao666"     // Redis密码

	go runWorker(redisAddr, redisPassword)
	time.Sleep(5 * time.Second)
	runClient(redisAddr, redisPassword)
}

// 运行客户端（任务生产者）
func runClient(redisAddr, redisPassword string) {
	log.Println("启动任务生产者...")

	// 创建客户端
	c, err := NewClient(redisAddr, redisPassword)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
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
