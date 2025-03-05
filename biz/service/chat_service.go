package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/your-username/Sweekar/biz/model"
)

// ChatService 聊天服务
type ChatService struct {
	coll *mongo.Collection
}

// NewChatService 创建新的聊天服务
func NewChatService(db *mongo.Database) *ChatService {
	return &ChatService{
		coll: db.Collection("chat_messages"),
	}
}

// SaveMessage 保存聊天消息
func (s *ChatService) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	msg.CreatedAt = time.Now()
	_, err := s.coll.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("保存聊天消息失败: %v", err)
	}
	return nil
}

// GetUserMessages 获取用户的聊天记录
func (s *ChatService) GetUserMessages(ctx context.Context, userID uint64, limit int64) ([]*model.ChatMessage, error) {
	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(limit)

	cursor, err := s.coll.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("查询聊天记录失败: %v", err)
	}
	defer cursor.Close(ctx)

	var messages []*model.ChatMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("解析聊天记录失败: %v", err)
	}

	return messages, nil
}

// GetParentMessages 获取家长的聊天记录
func (s *ChatService) GetParentMessages(ctx context.Context, parentID uint64, limit int64) ([]*model.ChatMessage, error) {
	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(limit)

	cursor, err := s.coll.Find(ctx, bson.M{"parent_id": parentID}, opts)
	if err != nil {
		return nil, fmt.Errorf("查询聊天记录失败: %v", err)
	}
	defer cursor.Close(ctx)

	var messages []*model.ChatMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("解析聊天记录失败: %v", err)
	}

	return messages, nil
}

// GetEmotionMessages 获取情绪消息记录
func (s *ChatService) GetEmotionMessages(ctx context.Context, userID uint64, startTime, endTime time.Time) ([]*model.ChatMessage, error) {
	filter := bson.M{
		"user_id": userID,
		"type": model.EmotionMessage,
		"created_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询情绪记录失败: %v", err)
	}
	defer cursor.Close(ctx)

	var messages []*model.ChatMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("解析情绪记录失败: %v", err)
	}

	return messages, nil
}