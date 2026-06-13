// 代码示例 09: 四种记忆类型深度解析
// 详细展示WorkingMemory、EpisodicMemory、SemanticMemory、PerceptualMemory的实现特点

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
		return fmt.Sprintf("📊 %s 统计: 已存储多条记忆", strings.Join(m.MemoryTypes, ", "))
	case "forget":
		return "🗑️ 已清理低重要性记忆"
	case "consolidate":
		fromType, _ := params["from_type"].(string)
		toType, _ := params["to_type"].(string)
		return fmt.Sprintf("🔄 整合完成: %s → %s", fromType, toType)
	default:
		return "操作完成"
	}
}

// MemoryTypesDeepDive 四种记忆类型深度解析演示类
type MemoryTypesDeepDive struct {
	WorkingMemoryTool   *MemoryTool
	EpisodicMemoryTool  *MemoryTool
	SemanticMemoryTool  *MemoryTool
	PerceptualMemoryTool *MemoryTool
}

// NewMemoryTypesDeepDive 创建演示实例
func NewMemoryTypesDeepDive() *MemoryTypesDeepDive {
	fmt.Println("🧠 四种记忆类型深度解析")
	fmt.Println(strings.Repeat("=", 60))

	demo := &MemoryTypesDeepDive{
		WorkingMemoryTool:    NewMemoryTool("working_memory_user", []string{"working"}),
		EpisodicMemoryTool:   NewMemoryTool("episodic_memory_user", []string{"episodic"}),
		SemanticMemoryTool:   NewMemoryTool("semantic_memory_user", []string{"semantic"}),
		PerceptualMemoryTool: NewMemoryTool("perceptual_memory_user", []string{"perceptual"}),
	}

	fmt.Println("✅ 四种记忆系统初始化完成")
	return demo
}

// DemonstrateWorkingMemory 演示工作记忆的特点
func (d *MemoryTypesDeepDive) DemonstrateWorkingMemory() {
	fmt.Println("\n💭 工作记忆 (Working Memory) 深度解析")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 工作记忆特点:")
	fmt.Println("• ⚡ 访问速度极快（纯内存存储）")
	fmt.Println("• 📏 容量有限（默认50条记忆）")
	fmt.Println("• ⏰ 自动过期（TTL机制）")
	fmt.Println("• 🔄 适合临时信息存储")

	// 演示容量限制
	fmt.Printf("\n1. 容量限制演示:\n")
	fmt.Println("添加大量临时记忆，观察容量管理...")

	for i := range 8 {
		content := fmt.Sprintf("临时工作记忆 %d: 当前正在处理任务步骤 %d", i+1, i+1)
		result := d.WorkingMemoryTool.Run(map[string]any{
			"action":      "add",
			"content":     content,
			"memory_type": "working",
			"importance":  0.3 + (float64(i) * 0.1),
			"task_step":   i + 1,
		})
		fmt.Printf("  添加记忆 %d: %s\n", i+1, result)
	}

	// 检查当前状态
	stats := d.WorkingMemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n当前工作记忆状态: %s\n", stats)

	// 演示TTL机制
	fmt.Printf("\n2. TTL（生存时间）机制演示:\n")

	// 模拟不同时间的记忆
	type timeMemory struct {
		Content        string
		MinutesAgo     int
		Importance     float64
	}

	timeMemories := []timeMemory{
		{"刚刚的想法", 0, 0.8},
		{"5分钟前的任务", 5, 0.6},
		{"10分钟前的提醒", 10, 0.4},
		{"很久以前的笔记", 30, 0.2},
	}

	for _, memory := range timeMemories {
		d.WorkingMemoryTool.Run(map[string]any{
			"action":                "add",
			"content":               memory.Content,
			"memory_type":           "working",
			"importance":            memory.Importance,
			"simulated_age_minutes": memory.MinutesAgo,
		})
		fmt.Printf("  添加记忆: %s (模拟 %d 分钟前)\n", memory.Content, memory.MinutesAgo)
	}

	// 演示快速检索
	fmt.Printf("\n3. 快速检索演示:\n")

	searchQueries := []string{"任务", "想法", "提醒"}
	for _, query := range searchQueries {
		startTime := time.Now()
		results := d.WorkingMemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       query,
			"memory_type": "working",
			"limit":       3,
		})
		searchTime := time.Since(startTime)
		fmt.Printf("  查询 '%s': %.4f秒\n", query, searchTime.Seconds())
		if len(results) > 100 {
			fmt.Printf("    结果: %s...\n", results[:100])
		} else {
			fmt.Printf("    结果: %s\n", results)
		}
	}

	// 演示自动清理
	fmt.Printf("\n4. 自动清理机制:\n")

	// 获取清理前的统计
	beforeStats := d.WorkingMemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("清理前: %s\n", beforeStats)

	// 触发清理（通过遗忘低重要性记忆）
	forgetResult := d.WorkingMemoryTool.Run(map[string]any{
		"action":    "forget",
		"strategy":  "importance_based",
		"threshold": 0.4,
	})
	fmt.Printf("清理结果: %s\n", forgetResult)

	// 获取清理后的统计
	afterStats := d.WorkingMemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("清理后: %s\n", afterStats)
}

