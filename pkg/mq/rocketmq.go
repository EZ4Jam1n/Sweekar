package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// RocketMQConfig RocketMQ配置
type RocketMQConfig struct {
	NameServers []string
	GroupID     string
	MaxRetries  int
}

// RocketMQClient RocketMQ客户端
type RocketMQClient struct {
	config    *RocketMQConfig
	producer  rocketmq.Producer
	consumers map[string]rocketmq.PushConsumer
	mutex     sync.RWMutex
}

// NewRocketMQClient 创建新的RocketMQ客户端
func NewRocketMQClient(config *RocketMQConfig) *RocketMQClient {
	return &RocketMQClient{
		config:    config,
		consumers: make(map[string]rocketmq.PushConsumer),
	}
}

// Start 启动RocketMQ客户端
func (c *RocketMQClient) Start(ctx context.Context) error {
	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(c.config.NameServers),
		producer.WithGroupName(c.config.GroupID),
		producer.WithRetry(c.config.MaxRetries),
	)
	if err != nil {
		return fmt.Errorf("create producer error: %v", err)
	}

	if err := p.Start(); err != nil {
		return fmt.Errorf("start producer error: %v", err)
	}

	c.producer = p
	return nil
}

// Stop 停止RocketMQ客户端
func (c *RocketMQClient) Stop(ctx context.Context) error {
	// 停止生产者
	if err := c.producer.Shutdown(); err != nil {
		return fmt.Errorf("shutdown producer error: %v", err)
	}

	// 停止所有消费者
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, cons := range c.consumers {
		if err := cons.Shutdown(); err != nil {
			return fmt.Errorf("shutdown consumer error: %v", err)
		}
	}

	return nil
}

// SendMessage 发送消息
func (c *RocketMQClient) SendMessage(ctx context.Context, topic string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message error: %v", err)
	}

	message := primitive.NewMessage(topic, body)
	message.WithProperty("timestamp", time.Now().String())

	_, err = c.producer.SendSync(ctx, message)
	if err != nil {
		return fmt.Errorf("send message error: %v", err)
	}

	return nil
}

// ConsumeMessage 消费消息
func (c *RocketMQClient) ConsumeMessage(ctx context.Context, topic string, numWorkers int, handler func(context.Context, []byte) error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 创建消费者
	cons, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(c.config.NameServers),
		consumer.WithGroupName(c.config.GroupID),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeMessageBatchMaxSize(1),
	)
	if err != nil {
		return fmt.Errorf("create consumer error: %v", err)
	}

	// 创建工作线程池
	workerPool := make(chan struct{}, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerPool <- struct{}{}
	}

	// 订阅主题
	err = cons.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		// 获取工作线程
		<-workerPool
		defer func() { workerPool <- struct{}{} }()

		for _, msg := range msgs {
			err := handler(ctx, msg.Body)
			if err != nil {
				return consumer.ConsumeRetryLater, err
			}
		}

		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		return fmt.Errorf("subscribe topic error: %v", err)
	}

	// 启动消费者
	if err := cons.Start(); err != nil {
		return fmt.Errorf("start consumer error: %v", err)
	}

	c.consumers[topic] = cons
	return nil
}