package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageType 消息类型
type MessageType string

const (
	TextMessage    MessageType = "text"     // 文本消息
	VoiceMessage   MessageType = "voice"    // 语音消息
	EmotionMessage MessageType = "emotion"  // 情绪消息
)

// ChatMessage 聊天消息记录
type ChatMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    uint64            `bson:"user_id" json:"user_id"`         // 发送者ID
	ParentID  uint64            `bson:"parent_id" json:"parent_id"`     // 家长ID
	Type      MessageType       `bson:"type" json:"type"`               // 消息类型
	Content   string            `bson:"content" json:"content"`         // 消息内容
	Emotion   *EmotionData     `bson:"emotion,omitempty" json:"emotion"` // 情绪数据
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`   // 创建时间
}