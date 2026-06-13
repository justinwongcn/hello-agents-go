// 代码示例 02: MemoryTool架构设计
// 展示MemoryTool和MemoryManager的分层架构

package main

import (
	"fmt"
	"strings"
)

// MemoryConfig 记忆配置
type MemoryConfig struct {
	WorkingMemoryCapacity    int
	WorkingMemoryTTLMinutes int
}

// NewMemoryConfig 创建默认配置
func NewMemoryConfig() *MemoryConfig {
	return &MemoryConfig{
		WorkingMemoryCapacity:    50,
		WorkingMemoryTTLMinutes: 60,
	}
}

// MemoryManager 记忆管理器
type MemoryManager struct {
	UserID      string
	Config      *MemoryConfig
	MemoryTypes map[string]string
}

// MemoryTool 记忆工具
type MemoryTool struct {
	MemoryManager *MemoryManager
	MemoryTypes   []string
	MemoryConfig  *MemoryConfig
}

// NewMemoryTool 初始化MemoryTool
func NewMemoryTool(userID string, config *MemoryConfig, memoryTypes []string) *MemoryTool {
	memoryTypeMap := make(map[string]string)
	typeNames := map[string]string{
		"working":    "WorkingMemory",
		"episodic":   "EpisodicMemory",
		"semantic":   "SemanticMemory",
		"perceptual": "PerceptualMemory",
	}
	for _, mt := range memoryTypes {
		if name, ok := typeNames[mt]; ok {
			memoryTypeMap[mt] = name
		}
	}

	return &MemoryTool{
		MemoryManager: &MemoryManager{
			UserID:      userID,
			Config:      config,
			MemoryTypes: memoryTypeMap,
		},
		MemoryTypes:  memoryTypes,
		MemoryConfig: config,
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
		return "📊 统计: 总计5条记忆"
	default:
		return "未知操作"
	}
}

// MemoryToolArchitectureDemo 架构演示类
type MemoryToolArchitectureDemo struct {
	MemoryConfig *MemoryConfig
	MemoryTypes  []string
}

// NewMemoryToolArchitectureDemo 创建演示实例
func NewMemoryToolArchitectureDemo() *MemoryToolArchitectureDemo {
	return &MemoryToolArchitectureDemo{
		MemoryConfig: NewMemoryConfig(),
		MemoryTypes:  []string{"working", "episodic", "semantic", "perceptual"},
	}
}

// DemonstrateMemoryToolInit 演示MemoryTool初始化过程
func (d *MemoryToolArchitectureDemo) DemonstrateMemoryToolInit() *MemoryTool {
	fmt.Println("🏗️ MemoryTool架构设计演示")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("📋 MemoryTool初始化过程:")
	fmt.Println("1. 创建MemoryConfig配置对象")
	fmt.Println("2. 指定启用的记忆类型")
	fmt.Println("3. 初始化MemoryManager管理器")
	fmt.Println("4. 根据配置启用不同记忆模块")

	// 演示MemoryTool的初始化
	memoryTool := NewMemoryTool(
		"architecture_demo_user",
		d.MemoryConfig,
		d.MemoryTypes,
	)

	fmt.Printf("\n✅ MemoryTool初始化完成\n")
	fmt.Printf("👤 用户ID: %s\n", memoryTool.MemoryManager.UserID)
	fmt.Printf("🧠 启用的记忆类型: %v\n", memoryTool.MemoryTypes)
	fmt.Printf("⚙️ 配置对象: MemoryConfig\n")

	return memoryTool
}

// DemonstrateMemoryManagerArchitecture 演示MemoryManager的组合模式架构
func (d *MemoryToolArchitectureDemo) DemonstrateMemoryManagerArchitecture(memoryTool *MemoryTool) {
	fmt.Println("\n🔧 MemoryManager架构设计")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("MemoryManager采用组合模式设计:")
	fmt.Println("- 统一的记忆操作接口")
	fmt.Println("- 独立的记忆类型组件")
	fmt.Println("- 灵活的配置和扩展能力")

	// 获取MemoryManager实例
	memoryManager := memoryTool.MemoryManager

	fmt.Printf("\n📊 MemoryManager状态:\n")
	fmt.Printf("用户ID: %s\n", memoryManager.UserID)
	fmt.Printf("配置类型: MemoryConfig\n")
	fmt.Printf("记忆类型数量: %d\n", len(memoryManager.MemoryTypes))

	// 显示各记忆类型的状态
	fmt.Printf("\n🧠 记忆类型组件:\n")
	for memoryType, memoryInstance := range memoryManager.MemoryTypes {
		fmt.Printf("  • %s: %s\n", memoryType, memoryInstance)
	}
}

