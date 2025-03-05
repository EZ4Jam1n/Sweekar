package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatedier/beego/logs"
	"github.com/sashabaranov/go-openai"

	"github.com/sweekar/biz/model"
	"github.com/sweekar/pkg/mq"
)

// VoiceProcessor 语音处理器
type VoiceProcessor struct {
	mqClient *mq.RocketMQClient

	// VAD配置
	vadModel   *whisper.Context
	vadConfig  *model.VADConfig
	vadWorkers int
	vadTopic   string

	// ASR配置
	asrClient  *funasr.Client
	asrWorkers int
	asrTopic   string

	// LLM配置
	llmClient  *openai.Client
	llmWorkers int
	llmTopic   string

	// TTS配置
	ttsClient  *tts.Client
	ttsWorkers int
	ttsTopic   string

	// WebSocket配置
	wsPool *websocket.Pool
}

// NewVoiceProcessor 创建语音处理器
func NewVoiceProcessor(config *VoiceProcessorConfig) (*VoiceProcessor, error) {
	// 初始化RocketMQ客户端
	mqClient := mq.NewRocketMQClient(&mq.RocketMQConfig{
		NameServers: config.MQNameServers,
		GroupID:     config.MQGroupID,
		MaxRetries:  config.MQMaxRetries,
	})

	// 初始化VAD模型
	vadModel, err := whisper.New(config.VADModelPath)
	if err != nil {
		return nil, fmt.Errorf("init vad model error: %v", err)
	}

	// 初始化ASR客户端
	asrClient, err := funasr.NewClient(config.ASRConfig)
	if err != nil {
		return nil, fmt.Errorf("init asr client error: %v", err)
	}

	// 初始化LLM客户端
	llmClient := openai.NewClient(config.LLMAPIKey)

	// 初始化TTS客户端
	ttsClient, err := tts.NewClient(config.TTSConfig)
	if err != nil {
		return nil, fmt.Errorf("init tts client error: %v", err)
	}

	return &VoiceProcessor{
		mqClient:   mqClient,
		vadModel:   vadModel,
		vadConfig:  config.VADConfig,
		vadWorkers: config.VADWorkers,
		vadTopic:   config.VADTopic,
		asrClient:  asrClient,
		asrWorkers: config.ASRWorkers,
		asrTopic:   config.ASRTopic,
		llmClient:  llmClient,
		llmWorkers: config.LLMWorkers,
		llmTopic:   config.LLMTopic,
		ttsClient:  ttsClient,
		ttsWorkers: config.TTSWorkers,
		ttsTopic:   config.TTSTopic,
		wsPool:     websocket.NewPool(),
	}, nil
}

// Start 启动语音处理器
func (p *VoiceProcessor) Start(ctx context.Context) error {
	// 启动RocketMQ客户端
	if err := p.mqClient.Start(ctx); err != nil {
		return err
	}

	// 启动VAD消费者
	if err := p.startVADConsumer(ctx); err != nil {
		return err
	}

	// 启动ASR消费者
	if err := p.startASRConsumer(ctx); err != nil {
		return err
	}

	// 启动LLM消费者
	if err := p.startLLMConsumer(ctx); err != nil {
		return err
	}

	// 启动TTS消费者
	if err := p.startTTSConsumer(ctx); err != nil {
		return err
	}

	return nil
}

// Stop 停止语音处理器
func (p *VoiceProcessor) Stop(ctx context.Context) error {
	return p.mqClient.Stop(ctx)
}

// ProcessVoice 处理语音消息
func (p *VoiceProcessor) ProcessVoice(ctx context.Context, msg *model.VoiceMessage) error {
	return p.mqClient.SendMessage(ctx, p.vadTopic, msg)
}

// startVADConsumer 启动VAD消费者
func (p *VoiceProcessor) startVADConsumer(ctx context.Context) error {
	return p.mqClient.ConsumeMessage(ctx, p.vadTopic, p.vadWorkers, func(ctx context.Context, data []byte) error {
		var msg model.VoiceMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}

		// 执行VAD处理
		result := p.processVAD(&msg)
		if result.IsSpeech {
			// 发送到ASR队列
			return p.mqClient.SendMessage(ctx, p.asrTopic, result)
		}
		return nil
	})
}

