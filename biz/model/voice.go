package model

import (
	"time"
)

// VoiceMessage 语音消息基础结构
type VoiceMessage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	RetryCount int      `json:"retry_count"`
}

// VADResult VAD处理结果
type VADResult struct {
	VoiceMessage
	IsSpeech     bool          `json:"is_speech"`
	SpeechStart  time.Duration `json:"speech_start"`
	SpeechEnd    time.Duration `json:"speech_end"`
	AudioSegment []byte        `json:"audio_segment"`
}

// ASRResult ASR识别结果
type ASRResult struct {
	VoiceMessage
	Text string `json:"text"`
}

// LLMResult LLM生成结果
type LLMResult struct {
	VoiceMessage
	Response string `json:"response"`
}

// TTSResult TTS转换结果
type TTSResult struct {
	VoiceMessage
	Audio []byte `json:"audio"`
}

// ProcessingStatus 处理状态
type ProcessingStatus string

const (
	StatusPending   ProcessingStatus = "pending"
	StatusRunning   ProcessingStatus = "running"
	StatusComplete  ProcessingStatus = "complete"
	StatusFailed    ProcessingStatus = "failed"
	StatusRetrying  ProcessingStatus = "retrying"
)