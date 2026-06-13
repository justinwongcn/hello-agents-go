// ContextBuilder 与 Agent 集成示例
//
// 展示如何将 ContextBuilder 集成到 Agent 中，实现：
// 1. 上下文感知的 Agent
// 2. 自动构建优化的上下文
// 3. 记忆管理与上下文构建的协同

package main

import (
	"fmt"
	"strings"
	"time"
)

// ContextConfig 上下文配置
type ContextConfig struct {
	MaxTokens         int
	ReserveRatio      float64
	MinRelevance      float64
	EnableCompression bool
}

// Message 消息
type Message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	MemoryTool any
	RAGTool    any
	Config     *ContextConfig
}

// NewContextBuilder 创建 ContextBuilder
func NewContextBuilder(memoryTool, ragTool any, config *ContextConfig) *ContextBuilder {
	return &ContextBuilder{
		MemoryTool: memoryTool,
		RAGTool:    ragTool,
		Config:     config,
	}
}

// Build 构建上下文
func (cb *ContextBuilder) Build(userQuery string, conversationHistory []Message, systemInstructions string) string {
	var sb strings.Builder
	sb.WriteString("=== System Instructions ===\n")
	sb.WriteString(systemInstructions)
	sb.WriteString("\n\n=== Conversation History ===\n")
	if len(conversationHistory) == 0 {
		sb.WriteString("(无历史记录)\n")
	} else {
		for _, msg := range conversationHistory {
			sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Role, msg.Content))
		}
	}
	sb.WriteString("\n=== User Query ===\n")
	sb.WriteString(userQuery)
	sb.WriteString("\n")
	return sb.String()
}

// HelloAgentsLLM 模拟 LLM
type HelloAgentsLLM struct{}

// Invoke 调用 LLM
func (llm *HelloAgentsLLM) Invoke(messages []map[string]string) string {
	// 模拟 LLM 回答
	for _, msg := range messages {
		if msg["role"] == "user" {
			content := msg["content"]
			if strings.Contains(content, "内存") {
				return "优化Pandas内存占用的方法: 1) 使用适当的数据类型 2) 使用chunksize分块读取 3) 使用category类型处理分类数据"
			}
			if strings.Contains(content, "代码示例") {
				return "```python\ndf = pd.read_csv('large.csv', dtype={'col1': 'int32', 'col2': 'category'})\n```"
			}
		}
	}
	return "请问有什么可以帮助您的?"
}

// ContextAwareAgent 具有上下文感知能力的 Agent
type ContextAwareAgent struct {
	Name              string
	LLM               *HelloAgentsLLM
	SystemPrompt      string
	ContextBuilder    *ContextBuilder
	ConversationHistory []Message
}

// NewContextAwareAgent 创建 ContextAwareAgent
func NewContextAwareAgent(name string, llm *HelloAgentsLLM, systemPrompt string) *ContextAwareAgent {
	config := &ContextConfig{MaxTokens: 4000}
	builder := NewContextBuilder(nil, nil, config)

	return &ContextAwareAgent{
		Name:              name,
		LLM:               llm,
		SystemPrompt:      systemPrompt,
		ContextBuilder:    builder,
		ConversationHistory: []Message{},
	}
}

// Run 运行 Agent,自动构建优化的上下文
func (a *ContextAwareAgent) Run(userInput string) string {
	// 1. 使用 ContextBuilder 构建优化的上下文
	optimizedContext := a.ContextBuilder.Build(
		userInput,
		a.ConversationHistory,
		a.SystemPrompt,
	)

	// 2. 使用优化后的上下文调用 LLM
	messages := []map[string]string{
		{"role": "system", "content": optimizedContext},
		{"role": "user", "content": userInput},
	}
	response := a.LLM.Invoke(messages)

	// 3. 更新对话历史
	now := time.Now()
	a.ConversationHistory = append(a.ConversationHistory, Message{Content: userInput, Role: "user", Timestamp: now})
	a.ConversationHistory = append(a.ConversationHistory, Message{Content: response, Role: "assistant", Timestamp: now})

	// 4. 将重要交互记录到记忆系统
	// a.memory_tool.Run({
	//     "action": "add",
	//     "content": fmt.Sprintf("Q: %s\nA: %s...", userInput, response[:200]),  // 摘要
	//     "memory_type": "episodic",
	//     "importance": 0.6
	// })

	return response
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("ContextBuilder 与 Agent 集成示例")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// 配置 LLM
	llm := &HelloAgentsLLM{}

	// 使用示例
	agent := NewContextAwareAgent(
		"数据分析顾问",
		llm,
		"你是一位资深的Python数据工程顾问。",
	)

	// 进行对话
	response := agent.Run("如何优化Pandas的内存占用?")
	fmt.Printf("助手回答:\n%s\n\n", response)

	// 继续对话
	response = agent.Run("能给出具体的代码示例吗?")
	fmt.Printf("助手回答:\n%s\n\n", response)

	fmt.Println(strings.Repeat("=", 80))
}
