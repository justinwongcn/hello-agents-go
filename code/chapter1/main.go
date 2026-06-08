package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// ---------------------------------------------------------------------------
// 1. 系统提示词
// ---------------------------------------------------------------------------
var AGENT_SYSTEM_PROMPT = `
你是一个智能旅行助手。你的任务是分析用户的请求，并使用可用工具一步步地解决问题。

# 可用工具:
- ` + "`get_weather(city: str)`" + `: 查询指定城市的实时天气。
- ` + "`get_attraction(city: str, weather: str)`" + `: 根据城市和天气搜索推荐的旅游景点。

# 输出格式要求:
你的每次回复必须严格遵循以下格式，包含一对Thought和Action：

Thought: [你的思考过程和下一步计划]
Action: [你要执行的具体行动]

Action的格式必须是以下之一：
1. 调用工具：function_name(arg_name="arg_value")
2. 结束任务：Finish[最终答案]

# 重要提示:
- 每次只输出一对Thought-Action
- Action必须在同一行，不要换行
- 当收集到足够信息可以回答用户问题时，必须使用 Action: Finish[最终答案] 格式结束

请开始吧！
`

// ---------------------------------------------------------------------------
// 2. HTTP 客户端（复用连接、超时）
// ---------------------------------------------------------------------------
var httpClient = &http.Client{Timeout: 30 * time.Second}

// ---------------------------------------------------------------------------
// 3. 天气查询工具（wttr.in API）
// ---------------------------------------------------------------------------

// wttrResponse 对应 wttr.in 返回的 JSON 结构
type wttrResponse struct {
	CurrentCondition []struct {
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		TempC string `json:"temp_C"`
	} `json:"current_condition"`
}

// getWeather 通过调用 wttr.in API 查询真实的天气信息
func getWeather(city string) string {
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", city)

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("错误：查询天气时遇到网络问题 - %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("错误：查询天气时遇到网络问题 - HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("错误：查询天气时遇到网络问题 - %v", err)
	}

	var data wttrResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Sprintf("错误：解析天气数据失败，可能是城市名称无效 - %v", err)
	}

	if len(data.CurrentCondition) == 0 || len(data.CurrentCondition[0].WeatherDesc) == 0 {
		return "错误：解析天气数据失败，可能是城市名称无效"
	}

	weatherDesc := data.CurrentCondition[0].WeatherDesc[0].Value
	tempC := data.CurrentCondition[0].TempC

	return fmt.Sprintf("%s当前天气：%s，气温%s摄氏度", city, weatherDesc, tempC)
}

// ---------------------------------------------------------------------------
// 4. 景点推荐工具（Tavily Search API）
// ---------------------------------------------------------------------------

// tavilyRequest Tavily Search API 请求体
type tavilyRequest struct {
	Query         string `json:"query"`
	SearchDepth   string `json:"search_depth"`
	IncludeAnswer bool   `json:"include_answer"`
}

// tavilyResponse Tavily Search API 响应体
type tavilyResponse struct {
	Answer  string `json:"answer"`
	Results []struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	} `json:"results"`
}

// getAttraction 根据城市和天气，使用 Tavily Search API 搜索并返回景点推荐
func getAttraction(city, weather string) string {
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return "错误：未配置TAVILY_API_KEY。"
	}

	query := fmt.Sprintf("'%s' 在'%s'天气下最值得去的旅游景点推荐及理由", city, weather)
	reqBody := tavilyRequest{
		Query:         query,
		SearchDepth:   "basic",
		IncludeAnswer: true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Sprintf("错误：执行Tavily搜索时出现问题 - %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.tavily.com/search", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Sprintf("错误：执行Tavily搜索时出现问题 - %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Sprintf("错误：执行Tavily搜索时出现问题 - %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("错误：执行Tavily搜索时出现问题 - %v", err)
	}

	var result tavilyResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Sprintf("错误：执行Tavily搜索时出现问题 - %v", err)
	}

	// Tavily 返回的 answer 是一个基于所有搜索结果的总结性回答
	if result.Answer != "" {
		return result.Answer
	}

	// 无综合回答时，格式化原始结果
	var formatted []string
	for _, r := range result.Results {
		formatted = append(formatted, fmt.Sprintf("- %s: %s", r.Title, r.Content))
	}

	if len(formatted) == 0 {
		return "抱歉，没有找到相关的旅游景点推荐。"
	}

	return "根据搜索，为您找到以下信息：\n" + strings.Join(formatted, "\n")
}

// ---------------------------------------------------------------------------
// 5. OpenAI 兼容客户端（使用官方 Go SDK）
// ---------------------------------------------------------------------------

// OpenAICompatibleClient 是一个用于调用任何兼容OpenAI接口的LLM服务的客户端
type OpenAICompatibleClient struct {
	model  string
	client *openai.Client
}

// NewOpenAICompatibleClient 创建新的 LLM 客户端
func NewOpenAICompatibleClient(model, apiKey, baseURL string) *OpenAICompatibleClient {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	return &OpenAICompatibleClient{
		model:  model,
		client: &client,
	}
}

// Generate 调用 LLM API 来生成回应
func (c *OpenAICompatibleClient) Generate(ctx context.Context, prompt, systemPrompt string) string {
	fmt.Println("正在调用大语言模型...")

	chat, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModel(c.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		fmt.Printf("调用LLM API时发生错误: %v\n", err)
		return "错误：调用语言模型服务时出错。"
	}

	if len(chat.Choices) == 0 {
		fmt.Println("调用LLM API时发生错误: 无返回结果")
		return "错误：调用语言模型服务时出错。"
	}

	answer := chat.Choices[0].Message.Content
	fmt.Println("大语言模型响应成功。")
	return answer
}

