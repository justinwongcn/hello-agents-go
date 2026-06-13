// 代码示例 06: 记忆整合机制演示
// 展示从短期记忆到长期记忆的智能转化过程

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
	case "summary":
		return "📋 记忆摘要: 工作记忆(8条), 情景记忆(2条), 语义记忆(1条)"
	case "stats":
		return "📊 统计: 总计11条记忆, 工作记忆8条, 情景记忆2条, 语义记忆1条"
	case "consolidate":
		fromType, _ := params["from_type"].(string)
		toType, _ := params["to_type"].(string)
		return fmt.Sprintf("🔄 整合完成: %s → %s", fromType, toType)
	default:
		return "操作完成"
	}
}

// MemoryConsolidationDemo 记忆整合演示类
type MemoryConsolidationDemo struct {
	MemoryTool *MemoryTool
}

// NewMemoryConsolidationDemo 创建演示实例
func NewMemoryConsolidationDemo() *MemoryConsolidationDemo {
	return &MemoryConsolidationDemo{
		MemoryTool: NewMemoryTool("consolidation_demo_user", []string{"working", "episodic", "semantic", "perceptual"}),
	}
}

// SetupInitialMemories 设置初始记忆数据
func (d *MemoryConsolidationDemo) SetupInitialMemories() {
	fmt.Println("📝 设置初始记忆数据")
	fmt.Println(strings.Repeat("=", 50))

	// 添加不同重要性的工作记忆
	type workingMemory struct {
		Content    string
		Importance float64
		Metadata   map[string]any
	}

	workingMemories := []workingMemory{
		{"学习了Transformer架构的基本原理", 0.9, map[string]any{"topic": "deep_learning", "session": "study_session_1"}},
		{"完成了Python代码调试任务", 0.8, map[string]any{"topic": "programming", "task_type": "debugging"}},
		{"参加了团队会议讨论项目进展", 0.7, map[string]any{"topic": "teamwork", "meeting_type": "progress_review"}},
		{"查看了今天的天气预报", 0.3, map[string]any{"topic": "daily_life", "category": "routine"}},
		{"阅读了关于注意力机制的论文", 0.85, map[string]any{"topic": "research", "paper_type": "technical"}},
		{"喝了一杯咖啡", 0.2, map[string]any{"topic": "daily_life", "category": "routine"}},
		{"解决了一个复杂的算法问题", 0.9, map[string]any{"topic": "problem_solving", "difficulty": "high"}},
		{"整理了桌面文件", 0.4, map[string]any{"topic": "organization", "category": "maintenance"}},
	}

	fmt.Println("添加工作记忆:")
	for i, memory := range workingMemories {
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

		contentPreview := memory.Content
		if len(contentPreview) > 40 {
			contentPreview = contentPreview[:40]
		}
		fmt.Printf("  %d. %s... (重要性: %.1f)\n", i+1, contentPreview, memory.Importance)
	}

	fmt.Printf("\n✅ 已添加 %d 条工作记忆\n", len(workingMemories))

	// 显示当前状态
	stats := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n📊 当前记忆统计:\n%s\n", stats)
}

// DemonstrateConsolidationCriteria 演示整合标准和筛选过程
func (d *MemoryConsolidationDemo) DemonstrateConsolidationCriteria() {
	fmt.Println("\n🎯 记忆整合标准演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("整合标准:")
	fmt.Println("• 重要性阈值筛选")
	fmt.Println("• 按重要性排序")
	fmt.Println("• 类型转换处理")
	fmt.Println("• 元数据更新")

	// 获取当前工作记忆摘要
	fmt.Println("\n📋 整合前的工作记忆状态:")
	summary := d.MemoryTool.Run(map[string]any{"action": "summary", "limit": 10})
	fmt.Println(summary)

	// 测试不同阈值的整合效果
	thresholds := []float64{0.5, 0.7, 0.8}

	for _, threshold := range thresholds {
		fmt.Printf("\n🔍 测试重要性阈值 %.1f:\n", threshold)
		fmt.Printf("  阈值 %.1f 下符合整合条件的记忆:\n", threshold)
		fmt.Printf("  • 重要性 >= %.1f 的记忆将被整合\n", threshold)
		fmt.Println("  • 整合后类型: working → episodic")
		fmt.Println("  • 重要性提升: importance × 1.1")
	}
}

// DemonstrateConsolidationProcess 演示实际的整合过程
func (d *MemoryConsolidationDemo) DemonstrateConsolidationProcess() {
	fmt.Println("\n🔄 记忆整合过程演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("整合过程步骤:")
	fmt.Println("1. 筛选符合条件的记忆")
	fmt.Println("2. 按重要性排序")
	fmt.Println("3. 创建新的记忆项")
	fmt.Println("4. 更新类型和元数据")
	fmt.Println("5. 添加整合标记")

	// 执行不同阈值的整合
	type consolidationTest struct {
		Threshold   float64
		Description string
	}

	consolidationTests := []consolidationTest{
		{0.6, "低阈值整合 - 整合更多记忆"},
		{0.8, "高阈值整合 - 只整合最重要的记忆"},
	}

	for _, ct := range consolidationTests {
		fmt.Printf("\n🔄 %s (阈值: %.1f):\n", ct.Description, ct.Threshold)

		// 获取整合前状态
		statsBefore := d.MemoryTool.Run(map[string]any{"action": "stats"})
		fmt.Printf("整合前状态: %s\n", statsBefore)

		// 执行整合
		startTime := time.Now()
		consolidationResult := d.MemoryTool.Run(map[string]any{
			"action":               "consolidate",
			"from_type":            "working",
			"to_type":              "episodic",
			"importance_threshold": ct.Threshold,
		})
		consolidationTime := time.Since(startTime)

		fmt.Printf("整合结果: %s\n", consolidationResult)
		fmt.Printf("整合耗时: %.3f秒\n", consolidationTime.Seconds())

		// 获取整合后状态
		statsAfter := d.MemoryTool.Run(map[string]any{"action": "stats"})
		fmt.Printf("整合后状态: %s\n", statsAfter)

		// 查看整合后的情景记忆
		fmt.Printf("\n📚 整合后的情景记忆:\n")
		episodicSearch := d.MemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       "",
			"memory_type": "episodic",
			"limit":       5,
		})
		fmt.Println(episodicSearch)
	}
}