// DemonstrateEpisodicMemory 演示情景记忆的特点
func (d *MemoryTypesDeepDive) DemonstrateEpisodicMemory() {
	fmt.Println("\n📖 情景记忆 (Episodic Memory) 深度解析")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 情景记忆特点:")
	fmt.Println("• 📅 完整的时间序列记录")
	fmt.Println("• 🎭 丰富的上下文信息")
	fmt.Println("• 🔗 支持记忆链条构建")
	fmt.Println("• 💾 持久化存储")

	// 演示完整事件记录
	fmt.Printf("\n1. 完整事件记录演示:\n")

	// 模拟一个完整的学习会话
	type learningEvent struct {
		Content    string
		Context    string
		Importance float64
		Metadata   map[string]any
	}

	learningSession := []learningEvent{
		{"开始学习Python机器学习", "学习开始", 0.7, map[string]any{"location": "家里书房", "mood": "专注"}},
		{"学习了线性回归的数学原理", "理论学习", 0.8, map[string]any{"chapter": "第3章", "difficulty": "中等"}},
		{"实现了第一个线性回归模型", "实践编程", 0.9, map[string]any{"code_lines": 45, "bugs_fixed": 2}},
		{"完成了课后练习题", "练习巩固", 0.6, map[string]any{"exercises_completed": 5, "accuracy": 0.8}},
		{"总结今天的学习收获", "学习总结", 0.8, map[string]any{"key_concepts": "线性回归, 梯度下降, 损失函数"}},
	}

	for i, event := range learningSession {
		params := map[string]any{
			"action":          "add",
			"content":         event.Content,
			"memory_type":     "episodic",
			"importance":      event.Importance,
			"session_id":      "learning_session_2024",
			"sequence_number": i + 1,
			"context":         event.Context,
		}
		for k, v := range event.Metadata {
			params[k] = v
		}
		result := d.EpisodicMemoryTool.Run(params)
		fmt.Printf("  事件 %d: %s\n", i+1, result)
	}

	// 演示时间序列检索
	fmt.Printf("\n2. 时间序列检索演示:\n")

	timelineSearch := d.EpisodicMemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "学习",
		"memory_type": "episodic",
		"limit":       10,
	})
	fmt.Printf("学习时间线: %s\n", timelineSearch)

	sessionSearch := d.EpisodicMemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "线性回归",
		"memory_type": "episodic",
		"limit":       5,
	})
	fmt.Printf("会话内容: %s\n", sessionSearch)

	// 演示上下文丰富性
	fmt.Printf("\n3. 上下文信息演示:\n")

	// 添加带有丰富上下文的记忆
	contextResult := d.EpisodicMemoryTool.Run(map[string]any{
		"action":              "add",
		"content":             "参加了AI技术分享会",
		"memory_type":         "episodic",
		"importance":          0.9,
		"event_type":          "conference",
		"location":            "北京国际会议中心",
		"speakers":            "张教授, 李博士, 王工程师",
		"topics":              "深度学习, 自然语言处理, 计算机视觉",
		"attendees_count":     200,
		"duration_hours":      6,
		"weather":             "晴朗",
		"transportation":      "地铁",
		"networking_contacts": 3,
		"key_insights":        "Transformer架构的演进, 多模态学习的前景",
		"follow_up_actions":   "阅读推荐论文, 尝试新框架",
		"satisfaction_rating":  9,
	})
	fmt.Printf("丰富上下文记忆: %s\n", contextResult)

	// 演示记忆链条
	fmt.Printf("\n4. 记忆链条构建:\n")

	type chainMemory struct {
		Content    string
		ChainType  string
		ParentType string
	}

	memoryChain := []chainMemory{
		{"看到一篇关于GPT的论文", "trigger", ""},
		{"决定深入研究Transformer架构", "decision", "trigger"},
		{"下载并阅读Attention is All You Need论文", "action", "decision"},
		{"实现了简化版的自注意力机制", "implementation", "action"},
		{"在项目中应用了学到的知识", "application", "implementation"},
	}

	chainMemories := make(map[string]string)
	for _, cm := range memoryChain {
		parentID := ""
		if cm.ParentType != "" {
			parentID = chainMemories[cm.ParentType]
		}

		d.EpisodicMemoryTool.Run(map[string]any{
			"action":        "add",
			"content":       cm.Content,
			"memory_type":   "episodic",
			"importance":    0.7,
			"chain_type":    cm.ChainType,
			"parent_memory": parentID,
			"chain_id":      "gpt_learning_chain",
		})

		chainMemories[cm.ChainType] = fmt.Sprintf("%s_memory", cm.ChainType)
		fmt.Printf("  链条记忆: %s (类型: %s)\n", cm.Content, cm.ChainType)
	}

	// 检索整个链条
	chainSearch := d.EpisodicMemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "GPT Transformer",
		"memory_type": "episodic",
		"limit":       8,
	})
	fmt.Printf("记忆链条检索: %s\n", chainSearch)
}