// DemonstrateMemoryTypesSpecialization 演示四种记忆类型的专业化特点
func (d *MemoryToolArchitectureDemo) DemonstrateMemoryTypesSpecialization(memoryTool *MemoryTool) {
	fmt.Println("\n🎯 四种记忆类型的专业化设计")
	fmt.Println(strings.Repeat("-", 40))

	type memoryTypeInfo struct {
		Name     string
		Features []string
		Storage  string
		TTL      string
	}

	memoryTypesInfo := map[string]memoryTypeInfo{
		"working": {
			Name:     "工作记忆",
			Features: []string{"容量有限", "访问速度快", "自动清理", "临时存储"},
			Storage:  "纯内存存储",
			TTL:      "60分钟TTL机制",
		},
		"episodic": {
			Name:     "情景记忆",
			Features: []string{"事件序列", "时间序列", "上下文丰富", "会话关联"},
			Storage:  "SQLite + Qdrant混合存储",
			TTL:      "持久化存储",
		},
		"semantic": {
			Name:     "语义记忆",
			Features: []string{"概念知识", "实体关系", "知识图谱", "语义推理"},
			Storage:  "Neo4j + Qdrant混合存储",
			TTL:      "长期存储",
		},
		"perceptual": {
			Name:     "感知记忆",
			Features: []string{"多模态", "跨模态检索", "感知数据", "内容生成"},
			Storage:  "分模态向量存储",
			TTL:      "按重要性管理",
		},
	}

	for _, memoryType := range d.MemoryTypes {
		info := memoryTypesInfo[memoryType]
		fmt.Printf("\n📚 %s (%s):\n", info.Name, memoryType)
		fmt.Printf("   特点: %s\n", strings.Join(info.Features, ", "))
		fmt.Printf("   存储: %s\n", info.Storage)
		fmt.Printf("   生命周期: %s\n", info.TTL)

		// 添加示例记忆来演示特点
		switch memoryType {
		case "working":
			memoryTool.Run(map[string]any{
				"action":        "add",
				"content":       fmt.Sprintf("演示%s的临时存储特性", info.Name),
				"memory_type":   memoryType,
				"importance":    0.6,
				"demo_feature":  "temporary_storage",
			})
		case "episodic":
			memoryTool.Run(map[string]any{
				"action":          "add",
				"content":         fmt.Sprintf("演示%s的事件记录特性", info.Name),
				"memory_type":     memoryType,
				"importance":      0.7,
				"event_type":      "demonstration",
				"session_context": "architecture_demo",
			})
		case "semantic":
			memoryTool.Run(map[string]any{
				"action":      "add",
				"content":     fmt.Sprintf("%s用于存储概念性知识和实体关系", info.Name),
				"memory_type": memoryType,
				"importance":  0.8,
				"concept":     "memory_architecture",
				"domain":      "cognitive_computing",
			})
		case "perceptual":
			memoryTool.Run(map[string]any{
				"action":      "add",
				"content":     fmt.Sprintf("演示%s的多模态数据处理", info.Name),
				"memory_type": memoryType,
				"importance":  0.6,
				"modality":    "text",
				"data_type":   "demonstration",
			})
		}
	}
}

