package api

import (
    "github.com/gin-gonic/gin"
    "github.com/sweekar/biz/handler"
    "github.com/sweekar/pkg/middleware"
)

func SetupRouter(userHandler *handler.UserHandler, chatHandler *handler.ChatHandler, emotionHandler *handler.EmotionHandler) *gin.Engine {
    router := gin.Default()

    // 用户服务API
    userGroup := router.Group("/api/v1/user")
    {
        userGroup.POST("/register", userHandler.Register)
        userGroup.POST("/login", userHandler.Login)
    }

    // 需要认证的API组
    authGroup := router.Group("/api/v1")
    authGroup.Use(middleware.AuthMiddleware())
    {
        // WebSocket连接
        authGroup.GET("/ws", chatHandler.HandleWebSocket)

        // 聊天服务API
        chatGroup := authGroup.Group("/chat")
        {
            chatGroup.GET("/history", chatHandler.GetChatHistory)
        }

        // 情绪分析服务API
        emotionGroup := authGroup.Group("/emotion")
        {
            emotionGroup.GET("/report", emotionHandler.GetEmotionReport)
            emotionGroup.GET("/trend", emotionHandler.GetEmotionTrend)
        }
    }

    return router
}