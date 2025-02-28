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
	userv1 "github.com/sweekar/api/proto/user/v1"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "用户服务端口")
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

	// 注册用户服务
	userv1.RegisterUserServiceServer(s, &userServer{})

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动 gRPC 服务器
	go func() {
		hlog.Infof("用户服务启动，监听端口: %d", *port)
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

// userServer 实现用户服务接口
type userServer struct {
	userv1.UnimplementedUserServiceServer
}

// ParentRegister 实现家长注册接口
func (s *userServer) ParentRegister(ctx context.Context, req *userv1.ParentRegisterRequest) (*userv1.ParentRegisterResponse, error) {
	// TODO: 实现家长注册逻辑
	return &userv1.ParentRegisterResponse{}, nil
}

// ParentLogin 实现家长登录接口
func (s *userServer) ParentLogin(ctx context.Context, req *userv1.ParentLoginRequest) (*userv1.ParentLoginResponse, error) {
	// TODO: 实现家长登录逻辑
	return &userv1.ParentLoginResponse{}, nil
}

// CreateChildAccount 实现创建儿童账户接口
func (s *userServer) CreateChildAccount(ctx context.Context, req *userv1.CreateChildAccountRequest) (*userv1.CreateChildAccountResponse, error) {
	// TODO: 实现创建儿童账户逻辑
	return &userv1.CreateChildAccountResponse{}, nil
}

// GetChildAccounts 实现获取儿童账户列表接口
func (s *userServer) GetChildAccounts(ctx context.Context, req *userv1.GetChildAccountsRequest) (*userv1.GetChildAccountsResponse, error) {
	// TODO: 实现获取儿童账户列表逻辑
	return &userv1.GetChildAccountsResponse{}, nil
}
