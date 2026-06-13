package llmclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// HelloAgentsLLM 为本书 "Hello Agents" 定制的LLM客户端。
// 它用于调用任何兼容OpenAI接口的服务，并默认使用流式响应。
type HelloAgentsLLM struct {
	Model    string
	APIKey   string
	BaseURL  string
	Timeout  int
	client   *http.Client
}

// Message 表示一条聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewHelloAgentsLLM 创建新的 LLM 客户端。
// 初始化客户端。优先使用传入参数，如果未提供，则从环境变量加载。
func NewHelloAgentsLLM(model, apiKey, baseURL string, timeout int) (*HelloAgentsLLM, error) {
	// 加载 .env 文件中的环境变量
	_ = godotenv.Load()

	if model == "" {
		model = os.Getenv("LLM_MODEL_ID")
	}
	if apiKey == "" {
		apiKey = os.Getenv("LLM_API_KEY")
	}
	if baseURL == "" {
		baseURL = os.Getenv("LLM_BASE_URL")
	}
	if timeout == 0 {
		timeout = 60
	}

	if model == "" || apiKey == "" || baseURL == "" {
		return nil, fmt.Errorf("模型ID、API密钥和服务地址必须被提供或在.env文件中定义。")
	}

	return &HelloAgentsLLM{
		Model:   model,
		APIKey:  apiKey,
		BaseURL: baseURL,
		Timeout: timeout,
		client:  &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}, nil
}

// Think 调用大语言模型进行思考，并返回其响应。
func (c *HelloAgentsLLM) Think(messages []Message, temperature float64) string {
	fmt.Printf("🧠 正在调用 %s 模型...\n", c.Model)

	// 构建请求体
	reqBody := struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float64   `json:"temperature"`
		Stream      bool      `json:"stream"`
	}{
		Model:       c.Model,
		Messages:    messages,
		Temperature: temperature,
		Stream:      true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("❌ 调用LLM API时发生错误: %v\n", err)
		return ""
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 调用LLM API时发生错误: %v\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("❌ 调用LLM API时发生错误: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	// 处理流式响应
	fmt.Println("✅ 大语言模型响应成功:")
	var collectedContent strings.Builder
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		content := chunk.Choices[0].Delta.Content
		fmt.Print(content)
		collectedContent.WriteString(content)
	}
	fmt.Println() // 在流式输出结束后换行

	return collectedContent.String()
}
