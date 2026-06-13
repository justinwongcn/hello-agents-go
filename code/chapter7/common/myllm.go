package common

import (
	"fmt"
	"os"

	"hello-agents-go/chapter4/llmclient"
)

// MyLLM 自定义LLM客户端，支持 ModelScope Provider
type MyLLM struct {
	*llmclient.HelloAgentsLLM
	Provider    string
	Temperature float64
	MaxTokens   int
}

// NewMyLLM 创建自定义LLM客户端实例。
// 支持 ModelScope Provider 和默认 Provider。
func NewMyLLM(model, apiKey, baseURL, provider string, temperature float64, maxTokens, timeout int) (*MyLLM, error) {
	if provider == "" {
		provider = "auto"
	}

	if provider == "modelscope" {
		fmt.Println("正在使用自定义的 ModelScope Provider")

		// 解析 ModelScope 的凭证
		if apiKey == "" {
			apiKey = os.Getenv("MODELSCOPE_API_KEY")
		}
		if baseURL == "" {
			baseURL = "https://api-inference.modelscope.cn/v1/"
		}

		// 验证凭证是否存在
		if apiKey == "" {
			return nil, fmt.Errorf("ModelScope API key not found. Please set MODELSCOPE_API_KEY environment variable.")
		}

		// 设置默认模型和其他参数
		if model == "" {
			model = os.Getenv("LLM_MODEL_ID")
		}
		if model == "" {
			model = "Qwen/Qwen2.5-VL-72B-Instruct"
		}
		if temperature == 0 {
			temperature = 0.7
		}
		if timeout == 0 {
			timeout = 60
		}

		// 使用获取的参数创建LLM客户端实例
		base, err := llmclient.NewHelloAgentsLLM(model, apiKey, baseURL, timeout)
		if err != nil {
			return nil, err
		}

		return &MyLLM{
			HelloAgentsLLM: base,
			Provider:       "modelscope",
			Temperature:    temperature,
			MaxTokens:      maxTokens,
		}, nil
	}

	// 如果不是 modelscope, 则完全使用默认逻辑来处理
	base, err := llmclient.NewHelloAgentsLLM(model, apiKey, baseURL, timeout)
	if err != nil {
		return nil, err
	}

	if temperature == 0 {
		temperature = 0.7
	}

	return &MyLLM{
		HelloAgentsLLM: base,
		Provider:       provider,
		Temperature:    temperature,
		MaxTokens:      maxTokens,
	}, nil
}
