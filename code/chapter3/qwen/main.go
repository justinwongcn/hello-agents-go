package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const systemPrompt = "You are a helpful assistant."

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model     string        `json:"model"`
	Messages  []chatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	// 增加HF_ENDPOINT，避免Connection aborted.
	os.Setenv("HF_ENDPOINT", "https://hf-mirror.com")

	// 指定模型ID
	modelID := "Qwen/Qwen1.5-0.5B-Chat"

	// 设置设备，优先使用GPU（在Go中通过环境变量控制）
	device := "cpu"
	fmt.Printf("Using device: %s\n", device)

	// 创建HTTP客户端连接到本地模型服务
	// 假设模型通过本地API服务暴露（如vLLM、Ollama等）
	apiBase := os.Getenv("OPENAI_API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost:8000/v1"
	}

	fmt.Println("模型和分词器加载完成！")

	// 准备对话输入
	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "你好，请介绍你自己。"},
	}

	// 使用模型生成回答
	// max_new_tokens 控制了模型最多能生成多少个新的Token
	reqBody := chatRequest{
		Model:     modelID,
		Messages:  messages,
		MaxTokens: 512,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("JSON编码失败: %v", err)
	}

	resp, err := http.Post(apiBase+"/chat/completions", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("模型生成失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		log.Fatalf("JSON解码失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		log.Fatal("模型未返回任何回答")
	}

	fmt.Println("\n模型的回答:")
	fmt.Println(chatResp.Choices[0].Message.Content)
}
