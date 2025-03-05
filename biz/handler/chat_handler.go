package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/sweekar/pkg/websocket"
    "github.com/sweekar/biz/service"
)

type ChatHandler struct {
    chatService *service.ChatService
    wsHandler  *websocket.Handler
}

func NewChatHandler(chatService *service.ChatService, wsHandler *websocket.Handler) *ChatHandler {
    return &ChatHandler{
        chatService: chatService,
        wsHandler:  wsHandler,
    }
}

// WebSocket连接处理
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "未授权"})
        return
    }

    // 升级HTTP连接为WebSocket连接
    if err := h.wsHandler.HandleConnection(c.Writer, c.Request, userID); err != nil {
        c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error()})
        return
    }
}

// 获取用户的聊天历史记录
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "未授权"})
        return
    }

    history, err := h.chatService.GetChatHistory(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error()})
        return
    }

    c.JSON(http.StatusOK, Response{Code: 200, Message: "获取成功", Data: history})
}