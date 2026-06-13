// 代码示例 08: Agent工具集成
// 展示如何在HelloAgents框架中集成MemoryTool和RAGTool

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
		return "📋 记忆摘要: 工作记忆(3条), 情景记忆(2条), 语义记忆(1条)"
	case "stats":
		return "📊 统计: 总计6条记忆"
	default:
		return "操作完成"
	}
}

// RAGTool 模拟RAGTool
type RAGTool struct {
	KnowledgeBasePath string
	RAGNamespace      string
}

// NewRAGTool 创建RAGTool
func NewRAGTool(kbPath, namespace string) *RAGTool {
	return &RAGTool{KnowledgeBasePath: kbPath, RAGNamespace: namespace}
}

// Run 执行RAG操作
func (r *RAGTool) Run(params map[string]any) string {
	action, _ := params["action"].(string)
	switch action {
	case "add_text":
		docID, _ := params["document_id"].(string)
		return fmt.Sprintf("✅ 文本 '%s' 已添加到知识库", docID)
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关文档片段", query)
	case "ask":
		question, _ := params["question"].(string)
		return fmt.Sprintf("💡 问答 '%s': 基于知识库的答案", question)
	case "stats":
		return "📊 统计: 文档数5, 分块数28, 向量维度384"
	default:
		return "操作完成"
	}
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools map[string]any
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]any)}
}

// RegisterTool 注册工具
func (tr *ToolRegistry) RegisterTool(name string, tool any) {
	tr.tools[name] = tool
}

// GetTool 获取工具
func (tr *ToolRegistry) GetTool(name string) any {
	return tr.tools[name]
}

// ListTools 列出工具
func (tr *ToolRegistry) ListTools() []string {
	var names []string
	for name := range tr.tools {
		names = append(names, name)
	}
	return names
}

// AgentIntegrationDemo Agent工具集成演示类
type AgentIntegrationDemo struct {
	MemoryTool   *MemoryTool
	RAGTool      *RAGTool
	ToolRegistry *ToolRegistry
	AgentName    string
	AgentPrompt  string
}

// NewAgentIntegrationDemo 创建演示实例
func NewAgentIntegrationDemo() *AgentIntegrationDemo {
	demo := &AgentIntegrationDemo{}
	demo.setupAgent()
	return demo
}

// setupAgent 设置Agent和工具
func (d *AgentIntegrationDemo) setupAgent() {
	fmt.Println("🤖 Agent工具集成设置")
	fmt.Println(strings.Repeat("=", 50))

	// 初始化工具
	fmt.Println("1. 初始化工具...")
	d.MemoryTool = NewMemoryTool(
		"agent_integration_user",
		[]string{"working", "episodic", "semantic", "perceptual"},
	)

	d.RAGTool = NewRAGTool("./agent_integration_kb", "agent_demo")

	fmt.Println("✅ MemoryTool和RAGTool初始化完成")

	// 创建Agent
	fmt.Println("\n2. 创建Agent...")
	d.AgentName = "智能学习助手"
	d.AgentPrompt = "集成记忆和RAG功能的智能助手"

	fmt.Println("✅ Agent创建完成")

	// 注册工具
	fmt.Println("\n3. 注册工具...")
	d.ToolRegistry = NewToolRegistry()
	d.ToolRegistry.RegisterTool("memory", d.MemoryTool)
	d.ToolRegistry.RegisterTool("rag", d.RAGTool)

	fmt.Println("✅ 工具注册完成")

	// 显示Agent状态
	fmt.Printf("\n📊 Agent状态:\n")
	fmt.Printf("  名称: %s\n", d.AgentName)
	fmt.Printf("  描述: %s\n", d.AgentPrompt)
	fmt.Printf("  可用工具: %v\n", d.ToolRegistry.ListTools())
}