// DemonstrateSemanticMemory 演示语义记忆的特点
func (d *MemoryTypesDeepDive) DemonstrateSemanticMemory() {
	fmt.Println("\n🧠 语义记忆 (Semantic Memory) 深度解析")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 语义记忆特点:")
	fmt.Println("• 🔗 知识图谱结构化存储")
	fmt.Println("• 🎯 概念和关系的抽象表示")
	fmt.Println("• 🔍 语义相似度检索")
	fmt.Println("• 🧮 支持推理和关联")

	// 演示概念存储
	fmt.Printf("\n1. 概念知识存储演示:\n")

	type concept struct {
		Content     string
		ConceptType string
		Domain      string
		Importance  float64
		Metadata    map[string]any
	}

	concepts := []concept{
		{"机器学习是人工智能的一个分支，通过算法让计算机从数据中学习模式", "definition", "artificial_intelligence", 0.9, map[string]any{"keywords": "机器学习, 人工智能, 算法, 数据, 模式"}},
		{"监督学习使用标记数据训练模型，包括分类和回归两大类任务", "category", "machine_learning", 0.8, map[string]any{"parent_concept": "机器学习", "subcategories": "分类, 回归"}},
		{"梯度下降是一种优化算法，通过迭代更新参数来最小化损失函数", "algorithm", "optimization", 0.8, map[string]any{"mathematical_basis": "微积分", "applications": "神经网络训练, 线性回归"}},
		{"过拟合是指模型在训练数据上表现很好，但在新数据上泛化能力差", "problem", "machine_learning", 0.7, map[string]any{"causes": "模型复杂度过高, 训练数据不足", "solutions": "正则化, 交叉验证, 早停"}},
	}

	for _, c := range concepts {
		params := map[string]any{
			"action":       "add",
			"content":      c.Content,
			"memory_type":  "semantic",
			"importance":   c.Importance,
			"concept_type": c.ConceptType,
			"domain":       c.Domain,
		}
		for k, v := range c.Metadata {
			params[k] = v
		}
		result := d.SemanticMemoryTool.Run(params)
		fmt.Printf("  概念存储: %s - %s\n", c.ConceptType, result)
	}

	// 演示关系推理
	fmt.Printf("\n2. 关系推理演示:\n")

	type relationship struct {
		Content      string
		RelationType string
		Subject      string
		Object       string
	}

	relationships := []relationship{
		{"深度学习是机器学习的子集，使用多层神经网络", "is_subset_of", "深度学习", "机器学习"},
		{"卷积神经网络特别适合处理图像数据", "suitable_for", "卷积神经网络", "图像处理"},
		{"反向传播算法用于训练神经网络", "used_for", "反向传播", "神经网络训练"},
	}

	for _, r := range relationships {
		d.SemanticMemoryTool.Run(map[string]any{
			"action":        "add",
			"content":       r.Content,
			"memory_type":   "semantic",
			"importance":    0.8,
			"relation_type": r.RelationType,
			"subject":       r.Subject,
			"object":        r.Object,
		})
		fmt.Printf("  关系存储: %s - ✅\n", r.RelationType)
	}

	// 演示语义检索
	fmt.Printf("\n3. 语义相似度检索:\n")

	semanticQueries := []string{
		"什么是人工智能？",
		"如何防止模型过拟合？",
		"神经网络的训练方法",
		"图像识别技术",
	}

	for _, query := range semanticQueries {
		startTime := time.Now()
		results := d.SemanticMemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       query,
			"memory_type": "semantic",
			"limit":       3,
		})
		searchTime := time.Since(startTime)
		fmt.Printf("  查询: '%s' (%.4f秒)\n", query, searchTime.Seconds())
		if len(results) > 150 {
			fmt.Printf("    结果: %s...\n", results[:150])
		} else {
			fmt.Printf("    结果: %s\n", results)
		}
	}

	// 演示知识图谱构建
	fmt.Printf("\n4. 知识图谱构建:\n")

	type entity struct {
		Content    string
		EntityType string
		Metadata   map[string]any
	}

	entities := []entity{
		{"TensorFlow是Google开发的深度学习框架", "framework", map[string]any{"developer": "Google", "domain": "deep_learning", "language": "Python", "year": 2015}},
		{"PyTorch是Facebook开发的深度学习框架，以动态图著称", "framework", map[string]any{"developer": "Facebook", "domain": "deep_learning", "feature": "dynamic_graph", "language": "Python"}},
		{"BERT是基于Transformer的预训练语言模型", "model", map[string]any{"architecture": "Transformer", "task": "natural_language_processing", "training_method": "pre_training"}},
	}

	for _, item := range entities {
		params := map[string]any{
			"action":      "add",
			"content":     item.Content,
			"memory_type": "semantic",
			"importance":  0.8,
			"entity_type": item.EntityType,
		}
		for k, v := range item.Metadata {
			params[k] = v
		}
		result := d.SemanticMemoryTool.Run(params)
		fmt.Printf("  实体关系: %s - %s\n", item.EntityType, result)
	}

	semanticStats := d.SemanticMemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n语义记忆统计: %s\n", semanticStats)
}

