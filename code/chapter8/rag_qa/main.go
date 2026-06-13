// 代码示例 07: RAGTool智能问答系统
// 展示完整的检索→上下文构建→答案生成流程

package main

import (
	"fmt"
	"strings"
	"time"
)

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
		return fmt.Sprintf("💡 问答 '%s': 基于知识库的参考来源生成答案", question)
	case "stats":
		return "📊 统计: 文档数3, 分块数24, 向量维度384"
	default:
		return "操作完成"
	}
}

// qualityResult 质量分析结果
type qualityResult struct {
	Question     string
	Difficulty   string
	AnswerLength int
	HasCitations bool
	ResponseTime float64
	QualityScore float64
}

// IntelligentQADemo 智能问答演示类
type IntelligentQADemo struct {
	RAGTool *RAGTool
}

// NewIntelligentQADemo 创建演示实例
func NewIntelligentQADemo() *IntelligentQADemo {
	demo := &IntelligentQADemo{
		RAGTool: NewRAGTool("./qa_demo_kb", "intelligent_qa_demo"),
	}
	demo.setupKnowledgeBase()
	return demo
}

// setupKnowledgeBase 设置知识库
func (d *IntelligentQADemo) setupKnowledgeBase() {
	fmt.Println("📚 设置智能问答知识库")
	fmt.Println(strings.Repeat("=", 50))

	// 添加技术知识文档
	type knowledgeDoc struct {
		ID      string
		Content string
	}

	knowledgeDocuments := []knowledgeDoc{
		{
			ID: "ai_fundamentals",
			Content: `# 人工智能基础

## 定义和历史
人工智能（Artificial Intelligence, AI）是计算机科学的一个分支，旨在创造能够执行通常需要人类智能的任务的机器。

## 主要分支
### 机器学习（Machine Learning）
机器学习是AI的核心分支，使计算机能够从数据中学习而无需明确编程。

#### 监督学习
- 分类：预测离散标签
- 回归：预测连续数值
- 常用算法：线性回归、决策树、随机森林、SVM

#### 无监督学习
- 聚类：发现数据中的群组
- 降维：减少数据维度
- 常用算法：K-means、PCA、t-SNE

#### 强化学习
通过与环境交互学习最优策略，应用于游戏AI、机器人控制等。

### 深度学习（Deep Learning）
基于人工神经网络的机器学习方法，在图像识别、自然语言处理等领域取得突破。

### 自然语言处理（NLP）
使计算机能够理解、解释和生成人类语言。
`,
		},
		{
			ID: "programming_best_practices",
			Content: `# 编程最佳实践

## 代码质量
高质量的代码应该具备可读性、可维护性和可扩展性。

### 命名规范
- 使用有意义的变量名和函数名
- 遵循一致的命名约定
- 避免使用缩写和模糊的名称

### 函数设计
- 单一职责原则：每个函数只做一件事
- 函数长度适中：通常不超过20-30行
- 参数数量合理：避免过多参数

## 测试策略
### 单元测试
- 测试单个函数或方法
- 使用断言验证预期结果
- 覆盖边界条件和异常情况

### 集成测试
- 测试模块间的交互
- 验证系统的整体功能

## 版本控制
### Git最佳实践
- 频繁提交，小步快跑
- 编写清晰的提交信息
- 使用分支管理功能开发
`,
		},
		{
			ID: "system_design",
			Content: `# 系统设计原则

## 设计模式
设计模式是软件设计中常见问题的典型解决方案。

### 创建型模式
- 单例模式：确保类只有一个实例
- 工厂模式：创建对象的接口
- 建造者模式：构建复杂对象

### 结构型模式
- 适配器模式：接口适配和转换
- 装饰器模式：动态添加功能
- 组合模式：树形结构的统一处理

### 行为型模式
- 观察者模式：对象间的一对多依赖
- 策略模式：算法的封装和切换
- 命令模式：请求的封装和参数化

## 架构原则
### SOLID原则
- 单一职责原则（SRP）
- 开闭原则（OCP）
- 里氏替换原则（LSP）
- 接口隔离原则（ISP）
- 依赖倒置原则（DIP）
`,
		},
	}

	// 批量添加知识文档
	for _, doc := range knowledgeDocuments {
		d.RAGTool.Run(map[string]any{
			"action":      "add_text",
			"text":        doc.Content,
			"document_id": doc.ID,
		})
		fmt.Printf("✅ 添加知识文档: %s\n", doc.ID)
	}

	fmt.Println("📊 知识库设置完成")
}