// DemonstrateToolRegistryPattern 演示工具注册模式
func (d *AgentIntegrationDemo) DemonstrateToolRegistryPattern() {
	fmt.Println("\n🔧 工具注册模式演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("工具注册模式特点:")
	fmt.Println("• 🔌 统一的工具接口")
	fmt.Println("• 📋 集中的工具管理")
	fmt.Println("• 🔄 动态工具加载")
	fmt.Println("• 🎯 工具能力发现")

	// 演示工具注册过程
	fmt.Printf("\n🔧 工具注册详情:\n")

	toolNames := d.ToolRegistry.ListTools()
	for _, toolName := range toolNames {
		fmt.Printf("\n工具: %s\n", toolName)

		// 显示工具的主要功能
		if toolName == "memory" {
			fmt.Printf("  类型: MemoryTool\n")
			fmt.Printf("  描述: 记忆管理工具\n")
			fmt.Printf("  主要功能: 记忆管理、搜索、整合、遗忘\n")
			fmt.Printf("  记忆类型: %v\n", d.MemoryTool.MemoryTypes)
		} else if toolName == "rag" {
			fmt.Printf("  类型: RAGTool\n")
			fmt.Printf("  描述: RAG检索增强生成工具\n")
			fmt.Printf("  主要功能: 文档处理、智能问答、知识检索\n")
			fmt.Printf("  命名空间: %s\n", d.RAGTool.RAGNamespace)
		}
	}

	// 演示工具发现机制
	fmt.Printf("\n🔍 工具能力发现:\n")
	availableTools := d.ToolRegistry.ListTools()
	fmt.Printf("可用工具列表: %v\n", availableTools)

	// 演示工具获取
	memoryTool := d.ToolRegistry.GetTool("memory")
	ragTool := d.ToolRegistry.GetTool("rag")

	fmt.Printf("\n✅ 工具获取成功:\n")
	fmt.Printf("  Memory工具: %T\n", memoryTool)
	fmt.Printf("  RAG工具: %T\n", ragTool)
}

// DemonstrateUnifiedInterface 演示统一接口模式
func (d *AgentIntegrationDemo) DemonstrateUnifiedInterface() {
	fmt.Println("\n🔗 统一接口模式演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("统一接口优势:")
	fmt.Println("• 🎯 一致的调用方式")
	fmt.Println("• 📝 标准化的参数传递")
	fmt.Println("• 🛡️ 统一的错误处理")
	fmt.Println("• 🔄 简化的工具切换")

	// 演示统一的run接口
	fmt.Printf("\n🔗 统一run接口演示:\n")

	// Memory工具操作
	fmt.Printf("\n1. Memory工具操作:\n")
	type memoryOperation struct {
		Action string
		Params map[string]any
	}

	memoryOperations := []memoryOperation{
		{"add", map[string]any{
			"content":     "学习了Agent工具集成模式",
			"memory_type": "episodic",
			"importance":  0.8,
			"topic":       "agent_integration",
		}},
		{"search", map[string]any{
			"query": "Agent集成",
			"limit": 2,
		}},
		{"stats", map[string]any{}},
	}

	for _, op := range memoryOperations {
		fmt.Printf("  操作: memory.Run('%s', %v)\n", op.Action, op.Params)
		params := map[string]any{"action": op.Action}
		for k, v := range op.Params {
			params[k] = v
		}
		result := d.MemoryTool.Run(params)
		resultStr := fmt.Sprintf("%v", result)
		if len(resultStr) > 100 {
			fmt.Printf("  结果: %s...\n", resultStr[:100])
		} else {
			fmt.Printf("  结果: %s\n", resultStr)
		}
	}

	// RAG工具操作
	fmt.Printf("\n2. RAG工具操作:\n")

	// 先添加一些内容
	d.RAGTool.Run(map[string]any{
		"action":      "add_text",
		"text":        "Agent工具集成是HelloAgents框架的核心特性，允许Agent使用多种工具来完成复杂任务。",
		"document_id": "agent_integration_guide",
	})

	type ragOperation struct {
		Action string
		Params map[string]any
	}

	ragOperations := []ragOperation{
		{"search", map[string]any{"query": "Agent工具集成", "limit": 2}},
		{"ask", map[string]any{"question": "什么是Agent工具集成？", "limit": 2}},
		{"stats", map[string]any{}},
	}

	for _, op := range ragOperations {
		fmt.Printf("  操作: rag.Run('%s', %v)\n", op.Action, op.Params)
		params := map[string]any{"action": op.Action}
		for k, v := range op.Params {
			params[k] = v
		}
		result := d.RAGTool.Run(params)
		resultStr := fmt.Sprintf("%v", result)
		if len(resultStr) > 100 {
			fmt.Printf("  结果: %s...\n", resultStr[:100])
		} else {
			fmt.Printf("  结果: %s\n", resultStr)
		}
	}
}

// DemonstrateCollaborativeWorkflow 演示协同工作流程
func (d *AgentIntegrationDemo) DemonstrateCollaborativeWorkflow() {
	fmt.Println("\n🤝 协同工作流程演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("协同工作场景:")
	fmt.Println("• 📚 学习新知识 → RAG存储 + Memory记录")
	fmt.Println("• 🔍 回顾学习历程 → Memory检索 + RAG补充")
	fmt.Println("• 💡 知识应用 → RAG查询 + Memory更新")
	fmt.Println("• 📊 学习分析 → 两工具统计整合")

	// 场景1：学习新知识
	fmt.Printf("\n📚 场景1：学习新知识\n")

	// 向RAG添加学习资料
	learningContent := `# 设计模式：观察者模式

## 定义
观察者模式定义了对象间的一对多依赖关系，当一个对象的状态发生改变时，所有依赖它的对象都会得到通知并自动更新。

## 结构
- Subject（主题）：维护观察者列表，提供注册和删除观察者的方法
- Observer（观察者）：定义更新接口
- ConcreteSubject（具体主题）：实现主题接口
- ConcreteObserver（具体观察者）：实现观察者接口

## 应用场景
- GUI事件处理
- 模型-视图架构
- 发布-订阅系统
`

	ragResult := d.RAGTool.Run(map[string]any{
		"action":      "add_text",
		"text":        learningContent,
		"document_id": "observer_pattern",
	})
	fmt.Printf("RAG添加结果: %s\n", ragResult)

	// 记录学习活动到记忆系统
	memoryResult := d.MemoryTool.Run(map[string]any{
		"action":        "add",
		"content":       "学习了观察者设计模式的定义、结构和应用场景",
		"memory_type":   "episodic",
		"importance":    0.8,
		"topic":         "design_patterns",
		"pattern_type":  "observer",
	})
	fmt.Printf("Memory记录结果: %s\n", memoryResult)

	// 场景2：回顾学习历程
	fmt.Printf("\n🔍 场景2：回顾学习历程\n")

	// 从记忆系统检索学习历史
	memorySearch := d.MemoryTool.Run(map[string]any{
		"action": "search",
		"query":  "设计模式学习",
		"limit":  3,
	})
	fmt.Printf("学习历史回顾: %s\n", memorySearch)

	// 从RAG获取相关知识补充
	ragSearch := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  "观察者模式",
		"limit":  2,
	})
	fmt.Printf("知识内容补充: %s\n", ragSearch)

	// 场景3：知识应用
	fmt.Printf("\n💡 场景3：知识应用\n")

	// 通过RAG查询应用方法
	applicationQuery := d.RAGTool.Run(map[string]any{
		"action":   "ask",
		"question": "观察者模式适用于什么场景？",
		"limit":    2,
	})
	fmt.Printf("应用场景查询: %s\n", applicationQuery)

	// 记录应用实践到记忆
	applicationMemory := d.MemoryTool.Run(map[string]any{
		"action":             "add",
		"content":            "查询了观察者模式的应用场景，准备在GUI项目中使用",
		"memory_type":        "working",
		"importance":         0.7,
		"application_context": "gui_project",
	})
	fmt.Printf("应用记录: %s\n", applicationMemory)

	// 场景4：学习分析
	fmt.Printf("\n📊 场景4：学习分析\n")

	// 获取记忆系统统计
	memoryStats := d.MemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("记忆统计: %s\n", memoryStats)

	// 获取RAG系统统计
	ragStats := d.RAGTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("知识库统计: %s\n", ragStats)

	// 生成学习摘要
	learningSummary := d.MemoryTool.Run(map[string]any{"action": "summary", "limit": 5})
	fmt.Printf("学习摘要: %s\n", learningSummary)
}

// DemonstrateAgentOrchestration 演示Agent编排能力
func (d *AgentIntegrationDemo) DemonstrateAgentOrchestration() {
	fmt.Println("\n🎭 Agent编排能力演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("Agent编排特点:")
	fmt.Println("• 🧠 智能工具选择")
	fmt.Println("• 🔄 工具链式调用")
	fmt.Println("• 📊 结果整合分析")
	fmt.Println("• 🎯 目标导向执行")

	// 模拟复杂任务的工具编排
	fmt.Printf("\n🎭 复杂任务编排示例:\n")
	fmt.Println("任务: 创建一个关于机器学习的学习计划")

	// 步骤1：从RAG获取机器学习知识结构
	fmt.Printf("\n步骤1: 获取知识结构\n")

	// 添加机器学习知识
	mlContent := `# 机器学习学习路径

## 基础阶段
1. 数学基础：线性代数、概率统计、微积分
2. 编程基础：Python、NumPy、Pandas
3. 机器学习概念：监督学习、无监督学习、强化学习

## 进阶阶段
1. 算法实现：从零实现经典算法
2. 深度学习：神经网络、CNN、RNN、Transformer
3. 实践项目：端到端机器学习项目

## 高级阶段
1. 模型优化：超参数调优、模型压缩
2. 部署运维：模型部署、监控、更新
3. 前沿技术：最新论文、开源项目
`

	d.RAGTool.Run(map[string]any{
		"action":      "add_text",
		"text":        mlContent,
		"document_id": "ml_learning_path",
	})

	knowledgeStructure := d.RAGTool.Run(map[string]any{
		"action":   "ask",
		"question": "机器学习的学习路径是什么？",
		"limit":    3,
	})
	if len(knowledgeStructure) > 200 {
		fmt.Printf("知识结构: %s...\n", knowledgeStructure[:200])
	} else {
		fmt.Printf("知识结构: %s\n", knowledgeStructure)
	}

	// 步骤2：记录学习计划到记忆系统
	fmt.Printf("\n步骤2: 记录学习计划\n")

	planMemory := d.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "制定了机器学习学习计划，包括基础、进阶、高级三个阶段",
		"memory_type": "episodic",
		"importance":  0.9,
		"plan_type":   "learning",
		"subject":     "machine_learning",
	})
	fmt.Printf("计划记录: %s\n", planMemory)

	// 步骤3：检索相关学习经验
	fmt.Printf("\n步骤3: 检索学习经验\n")

	experienceSearch := d.MemoryTool.Run(map[string]any{
		"action": "search",
		"query":  "学习计划 学习经验",
		"limit":  3,
	})
	fmt.Printf("相关经验: %s\n", experienceSearch)

	// 步骤4：整合生成最终建议
	fmt.Printf("\n步骤4: 生成最终建议\n")

	finalAdvice := d.RAGTool.Run(map[string]any{
		"action":   "ask",
		"question": "如何制定有效的机器学习学习计划？",
		"limit":    4,
	})
	if len(finalAdvice) > 300 {
		fmt.Printf("最终建议: %s...\n", finalAdvice[:300])
	} else {
		fmt.Printf("最终建议: %s\n", finalAdvice)
	}

	// 记录编排过程
	orchestrationMemory := d.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "完成了复杂的学习计划制定任务，使用了RAG和Memory的协同编排",
		"memory_type": "working",
		"importance":  0.8,
		"task_type":   "orchestration",
	})
	fmt.Printf("\n编排记录: %s\n", orchestrationMemory)
}

