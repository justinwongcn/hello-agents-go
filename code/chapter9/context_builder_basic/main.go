// ContextBuilder 基础使用示例
//
// 展示如何使用 ContextBuilder 构建优化的上下文，包括：
// 1. 初始化 ContextBuilder
// 2. 准备对话历史
// 3. 添加记忆
// 4. 构建结构化上下文

package main

import (
	"fmt"
	"strings"
	"time"
)

// ContextConfig 上下文配置
type ContextConfig struct {
	MaxTokens        int
	ReserveRatio     float64
	MinRelevance     float64
	EnableCompression bool
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

// Message 消息
type Message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

// Build 构建上下文
func (cb *ContextBuilder) Build(userQuery string, conversationHistory []Message, systemInstructions string) string {
	var sb strings.Builder

	sb.WriteString("=== System Instructions ===\n")
	sb.WriteString(systemInstructions)
	sb.WriteString("\n\n")

	sb.WriteString("=== Conversation History ===\n")
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

	sb.WriteString(fmt.Sprintf("\n=== Context Stats ===\n"))
	sb.WriteString(fmt.Sprintf("Max Tokens: %d\n", cb.Config.MaxTokens))
	sb.WriteString(fmt.Sprintf("Reserve Ratio: %.1f\n", cb.Config.ReserveRatio))
	sb.WriteString(fmt.Sprintf("Compression: %v\n", cb.Config.EnableCompression))

	return sb.String()
}

// HelloAgentsLLM 模拟 LLM
type HelloAgentsLLM struct{}

// Invoke 调用 LLM
func (llm *HelloAgentsLLM) Invoke(messages []map[string]string) string {
	return "LLM 回答: 建议使用 chunksize 参数分块读取数据，使用 df.astype() 优化数据类型以减少内存占用。"
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("ContextBuilder 基础使用示例")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// 1. 初始化工具（Optional）
	fmt.Println("1. 初始化工具...")
	// memory_tool = MemoryTool(user_id="user123")
	// rag_tool = RAGTool(knowledge_base_path="./knowledge_base")

	// 2. 创建 ContextBuilder
	fmt.Println("2. 创建 ContextBuilder...")
	config := &ContextConfig{
		MaxTokens:         3000,
		ReserveRatio:      0.2,
		MinRelevance:      0, // 最小相关性阈值，0代表所有历史信息会被保留
		EnableCompression: true,
	}

	builder := NewContextBuilder(
		// memory_tool,
		// rag_tool,
		nil, nil,
		config,
	)

	// 3. 准备对话历史
	fmt.Println("3. 准备对话历史...")
	now := time.Now()
	conversationHistory := []Message{
		{Content: "我正在开发一个数据分析工具", Role: "user", Timestamp: now},
		{Content: "很好!数据分析工具通常需要处理大量数据。您计划使用什么技术栈?", Role: "assistant", Timestamp: now},
		{Content: "我打算使用Python和Pandas,已经完成了CSV读取模块", Role: "user", Timestamp: now},
		{Content: "不错的选择!Pandas在数据处理方面非常强大。接下来您可能需要考虑数据清洗和转换。", Role: "assistant", Timestamp: now},
	}

	// 4. 添加一些记忆
	fmt.Println("4. 添加记忆...")
	// memory_tool.Run({
	//     "action": "add",
	//     "content": "用户正在开发数据分析工具,使用Python和Pandas",
	//     "memory_type": "semantic",
	//     "importance": 0.8
	// })

	// memory_tool.Run({
	//     "action": "add",
	//     "content": "已完成CSV读取模块的开发",
	//     "memory_type": "episodic",
	//     "importance": 0.7
	// })

	// 5. 构建上下文
	fmt.Println("5. 构建上下文...\n")
	contextStr := builder.Build(
		"如何优化Pandas的内存占用?",
		conversationHistory,
		"你是一位资深的Python数据工程顾问。你的回答需要:1) 提供具体可行的建议 2) 解释技术原理 3) 给出代码示例",
	)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("构建的上下文 (结构化字符串):")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(contextStr)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// 6. 将上下文字符串转换为消息格式供 LLM 使用
	fmt.Println("6. 将上下文传给 LLM...")
	messages := []map[string]string{
		{"role": "system", "content": contextStr},
		{"role": "user", "content": "请回答"},
	}

	llm := &HelloAgentsLLM{}
	// 注意: 实际使用时需要配置 LLM
	response := llm.Invoke(messages)
	fmt.Printf("LLM 回答: %s\n", response)

	fmt.Println("✅ ContextBuilder 演示完成!")
	fmt.Println("\n提示: ContextBuilder 返回的是结构化的上下文字符串,")
	fmt.Println("      可以直接作为 system message 传给 LLM。")
}
