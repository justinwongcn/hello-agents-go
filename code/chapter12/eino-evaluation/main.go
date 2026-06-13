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

	einoOpenai "github.com/cloudwego/eino-ext/components/model/openai"
)

func main() {
	fmt.Println("📊 第十二章：智能体性能评估 - Eino 框架实现")
	fmt.Println("支持 BFCL 评估、GAIA 评估、LLM Judge、Win Rate 等评估工具")
	fmt.Println("(输入 'quit' 退出)")
	fmt.Println()

	ctx := context.Background()

	// 检查 API 配置
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")
	modelID := os.Getenv("LLM_MODEL_ID")

	if apiKey == "" || apiKey == "your-api-key-here" {
		log.Println("⚠️ 未配置有效的 LLM_API_KEY，请设置环境变量")
		log.Println("   export LLM_API_KEY=\"your-api-key\"")
		log.Println("   export LLM_BASE_URL=\"https://api.openai.com/v1\"")
		log.Println("   export LLM_MODEL_ID=\"gpt-4o-mini\"")
		log.Println()
	}

	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if modelID == "" {
		modelID = "gpt-4o-mini"
	}

	// 创建评估工具列表
	evaluationTools := []tool.BaseTool{
		NewBFCLEvaluationTool(),
		NewGAIAEvaluationTool(),
		NewLLMJudgeTool(),
		NewWinRateTool(),
	}

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
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "EvaluationAssistant",
		Description: "智能体性能评估助手，支持 BFCL、GAIA、LLM Judge、Win Rate 等评估",
		Instruction: `你是一个智能体性能评估助手。你可以帮助用户：

1. **BFCL 评估**：评估智能体的函数调用能力（Function Calling）
   - 支持多种类别：simple_python, multiple, parallel 等
   - 输出准确率和详细评估结果

2. **GAIA 评估**：评估智能体解决真实世界问题的能力
   - 支持三个难度级别：Level 1（简单）、Level 2（中等）、Level 3（困难）
   - 需要 HuggingFace 访问权限

3. **LLM Judge 评估**：使用 LLM 评估生成内容质量
   - 从正确性、清晰度、难度匹配、完整性四个维度评分
   - 每个维度 1-5 分

4. **Win Rate 评估**：通过对比评估生成质量
   - 与参考真题对比计算胜率
   - Win Rate ≈ 50% 表示质量与真题相当

请用中文回答，并根据用户需求选择合适的评估工具。`,
		Model: chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: evaluationTools,
			},
		},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
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