// startASRConsumer 启动ASR消费者
func (p *VoiceProcessor) startASRConsumer(ctx context.Context) error {
	return p.mqClient.ConsumeMessage(ctx, p.asrTopic, p.asrWorkers, func(ctx context.Context, data []byte) error {
		var vadResult model.VADResult
		if err := json.Unmarshal(data, &vadResult); err != nil {
			return err
		}

		// 执行ASR处理
		result := p.processASR(&vadResult)
		// 发送到LLM队列
		return p.mqClient.SendMessage(ctx, p.llmTopic, result)
	})
}

// startLLMConsumer 启动LLM消费者
func (p *VoiceProcessor) startLLMConsumer(ctx context.Context) error {
	return p.mqClient.ConsumeMessage(ctx, p.llmTopic, p.llmWorkers, func(ctx context.Context, data []byte) error {
		var asrResult model.ASRResult
		if err := json.Unmarshal(data, &asrResult); err != nil {
			return err
		}

		// 执行LLM处理
		result := p.processLLM(&asrResult)
		// 发送到TTS队列
		return p.mqClient.SendMessage(ctx, p.ttsTopic, result)
	})
}

// startTTSConsumer 启动TTS消费者
func (p *VoiceProcessor) startTTSConsumer(ctx context.Context) error {
	return p.mqClient.ConsumeMessage(ctx, p.ttsTopic, p.ttsWorkers, func(ctx context.Context, data []byte) error {
		var llmResult model.LLMResult
		if err := json.Unmarshal(data, &llmResult); err != nil {
			return err
		}

		// 执行TTS处理
		result := p.processTTS(&llmResult)
		// TODO: 发送结果到WebSocket
		return nil
	})
}

// processVAD 执行VAD处理
func (p *VoiceProcessor) processVAD(msg *model.VoiceMessage) *model.VADResult {
	// TODO: 实现VAD处理逻辑
	return &model.VADResult{
		VoiceMessage: *msg,
		IsSpeech:     true,
		SpeechStart:  0,
		SpeechEnd:    time.Duration(len(msg.Data)) * time.Millisecond,
		AudioSegment: msg.Data,
	}
}

// processASR 执行ASR处理
func (p *VoiceProcessor) processASR(vadResult *model.VADResult) *model.ASRResult {
	// 调用FunASR进行语音识别
	text, err := p.asrClient.Recognize(vadResult.AudioSegment)
	if err != nil {
		logs.Error("ASR recognition error: %v", err)
		return &model.ASRResult{
			VoiceMessage: vadResult.VoiceMessage,
			Text:         "",
		}
	}

	return &model.ASRResult{
		VoiceMessage: vadResult.VoiceMessage,
		Text:         text,
	}
}

// processLLM 执行LLM处理
func (p *VoiceProcessor) processLLM(asrResult *model.ASRResult) *model.LLMResult {
	// 调用OpenAI API生成响应
	resp, err := p.llmClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: asrResult.Text,
				},
			},
		},
	)

	if err != nil {
		logs.Error("LLM generation error: %v", err)
		return &model.LLMResult{
			VoiceMessage: asrResult.VoiceMessage,
			Response:     "",
		}
	}

	return &model.LLMResult{
		VoiceMessage: asrResult.VoiceMessage,
		Response:     resp.Choices[0].Message.Content,
	}
}

// processTTS 执行TTS处理
func (p *VoiceProcessor) processTTS(llmResult *model.LLMResult) *model.TTSResult {
	// 调用TTS服务进行语音合成
	audio, err := p.ttsClient.Synthesize(llmResult.Response)
	if err != nil {
		logs.Error("TTS synthesis error: %v", err)
		return &model.TTSResult{
			VoiceMessage: llmResult.VoiceMessage,
			Audio:        nil,
		}
	}

	result := &model.TTSResult{
		VoiceMessage: llmResult.VoiceMessage,
		Audio:        audio,
	}

	// 通过WebSocket推送TTS结果给客户端
	if client := p.wsPool.GetClient(llmResult.UserID); client != nil {
		msg := websocket.Message{
			Type:    "voice_response",
			Payload: result,
		}
		msgData, err := json.Marshal(msg)
		if err != nil {
			logs.Error("序列化TTS结果失败: %v", err)
		} else {
			client.Mu.Lock()
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, msgData); err != nil {
				logs.Error("推送TTS结果失败: %v", err)
			}
			client.Mu.Unlock()
		}
	}

	return result
}