// ---------------------------------------------------------------------------
// 6. 工具注册表
// ---------------------------------------------------------------------------

// ToolFunc 工具函数签名：接收解析好的参数字典，返回观察结果
type ToolFunc func(args map[string]string) string

// availableTools 将所有工具函数放入一个字典，方便后续调用
var availableTools = map[string]ToolFunc{
	"get_weather": func(args map[string]string) string {
		return getWeather(args["city"])
	},
	"get_attraction": func(args map[string]string) string {
		return getAttraction(args["city"], args["weather"])
	},
}

// ---------------------------------------------------------------------------
// 7. 正则表达式（一次性编译）
// ---------------------------------------------------------------------------
var (
	// 用于截断多余 Thought-Action 对：匹配第一对，边界是下一个标签或字符串结尾
	reTruncate = regexp.MustCompile(`(?s)(Thought:.*?Action:.*?)(?:\n\s*(?:Thought:|Action:|Observation:)|\z)`)
	// 提取 Action 行（贪婪匹配到结尾，后续进一步解析）
	reAction = regexp.MustCompile(`(?s)Action: (.*)`)
	// 提取 Finish 中的最终答案
	reFinish = regexp.MustCompile(`^Finish\[(.*)\]`)
	// 提取工具名（"funcName(" 中的 funcName）
	reToolName = regexp.MustCompile(`(\w+)\(`)
	// 提取键值对参数（key="value"）
	reArgs = regexp.MustCompile(`(\w+)="([^"]*)"`)
)

// ---------------------------------------------------------------------------
// 8. Main — ReAct 主循环
// ---------------------------------------------------------------------------
func main() {
	// --- 1. 加载环境变量 ---
	_ = godotenv.Load() // 忽略 .env 文件不存在的错误

	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")
	modelID := os.Getenv("LLM_MODEL_ID")

	llm := NewOpenAICompatibleClient(modelID, apiKey, baseURL)
	ctx := context.Background()

	// --- 2. 初始化 ---
	userPrompt := "你好，请帮我查询一下今天北京的天气，然后根据天气推荐一个合适的旅游景点。"
	promptHistory := []string{fmt.Sprintf("用户请求: %s", userPrompt)}

	fmt.Printf("用户输入: %s\n", userPrompt)
	fmt.Println(strings.Repeat("=", 40))

	// --- 3. 运行主循环 ---
	for i := range 5 { // 设置最大循环次数
		fmt.Printf("--- 循环 %d ---\n\n", i+1)

		// 3.1. 构建 Prompt
		fullPrompt := strings.Join(promptHistory, "\n")

		// 3.2. 调用 LLM 进行思考
		llmOutput := llm.Generate(ctx, fullPrompt, AGENT_SYSTEM_PROMPT)

		// 模型可能会输出多余的 Thought-Action，需要截断
		match := reTruncate.FindStringSubmatch(llmOutput)
		if match != nil {
			truncated := strings.TrimSpace(match[1])
			if truncated != strings.TrimSpace(llmOutput) {
				llmOutput = truncated
				fmt.Println("已截断多余的 Thought-Action 对")
			}
		}
		fmt.Printf("模型输出:\n%s\n\n", llmOutput)
		promptHistory = append(promptHistory, llmOutput)

		// 3.3. 解析并执行行动
		actionMatch := reAction.FindStringSubmatch(llmOutput)
		if actionMatch == nil {
			observation := "错误: 未能解析到 Action 字段。请确保你的回复严格遵循 'Thought: ... Action: ...' 的格式。"
			observationStr := fmt.Sprintf("Observation: %s", observation)
			fmt.Printf("%s\n%s\n", observationStr, strings.Repeat("=", 40))
			promptHistory = append(promptHistory, observationStr)
			continue
		}
		actionStr := strings.TrimSpace(actionMatch[1])

		if strings.HasPrefix(actionStr, "Finish") {
			finishMatch := reFinish.FindStringSubmatch(actionStr)
			if finishMatch != nil {
				finalAnswer := finishMatch[1]
				fmt.Printf("任务完成，最终答案: %s\n", finalAnswer)
			}
			break
		}

		toolMatch := reToolName.FindStringSubmatch(actionStr)
		if toolMatch == nil {
			observation := fmt.Sprintf("错误：无法解析工具调用 - %s", actionStr)
			observationStr := fmt.Sprintf("Observation: %s", observation)
			fmt.Printf("%s\n%s\n", observationStr, strings.Repeat("=", 40))
			promptHistory = append(promptHistory, observationStr)
			continue
		}
		toolName := toolMatch[1]

		// 提取参数
		argsMatch := reArgs.FindAllStringSubmatch(actionStr, -1)
		kwargs := make(map[string]string, len(argsMatch))
		for _, m := range argsMatch {
			kwargs[m[1]] = m[2]
		}

		var observation string
		if fn, ok := availableTools[toolName]; ok {
			observation = fn(kwargs)
		} else {
			observation = fmt.Sprintf("错误：未定义的工具 '%s'", toolName)
		}

		// 3.4. 记录观察结果
		observationStr := fmt.Sprintf("Observation: %s", observation)
		fmt.Printf("%s\n%s\n", observationStr, strings.Repeat("=", 40))
		promptHistory = append(promptHistory, observationStr)
	}
}