// DemonstratePerformanceAnalysis 演示性能分析
func (d *AgentIntegrationDemo) DemonstratePerformanceAnalysis() {
	fmt.Println("\n📊 性能分析演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("性能分析指标:")
	fmt.Println("• ⏱️ 工具响应时间")
	fmt.Println("• 🔄 工具切换开销")
	fmt.Println("• 💾 内存使用情况")
	fmt.Println("• 🎯 任务完成效率")

	// 性能测试
	fmt.Printf("\n📊 性能测试:\n")

	// 单工具性能测试
	fmt.Printf("\n1. 单工具性能:\n")

	// Memory工具性能
	startTime := time.Now()
	for i := range 5 {
		d.MemoryTool.Run(map[string]any{
			"action":      "add",
			"content":     fmt.Sprintf("性能测试记忆 %d", i+1),
			"memory_type": "working",
			"importance":  0.5,
		})
	}
	memoryTime := time.Since(startTime)
	fmt.Printf("Memory工具 - 5次添加操作: %.3f秒\n", memoryTime.Seconds())

	// RAG工具性能
	startTime = time.Now()
	for i := range 3 {
		d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  fmt.Sprintf("测试查询 %d", i+1),
			"limit":  2,
		})
	}
	ragTime := time.Since(startTime)
	fmt.Printf("RAG工具 - 3次搜索操作: %.3f秒\n", ragTime.Seconds())

	// 协同工作性能测试
	fmt.Printf("\n2. 协同工作性能:\n")

	startTime = time.Now()

	// 模拟协同工作流程
	d.RAGTool.Run(map[string]any{
		"action":      "add_text",
		"text":        "这是一个性能测试文档",
		"document_id": "perf_test",
	})

	d.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     "执行了性能测试",
		"memory_type": "working",
		"importance":  0.6,
	})

	d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  "性能测试",
		"limit":  1,
	})

	d.MemoryTool.Run(map[string]any{
		"action": "search",
		"query":  "性能测试",
		"limit":  1,
	})

	collaborativeTime := time.Since(startTime)
	fmt.Printf("协同工作流程: %.3f秒\n", collaborativeTime.Seconds())

	// 性能分析总结
	fmt.Printf("\n📈 性能分析总结:\n")
	fmt.Printf("Memory工具平均响应: %.3f秒/操作\n", memoryTime.Seconds()/5)
	fmt.Printf("RAG工具平均响应: %.3f秒/操作\n", ragTime.Seconds()/3)
	fmt.Printf("协同工作效率: %.3f秒/流程\n", collaborativeTime.Seconds())

	// 获取最终统计
	finalMemoryStats := d.MemoryTool.Run(map[string]any{"action": "stats"})
	finalRAGStats := d.RAGTool.Run(map[string]any{"action": "stats"})

	fmt.Printf("\n📊 最终系统状态:\n")
	fmt.Printf("Memory系统: %s\n", finalMemoryStats)
	fmt.Printf("RAG系统: %s\n", finalRAGStats)
}

