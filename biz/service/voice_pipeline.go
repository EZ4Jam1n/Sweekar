package service

import (
	"context"
	"fmt"
	"sync"
)

// VoicePipelineService 语音处理流水线服务
type VoicePipelineService struct {
	// 语音处理器
	processor *VoiceProcessor

	// 服务状态
	isRunning bool
	mutex     sync.RWMutex
}

// NewVoicePipelineService 创建新的语音处理流水线服务
func NewVoicePipelineService(config *VoiceProcessorConfig) (*VoicePipelineService, error) {
	// 创建语音处理器
	processor, err := NewVoiceProcessor(config)
	if err != nil {
		return nil, fmt.Errorf("create voice processor error: %v", err)
	}

	return &VoicePipelineService{
		processor: processor,
		isRunning: false,
	}, nil
}

// Start 启动语音处理流水线服务
func (s *VoicePipelineService) Start(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return nil
	}

	// 启动语音处理器
	if err := s.processor.Start(ctx); err != nil {
		return fmt.Errorf("start voice processor error: %v", err)
	}

	s.isRunning = true
	return nil
}

// Stop 停止语音处理流水线服务
func (s *VoicePipelineService) Stop(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil
	}

	// 停止语音处理器
	if err := s.processor.Stop(ctx); err != nil {
		return fmt.Errorf("stop voice processor error: %v", err)
	}

	s.isRunning = false
	return nil
}
