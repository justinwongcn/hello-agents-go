package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"hello-agents-go/chapter4/llmclient"
	"hello-agents-go/chapter4/toolexecutor"
)

// REACT_PROMPT_TEMPLATE ReAct 智能体的提示词模板
const REACT_PROMPT_TEMPLATE = `请注意，你是一个有能力调用外部工具的智能助手。

可用工具如下：
{tools}

请严格按照以下格式进行回应：

Thought: 你的思考过程，用于分析问题、拆解任务和规划下一步行动。
Action: 你决定采取的行动，必须是以下格式之一：
- ` + "`{tool_name}[{tool_input}]`" + `：调用一个可用工具。
- ` + "`Finish[最终答案]`" + `：当你认为已经获得最终答案时。
- 当你收集到足够的信息，能够回答用户的最终问题时，你必须在` + "`Action:`" + `字段后使用 ` + "`Finish[最终答案]`" + `来输出最终答案。


现在，请开始解决以下问题：
Question: {question}
History: {history}
`

// 正则表达式（一次性编译）
var (
	reThought = regexp.MustCompile(`(?s)Thought:\s*(.*?)(?:\nAction:|$)`)
	reAction  = regexp.MustCompile(`(?s)Action:\s*(.*?)$`)
	reActionName = regexp.MustCompile(`(\w+)\[`)
	reActionInput = regexp.MustCompile(`\w+\[(.*)\]`)
)

// ReActAgent ReAct 智能体
type ReActAgent struct {
	llmClient    *llmclient.HelloAgentsLLM
	toolExecutor *toolexecutor.ToolExecutor
	maxSteps     int
	history      []string
}

// NewReActAgent 创建新的 ReAct 智能体
func NewReActAgent(llmClient *llmclient.HelloAgentsLLM, toolExecutor *toolexecutor.ToolExecutor, maxSteps int) *ReActAgent {
	return &ReActAgent{
		llmClient:    llmClient,
		toolExecutor: toolExecutor,
		maxSteps:     maxSteps,
	}
}

// Run 运行 ReAct 智能体
func (a *ReActAgent) Run(question string) string {
	a.history = nil
	currentStep := 0

	for currentStep < a.maxSteps {
		currentStep++
		fmt.Printf("\n--- 第 %d 步 ---\n", currentStep)

		toolsDesc := a.toolExecutor.GetAvailableTools()
		historyStr := strings.Join(a.history, "\n")
		prompt := REACT_PROMPT_TEMPLATE
		prompt = strings.ReplaceAll(prompt, "{tools}", toolsDesc)
		prompt = strings.ReplaceAll(prompt, "{question}", question)
		prompt = strings.ReplaceAll(prompt, "{history}", historyStr)

		messages := []llmclient.Message{{Role: "user", Content: prompt}}
		responseText := a.llmClient.Think(messages, 0)
		if responseText == "" {
			fmt.Println("错误：LLM未能返回有效响应。")
			break
		}

		thought, action := a.parseOutput(responseText)
		if thought != "" {
			fmt.Printf("🤔 思考: %s\n", thought)
		}
		if action == "" {
			fmt.Println("警告：未能解析出有效的Action，流程终止。")
			break
		}

		if strings.HasPrefix(action, "Finish") {
			finalAnswer := a.parseActionInput(action)
			fmt.Printf("🎉 最终答案: %s\n", finalAnswer)
			return finalAnswer
		}

		toolName, toolInput := a.parseAction(action)
		if toolName == "" || toolInput == "" {
			a.history = append(a.history, "Observation: 无效的Action格式，请检查。")
			continue
		}

		fmt.Printf("🎬 行动: %s[%s]\n", toolName, toolInput)
		toolFunction := a.toolExecutor.GetTool(toolName)
		var observation string
		if toolFunction != nil {
			observation = toolFunction(toolInput)
		} else {
			observation = fmt.Sprintf("错误：未找到名为 '%s' 的工具。", toolName)
		}

		fmt.Printf("👀 观察: %s\n", observation)
		a.history = append(a.history, fmt.Sprintf("Action: %s", action))
		a.history = append(a.history, fmt.Sprintf("Observation: %s", observation))
	}

	fmt.Println("已达到最大步数，流程终止。")
	return ""
}

// parseOutput 解析 LLM 输出，提取 Thought 和 Action
func (a *ReActAgent) parseOutput(text string) (string, string) {
	var thought, action string

	thoughtMatch := reThought.FindStringSubmatch(text)
	if thoughtMatch != nil {
		thought = strings.TrimSpace(thoughtMatch[1])
	}

	actionMatch := reAction.FindStringSubmatch(text)
	if actionMatch != nil {
		action = strings.TrimSpace(actionMatch[1])
	}

	return thought, action
}

// parseAction 解析 Action 中的工具名和输入
func (a *ReActAgent) parseAction(actionText string) (string, string) {
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
func (a *ReActAgent) parseActionInput(actionText string) string {
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

	toolExecutor := toolexecutor.NewToolExecutor()
	searchDesc := "一个网页搜索引擎。当你需要回答关于时事、事实以及在你的知识库中找不到的信息时，应使用此工具。"
	toolExecutor.RegisterTool("Search", searchDesc, toolexecutor.Search)

	agent := NewReActAgent(llm, toolExecutor, 5)
	question := "华为最新的手机是哪一款？它的主要卖点是什么？"
	agent.Run(question)
}
