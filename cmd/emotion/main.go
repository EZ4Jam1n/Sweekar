package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"google.golang.org/grpc"
	emotionv1 "github.com/sweekar/api/proto/emotion/v1"
)

var (
	port = flag.Int("port", 50053, "情绪分析服务端口")
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

	// 创建情绪分析服务实例
	emotionServer := &emotionServer{}

	// 注册情绪分析服务
	emotionv1.RegisterEmotionServiceServer(s, emotionServer)

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动定时推送任务
	go emotionServer.startEmotionReportScheduler(ctx)

	// 处理系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动 gRPC 服务器
	go func() {
		hlog.Infof("情绪分析服务启动，监听端口: %d", *port)
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

// emotionServer 实现情绪分析服务接口
type emotionServer struct {
	emotionv1.UnimplementedEmotionServiceServer
}

// AnalyzeEmotion 实现情绪分析接口
func (s *emotionServer) AnalyzeEmotion(ctx context.Context, req *emotionv1.AnalyzeEmotionRequest) (*emotionv1.AnalyzeEmotionResponse, error) {
	// TODO: 实现情绪分析逻辑
	return &emotionv1.AnalyzeEmotionResponse{}, nil
}

// GetEmotionReport 实现获取情绪报告接口
func (s *emotionServer) GetEmotionReport(ctx context.Context, req *emotionv1.GetEmotionReportRequest) (*emotionv1.GetEmotionReportResponse, error) {
	// TODO: 实现获取情绪报告逻辑
	return &emotionv1.GetEmotionReportResponse{}, nil
}

// PushEmotionReport 实现推送情绪报告接口
func (s *emotionServer) PushEmotionReport(ctx context.Context, req *emotionv1.PushEmotionReportRequest) (*emotionv1.PushEmotionReportResponse, error) {
	// TODO: 实现推送情绪报告逻辑
	return &emotionv1.PushEmotionReportResponse{}, nil
}

// startEmotionReportScheduler 启动情绪报告定时推送任务
func (s *emotionServer) startEmotionReportScheduler(ctx context.Context) {
	ticker := time.NewTicker(time.Hour) // 每小时检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			// 检查是否到达推送时间（每天晚上 8 点）
			if t.Hour() == 20 {
				// TODO: 实现批量推送情绪报告的逻辑
				hlog.Info("开始执行情绪报告推送任务")
			}
		}
	}
}