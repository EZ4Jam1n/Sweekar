package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/sweekar/biz/model"
    "github.com/sweekar/biz/service"
)

type EmotionHandler struct {
    emotionProcessor *service.EmotionProcessor
}

func NewEmotionHandler(emotionProcessor *service.EmotionProcessor) *EmotionHandler {
    return &EmotionHandler{emotionProcessor: emotionProcessor}
}

// 获取用户的情绪分析报告
func (h *EmotionHandler) GetEmotionReport(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "未授权"})
        return
    }

    report, err := h.emotionProcessor.GetEmotionReport(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error()})
        return
    }

    c.JSON(http.StatusOK, Response{Code: 200, Message: "获取成功", Data: report})
}

// 获取用户的情绪趋势分析
func (h *EmotionHandler) GetEmotionTrend(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "未授权"})
        return
    }

    startTime := c.Query("start_time")
    endTime := c.Query("end_time")

    trend, err := h.emotionProcessor.GetEmotionTrend(c.Request.Context(), userID, startTime, endTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error()})
        return
    }

    c.JSON(http.StatusOK, Response{Code: 200, Message: "获取成功", Data: trend})
}