// DemonstrateUnifiedInterface 演示统一接口的设计优势
func (d *MemoryToolArchitectureDemo) DemonstrateUnifiedInterface(memoryTool *MemoryTool) {
	fmt.Println("\n🔗 统一接口设计优势")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("统一的execute方法提供:")
	fmt.Println("• 一致的调用方式")
	fmt.Println("• 灵活的参数传递")
	fmt.Println("• 统一的错误处理")
	fmt.Println("• 简化的用户体验")

	// 演示统一接口的使用
	type operation struct {
		name   string
		params map[string]any
	}

	operations := []operation{
		{"search", map[string]any{"query": "演示", "limit": 2}},
		{"summary", map[string]any{"limit": 3}},
		{"stats", map[string]any{}},
	}

	fmt.Printf("\n🔧 统一接口操作演示:\n")
	for _, op := range operations {
		fmt.Printf("\n操作: %s\n", op.name)
		fmt.Printf("参数: %v\n", op.params)
		params := map[string]any{"action": op.name}
		for k, v := range op.params {
			params[k] = v
		}
		result := memoryTool.Run(params)
		resultStr := fmt.Sprintf("%v", result)
		if len(resultStr) > 100 {
			fmt.Printf("结果: %s...\n", resultStr[:100])
		} else {
			fmt.Printf("结果: %s\n", resultStr)
		}
	}
}

// DemonstrateExtensibility 演示系统的扩展性设计
func (d *MemoryToolArchitectureDemo) DemonstrateExtensibility() {
	fmt.Println("\n🚀 系统扩展性设计")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("扩展性特点:")
	fmt.Println("• 插件化的记忆类型")
	fmt.Println("• 可配置的存储后端")
	fmt.Println("• 灵活的记忆策略")
	fmt.Println("• 模块化的组件设计")

	// 演示自定义配置
	customConfig := NewMemoryConfig()
	customConfig.WorkingMemoryCapacity = 100
	customConfig.WorkingMemoryTTLMinutes = 120

	fmt.Printf("\n⚙️ 自定义配置示例:\n")
	fmt.Printf("工作记忆容量: %d\n", customConfig.WorkingMemoryCapacity)
	fmt.Printf("工作记忆TTL: %d分钟\n", customConfig.WorkingMemoryTTLMinutes)

	// 演示选择性启用记忆类型
	selectiveMemoryTool := NewMemoryTool(
		"selective_user",
		customConfig,
		[]string{"working", "semantic"}, // 只启用部分类型
	)

	fmt.Printf("\n🎯 选择性启用示例:\n")
	fmt.Printf("启用的记忆类型: %v\n", selectiveMemoryTool.MemoryTypes)
	fmt.Println("✅ 系统支持根据需求灵活配置")
}

func main() {
	fmt.Println("🏗️ MemoryTool架构设计完整演示")
	fmt.Println("展示记忆系统的分层架构和设计模式")
	fmt.Println(strings.Repeat("=", 60))

	demo := NewMemoryToolArchitectureDemo()

	// 1. MemoryTool初始化演示
	memoryTool := demo.DemonstrateMemoryToolInit()

	// 2. MemoryManager架构演示
	demo.DemonstrateMemoryManagerArchitecture(memoryTool)

	// 3. 记忆类型专业化演示
	demo.DemonstrateMemoryTypesSpecialization(memoryTool)

	// 4. 统一接口演示
	demo.DemonstrateUnifiedInterface(memoryTool)

	// 5. 扩展性演示
	demo.DemonstrateExtensibility()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 MemoryTool架构演示完成！")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n✨ 架构设计亮点:")
	fmt.Println("1. 🏗️ 分层架构 - 关注点分离，职责清晰")
	fmt.Println("2. 🔧 组合模式 - 灵活组合，独立管理")
	fmt.Println("3. 🎯 专业化设计 - 各记忆类型特点鲜明")
	fmt.Println("4. 🔗 统一接口 - 简化使用，一致体验")
	fmt.Println("5. 🚀 高扩展性 - 插件化设计，灵活配置")

	fmt.Println("\n🎯 设计原则:")
	fmt.Println("• 单一职责原则 - 每个组件专注特定功能")
	fmt.Println("• 开闭原则 - 对扩展开放，对修改封闭")
	fmt.Println("• 依赖倒置原则 - 依赖抽象，不依赖具体")
	fmt.Println("• 组合优于继承 - 灵活组合，避免复杂继承")
}
