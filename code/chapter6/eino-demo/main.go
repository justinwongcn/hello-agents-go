package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	einoOpenai "github.com/cloudwego/eino-ext/components/model/openai"
)

// WeatherTool 天气查询工具
type WeatherTool struct{}

func (t *WeatherTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_weather",
		Desc: "获取指定城市的当前天气信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"city": {
				Type:     "string",
				Desc:     "城市名称，如：北京、上海、广州",
				Required: true,
			},
		}),
	}, nil
}

func (t *WeatherTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	city := extractStringArg(argumentsInJSON, "city")
	if city == "" {
		city = "未知城市"
	}

	weatherData := map[string]string{
		"北京": "晴天，气温 25°C，空气质量良好",
		"上海": "多云，气温 28°C，湿度较高",
		"广州": "阵雨，气温 32°C，注意带伞",
		"深圳": "晴转多云，气温 31°C，海风较大",
		"杭州": "阴天，气温 26°C，适合出游",
	}

	if weather, ok := weatherData[city]; ok {
		return fmt.Sprintf("%s的天气：%s", city, weather), nil
	}
	return fmt.Sprintf("%s的天气：晴天，气温 22°C，天气宜人", city), nil
}

// CalculatorTool 计算器工具
type CalculatorTool struct{}

func (t *CalculatorTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "calculator",
		Desc: "执行简单的数学计算",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"expression": {
				Type:     "string",
				Desc:     "数学表达式，如：2+3*4",
				Required: true,
			},
		}),
	}, nil
}

func (t *CalculatorTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	expression := extractStringArg(argumentsInJSON, "expression")
	if expression == "" {
		return "错误：未提供表达式", nil
	}

	switch expression {
	case "2+3":
		return "2+3 = 5", nil
	case "10*5":
		return "10*5 = 50", nil
	case "100/4":
		return "100/4 = 25", nil
	default:
		return fmt.Sprintf("计算结果：%s = [需要实际计算]", expression), nil
	}
}

// SearchTool 搜索工具
type SearchTool struct{}

func (t *SearchTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "web_search",
		Desc: "搜索互联网上的最新信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     "string",
				Desc:     "搜索关键词",
				Required: true,
			},
		}),
	}, nil
}

func (t *SearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	query := extractStringArg(argumentsInJSON, "query")
	if query == "" {
		return "错误：未提供搜索关键词", nil
	}

	return fmt.Sprintf("搜索「%s」的结果：\n1. 相关文章标题 - 这是一篇关于%s的详细文章\n2. 最新新闻 - %s的最新动态\n3. 技术文档 - %s的技术指南", query, query, query, query), nil
}

// extractStringArg 从 JSON 参数中提取字符串值
func extractStringArg(argumentsInJSON, key string) string {
	searchStr := fmt.Sprintf(`"%s":"`, key)
	start := strings.Index(argumentsInJSON, searchStr)
	if start == -1 {
		return ""
	}
	start += len(searchStr)
	end := strings.Index(argumentsInJSON[start:], `"`)
	if end == -1 {
		return ""
	}
	return argumentsInJSON[start : start+end]
}

func main() {
	fmt.Println("🤖 Eino 智能助手 - 基于 CloudWeGo Eino 框架")
	fmt.Println("支持天气查询、数学计算、网络搜索等功能")
	fmt.Println("(输入 'quit' 退出)")
	fmt.Println()

	ctx := context.Background()

	// 检查 API 配置
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")
	modelID := os.Getenv("LLM_MODEL_ID")

	useMock := apiKey == "" || apiKey == "your-api-key-here"
	if useMock {
		log.Println("⚠️ 未配置有效的 LLM_API_KEY，将使用模拟模式运行")
		log.Println("   设置环境变量 LLM_API_KEY 以启用完整功能")
		log.Println()
	}

	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if modelID == "" {
		modelID = "gpt-4o-mini"
	}

	// 创建工具列表
	tools := []tool.BaseTool{
		&WeatherTool{},
		&CalculatorTool{},
		&SearchTool{},
	}

	// 创建 Agent
	var agent adk.Agent
	if useMock {
		// 模拟模式
		agent = &MockAgent{tools: tools}
	} else {
		// 创建 OpenAI 兼容的聊天模型
		chatModel, err := einoOpenai.NewChatModel(ctx, &einoOpenai.ChatModelConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Model:   modelID,
		})
		if err != nil {
			log.Fatalf("创建聊天模型失败: %v", err)
		}

		// 创建 ADK Agent
		agent, err = adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
			Name:        "SmartAssistant",
			Description: "一个支持天气查询、数学计算和网络搜索的智能助手",
			Instruction: "你是一个智能助手，可以帮助用户查询天气、进行数学计算和搜索信息。请用中文回答。",
			Model:       chatModel,
			ToolsConfig: adk.ToolsConfig{
				ToolsNodeConfig: compose.ToolsNodeConfig{
					Tools: tools,
				},
			},
		})
		if err != nil {
			log.Fatalf("创建 Agent 失败: %v", err)
		}
	}

	// 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	// 交互式对话循环
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("🤔 您的问题: ")
		if !scanner.Scan() {
			break
		}
		query := strings.TrimSpace(scanner.Text())

		if strings.ToLower(query) == "quit" || strings.ToLower(query) == "q" || query == "退出" {
			fmt.Println("感谢使用！再见！👋")
			break
		}

		if query == "" {
			continue
		}

		fmt.Printf("\n%s\n", strings.Repeat("=", 60))

		// 执行查询
		iter := runner.Query(ctx, query)
		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				fmt.Printf("❌ 错误: %v\n", event.Err)
				continue
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				msg, err := event.Output.MessageOutput.GetMessage()
				if err != nil {
					fmt.Printf("❌ 获取消息失败: %v\n", err)
					continue
				}
				if msg != nil && msg.Content != "" {
					fmt.Printf("💡 回答: %s\n", msg.Content)
				}
			}
		}

		fmt.Printf("\n%s\n\n", strings.Repeat("=", 60))
	}
}

// MockAgent 模拟 Agent（当没有配置 API 密钥时使用）
type MockAgent struct {
	tools []tool.BaseTool
}

func (a *MockAgent) Name(ctx context.Context) string {
	return "MockAgent"
}

func (a *MockAgent) Description(ctx context.Context) string {
	return "一个模拟的智能助手"
}

func (a *MockAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	go func() {
		defer gen.Close()

		if len(input.Messages) == 0 {
			gen.Send(&adk.AgentEvent{
				Output: &adk.AgentOutput{
					MessageOutput: &adk.MessageVariant{
						Message: schema.AssistantMessage("请输入您的问题", nil),
					},
				},
			})
			return
		}

		query := input.Messages[len(input.Messages)-1].Content

		// 简单的模拟响应
		response := fmt.Sprintf("收到您的问题：%s\n\n"+
			"这是一个使用 Eino 框架的智能助手演示。\n"+
			"在实际应用中，这里会调用 LLM 并使用工具来回答您的问题。\n\n"+
			"可用工具：\n"+
			"- get_weather: 查询天气\n"+
			"- calculator: 数学计算\n"+
			"- web_search: 网络搜索\n\n"+
			"请配置 LLM_API_KEY 环境变量以启用完整功能。", query)

		gen.Send(&adk.AgentEvent{
			Output: &adk.AgentOutput{
				MessageOutput: &adk.MessageVariant{
					Message: schema.AssistantMessage(response, nil),
				},
			},
		})
	}()

	return iter
}
