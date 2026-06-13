package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"hello-agents-go/chapter4/llmclient"
	"hello-agents-go/chapter7/common"
)

// MY_REACT_PROMPT ReAct 智能体的提示词模板
const MY_REACT_PROMPT = `你是一个具备推理和行动能力的AI助手。你可以通过思考分析问题，然后调用合适的工具来获取信息，最终给出准确的答案。

## 可用工具
{tools}

## 工作流程
请严格按照以下格式进行回应，每次只能执行一个步骤：

Thought: 你的思考过程，用于分析问题、拆解任务和规划下一步行动。
Action: 你决定采取的行动，必须是以下格式之一：
- ` + "`{tool_name}[{tool_input}]`" + ` - 调用指定工具
- ` + "`Finish[最终答案]`" + ` - 当你有足够信息给出最终答案时

## 重要提醒
1. 每次回应必须包含Thought和Action两部分
2. 工具调用的格式必须严格遵循：工具名[参数]
3. 只有当你确信有足够信息回答问题时，才使用Finish
4. 如果工具返回的信息不够，继续使用其他工具或相同工具的不同参数

## 当前任务
**Question:** {question}

## 执行历史
{history}

现在开始你的推理和行动：
`

// 正则表达式（预编译为包级变量）
var (
	reAction      = regexp.MustCompile(`(?s)Action:\s*(.*?)$`)
	reActionName  = regexp.MustCompile(`(\w+)\[`)
	reActionInput = regexp.MustCompile(`\w+\[(.*)\]`)
)

// MyReActAgent 重写的 ReAct Agent - 推理与行动结合的智能体
type MyReActAgent struct {
	Name           string
	LLM            *llmclient.HelloAgentsLLM
	ToolRegistry   *common.ToolRegistry
	SystemPrompt   string
	MaxSteps       int
	History        []string
	PromptTemplate string
}

// NewMyReActAgent 创建新的 ReAct Agent
func NewMyReActAgent(
	name string,
	llm *llmclient.HelloAgentsLLM,
	toolRegistry *common.ToolRegistry,
	systemPrompt string,
	maxSteps int,
	customPrompt string,
) *MyReActAgent {
	if maxSteps == 0 {
		maxSteps = 5
	}

	promptTemplate := MY_REACT_PROMPT
	if customPrompt != "" {
		promptTemplate = customPrompt
	}

	fmt.Printf("✅ %s 初始化完成，最大步数: %d\n", name, maxSteps)

	return &MyReActAgent{
		Name:           name,
		LLM:            llm,
		ToolRegistry:   toolRegistry,
		SystemPrompt:   systemPrompt,
		MaxSteps:       maxSteps,
		PromptTemplate: promptTemplate,
	}
}

// Run 运行 ReAct Agent
func (a *MyReActAgent) Run(inputText string) string {
	a.History = nil
	currentStep := 0

	fmt.Printf("\n🤖 %s 开始处理问题: %s\n", a.Name, inputText)

	for currentStep < a.MaxSteps {
		currentStep++
		fmt.Printf("\n--- 第 %d 步 ---\n", currentStep)

		// 1. 构建提示词
		toolsDesc := a.ToolRegistry.GetToolsDescription()
		historyStr := strings.Join(a.History, "\n")
		prompt := a.PromptTemplate
		prompt = strings.ReplaceAll(prompt, "{tools}", toolsDesc)
		prompt = strings.ReplaceAll(prompt, "{question}", inputText)
		prompt = strings.ReplaceAll(prompt, "{history}", historyStr)

		// 2. 调用LLM
		messages := []llmclient.Message{{Role: "user", Content: prompt}}
		responseText := a.LLM.Think(messages, 0)

		// 3. 解析输出
		_, action := a.parseOutput(responseText)

		// 4. 检查完成条件
		if action != "" && strings.HasPrefix(action, "Finish") {
			finalAnswer := a.parseActionInput(action)
			fmt.Printf("🎉 最终答案: %s\n", finalAnswer)
			return finalAnswer
		}

		// 5. 执行工具调用
		if action != "" {
			toolName, toolInput := a.parseAction(action)
			var observation string
			if toolName != "" {
				observation = a.ToolRegistry.ExecuteTool(toolName, toolInput)
			} else {
				observation = "无效的Action格式，请检查。"
			}
			fmt.Printf("👀 观察: %s\n", observation)
			a.History = append(a.History, fmt.Sprintf("Action: %s", action))
			a.History = append(a.History, fmt.Sprintf("Observation: %s", observation))
		}
	}

	// 达到最大步数
	finalAnswer := "抱歉，我无法在限定步数内完成这个任务。"
	fmt.Println("已达到最大步数，流程终止。")
	return finalAnswer
}

// parseOutput 解析 LLM 输出，提取 Thought 和 Action
func (a *MyReActAgent) parseOutput(text string) (string, string) {
	var action string

	actionMatch := reAction.FindStringSubmatch(text)
	if actionMatch != nil {
		action = strings.TrimSpace(actionMatch[1])
	}

	return "", action
}

// parseAction 解析 Action 中的工具名和输入
func (a *MyReActAgent) parseAction(actionText string) (string, string) {
	match := reActionName.FindStringSubmatch(actionText)
	if match == nil {
		return "", ""
	}
	toolName := match[1]

	inputMatch := reActionInput.FindStringSubmatch(actionText)
	if inputMatch == nil {
		return "", ""
	}
	return toolName, inputMatch[1]
}

// parseActionInput 从 Action 中提取输入（用于 Finish）
func (a *MyReActAgent) parseActionInput(actionText string) string {
	match := reActionInput.FindStringSubmatch(actionText)
	if match != nil {
		return match[1]
	}
	return ""
}

func main() {
	_ = godotenv.Load()

	llm, err := llmclient.NewHelloAgentsLLM("", "", "", 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 创建工具注册表
	registry := common.NewToolRegistry()

	// 注册计算器工具
	calculator := common.NewCalculator()
	registry.RegisterFunction(
		"my_calculator",
		"简单的数学计算工具，支持基本运算(+,-,*,/)和sqrt函数",
		calculator.Run,
	)

	// 创建 ReAct Agent
	agent := NewMyReActAgent("MyReActAgent", llm, registry, "", 5, "")

	// 运行Agent
	question := "请帮我计算 (sqrt(16) + 3) * 2 的结果"
	agent.Run(question)
}
