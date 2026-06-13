// 代码示例 03: WorkingMemory实现详解
// 展示工作记忆的混合检索策略和TTL机制

package main

import (
	"fmt"
	"strings"
	"time"
)

// MemoryTool 模拟MemoryTool
type MemoryTool struct {
	UserID      string
	MemoryTypes []string
}

// NewMemoryTool 创建MemoryTool
func NewMemoryTool(userID string, memoryTypes []string) *MemoryTool {
	return &MemoryTool{UserID: userID, MemoryTypes: memoryTypes}
}

// Run 执行记忆操作
func (m *MemoryTool) Run(params map[string]any) string {
	action, _ := params["action"].(string)
	switch action {
	case "add":
		return fmt.Sprintf("✅ 已添加到 %s", params["memory_type"])
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关记忆", query)
	case "stats":
		return "📊 工作记忆: 15条, 容量: 50, TTL: 60分钟"
	case "forget":
		return "🗑️ 已清理低重要性记忆"
	default:
		return "操作完成"
	}
}

// WorkingMemoryDemo 工作记忆演示类
type WorkingMemoryDemo struct {
	MemoryTool *MemoryTool
}

// NewWorkingMemoryDemo 创建演示实例
func NewWorkingMemoryDemo() *WorkingMemoryDemo {
	return &WorkingMemoryDemo{
		MemoryTool: NewMemoryTool("working_memory_demo", []string{"working"}),
	}
}

// DemonstrateCapacityManagement 演示容量管理和TTL机制
func (d *WorkingMemoryDemo) DemonstrateCapacityManagement() {
	fmt.Println("🧠 工作记忆容量管理演示")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("工作记忆特点:")
	fmt.Println("• 容量有限（默认50条）")
	fmt.Println("• TTL机制（默认60分钟）")
	fmt.Println("• 自动清理过期记忆")
	fmt.Println("• 优先级管理（重要性排序）")

	// 添加多条记忆来演示容量管理
	fmt.Printf("\n📝 添加测试记忆...\n")
	for i := range 10 {
		importance := 0.3 + (float64(i) * 0.07) // 递增重要性
		d.MemoryTool.Run(map[string]any{
			"action":      "add",
			"content":     fmt.Sprintf("工作记忆测试项目 %d - 重要性 %.2f", i+1, importance),
			"memory_type": "working",
			"importance":  importance,
			"test_id":     i + 1,
			"category":    "capacity_test",
		})
	}

	// 查看当前状态
	stats := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("当前状态: %s\n", stats)

	// 演示重要性排序
	fmt.Printf("\n🔍 按重要性搜索:\n")
	result := d.MemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "测试项目",
		"memory_type": "working",
		"limit":       5,
	})
	fmt.Println(result)
}

// DemonstrateMixedRetrievalStrategy 演示混合检索策略
func (d *WorkingMemoryDemo) DemonstrateMixedRetrievalStrategy() {
	fmt.Println("\n🔍 混合检索策略演示")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("混合检索策略包括:")
	fmt.Println("• TF-IDF向量化语义检索")
	fmt.Println("• 关键词匹配检索")
	fmt.Println("• 时间衰减因子")
	fmt.Println("• 重要性权重调整")

	// 添加不同类型的记忆用于检索测试
	type testMemory struct {
		Content    string
		Importance float64
		Metadata   map[string]any
	}

	testMemories := []testMemory{
		{
			Content:    "Python是一种高级编程语言，语法简洁清晰",
			Importance: 0.8,
			Metadata:   map[string]any{"topic": "programming", "language": "python"},
		},
		{
			Content:    "机器学习是人工智能的重要分支，包括监督学习和无监督学习",
			Importance: 0.9,
			Metadata:   map[string]any{"topic": "ai", "domain": "machine_learning"},
		},
		{
			Content:    "数据结构包括数组、链表、栈、队列等基本结构",
			Importance: 0.7,
			Metadata:   map[string]any{"topic": "computer_science", "category": "data_structures"},
		},
		{
			Content:    "算法复杂度分析使用大O记号来描述时间和空间复杂度",
			Importance: 0.8,
			Metadata:   map[string]any{"topic": "algorithms", "analysis": "complexity"},
		},
	}

	fmt.Printf("\n📝 添加测试记忆...\n")
	for _, memory := range testMemories {
		params := map[string]any{
			"action":      "add",
			"content":     memory.Content,
			"memory_type": "working",
			"importance":  memory.Importance,
		}
		for k, v := range memory.Metadata {
			params[k] = v
		}
		d.MemoryTool.Run(params)
	}

	// 测试不同类型的检索
	type searchTest struct {
		Query       string
		Description string
	}

	searchTests := []searchTest{
		{"Python编程", "测试语义匹配"},
		{"学习", "测试关键词匹配"},
		{"复杂度", "测试部分匹配"},
		{"人工智能机器学习", "测试多词匹配"},
	}

	fmt.Printf("\n🔍 混合检索测试:\n")
	for _, test := range searchTests {
		fmt.Printf("\n查询: '%s' (%s)\n", test.Query, test.Description)
		result := d.MemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       test.Query,
			"memory_type": "working",
			"limit":       2,
		})
		fmt.Printf("结果: %s\n", result)
	}
}