// DemonstrateQuestionUnderstanding 演示问题理解和分类
func (d *IntelligentQADemo) DemonstrateQuestionUnderstanding() {
	fmt.Println("\n🧠 问题理解和分类演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("问题类型分析:")
	fmt.Println("• 📖 概念定义类 - '什么是...？'")
	fmt.Println("• 🔍 方法询问类 - '如何...？'")
	fmt.Println("• ⚖️ 对比分析类 - '...和...的区别？'")
	fmt.Println("• 💡 应用场景类 - '...用于什么？'")
	fmt.Println("• 🔧 实现细节类 - '...是怎么实现的？'")

	type questionCategory struct {
		Category  string
		Questions []string
	}

	questionCategories := []questionCategory{
		{
			Category:  "概念定义",
			Questions: []string{"什么是人工智能？", "什么是深度学习？", "什么是Transformer架构？"},
		},
		{
			Category:  "方法询问",
			Questions: []string{"如何提高代码质量？", "如何进行系统设计？", "如何优化算法性能？"},
		},
		{
			Category:  "对比分析",
			Questions: []string{"监督学习和无监督学习的区别是什么？", "CNN和RNN有什么不同？", "单元测试和集成测试的区别？"},
		},
		{
			Category:  "应用场景",
			Questions: []string{"强化学习主要用于什么场景？", "设计模式在什么情况下使用？", "缓存策略适用于哪些场景？"},
		},
	}

	// 测试不同类型问题的处理效果
	for _, categoryInfo := range questionCategories {
		fmt.Printf("\n📋 %s问题测试:\n", categoryInfo.Category)

		questions := categoryInfo.Questions
		if len(questions) > 2 {
			questions = questions[:2] // 每类测试2个问题
		}

		for _, question := range questions {
			fmt.Printf("\n❓ 问题: %s\n", question)

			startTime := time.Now()
			answer := d.RAGTool.Run(map[string]any{
				"action":            "ask",
				"question":          question,
				"limit":             3,
				"include_citations": true,
			})
			qaTime := time.Since(startTime)

			fmt.Printf("⏱️ 响应时间: %.3f秒\n", qaTime.Seconds())
			if len(answer) > 300 {
				fmt.Printf("🤖 回答: %s...\n", answer[:300])
			} else {
				fmt.Printf("🤖 回答: %s\n", answer)
			}
			fmt.Println(strings.Repeat("-", 40))
		}
	}
}

// DemonstrateContextConstruction 演示上下文构建过程
func (d *IntelligentQADemo) DemonstrateContextConstruction() {
	fmt.Println("\n🏗️ 上下文构建过程演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("上下文构建步骤:")
	fmt.Println("1. 🔍 检索相关文档片段")
	fmt.Println("2. 📊 按相关性排序")
	fmt.Println("3. 🧹 清理和格式化内容")
	fmt.Println("4. ✂️ 智能截断保持完整性")
	fmt.Println("5. 🔗 添加引用信息")

	// 使用复杂问题演示上下文构建
	complexQuestion := "如何设计一个高质量的机器学习系统？"

	fmt.Printf("\n❓ 复杂问题: %s\n", complexQuestion)
	fmt.Println("这个问题需要整合多个文档的信息...")

	// 先进行搜索，查看检索到的片段
	fmt.Printf("\n🔍 第一步：检索相关片段\n")
	searchResult := d.RAGTool.Run(map[string]any{
		"action":                 "search",
		"query":                  complexQuestion,
		"limit":                  4,
		"enable_advanced_search": true,
	})
	fmt.Printf("检索片段: %s\n", searchResult)

	// 然后进行智能问答，查看完整的上下文构建
	fmt.Printf("\n🤖 第二步：构建上下文并生成答案\n")
	startTime := time.Now()
	qaResult := d.RAGTool.Run(map[string]any{
		"action":                 "ask",
		"question":               complexQuestion,
		"limit":                  4,
		"enable_advanced_search": true,
		"include_citations":      true,
		"max_chars":              1500,
	})
	qaTime := time.Since(startTime)

	fmt.Printf("问答耗时: %.3f秒\n", qaTime.Seconds())
	fmt.Printf("完整回答: %s\n", qaResult)
}

// DemonstrateAnswerQualityAnalysis 演示答案质量分析
func (d *IntelligentQADemo) DemonstrateAnswerQualityAnalysis() {
	fmt.Println("\n📊 答案质量分析演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("质量评估指标:")
	fmt.Println("• 🎯 相关性得分 - 检索内容与问题的匹配度")
	fmt.Println("• 📚 引用完整性 - 答案来源的可追溯性")
	fmt.Println("• 💡 答案完整性 - 回答的全面性和准确性")
	fmt.Println("• ⚡ 响应速度 - 系统的响应时间")

	// 质量测试问题集
	type qualityTestCase struct {
		Question         string
		ExpectedAspects  []string
		Difficulty       string
	}

	qualityTestQuestions := []qualityTestCase{
		{
			Question:        "什么是机器学习？",
			ExpectedAspects: []string{"定义", "分类", "应用"},
			Difficulty:      "简单",
		},
		{
			Question:        "如何选择合适的机器学习算法？",
			ExpectedAspects: []string{"数据特点", "问题类型", "性能要求"},
			Difficulty:      "中等",
		},
		{
			Question:        "在设计大规模系统时如何平衡性能和可维护性？",
			ExpectedAspects: []string{"架构设计", "性能优化", "代码质量"},
			Difficulty:      "复杂",
		},
	}

	fmt.Printf("\n📊 答案质量测试:\n")

	var qualityResults []qualityResult

	for _, testCase := range qualityTestQuestions {
		fmt.Printf("\n❓ 问题: %s\n", testCase.Question)
		fmt.Printf("🎯 难度: %s\n", testCase.Difficulty)
		fmt.Printf("📋 期望涵盖: %s\n", strings.Join(testCase.ExpectedAspects, ", "))

		// 执行问答
		startTime := time.Now()
		answer := d.RAGTool.Run(map[string]any{
			"action":                 "ask",
			"question":               testCase.Question,
			"limit":                  4,
			"enable_advanced_search": true,
			"include_citations":      true,
		})
		qaTime := time.Since(startTime)

		// 分析答案质量
		answerLength := len(answer)
		hasCitations := strings.Contains(answer, "参考来源")
		responseTime := qaTime.Seconds()

		qualityScore := calculateQualityScore(answer, testCase.ExpectedAspects, responseTime)

		qualityResults = append(qualityResults, qualityResult{
			Question:     testCase.Question,
			Difficulty:   testCase.Difficulty,
			AnswerLength: answerLength,
			HasCitations: hasCitations,
			ResponseTime: responseTime,
			QualityScore: qualityScore,
		})

		fmt.Printf("⏱️ 响应时间: %.3f秒\n", responseTime)
		fmt.Printf("📏 答案长度: %d字符\n", answerLength)
		citationStr := "否"
		if hasCitations {
			citationStr = "是"
		}
		fmt.Printf("📚 包含引用: %s\n", citationStr)
		fmt.Printf("⭐ 质量评分: %.2f/10\n", qualityScore)
		if len(answer) > 200 {
			fmt.Printf("🤖 答案预览: %s...\n", answer[:200])
		} else {
			fmt.Printf("🤖 答案预览: %s\n", answer)
		}
		fmt.Println(strings.Repeat("-", 50))
	}

	// 质量分析总结
	analyzeQualityResults(qualityResults)
}

// calculateQualityScore 计算答案质量评分
func calculateQualityScore(answer string, expectedAspects []string, responseTime float64) float64 {
	// 内容完整性评分 (40%)
	contentScore := 0
	for _, aspect := range expectedAspects {
		if strings.Contains(strings.ToLower(answer), strings.ToLower(aspect)) {
			contentScore++
		}
	}
	contentScoreVal := (float64(contentScore) / float64(len(expectedAspects))) * 4.0

	// 答案长度评分 (30%)
	lengthScore := min(float64(len(answer))/500, 1.0) * 3.0

	// 引用完整性评分 (20%)
	citationScore := 0.0
	if strings.Contains(answer, "参考来源") {
		citationScore = 2.0
	}

	// 响应速度评分 (10%)
	speedScore := max(0, 1.0-(responseTime-1.0)/5.0) * 1.0

	totalScore := contentScoreVal + lengthScore + citationScore + speedScore
	return min(totalScore, 10.0)
}

// analyzeQualityResults 分析质量测试结果
func analyzeQualityResults(results []qualityResult) {
	fmt.Printf("\n📈 质量分析总结:\n")

	var totalScore, totalTime float64
	citationCount := 0
	for _, r := range results {
		totalScore += r.QualityScore
		totalTime += r.ResponseTime
		if r.HasCitations {
			citationCount++
		}
	}

	avgScore := totalScore / float64(len(results))
	avgTime := totalTime / float64(len(results))
	citationRate := float64(citationCount) / float64(len(results))

	fmt.Printf("平均质量评分: %.2f/10\n", avgScore)
	fmt.Printf("平均响应时间: %.3f秒\n", avgTime)
	fmt.Printf("引用完整率: %.1f%%\n", citationRate*100)

	// 按难度分析
	difficultyAnalysis := make(map[string][]float64)
	for _, r := range results {
		difficultyAnalysis[r.Difficulty] = append(difficultyAnalysis[r.Difficulty], r.QualityScore)
	}

	fmt.Printf("\n📊 按难度分析:\n")
	for difficulty, scores := range difficultyAnalysis {
		var total float64
		for _, s := range scores {
			total += s
		}
		avg := total / float64(len(scores))
		fmt.Printf("  %s: %.2f/10\n", difficulty, avg)
	}
}

// DemonstratePromptEngineering 演示提示词工程
func (d *IntelligentQADemo) DemonstratePromptEngineering() {
	fmt.Println("\n🎨 提示词工程演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("提示词设计要素:")
	fmt.Println("• 🎯 系统角色定义")
	fmt.Println("• 📋 任务明确描述")
	fmt.Println("• 🔍 上下文信息注入")
	fmt.Println("• 📝 输出格式要求")
	fmt.Println("• 🚫 限制和约束条件")

	// 演示不同的提示词策略
	type promptStrategy struct {
		Name           string
		SystemPrompt   string
		Description    string
	}

	promptStrategies := []promptStrategy{
		{
			Name:         "基础提示",
			SystemPrompt: "你是一个AI助手，请回答用户的问题。",
			Description:  "简单直接的角色定义",
		},
		{
			Name: "专业提示",
			SystemPrompt: "你是一个专业的技术顾问，具备以下能力：\n" +
				"1. 深入理解技术概念和原理\n" +
				"2. 提供准确可靠的技术建议\n" +
				"3. 用清晰简洁的语言解释复杂概念\n" +
				"4. 基于提供的上下文信息回答问题",
			Description: "详细的能力描述和要求",
		},
		{
			Name: "结构化提示",
			SystemPrompt: "你是一个专业的知识助手，请按以下要求回答：\n" +
				"【理解】仔细分析问题的核心意图\n" +
				"【检索】基于提供的上下文信息\n" +
				"【整合】从多个片段提取关键信息\n" +
				"【回答】用结构化格式清晰表达\n" +
				"【引用】标注信息来源和依据",
			Description: "结构化的处理流程",
		},
	}

	testQuestion := "什么是深度学习，它有哪些主要应用？"

	fmt.Printf("\n🧪 提示词策略对比测试:\n")
	fmt.Printf("测试问题: %s\n", testQuestion)

	for _, strategy := range promptStrategies {
		fmt.Printf("\n📝 %s (%s):\n", strategy.Name, strategy.Description)

		startTime := time.Now()
		answer := d.RAGTool.Run(map[string]any{
			"action":   "ask",
			"question": testQuestion,
			"limit":    3,
		})
		responseTime := time.Since(startTime)

		fmt.Printf("⏱️ 响应时间: %.3f秒\n", responseTime.Seconds())
		fmt.Printf("🤖 回答长度: %d字符\n", len(answer))
		if len(answer) > 250 {
			fmt.Printf("📄 回答预览: %s...\n", answer[:250])
		} else {
			fmt.Printf("📄 回答预览: %s\n", answer)
		}
	}
}

// DemonstrateCitationSystem 演示引用系统
func (d *IntelligentQADemo) DemonstrateCitationSystem() {
	fmt.Println("\n📚 引用系统演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("引用系统特点:")
	fmt.Println("• 🔗 自动标注信息来源")
	fmt.Println("• 📊 显示相似度得分")
	fmt.Println("• 📄 提供文档定位")
	fmt.Println("• ✅ 确保答案可追溯性")

	citationTestQuestions := []string{
		"机器学习有哪些主要类型？",
		"如何进行代码质量管理？",
		"系统设计中的SOLID原则是什么？",
	}

	fmt.Printf("\n📚 引用系统测试:\n")

	for _, question := range citationTestQuestions {
		fmt.Printf("\n❓ 问题: %s\n", question)

		// 启用引用的问答
		answerWithCitations := d.RAGTool.Run(map[string]any{
			"action":            "ask",
			"question":          question,
			"limit":             3,
			"include_citations": true,
		})

		// 禁用引用的问答对比
		answerWithoutCitations := d.RAGTool.Run(map[string]any{
			"action":            "ask",
			"question":          question,
			"limit":             3,
			"include_citations": false,
		})

		if len(answerWithCitations) > 400 {
			fmt.Printf("🔗 带引用回答: %s...\n", answerWithCitations[:400])
		} else {
			fmt.Printf("🔗 带引用回答: %s\n", answerWithCitations)
		}
		if len(answerWithoutCitations) > 200 {
			fmt.Printf("📝 无引用回答: %s...\n", answerWithoutCitations[:200])
		} else {
			fmt.Printf("📝 无引用回答: %s\n", answerWithoutCitations)
		}

		// 分析引用信息
		citationCount := strings.Count(answerWithCitations, "参考来源")
		fmt.Printf("📊 引用分析: 包含 %d 个引用来源\n", citationCount)
	}
}

func main() {
	fmt.Println("🤖 RAGTool智能问答系统演示")
	fmt.Println("展示完整的检索→上下文构建→答案生成流程")
	fmt.Println(strings.Repeat("=", 70))

	demo := NewIntelligentQADemo()

	// 1. 问题理解和分类演示
	demo.DemonstrateQuestionUnderstanding()

	// 2. 上下文构建过程演示
	demo.DemonstrateContextConstruction()

	// 3. 答案质量分析演示
	demo.DemonstrateAnswerQualityAnalysis()

	// 4. 提示词工程演示
	demo.DemonstratePromptEngineering()

	// 5. 引用系统演示
	demo.DemonstrateCitationSystem()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎉 智能问答系统演示完成！")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n✨ 智能问答核心能力:")
	fmt.Println("1. 🧠 问题理解 - 准确识别问题类型和意图")
	fmt.Println("2. 🔍 智能检索 - 多策略检索相关内容")
	fmt.Println("3. 🏗️ 上下文构建 - 智能整合检索结果")
	fmt.Println("4. 🤖 答案生成 - 基于上下文的准确回答")
	fmt.Println("5. 📚 引用标注 - 完整的来源追溯")

	fmt.Println("\n🎯 技术优势:")
	fmt.Println("• 语义理解 - 深度理解问题语义和意图")
	fmt.Println("• 上下文感知 - 充分利用检索上下文")
	fmt.Println("• 质量保证 - 多层次的质量控制机制")
	fmt.Println("• 可追溯性 - 完整的答案来源追溯")

	fmt.Println("\n💡 应用场景:")
	fmt.Println("• 技术支持 - 自动回答技术问题")
	fmt.Println("• 知识问答 - 企业内部知识查询")
	fmt.Println("• 学习辅导 - 个性化学习问答")
	fmt.Println("• 文档助手 - 快速理解复杂文档")
}
