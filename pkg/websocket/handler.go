package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/sweekar/biz/model"
	"github.com/sweekar/biz/service"
)

// MessageType 消息类型
type MessageType string

const (
	VoiceChat    MessageType = "voice_chat"    // 语音聊天消息
	VoiceResponse MessageType = "voice_response" // 语音响应消息
)

// Message WebSocket消息结构
type Message struct {
	Type    MessageType     `json:"type"`
	Payload interface{}    `json:"payload"`
}

// Handler WebSocket消息处理器
type Handler struct {
	pool *Pool
	voiceProcessor *service.VoiceProcessor
}

// NewHandler 创建新的消息处理器
func NewHandler(pool *Pool, voiceProcessor *service.VoiceProcessor) *Handler {
	return &Handler{
		pool:           pool,
		voiceProcessor: voiceProcessor,
	}
}

// HandleConnection 处理新的WebSocket连接
func (h *Handler) HandleConnection(conn *websocket.Conn, userID, parentID uint64) {
	client := &Client{
		Conn:     conn,
		UserID:   userID,
		ParentID: parentID,
	}

	// 注册客户端
	h.pool.Register(client)
	defer h.pool.Unregister(userID)

	// 开始接收消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket连接异常关闭: %v", err)
			}
			break
		}

		// 处理接收到的消息
		h.handleMessage(message, client)
	}
}

// handleMessage 处理接收到的消息
func (h *Handler) handleMessage(data []byte, client *Client) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("解析消息失败: %v", err)
		return
	}

	switch msg.Type {
	case VoiceChat:
		h.handleVoiceChat(msg.Payload, client)
	default:
		log.Printf("未知的消息类型: %s", msg.Type)
	}
}



// handleVoiceChat 处理语音聊天消息
func (h *Handler) handleVoiceChat(payload interface{}, client *Client) {
	// 将语音消息转换为VoiceMessage结构
	voiceData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("序列化语音数据失败: %v", err)
		return
	}

	// 创建语音消息对象
	voiceMsg := &model.VoiceMessage{
		UserID: client.UserID,
		Data:   voiceData,
	}

	// 调用语音处理服务
	if err := h.voiceProcessor.ProcessVoice(context.Background(), voiceMsg); err != nil {
		log.Printf("处理语音消息失败: %v", err)
	}
}