func main() {
	fmt.Println("🤖 Agent工具集成演示")
	fmt.Println("展示如何在HelloAgents框架中集成MemoryTool和RAGTool")
	fmt.Println(strings.Repeat("=", 70))

	demo := NewAgentIntegrationDemo()

	// 1. 工具注册模式演示
	demo.DemonstrateToolRegistryPattern()

	// 2. 统一接口模式演示
	demo.DemonstrateUnifiedInterface()

	// 3. 协同工作流程演示
	demo.DemonstrateCollaborativeWorkflow()

	// 4. Agent编排能力演示
	demo.DemonstrateAgentOrchestration()

	// 5. 性能分析演示
	demo.DemonstratePerformanceAnalysis()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎉 Agent工具集成演示完成！")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n✨ Agent集成核心特性:")
	fmt.Println("1. 🔧 工具注册模式 - 统一的工具管理和发现")
	fmt.Println("2. 🔗 统一接口设计 - 一致的工具调用方式")
	fmt.Println("3. 🤝 协同工作流程 - 工具间的智能协作")
	fmt.Println("4. 🎭 智能编排能力 - 复杂任务的自动分解")
	fmt.Println("5. 📊 性能监控分析 - 全面的性能评估")

	fmt.Println("\n🎯 设计优势:")
	fmt.Println("• 模块化 - 工具独立开发，灵活组合")
	fmt.Println("• 可扩展 - 支持动态添加新工具")
	fmt.Println("• 高内聚 - 每个工具专注特定功能")
	fmt.Println("• 低耦合 - 工具间依赖关系最小")

	fmt.Println("\n💡 应用价值:")
	fmt.Println("• 智能助手 - 构建多功能智能助手")
	fmt.Println("• 知识管理 - 企业级知识管理系统")
	fmt.Println("• 学习平台 - 个性化学习支持系统")
	fmt.Println("• 决策支持 - 基于知识和经验的决策")
}
