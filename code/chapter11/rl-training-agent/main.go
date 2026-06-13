package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	einoOpenai "github.com/cloudwego/eino-ext/components/model/openai"
)

func main() {
	fmt.Println("🤖 RL Training Agent - 基于 Eino 框架的强化学习训练助手")
	fmt.Println("支持: 加载数据集、SFT/GRPO训练、模型评估、奖励函数创建")
	fmt.Println("(输入 'quit' 退出)")
	fmt.Println()

	ctx := context.Background()

	// 获取脚本目录（默认为上级目录的Python脚本）
	scriptDir := os.Getenv("SCRIPT_DIR")
	if scriptDir == "" {
		// 默认使用chapter11目录（Python脚本所在目录）
		scriptDir = filepath.Join("..")
	}

	pythonPath := os.Getenv("PYTHON_PATH")
	if pythonPath == "" {
		pythonPath = "python3"
	}

	// 检查 API 配置
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")
	modelID := os.Getenv("LLM_MODEL_ID")

	useMock := apiKey == "" || apiKey == "your-api-key-here"
	if useMock {
		log.Println("⚠️  未配置有效的 LLM_API_KEY，将使用模拟模式运行")
		log.Println("   设置环境变量 LLM_API_KEY 以启用完整功能")
		log.Println()
	}

	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if modelID == "" {
		modelID = "gpt-4o-mini"
	}

	// 创建RL训练工具
	rlTool := NewRLTrainingTool(RLTrainingToolConfig{
		ScriptDir: scriptDir,
		PythonPath: pythonPath,
	})

	// 工具列表
	tools := []tool.BaseTool{rlTool}

	// 创建 Agent
	var agent adk.Agent
	if useMock {
		agent = &MockRLAgent{tools: tools}
	} else {
		chatModel, err := einoOpenai.NewChatModel(ctx, &einoOpenai.ChatModelConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Model:   modelID,
		})
		if err != nil {
			log.Fatalf("创建聊天模型失败: %v", err)
		}

		agent, err = adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
			Name:        "RLTrainingAgent",
			Description: "强化学习训练助手，支持SFT/GRPO训练、数据集加载、模型评估",
			Instruction: `你是一个强化学习训练助手。你可以帮助用户：
1. 加载数据集 (action: load_dataset) - 加载GSM8K等数据集
2. 训练模型 (action: train) - 进行SFT或GRPO训练
3. 评估模型 (action: evaluate) - 评估训练后的模型性能
4. 创建奖励函数 (action: create_reward) - 创建准确性/长度惩罚/步骤奖励

常用参数：
- model_name: 模型名称，如 "Qwen/Qwen3-0.6B"
- algorithm: 训练算法 "sft" 或 "grpo"
- max_samples: 最大样本数
- num_epochs: 训练轮数
- batch_size: 批量大小
- use_lora: 是否使用LoRA
- lora_r: LoRA秩
- lora_alpha: LoRA缩放因子

请用中文回答，帮助用户完成强化学习训练任务。`,
			Model: chatModel,
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

// ============================================================================
// MockRLAgent 模拟 Agent（当没有配置API密钥时使用）
// ============================================================================

type MockRLAgent struct {
	tools []tool.BaseTool
}

func (a *MockRLAgent) Name(ctx context.Context) string {
	return "MockRLAgent"
}

func (a *MockRLAgent) Description(ctx context.Context) string {
	return "强化学习训练模拟助手"
}

func (a *MockRLAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	go func() {
		defer gen.Close()

		if len(input.Messages) == 0 {
			gen.Send(&adk.AgentEvent{
				Output: &adk.AgentOutput{
					MessageOutput: &adk.MessageVariant{
						Message: schema.AssistantMessage("请输入您的问题，例如：'加载GSM8K数据集' 或 '开始GRPO训练'", nil),
					},
				},
			})
			return
		}

		query := input.Messages[len(input.Messages)-1].Content

		response := fmt.Sprintf("收到您的请求：%s\n\n"+
			"这是一个使用 Eino 框架的 RL Training Agent 演示。\n"+
			"在实际应用中，这里会调用 RLTrainingTool 来执行操作。\n\n"+
			"可用操作:\n"+
			"- load_dataset: 加载数据集 (GSM8K)\n"+
			"- train: 训练模型 (SFT/GRPO)\n"+
			"- evaluate: 评估模型\n"+
			"- create_reward: 创建奖励函数\n\n"+
			"常用参数:\n"+
			"- model_name: Qwen/Qwen3-0.6B\n"+
			"- algorithm: sft / grpo\n"+
			"- max_samples: 样本数\n"+
			"- use_lora: true/false\n\n"+
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