// DemonstrateTimeDecayMechanism 演示时间衰减机制
func (d *WorkingMemoryDemo) DemonstrateTimeDecayMechanism() {
	fmt.Println("\n⏰ 时间衰减机制演示")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("时间衰减机制:")
	fmt.Println("• 新记忆权重更高")
	fmt.Println("• 旧记忆权重衰减")
	fmt.Println("• 模拟人类记忆特点")
	fmt.Println("• 平衡新旧信息重要性")

	// 添加不同时间的记忆（模拟）
	type timeTestMemory struct {
		Content     string
		Importance  float64
		AgeCategory string
	}

	timeTestMemories := []timeTestMemory{
		{"最新的重要信息 - 刚刚学习的概念", 0.7, "newest"},
		{"较新的信息 - 昨天学习的内容", 0.7, "recent"},
		{"较旧的信息 - 上周学习的内容", 0.7, "older"},
		{"最旧的信息 - 很久以前的内容", 0.7, "oldest"},
	}

	fmt.Printf("\n📝 添加不同时期的记忆...\n")
	for _, memory := range timeTestMemories {
		d.MemoryTool.Run(map[string]any{
			"action":             "add",
			"content":            memory.Content,
			"memory_type":        "working",
			"importance":         memory.Importance,
			"age_category":       memory.AgeCategory,
			"timestamp_category": memory.AgeCategory,
		})
	}

	// 搜索测试时间衰减效果
	fmt.Printf("\n🔍 时间衰减效果测试:\n")
	result := d.MemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "学习的内容",
		"memory_type": "working",
		"limit":       4,
	})
	fmt.Println("搜索结果（注意时间因素对排序的影响）:")
	fmt.Println(result)
}

// DemonstrateAutomaticCleanup 演示自动清理机制
func (d *WorkingMemoryDemo) DemonstrateAutomaticCleanup() {
	fmt.Println("\n🧹 自动清理机制演示")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("自动清理机制:")
	fmt.Println("• 过期记忆自动清理")
	fmt.Println("• 容量超限时清理低优先级记忆")
	fmt.Println("• 保持系统性能和响应速度")
	fmt.Println("• 模拟工作记忆的有限容量")

	// 获取清理前的状态
	statsBefore := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n清理前状态: %s\n", statsBefore)

	// 添加一些低重要性的记忆
	fmt.Printf("\n📝 添加低重要性记忆...\n")
	for i := range 5 {
		d.MemoryTool.Run(map[string]any{
			"action":       "add",
			"content":      fmt.Sprintf("低重要性临时记忆 %d", i+1),
			"memory_type":  "working",
			"importance":   0.1 + float64(i)*0.05,
			"temporary":    true,
			"cleanup_test": true,
		})
	}

	// 触发基于重要性的清理
	fmt.Printf("\n🧹 执行基于重要性的清理...\n")
	cleanupResult := d.MemoryTool.Run(map[string]any{
		"action":    "forget",
		"strategy":  "importance_based",
		"threshold": 0.3,
	})
	fmt.Printf("清理结果: %s\n", cleanupResult)

	// 获取清理后的状态
	statsAfter := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n清理后状态: %s\n", statsAfter)
}

// DemonstratePerformanceCharacteristics 演示性能特征
func (d *WorkingMemoryDemo) DemonstratePerformanceCharacteristics() {
	fmt.Println("\n⚡ 性能特征演示")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("工作记忆性能特点:")
	fmt.Println("• 纯内存存储，访问速度极快")
	fmt.Println("• 无需磁盘I/O，响应时间短")
	fmt.Println("• 适合频繁访问的临时数据")
	fmt.Println("• 系统重启后数据丢失（符合设计）")

	// 性能测试
	fmt.Printf("\n⏱️ 性能测试:\n")

	// 批量添加测试
	startTime := time.Now()
	for i := range 20 {
		d.MemoryTool.Run(map[string]any{
			"action":           "add",
			"content":          fmt.Sprintf("性能测试记忆 %d", i+1),
			"memory_type":      "working",
			"importance":       0.5,
			"performance_test": true,
		})
	}
	addTime := time.Since(startTime)
	fmt.Printf("批量添加20条记忆耗时: %.3f秒\n", addTime.Seconds())

	// 批量搜索测试
	startTime = time.Now()
	for range 10 {
		d.MemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       "性能测试",
			"memory_type": "working",
			"limit":       3,
		})
	}
	searchTime := time.Since(startTime)
	fmt.Printf("批量搜索10次耗时: %.3f秒\n", searchTime.Seconds())

	// 获取最终统计
	finalStats := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n📊 最终统计: %s\n", finalStats)
}

func main() {
	fmt.Println("🧠 WorkingMemory实现详解")
	fmt.Println("展示工作记忆的核心特性和实现机制")
	fmt.Println(strings.Repeat("=", 60))

	demo := NewWorkingMemoryDemo()

	// 1. 容量管理演示
	demo.DemonstrateCapacityManagement()

	// 2. 混合检索策略演示
	demo.DemonstrateMixedRetrievalStrategy()

	// 3. 时间衰减机制演示
	demo.DemonstrateTimeDecayMechanism()

	// 4. 自动清理机制演示
	demo.DemonstrateAutomaticCleanup()

	// 5. 性能特征演示
	demo.DemonstratePerformanceCharacteristics()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 WorkingMemory实现演示完成！")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n✨ 工作记忆核心特性:")
	fmt.Println("1. 🧠 有限容量 - 模拟人类工作记忆限制")
	fmt.Println("2. ⚡ 高速访问 - 纯内存存储，响应迅速")
	fmt.Println("3. 🔍 混合检索 - 语义+关键词+时间+重要性")
	fmt.Println("4. ⏰ 时间衰减 - 新信息优先，旧信息衰减")
	fmt.Println("5. 🧹 自动清理 - TTL机制+优先级管理")

	fmt.Println("\n🎯 设计理念:")
	fmt.Println("• 临时性 - 存储当前会话的临时信息")
	fmt.Println("• 高效性 - 快速访问和处理能力")
	fmt.Println("• 智能性 - 自动管理和优化策略")
	fmt.Println("• 仿生性 - 模拟人类工作记忆特点")
}