// DemonstrateConsolidationMetadata 演示整合过程中的元数据处理
func (d *MemoryConsolidationDemo) DemonstrateConsolidationMetadata() {
	fmt.Println("\n📋 整合元数据处理演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("元数据处理:")
	fmt.Println("• 保留原始元数据")
	fmt.Println("• 添加整合标记")
	fmt.Println("• 记录整合时间")
	fmt.Println("• 保存原始ID引用")

	// 添加一个特殊的工作记忆用于演示
	specialMemoryResult := d.MemoryTool.Run(map[string]any{
		"action":           "add",
		"content":          "这是一个用于演示整合元数据处理的特殊记忆",
		"memory_type":      "working",
		"importance":       0.85,
		"special_tag":      "metadata_demo",
		"original_context": "demonstration",
		"creation_purpose": "show_consolidation_metadata",
	})

	fmt.Printf("添加特殊记忆: %s\n", specialMemoryResult)

	// 执行整合
	fmt.Printf("\n🔄 执行整合...\n")
	consolidationResult := d.MemoryTool.Run(map[string]any{
		"action":               "consolidate",
		"from_type":            "working",
		"to_type":              "episodic",
		"importance_threshold": 0.8,
	})

	fmt.Printf("整合结果: %s\n", consolidationResult)

	// 搜索整合后的记忆查看元数据
	fmt.Printf("\n🔍 查看整合后的记忆元数据:\n")
	searchResult := d.MemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "特殊记忆",
		"memory_type": "episodic",
		"limit":       1,
	})
	fmt.Println(searchResult)
}

