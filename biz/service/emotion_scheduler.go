package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"sweekar/biz/model"
)

// EmotionScheduler 情绪报告调度器
type EmotionScheduler struct {
	db         *gorm.DB
	mqProducer rocketmq.Producer
	cron       *cron.Cron
	processor  *EmotionProcessor
}

// NewEmotionScheduler 创建情绪报告调度器
func NewEmotionScheduler(db *gorm.DB, producer rocketmq.Producer, processor *EmotionProcessor) *EmotionScheduler {
	scheduler := &EmotionScheduler{
		db:         db,
		mqProducer: producer,
		cron:       cron.New(cron.WithSeconds()),
		processor:  processor,
	}

	// 每天19:00生成情绪报告
	_, err := scheduler.cron.AddFunc("0 0 19 * * ?", scheduler.generateDailyReports)
	if err != nil {
		panic(fmt.Sprintf("添加生成报告定时任务失败: %v", err))
	}

	// 每天20:00推送情绪报告
	_, err = scheduler.cron.AddFunc("0 0 20 * * ?", scheduler.pushDailyReports)
	if err != nil {
		panic(fmt.Sprintf("添加推送报告定时任务失败: %v", err))
	}

	return scheduler
}

// Start 启动调度器
func (s *EmotionScheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *EmotionScheduler) Stop() {
	s.cron.Stop()
}

// generateDailyReports 生成所有用户的每日情绪报告
func (s *EmotionScheduler) generateDailyReports() {
	// 获取所有有聊天记录的用户
	var userIDs []uint64
	today := time.Now().Truncate(24 * time.Hour)

	err := s.db.Model(&model.EmotionRecord{}).
		Where("DATE(created_at) = DATE(?)", today).
		Distinct().
		Pluck("user_id", &userIDs).Error

	if err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
		return
	}

	// 为每个用户生成报告
	for _, userID := range userIDs {
		err := s.processor.GenerateDailyReport(context.Background(), userID, today)
		if err != nil {
			fmt.Printf("生成用户 %d 的情绪报告失败: %v\n", userID, err)
			continue
		}
	}
}

// pushDailyReports 推送所有用户的每日情绪报告
func (s *EmotionScheduler) pushDailyReports() {
	// 获取今天生成但未推送的报告
	var reports []model.EmotionReport
	today := time.Now().Truncate(24 * time.Hour)

	err := s.db.Where("DATE(date) = DATE(?) AND pushed_at IS NULL", today).Find(&reports).Error
	if err != nil {
		fmt.Printf("获取待推送报告失败: %v\n", err)
		return
	}

	// 推送每个报告
	for _, report := range reports {
		// 序列化报告数据
		reportData, err := json.Marshal(report)
		if err != nil {
			fmt.Printf("序列化报告 %d 失败: %v\n", report.ID, err)
			continue
		}

		// 发送推送消息
		msg := primitive.NewMessage("emotion_report_push", reportData)
		msg.WithKeys([]string{fmt.Sprintf("user_%d", report.UserID)})

		_, err = s.mqProducer.SendSync(context.Background(), msg)
		if err != nil {
			fmt.Printf("推送报告 %d 失败: %v\n", report.ID, err)
			continue
		}

		// 更新推送时间
		report.PushedAt = time.Now()
		if err := s.db.Save(&report).Error; err != nil {
			fmt.Printf("更新报告 %d 推送时间失败: %v\n", report.ID, err)
		}
	}
}