// DemonstratePerceptualMemory 演示感知记忆的特点
func (d *MemoryTypesDeepDive) DemonstratePerceptualMemory() {
	fmt.Println("\n👁️ 感知记忆 (Perceptual Memory) 深度解析")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 感知记忆特点:")
	fmt.Println("• 🎨 多模态数据支持")
	fmt.Println("• 🔄 跨模态相似性搜索")
	fmt.Println("• 📊 感知数据的语义理解")
	fmt.Println("• 🎯 内容生成和检索")

	// 演示文本感知记忆
	fmt.Printf("\n1. 文本感知记忆:\n")

	type textPerception struct {
		Content  string
		Genre    string
		Metadata map[string]any
	}

	textPerceptions := []textPerception{
		{"这是一段优美的诗歌：春江潮水连海平，海上明月共潮生", "poetry", map[string]any{"modality": "text", "emotion": "peaceful", "language": "chinese", "aesthetic_value": 0.9}},
		{"技术文档：API接口返回JSON格式数据，包含状态码和响应体", "technical", map[string]any{"modality": "text", "complexity": "medium", "language": "chinese", "practical_value": 0.8}},
	}

	for _, perception := range textPerceptions {
		params := map[string]any{
			"action":      "add",
			"content":     perception.Content,
			"memory_type": "perceptual",
			"importance":  0.7,
			"genre":       perception.Genre,
		}
		for k, v := range perception.Metadata {
			params[k] = v
		}
		result := d.PerceptualMemoryTool.Run(params)
		fmt.Printf("  文本感知: %s - %s\n", perception.Genre, result)
	}

	// 演示图像感知记忆（模拟）
	fmt.Printf("\n2. 图像感知记忆（模拟）:\n")

	type imagePerception struct {
		Content  string
		FilePath string
		Metadata map[string]any
	}

	imagePerceptions := []imagePerception{
		{"一张美丽的日落风景照片", "/simulated/sunset.jpg", map[string]any{"modality": "image", "scene_type": "landscape", "colors": "orange, red, purple", "objects": "sun, clouds, horizon", "mood": "serene", "quality": "high"}},
		{"技术架构图展示了微服务系统设计", "/simulated/architecture.png", map[string]any{"modality": "image", "diagram_type": "technical", "components": "API Gateway, Services, Database", "complexity": "high", "purpose": "documentation"}},
	}

	for _, perception := range imagePerceptions {
		params := map[string]any{
			"action":      "add",
			"content":     perception.Content,
			"memory_type": "perceptual",
			"importance":  0.8,
			"file_path":   perception.FilePath,
		}
		for k, v := range perception.Metadata {
			params[k] = v
		}
		result := d.PerceptualMemoryTool.Run(params)
		fmt.Printf("  图像感知: %s - %s\n", perception.Content, result)
	}

	// 演示音频感知记忆（模拟）
	fmt.Printf("\n3. 音频感知记忆（模拟）:\n")

	type audioPerception struct {
		Content  string
		FilePath string
		Metadata map[string]any
	}

	audioPerceptions := []audioPerception{
		{"一段优美的古典音乐演奏", "/simulated/classical.mp3", map[string]any{"modality": "audio", "genre": "classical", "instruments": "piano, violin, cello", "tempo": "andante", "emotion": "elegant", "duration_seconds": 240}},
		{"技术会议的录音，讨论AI发展趋势", "/simulated/conference.wav", map[string]any{"modality": "audio", "content_type": "speech", "topic": "artificial_intelligence", "speakers": 3, "language": "chinese", "duration_seconds": 1800}},
	}

	for _, perception := range audioPerceptions {
		params := map[string]any{
			"action":      "add",
			"content":     perception.Content,
			"memory_type": "perceptual",
			"importance":  0.7,
			"file_path":   perception.FilePath,
		}
		for k, v := range perception.Metadata {
			params[k] = v
		}
		result := d.PerceptualMemoryTool.Run(params)
		fmt.Printf("  音频感知: %s - %s\n", perception.Content, result)
	}

	// 演示跨模态检索
	fmt.Printf("\n4. 跨模态检索演示:\n")

	type crossModalQuery struct {
		Query       string
		Description string
	}

	crossModalQueries := []crossModalQuery{
		{"美丽的风景", "寻找视觉美感相关内容"},
		{"技术文档", "查找技术相关的多模态内容"},
		{"音乐和艺术", "检索艺术相关的感知记忆"},
		{"会议和讨论", "查找交流相关的内容"},
	}

	for _, q := range crossModalQueries {
		results := d.PerceptualMemoryTool.Run(map[string]any{
			"action":      "search",
			"query":       q.Query,
			"memory_type": "perceptual",
			"limit":       3,
		})
		fmt.Printf("  跨模态查询: '%s' (%s)\n", q.Query, q.Description)
		if len(results) > 120 {
			fmt.Printf("    结果: %s...\n", results[:120])
		} else {
			fmt.Printf("    结果: %s\n", results)
		}
	}

	// 演示感知特征分析
	fmt.Printf("\n5. 感知特征分析:\n")

	perceptualStats := d.PerceptualMemoryTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("感知记忆统计: %s\n", perceptualStats)

	modalityAnalysis := d.PerceptualMemoryTool.Run(map[string]any{
		"action":      "search",
		"query":       "模态分析",
		"memory_type": "perceptual",
		"limit":       10,
	})
	fmt.Printf("模态分布分析: %s\n", modalityAnalysis)
}

