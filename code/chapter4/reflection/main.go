package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"hello-agents-go/chapter4/llmclient"
)

// --- 模块 1: 记忆模块 ---

// Record 一条记忆记录
type Record struct {
	Type    string
	Content string
}

// Memory 一个简单的短期记忆模块，用于存储智能体的行动与反思轨迹。
type Memory struct {
	Records []Record
}

// NewMemory 创建新的记忆模块
func NewMemory() *Memory {
	return &Memory{}
}

// AddRecord 向记忆中添加一条新记录。
// 参数:
// - recordType: 记录的类型 ('execution' 或 'reflection')
// - content: 记录的具体内容 (例如，生成的代码或反思的反馈)
func (m *Memory) AddRecord(recordType, content string) {
	m.Records = append(m.Records, Record{Type: recordType, Content: content})
	fmt.Printf("📝 记忆已更新，新增一条 '%s' 记录。\n", recordType)
}

// GetTrajectory 将所有记忆记录格式化为一个连贯的字符串文本，用于构建提示词。
func (m *Memory) GetTrajectory() string {
	var parts []string
	for _, record := range m.Records {
		if record.Type == "execution" {
			parts = append(parts, fmt.Sprintf("--- 上一轮尝试 (代码) ---\n%s\n", record.Content))
		} else if record.Type == "reflection" {
			parts = append(parts, fmt.Sprintf("--- 评审员反馈 ---\n%s\n", record.Content))
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

// GetLastExecution 获取最近一次的执行结果 (例如，最新生成的代码)。
func (m *Memory) GetLastExecution() string {
	for i := len(m.Records) - 1; i >= 0; i-- {
		if m.Records[i].Type == "execution" {
			return m.Records[i].Content
		}
	}
	return ""
}

// --- 模块 2: Reflection 智能体 ---

// 1. 初始执行提示词
const INITIAL_PROMPT_TEMPLATE = `你是一位资深的Python程序员。请根据以下要求，编写一个Python函数。
你的代码必须包含完整的函数签名、文档字符串，并遵循PEP 8编码规范。

要求: {task}

请直接输出代码，不要包含任何额外的解释。
`

// 2. 反思提示词
const REFLECT_PROMPT_TEMPLATE = `你是一位极其严格的代码评审专家和资深算法工程师，对代码的性能有极致的要求。
你的任务是审查以下Python代码，并专注于找出其在**算法效率**上的主要瓶颈。

# 原始任务:
{task}

# 待审查的代码:
` + "```python" + `
{code}
` + "```" + `

请分析该代码的时间复杂度，并思考是否存在一种**算法上更优**的解决方案来显著提升性能。
如果存在，请清晰地指出当前算法的不足，并提出具体的、可行的改进算法建议（例如，使用筛法替代试除法）。
如果代码在算法层面已经达到最优，才能回答"无需改进"。

请直接输出你的反馈，不要包含任何额外的解释。
`

// 3. 优化提示词
const REFINE_PROMPT_TEMPLATE = `你是一位资深的Python程序员。你正在根据一位代码评审专家的反馈来优化你的代码。

# 原始任务:
{task}

# 你上一轮尝试的代码:
{last_code_attempt}

# 评审员的反馈:
{feedback}

请根据评审员的反馈，生成一个优化后的新版本代码。
你的代码必须包含完整的函数签名、文档字符串，并遵循PEP 8编码规范。
请直接输出优化后的代码，不要包含任何额外的解释。
`

// ReflectionAgent Reflection 智能体
type ReflectionAgent struct {
	llmClient      *llmclient.HelloAgentsLLM
	memory         *Memory
	maxIterations  int
}

// NewReflectionAgent 创建新的 Reflection 智能体
func NewReflectionAgent(llmClient *llmclient.HelloAgentsLLM, maxIterations int) *ReflectionAgent {
	return &ReflectionAgent{
		llmClient:     llmClient,
		memory:        NewMemory(),
		maxIterations: maxIterations,
	}
}

// Run 运行 Reflection 智能体
func (a *ReflectionAgent) Run(task string) string {
	fmt.Printf("\n--- 开始处理任务 ---\n任务: %s\n", task)

	// --- 1. 初始执行 ---
	fmt.Println("\n--- 正在进行初始尝试 ---")
	initialPrompt := strings.ReplaceAll(INITIAL_PROMPT_TEMPLATE, "{task}", task)
	initialCode := a.getLLMResponse(initialPrompt)
	a.memory.AddRecord("execution", initialCode)

	// --- 2. 迭代循环：反思与优化 ---
	for i := 0; i < a.maxIterations; i++ {
		fmt.Printf("\n--- 第 %d/%d 轮迭代 ---\n", i+1, a.maxIterations)

		// a. 反思
		fmt.Println("\n-> 正在进行反思...")
		lastCode := a.memory.GetLastExecution()
		reflectPrompt := REFLECT_PROMPT_TEMPLATE
		reflectPrompt = strings.ReplaceAll(reflectPrompt, "{task}", task)
		reflectPrompt = strings.ReplaceAll(reflectPrompt, "{code}", lastCode)
		feedback := a.getLLMResponse(reflectPrompt)
		a.memory.AddRecord("reflection", feedback)

		// b. 检查是否需要停止
		if strings.Contains(feedback, "无需改进") || strings.Contains(strings.ToLower(feedback), "no need for improvement") {
			fmt.Println("\n✅ 反思认为代码已无需改进，任务完成。")
			break
		}

		// c. 优化
		fmt.Println("\n-> 正在进行优化...")
		refinePrompt := REFINE_PROMPT_TEMPLATE
		refinePrompt = strings.ReplaceAll(refinePrompt, "{task}", task)
		refinePrompt = strings.ReplaceAll(refinePrompt, "{last_code_attempt}", lastCode)
		refinePrompt = strings.ReplaceAll(refinePrompt, "{feedback}", feedback)
		refinedCode := a.getLLMResponse(refinePrompt)
		a.memory.AddRecord("execution", refinedCode)
	}

	finalCode := a.memory.GetLastExecution()
	fmt.Printf("\n--- 任务完成 ---\n最终生成的代码:\n%s\n", finalCode)
	return finalCode
}

// getLLMResponse 一个辅助方法，用于调用LLM并获取完整的流式响应。
func (a *ReflectionAgent) getLLMResponse(prompt string) string {
	messages := []llmclient.Message{{Role: "user", Content: prompt}}
	responseText := a.llmClient.Think(messages, 0)
	return responseText
}

func main() {
	_ = godotenv.Load()

	llm, err := llmclient.NewHelloAgentsLLM("", "", "", 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 初始化 Reflection 智能体，设置最多迭代2轮
	agent := NewReflectionAgent(llm, 2)

	// 定义任务并运行智能体
	task := "编写一个Python函数，找出1到n之间所有的素数 (prime numbers)。"
	agent.Run(task)
}
