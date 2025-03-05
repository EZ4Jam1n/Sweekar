package model

import (
	"time"
)

// EmotionType 情绪类型
type EmotionType string

const (
	EmotionHappy   EmotionType = "happy"   // 开心
	EmotionSad     EmotionType = "sad"     // 难过
	EmotionAngry   EmotionType = "angry"   // 生气
	EmotionNeutral EmotionType = "neutral" // 平静
)

// EmotionRecord 情绪记录
type EmotionRecord struct {
	ID        uint64      `json:"id" gorm:"primaryKey"`
	UserID    uint64      `json:"user_id" gorm:"index"`              // 用户ID
	ChatID    uint64      `json:"chat_id" gorm:"index"`              // 聊天记录ID
	Emotion   EmotionType `json:"emotion"`                            // 情绪类型
	Confidence float64     `json:"confidence"`                         // 情绪判断的置信度
	CreatedAt time.Time   `json:"created_at" gorm:"index"`           // 创建时间
}

// EmotionReport 情绪报告
type EmotionReport struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	UserID       uint64    `json:"user_id" gorm:"index"`                     // 用户ID
	Date         time.Time `json:"date" gorm:"index;type:date"`              // 报告日期
	ChatCount    int       `json:"chat_count"`                                // 当天聊天次数
	EmotionStats map[EmotionType]int `json:"emotion_stats" gorm:"type:json"` // 情绪统计
	Summary      string    `json:"summary"`                                   // 情绪总结
	CreatedAt    time.Time `json:"created_at"`                                // 创建时间
	PushedAt     time.Time `json:"pushed_at"`                                 // 推送时间
}