// 代码示例 01: MemoryTool基础操作
// 展示MemoryTool的核心execute方法和基本操作

package main

import (
	"fmt"
	"strings"
)

// MemoryTool 模拟MemoryTool结构
type MemoryTool struct {
	UserID      string
	MemoryTypes []string
}

// NewMemoryTool 初始化MemoryTool
func NewMemoryTool(userID string, memoryTypes []string) *MemoryTool {
	return &MemoryTool{
		UserID:      userID,
		MemoryTypes: memoryTypes,
	}
}

// Run 执行记忆操作
func (m *MemoryTool) Run(params map[string]any) string {
	action, _ := params["action"].(string)
	switch action {
	case "add":
		return fmt.Sprintf("✅ 已添加记忆到 %s 类型", params["memory_type"])
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s' 的结果: 找到相关记忆条目", query)
	case "summary":
		return "📋 记忆摘要: 工作记忆(2条), 情景记忆(1条), 语义记忆(1条), 感知记忆(1条)"
	case "stats":
		return "📊 统计: 总计5条记忆, 工作记忆2条, 情景记忆1条, 语义记忆1条, 感知记忆1条"
	case "forget":
		return "🗑️ 已遗忘低重要性记忆"
	case "consolidate":
		return "🔄 记忆整合完成: working → episodic"
	default:
		return "未知操作"
	}
}

// memoryToolExecuteDemo MemoryTool execute方法演示
func memoryToolExecuteDemo() *MemoryTool {
	fmt.Println("🧠 MemoryTool基础操作演示")
	fmt.Println(strings.Repeat("=", 50))

	// 初始化MemoryTool
	memoryTool := NewMemoryTool(
		"demo_user",
		[]string{"working", "episodic", "semantic", "perceptual"},
	)

	fmt.Println("✅ MemoryTool初始化完成")
	fmt.Println("📋 支持的操作: add, search, summary, stats, update, remove, forget, consolidate, clear_all")

	return memoryTool
}

// addMemoryDemo 添加记忆演示 - 模拟人类记忆编码过程
func addMemoryDemo(memoryTool *MemoryTool) {
	fmt.Println("\n📝 添加记忆演示")
	fmt.Println(strings.Repeat("-", 30))

	// 添加工作记忆
	result := memoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "正在学习HelloAgents框架的记忆系统",
		"memory_type": "working",
		"importance":  0.7,
		"task_type":   "learning",
	})
	fmt.Printf("工作记忆: %s\n", result)

	// 添加情景记忆
	result = memoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "2024年开始深入研究AI Agent技术",
		"memory_type": "episodic",
		"importance":  0.8,
		"event_type":  "milestone",
		"location":    "研发中心",
	})
	fmt.Printf("情景记忆: %s\n", result)

	// 添加语义记忆
	result = memoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "记忆系统包括工作记忆、情景记忆、语义记忆和感知记忆四种类型",
		"memory_type": "semantic",
		"importance":  0.9,
		"concept":     "memory_types",
		"domain":      "cognitive_science",
	})
	fmt.Printf("语义记忆: %s\n", result)

	// 添加感知记忆
	result = memoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "查看了记忆系统的架构图和实现代码",
		"memory_type": "perceptual",
		"importance":  0.6,
		"modality":    "document",
		"source":      "technical_documentation",
	})
	fmt.Printf("感知记忆: %s\n", result)
}

// searchMemoryDemo 搜索记忆演示 - 实现语义理解的检索
func searchMemoryDemo(memoryTool *MemoryTool) {
	fmt.Println("\n🔍 搜索记忆演示")
	fmt.Println(strings.Repeat("-", 30))

	// 基础搜索
	fmt.Println("基础搜索 - '记忆系统':")
	result := memoryTool.Run(map[string]any{"action": "search", "query": "记忆系统", "limit": 3})
	fmt.Println(result)

	// 按类型搜索
	fmt.Println("\n按类型搜索 - 语义记忆中的'记忆':")
	result = memoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "记忆",
		"memory_type": "semantic",
		"limit":       2,
	})
	fmt.Println(result)

	// 设置重要性阈值
	fmt.Println("\n高重要性记忆搜索:")
	result = memoryTool.Run(map[string]any{
		"action":         "search",
		"query":          "AI Agent",
		"min_importance": 0.7,
		"limit":          3,
	})
	fmt.Println(result)
}

// memorySummaryDemo 记忆摘要演示 - 提供系统全貌
func memorySummaryDemo(memoryTool *MemoryTool) {
	fmt.Println("\n📋 记忆摘要演示")
	fmt.Println(strings.Repeat("-", 30))

	// 获取记忆摘要
	result := memoryTool.Run(map[string]any{"action": "summary", "limit": 5})
	fmt.Println("记忆摘要:")
	fmt.Println(result)

	// 获取统计信息
	fmt.Println("\n📊 统计信息:")
	result = memoryTool.Run(map[string]any{"action": "stats"})
	fmt.Println(result)
}

// memoryManagementDemo 记忆管理演示 - 遗忘和整合
func memoryManagementDemo(memoryTool *MemoryTool) {
	fmt.Println("\n⚙️ 记忆管理演示")
	fmt.Println(strings.Repeat("-", 30))

	// 添加一个低重要性记忆用于遗忘测试
	memoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "这是一个临时的测试记忆，重要性很低",
		"memory_type": "working",
		"importance":  0.1,
	})

	// 基于重要性的遗忘
	fmt.Println("基于重要性的遗忘 (阈值=0.2):")
	result := memoryTool.Run(map[string]any{
		"action":    "forget",
		"strategy":  "importance_based",
		"threshold": 0.2,
	})
	fmt.Println(result)

	// 记忆整合 - 将重要的工作记忆转为情景记忆
	fmt.Println("\n记忆整合 (working → episodic):")
	result = memoryTool.Run(map[string]any{
		"action":               "consolidate",
		"from_type":            "working",
		"to_type":              "episodic",
		"importance_threshold": 0.6,
	})
	fmt.Println(result)
}

func main() {
	fmt.Println("🚀 MemoryTool基础操作完整演示")
	fmt.Println("展示记忆系统的核心功能和操作方法")
	fmt.Println(strings.Repeat("=", 60))

	// 1. 初始化MemoryTool
	memoryTool := memoryToolExecuteDemo()

	// 2. 添加记忆演示
	addMemoryDemo(memoryTool)

	// 3. 搜索记忆演示
	searchMemoryDemo(memoryTool)

	// 4. 记忆摘要演示
	memorySummaryDemo(memoryTool)

	// 5. 记忆管理演示
	memoryManagementDemo(memoryTool)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 MemoryTool基础操作演示完成！")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n✨ 演示的核心功能:")
	fmt.Println("1. 🧠 四种记忆类型的添加和管理")
	fmt.Println("2. 🔍 智能语义搜索和过滤")
	fmt.Println("3. 📋 记忆摘要和统计分析")
	fmt.Println("4. ⚙️ 记忆整合和选择性遗忘")

	fmt.Println("\n🎯 设计特点:")
	fmt.Println("• 统一的execute接口，操作简洁一致")
	fmt.Println("• 丰富的元数据支持，便于分类和检索")
	fmt.Println("• 智能的重要性评估和时间衰减机制")
	fmt.Println("• 模拟人类认知的记忆管理策略")
}