// DemonstrateMemoryInteractions 演示四种记忆类型的交互
func (d *MemoryTypesDeepDive) DemonstrateMemoryInteractions() {
	fmt.Println("\n🔄 四种记忆类型交互演示")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 记忆交互模式:")
	fmt.Println("• 🔄 工作记忆 → 情景记忆（重要事件固化）")
	fmt.Println("• 📚 情景记忆 → 语义记忆（经验抽象化）")
	fmt.Println("• 👁️ 感知记忆 → 其他记忆（多模态信息整合）")
	fmt.Println("• 🧠 语义记忆 → 工作记忆（知识激活）")

	// 模拟一个完整的学习过程
	fmt.Printf("\n完整学习过程模拟:\n")

	// 1. 感知阶段：接收多模态信息
	fmt.Printf("\n1. 感知阶段 - 接收信息:\n")

	perceptualInput := d.PerceptualMemoryTool.Run(map[string]any{
		"action":           "add",
		"content":          "观看了一个关于深度学习的视频教程",
		"memory_type":      "perceptual",
		"importance":       0.8,
		"modality":         "video",
		"topic":            "deep_learning",
		"duration_minutes": 45,
		"quality":          "high",
	})
	fmt.Printf("感知记忆: %s\n", perceptualInput)

	// 2. 工作记忆阶段：临时处理和思考
	fmt.Printf("\n2. 工作记忆阶段 - 临时处理:\n")

	workingThoughts := []string{
		"理解了卷积神经网络的基本原理",
		"需要记住反向传播的计算步骤",
		"想到了之前学过的线性代数知识",
		"计划实现一个简单的CNN模型",
	}

	for _, thought := range workingThoughts {
		result := d.WorkingMemoryTool.Run(map[string]any{
			"action":            "add",
			"content":           thought,
			"memory_type":       "working",
			"importance":        0.6,
			"processing_stage":  "active_thinking",
		})
		if len(thought) > 30 {
			fmt.Printf("  工作记忆: %s... - %s\n", thought[:30], result)
		} else {
			fmt.Printf("  工作记忆: %s - %s\n", thought, result)
		}
	}

	// 3. 情景记忆阶段：记录完整学习事件
	fmt.Printf("\n3. 情景记忆阶段 - 事件记录:\n")

	episodicEvent := d.EpisodicMemoryTool.Run(map[string]any{
		"action":           "add",
		"content":          "完成了深度学习视频教程的学习，理解了CNN的核心概念",
		"memory_type":      "episodic",
		"importance":       0.9,
		"event_type":       "learning_session",
		"duration_minutes": 45,
		"location":         "家里",
		"learning_outcome": "理解CNN原理",
		"next_action":      "实践编程",
	})
	fmt.Printf("情景记忆: %s\n", episodicEvent)

	// 4. 语义记忆阶段：抽象知识存储
	fmt.Printf("\n4. 语义记忆阶段 - 知识抽象:\n")

	type semanticKnowledge struct {
		Content  string
		Concept  string
		Metadata map[string]any
	}

	semanticKnowledgeItems := []semanticKnowledge{
		{"卷积神经网络通过卷积层提取图像特征，适合计算机视觉任务", "CNN", map[string]any{"domain": "deep_learning", "application": "computer_vision"}},
		{"反向传播算法通过链式法则计算梯度，用于更新网络参数", "backpropagation", map[string]any{"domain": "optimization", "mathematical_basis": "chain_rule"}},
	}

	for _, knowledge := range semanticKnowledgeItems {
		params := map[string]any{
			"action":      "add",
			"content":     knowledge.Content,
			"memory_type": "semantic",
			"importance":  0.8,
			"concept":     knowledge.Concept,
		}
		for k, v := range knowledge.Metadata {
			params[k] = v
		}
		result := d.SemanticMemoryTool.Run(params)
		fmt.Printf("  语义记忆: %s - %s\n", knowledge.Concept, result)
	}

	// 5. 记忆整合演示
	fmt.Printf("\n5. 记忆整合演示:\n")

	consolidationResult := d.WorkingMemoryTool.Run(map[string]any{
		"action":               "consolidate",
		"from_type":            "working",
		"to_type":              "episodic",
		"importance_threshold": 0.6,
	})
	fmt.Printf("工作记忆整合: %s\n", consolidationResult)

	// 跨记忆类型检索
	fmt.Printf("\n6. 跨记忆类型检索:\n")

	query := "深度学习CNN"

	// 在所有记忆类型中搜索
	type memorySearch struct {
		Name string
		Tool *MemoryTool
	}

	memoryTools := []memorySearch{
		{"工作记忆", d.WorkingMemoryTool},
		{"情景记忆", d.EpisodicMemoryTool},
		{"语义记忆", d.SemanticMemoryTool},
		{"感知记忆", d.PerceptualMemoryTool},
	}

	for _, ms := range memoryTools {
		results := ms.Tool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  2,
		})
		if len(results) > 80 {
			fmt.Printf("  %s检索: %s...\n", ms.Name, results[:80])
		} else {
			fmt.Printf("  %s检索: %s\n", ms.Name, results)
		}
	}

	// 获取所有记忆系统的统计
	fmt.Printf("\n7. 系统整体状态:\n")

	for _, ms := range memoryTools {
		stats := ms.Tool.Run(map[string]any{"action": "stats"})
		fmt.Printf("  %s: %s\n", ms.Name, stats)
	}
}

