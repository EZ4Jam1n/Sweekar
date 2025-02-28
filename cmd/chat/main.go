package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"google.golang.org/grpc"
	chatv1 "github.com/sweekar/api/proto/chat/v1"
	emotionv1 "github.com/sweekar/api/proto/emotion/v1"
)

var (
	port = flag.Int("port", 50052, "聊天服务端口")
)

func main() {
	flag.Parse()

	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		hlog.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器实例
	s := grpc.NewServer()

	// 注册聊天服务
	chatv1.RegisterChatServiceServer(s, &chatServer{})

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动 gRPC 服务器
	go func() {
		hlog.Infof("聊天服务启动，监听端口: %d", *port)
		if err := s.Serve(lis); err != nil {
			hlog.Fatalf("failed to serve: %v", err)
		}
	}()

	// 等待系统信号
	<-sigCh
	hlog.Info("正在关闭服务器...")

	// 优雅关闭服务器
	s.GracefulStop()
	hlog.Info("服务器已关闭")
}

// chatServer 实现聊天服务接口
type chatServer struct {
	chatv1.UnimplementedChatServiceServer
	emotionClient emotionv1.EmotionServiceClient
}

// StartChatSession 实现开始聊天会话接口
func (s *chatServer) StartChatSession(ctx context.Context, req *chatv1.StartChatSessionRequest) (*chatv1.StartChatSessionResponse, error) {
	// TODO: 实现开始聊天会话逻辑
	return &chatv1.StartChatSessionResponse{}, nil
}

// SendVoiceMessage 实现发送语音消息接口
func (s *chatServer) SendVoiceMessage(stream chatv1.ChatService_SendVoiceMessageServer) error {
	// TODO: 实现语音消息处理逻辑
	return nil
}

// GetSystemRoles 实现获取系统角色列表接口
func (s *chatServer) GetSystemRoles(ctx context.Context, req *chatv1.GetSystemRolesRequest) (*chatv1.GetSystemRolesResponse, error) {
	// TODO: 实现获取系统角色列表逻辑
	return &chatv1.GetSystemRolesResponse{}, nil
}

// GetChatHistory 实现获取聊天历史记录接口
func (s *chatServer) GetChatHistory(ctx context.Context, req *chatv1.GetChatHistoryRequest) (*chatv1.GetChatHistoryResponse, error) {
	// TODO: 实现获取聊天历史记录逻辑
	return &chatv1.GetChatHistoryResponse{}, nil
}