// DemonstrateMultiTypeConsolidation 演示多类型记忆整合
func (d *MemoryConsolidationDemo) DemonstrateMultiTypeConsolidation() {
	fmt.Println("\n🔀 多类型记忆整合演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("多类型整合场景:")
	fmt.Println("• working → episodic (经历记录)")
	fmt.Println("• working → semantic (知识提取)")
	fmt.Println("• episodic → semantic (经验总结)")

	// 添加一些适合不同整合路径的记忆
	type consolidationCandidate struct {
		Content      string
		MemoryType   string
		Importance   float64
		SuitableFor  string
		Metadata     map[string]any
	}

	candidates := []consolidationCandidate{
		{
			Content:     "学习了深度学习中的反向传播算法原理",
			MemoryType:  "working",
			Importance:  0.9,
			SuitableFor: "semantic",
			Metadata:    map[string]any{"learning_type": "concept"},
		},
		{
			Content:     "今天下午参加了AI技术分享会",
			MemoryType:  "working",
			Importance:  0.8,
			SuitableFor: "episodic",
			Metadata:    map[string]any{"event_type": "meeting"},
		},
		{
			Content:     "通过多次实践掌握了Transformer的实现技巧",
			MemoryType:  "episodic",
			Importance:  0.85,
			SuitableFor: "semantic",
			Metadata:    map[string]any{"experience_type": "skill"},
		},
	}

	fmt.Printf("\n📝 添加整合候选记忆:\n")
	for _, memory := range candidates {
		params := map[string]any{
			"action":      "add",
			"content":     memory.Content,
			"memory_type": memory.MemoryType,
			"importance":  memory.Importance,
		}
		for k, v := range memory.Metadata {
			params[k] = v
		}

		d.MemoryTool.Run(params)

		contentPreview := memory.Content
		if len(contentPreview) > 50 {
			contentPreview = contentPreview[:50]
		}
		fmt.Printf("  • %s... → 适合整合为%s\n", contentPreview, memory.SuitableFor)
	}

	// 执行不同类型的整合
	type consolidationPath struct {
		FromType    string
		ToType      string
		Threshold   float64
		Description string
	}

	consolidationPaths := []consolidationPath{
		{"working", "episodic", 0.75, "经历记录整合"},
		{"working", "semantic", 0.85, "知识提取整合"},
		{"episodic", "semantic", 0.8, "经验总结整合"},
	}

	for _, path := range consolidationPaths {
		fmt.Printf("\n🔄 %s (%s → %s):\n", path.Description, path.FromType, path.ToType)

		result := d.MemoryTool.Run(map[string]any{
			"action":               "consolidate",
			"from_type":            path.FromType,
			"to_type":              path.ToType,
			"importance_threshold": path.Threshold,
		})

		fmt.Printf("整合结果: %s\n", result)
	}
}

// DemonstrateConsolidationBenefits 演示记忆整合的益处
func (d *MemoryConsolidationDemo) DemonstrateConsolidationBenefits() {
	fmt.Println("\n✨ 记忆整合益处演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("整合益处:")
	fmt.Println("• 长期保存重要信息")
	fmt.Println("• 释放工作记忆空间")
	fmt.Println("• 形成知识体系")
	fmt.Println("• 提升检索效率")

	// 获取最终的记忆系统状态
	fmt.Printf("\n📊 最终记忆系统状态:\n")
	finalStats := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Println(finalStats)

	// 获取各类型记忆的摘要
	fmt.Printf("\n📋 各类型记忆摘要:\n")

	memoryTypes := []string{"working", "episodic", "semantic"}
	for _, memoryType := range memoryTypes {
		fmt.Printf("\n%s记忆:\n", strings.ToUpper(memoryType))
		typeSummary := d.MemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       "",
			"memory_type": memoryType,
			"limit":       3,
		})
		fmt.Println(typeSummary)
	}

	// 演示整合后的检索效果
	fmt.Printf("\n🔍 整合后的检索效果测试:\n")
	type searchQuery struct {
		Query       string
		Description string
	}

	searchQueries := []searchQuery{
		{"深度学习", "测试跨类型检索"},
		{"学习经历", "测试整合记忆检索"},
		{"重要概念", "测试语义记忆检索"},
	}

	for _, sq := range searchQueries {
		fmt.Printf("\n查询: '%s' (%s)\n", sq.Query, sq.Description)
		result := d.MemoryTool.Run(map[string]any{
			"action": "search",
			"query":  sq.Query,
			"limit":  3,
		})
		fmt.Println(result)
	}
}

func main() {
	fmt.Println("🔄 记忆整合机制演示")
	fmt.Println("展示从短期记忆到长期记忆的智能转化过程")
	fmt.Println(strings.Repeat("=", 60))

	demo := NewMemoryConsolidationDemo()

	// 1. 设置初始记忆数据
	demo.SetupInitialMemories()

	// 2. 演示整合标准
	demo.DemonstrateConsolidationCriteria()

	// 3. 演示整合过程
	demo.DemonstrateConsolidationProcess()

	// 4. 演示元数据处理
	demo.DemonstrateConsolidationMetadata()

	// 5. 演示多类型整合
	demo.DemonstrateMultiTypeConsolidation()

	// 6. 演示整合益处
	demo.DemonstrateConsolidationBenefits()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 记忆整合机制演示完成！")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n✨ 记忆整合核心特性:")
	fmt.Println("1. 🎯 智能筛选 - 基于重要性阈值的自动筛选")
	fmt.Println("2. 🔄 类型转换 - 灵活的记忆类型转换机制")
	fmt.Println("3. 📋 元数据保持 - 完整保留原始上下文信息")
	fmt.Println("4. ⚡ 自动化处理 - 无需人工干预的自动整合")
	fmt.Println("5. 🔀 多路径支持 - 支持多种整合路径")

	fmt.Println("\n🎯 设计理念:")
	fmt.Println("• 仿生性 - 模拟人类大脑的记忆固化过程")
	fmt.Println("• 智能性 - 自动识别和处理重要信息")
	fmt.Println("• 灵活性 - 支持多种整合策略和路径")
	fmt.Println("• 完整性 - 保持记忆的完整性和可追溯性")

	fmt.Println("\n💡 应用价值:")
	fmt.Println("• 知识管理 - 将临时学习转化为长期知识")
	fmt.Println("• 经验积累 - 保存重要的实践经验")
	fmt.Println("• 系统优化 - 释放短期记忆空间")
	fmt.Println("• 智能决策 - 基于历史经验的决策支持")
}