func main() {
	fmt.Println("🧠 四种记忆类型深度解析演示")
	fmt.Println("详细展示WorkingMemory、EpisodicMemory、SemanticMemory、PerceptualMemory")
	fmt.Println(strings.Repeat("=", 80))

	demo := NewMemoryTypesDeepDive()

	// 1. 工作记忆演示
	demo.DemonstrateWorkingMemory()

	// 2. 情景记忆演示
	demo.DemonstrateEpisodicMemory()

	// 3. 语义记忆演示
	demo.DemonstrateSemanticMemory()

	// 4. 感知记忆演示
	demo.DemonstratePerceptualMemory()

	// 5. 记忆交互演示
	demo.DemonstrateMemoryInteractions()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎉 四种记忆类型深度解析完成！")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n✨ 记忆类型特性总结:")
	fmt.Println("1. 💭 工作记忆 - 快速临时存储，容量有限，自动过期")
	fmt.Println("2. 📖 情景记忆 - 完整事件记录，时间序列，丰富上下文")
	fmt.Println("3. 🧠 语义记忆 - 抽象知识存储，概念关系，语义推理")
	fmt.Println("4. 👁️ 感知记忆 - 多模态支持，跨模态检索，感知理解")

	fmt.Println("\n🔄 记忆交互模式:")
	fmt.Println("• 感知 → 工作 → 情景 → 语义（信息处理流程）")
	fmt.Println("• 语义 → 工作（知识激活和应用）")
	fmt.Println("• 跨类型检索和整合（智能记忆管理）")

	fmt.Println("\n💡 设计价值:")
	fmt.Println("• 模拟人类认知过程")
	fmt.Println("• 支持多层次信息处理")
	fmt.Println("• 实现智能记忆管理")
	fmt.Println("• 提供丰富的检索能力")
}
