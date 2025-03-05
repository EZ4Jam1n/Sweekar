package service

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"

	"sweekar/biz/model"
	"sweekar/pkg/websocket"
)

// EmotionProcessor 情绪处理器
type EmotionProcessor struct {
	db           *gorm.DB
	mqProducer   rocketmq.Producer
	mqConsumer   rocketmq.PushConsumer
	wsHandler    *websocket.Handler
}

// NewEmotionProcessor 创建情绪处理器
func NewEmotionProcessor(db *gorm.DB, producer rocketmq.Producer, consumer rocketmq.PushConsumer, wsHandler *websocket.Handler) *EmotionProcessor {
	return &EmotionProcessor{
		db:         db,
		mqProducer: producer,
		mqConsumer: consumer,
		wsHandler:  wsHandler,
	}
}

// AnalyzeEmotion 分析聊天内容的情绪
func (p *EmotionProcessor) AnalyzeEmotion(ctx context.Context, chatID uint64, userID uint64, content string) error {
	// TODO: 接入情绪分析AI模型
	// 这里模拟情绪分析结果
	emotion := model.EmotionRecord{
		UserID:    userID,
		ChatID:    chatID,
		Emotion:   model.EmotionHappy,
		Confidence: 0.85,
		CreatedAt: time.Now(),
	}

	// 保存情绪记录
	if err := p.db.Create(&emotion).Error; err != nil {
		return err
	}

	return nil
}

// GenerateDailyReport 生成每日情绪报告
func (p *EmotionProcessor) GenerateDailyReport(ctx context.Context, userID uint64, date time.Time) error {
	// 获取当天的聊天情绪记录
	var records []model.EmotionRecord
	if err := p.db.Where("user_id = ? AND DATE(created_at) = DATE(?)", userID, date).Find(&records).Error; err != nil {
		return err
	}

	// 如果当天没有聊天记录，则不生成报告
	if len(records) == 0 {
		return nil
	}

	// 统计情绪分布
	emotionStats := make(map[model.EmotionType]int)
	for _, record := range records {
		emotionStats[record.Emotion]++
	}

	// 生成情绪报告
	report := model.EmotionReport{
		UserID:       userID,
		Date:         date,
		ChatCount:    len(records),
		EmotionStats: emotionStats,
		Summary:      generateEmotionSummary(emotionStats), // 生成情绪总结
		CreatedAt:    time.Now(),
	}

	// 保存情绪报告
	if err := p.db.Create(&report).Error; err != nil {
		return err
	}

	// 发送报告生成消息到消息队列
	msg := primitive.NewMessage("emotion_report", []byte(report.ID))
	_, err := p.mqProducer.SendSync(context.Background(), msg)
	if err != nil {
		return err
	}

	// 通过WebSocket实时推送报告给家长
	return p.wsHandler.BroadcastEmotionReport(userID, report)
}

// generateEmotionSummary 生成情绪总结
func generateEmotionSummary(stats map[model.EmotionType]int) string {
	// TODO: 根据情绪统计生成更智能的总结
	return "今天整体心情不错，多数时候比较开心"
}