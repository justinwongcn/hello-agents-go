package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"hello-agents-go/chapter4/llmclient"
)

// --- 1. 规划器 (Planner) 定义 ---

// PLANNER_PROMPT_TEMPLATE 规划器提示词
const PLANNER_PROMPT_TEMPLATE = `你是一个顶级的AI规划专家。你的任务是将用户提出的复杂问题分解成一个由多个简单步骤组成的行动计划。
请确保计划中的每个步骤都是一个独立的、可执行的子任务，并且严格按照逻辑顺序排列。
你的输出必须是一个Python列表，其中每个元素都是一个描述子任务的字符串。

问题: {question}

请严格按照以下格式输出你的计划，` + "```python" + `与` + "```" + `作为前后缀是必要的:
` + "```python" + `
["步骤1", "步骤2", "步骤3", ...]
` + "```" + `
`

// Planner 规划器
type Planner struct {
	llmClient *llmclient.HelloAgentsLLM
}

// NewPlanner 创建新的规划器
func NewPlanner(llmClient *llmclient.HelloAgentsLLM) *Planner {
	return &Planner{llmClient: llmClient}
}

// Plan 生成行动计划
func (p *Planner) Plan(question string) []string {
	prompt := strings.ReplaceAll(PLANNER_PROMPT_TEMPLATE, "{question}", question)
	messages := []llmclient.Message{{Role: "user", Content: prompt}}

	fmt.Println("--- 正在生成计划 ---")
	responseText := p.llmClient.Think(messages, 0)
	fmt.Printf("✅ 计划已生成:\n%s\n", responseText)

	// 解析计划：提取 ```python ... ``` 中的内容
	planStr := ""
	if parts := strings.Split(responseText, "```python"); len(parts) > 1 {
		if endParts := strings.Split(parts[1], "```"); len(endParts) > 0 {
			planStr = strings.TrimSpace(endParts[0])
		}
	}

	if planStr == "" {
		fmt.Println("❌ 解析计划时出错: 未找到计划内容")
		return nil
	}

	// 解析 JSON 数组
	var plan []string
	if err := json.Unmarshal([]byte(planStr), &plan); err != nil {
		fmt.Printf("❌ 解析计划时出错: %v\n", err)
		fmt.Printf("原始响应: %s\n", responseText)
		return nil
	}

	return plan
}

// --- 2. 执行器 (Executor) 定义 ---

// EXECUTOR_PROMPT_TEMPLATE 执行器提示词
const EXECUTOR_PROMPT_TEMPLATE = `你是一位顶级的AI执行专家。你的任务是严格按照给定的计划，一步步地解决问题。
你将收到原始问题、完整的计划、以及到目前为止已经完成的步骤和结果。
请你专注于解决"当前步骤"，并仅输出该步骤的最终答案，不要输出任何额外的解释或对话。

# 原始问题:
{question}

# 完整计划:
{plan}

# 历史步骤与结果:
{history}

# 当前步骤:
{current_step}

请仅输出针对"当前步骤"的回答:
`

// Executor 执行器
type Executor struct {
	llmClient *llmclient.HelloAgentsLLM
}

// NewExecutor 创建新的执行器
func NewExecutor(llmClient *llmclient.HelloAgentsLLM) *Executor {
	return &Executor{llmClient: llmClient}
}

// Execute 执行计划
func (e *Executor) Execute(question string, plan []string) string {
	history := ""
	finalAnswer := ""

	fmt.Println("\n--- 正在执行计划 ---")
	for i, step := range plan {
		fmt.Printf("\n-> 正在执行步骤 %d/%d: %s\n", i+1, len(plan), step)

		planStr := fmt.Sprintf("%v", plan)
		historyDisplay := history
		if historyDisplay == "" {
			historyDisplay = "无"
		}

		prompt := EXECUTOR_PROMPT_TEMPLATE
		prompt = strings.ReplaceAll(prompt, "{question}", question)
		prompt = strings.ReplaceAll(prompt, "{plan}", planStr)
		prompt = strings.ReplaceAll(prompt, "{history}", historyDisplay)
		prompt = strings.ReplaceAll(prompt, "{current_step}", step)

		messages := []llmclient.Message{{Role: "user", Content: prompt}}
		responseText := e.llmClient.Think(messages, 0)

		history += fmt.Sprintf("步骤 %d: %s\n结果: %s\n\n", i+1, step, responseText)
		finalAnswer = responseText
		fmt.Printf("✅ 步骤 %d 已完成，结果: %s\n", i+1, finalAnswer)
	}

	return finalAnswer
}

// --- 3. 智能体 (Agent) 整合 ---

// PlanAndSolveAgent Plan-and-Solve 智能体
type PlanAndSolveAgent struct {
	llmClient *llmclient.HelloAgentsLLM
	planner   *Planner
	executor  *Executor
}

// NewPlanAndSolveAgent 创建新的 Plan-and-Solve 智能体
func NewPlanAndSolveAgent(llmClient *llmclient.HelloAgentsLLM) *PlanAndSolveAgent {
	return &PlanAndSolveAgent{
		llmClient: llmClient,
		planner:   NewPlanner(llmClient),
		executor:  NewExecutor(llmClient),
	}
}

// Run 运行 Plan-and-Solve 智能体
func (a *PlanAndSolveAgent) Run(question string) {
	fmt.Printf("\n--- 开始处理问题 ---\n问题: %s\n", question)
	plan := a.planner.Plan(question)
	if len(plan) == 0 {
		fmt.Println("\n--- 任务终止 --- \n无法生成有效的行动计划。")
		return
	}
	finalAnswer := a.executor.Execute(question, plan)
	fmt.Printf("\n--- 任务完成 ---\n最终答案: %s\n", finalAnswer)
}

// --- 4. 主函数入口 ---

func main() {
	_ = godotenv.Load()

	llm, err := llmclient.NewHelloAgentsLLM("", "", "", 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	agent := NewPlanAndSolveAgent(llm)
	question := "一个水果店周一卖出了15个苹果。周二卖出的苹果数量是周一的两倍。周三卖出的数量比周二少了5个。请问这三天总共卖出了多少个苹果？"
	agent.Run(question)
}
