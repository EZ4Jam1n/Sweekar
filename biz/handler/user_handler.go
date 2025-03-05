package handler

import (
    "encoding/json"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/sweekar/biz/model"
    "github.com/sweekar/biz/service"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
        return
    }

    user := &model.User{
        Username: req.Username,
        Password: req.Password,
        Email:    req.Email,
    }

    if err := h.userService.Register(c.Request.Context(), user); err != nil {
        c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error()})
        return
    }

    c.JSON(http.StatusOK, Response{Code: 200, Message: "注册成功"})
}

func (h *UserHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
        return
    }

    token, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: err.Error()})
        return
    }

    c.JSON(http.StatusOK, Response{Code: 200, Message: "登录成功", Data: gin.H{"token": token}})
}