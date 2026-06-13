// 代码示例 10: RAG完整处理管道
// 展示从文档处理到智能问答的完整RAG流程

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
		return fmt.Sprintf("✅ 文本 '%s' 已分块并添加到知识库", docID)
	case "add_document":
		return "✅ 文档已处理并添加到知识库"
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关文档片段并按相关性排序", query)
	case "ask":
		question, _ := params["question"].(string)
		return fmt.Sprintf("💡 问答 '%s': 基于检索上下文生成的答案", question)
	case "stats":
		return "📊 统计: 文档数15, 分块数68, 向量维度384"
	default:
		return "操作完成"
	}
}

// RAGPipelineComplete RAG完整处理管道演示类
type RAGPipelineComplete struct {
	RAGTool *RAGTool
}

// NewRAGPipelineComplete 创建演示实例
func NewRAGPipelineComplete() *RAGPipelineComplete {
	fmt.Println("📚 RAG完整处理管道演示")
	fmt.Println(strings.Repeat("=", 60))

	demo := &RAGPipelineComplete{
		RAGTool: NewRAGTool("./rag_pipeline_kb", "complete_pipeline"),
	}

	fmt.Println("✅ RAG系统初始化完成")

	fmt.Printf("\n📊 系统配置:\n")
	fmt.Println("  知识库路径: ./rag_pipeline_kb")
	fmt.Println("  命名空间: complete_pipeline")
	fmt.Println("  支持格式: PDF, DOCX, TXT, MD, HTML, JSON")

	return demo
}

// DemonstrateDocumentIngestion 演示文档摄取过程
func (d *RAGPipelineComplete) DemonstrateDocumentIngestion() {
	fmt.Println("\n📥 文档摄取过程演示")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 文档摄取特点:")
	fmt.Println("• 📄 多格式文档支持")
	fmt.Println("• 🔄 MarkItDown格式转换")
	fmt.Println("• ✂️ 智能文档分块")
	fmt.Println("• 🎯 元数据提取")

	// 演示不同类型文档的处理
	fmt.Printf("\n1. 多格式文档处理:\n")

	type document struct {
		Content    string
		DocumentID string
		Format     string
		Metadata   map[string]any
	}

	documents := []document{
		{
			Content: `# 机器学习基础教程

## 第一章：机器学习概述

机器学习是人工智能的一个重要分支，它使计算机能够在没有明确编程的情况下学习和改进。

### 1.1 机器学习的定义

机器学习是一种数据分析方法，它自动化分析模型的构建。

### 1.2 机器学习的类型

1. **监督学习**：使用标记的训练数据来学习映射函数
2. **无监督学习**：从未标记的数据中发现隐藏的模式
3. **强化学习**：通过与环境交互来学习最优行为

### 1.3 常见算法

- 线性回归
- 逻辑回归
- 决策树
- 随机森林
- 支持向量机
- 神经网络
`,
			DocumentID: "ml_tutorial_chapter1",
			Format:     "markdown",
			Metadata: map[string]any{
				"title":                   "机器学习基础教程",
				"chapter":                 1,
				"author":                  "AI教学团队",
				"difficulty":              "beginner",
				"estimated_reading_time":  15,
			},
		},
		{
			Content: `深度学习技术报告

执行摘要：
本报告分析了深度学习在计算机视觉领域的最新进展。通过对比不同架构的性能，我们发现Transformer架构在多个任务上都表现出色。

主要发现：
1. Vision Transformer (ViT) 在图像分类任务上超越了传统CNN
2. CLIP模型实现了图像和文本的统一表示
3. 自监督学习方法显著减少了对标注数据的依赖

结论：
深度学习技术在计算机视觉领域持续快速发展，Transformer架构的引入为该领域带来了新的突破。
`,
			DocumentID: "deep_learning_report",
			Format:     "text",
			Metadata: map[string]any{
				"title":           "深度学习技术报告",
				"type":            "technical_report",
				"date":            "2024-01-15",
				"department":      "AI研究部",
				"confidentiality": "internal",
			},
		},
		{
			Content: `{
    "api_documentation": {
        "title": "机器学习API文档",
        "version": "v2.1",
        "endpoints": [
            {
                "path": "/models",
                "method": "GET",
                "description": "获取可用模型列表"
            },
            {
                "path": "/predict",
                "method": "POST",
                "description": "使用指定模型进行预测"
            }
        ]
    }
}`,
			DocumentID: "ml_api_docs",
			Format:     "json",
			Metadata: map[string]any{
				"title":         "机器学习API文档",
				"version":       "v2.1",
				"type":          "api_documentation",
				"last_updated":  "2024-01-20",
			},
		},
	}

	// 处理每个文档
	for _, doc := range documents {
		fmt.Printf("\n处理文档: %s (%s)\n", doc.DocumentID, doc.Format)

		params := map[string]any{
			"action":      "add_text",
			"text":        doc.Content,
			"document_id": doc.DocumentID,
		}
		for k, v := range doc.Metadata {
			params[k] = v
		}

		result := d.RAGTool.Run(params)
		fmt.Printf("  摄取结果: %s\n", result)

		// 显示文档统计
		docStats := map[string]any{
			"字符数":     len(doc.Content),
			"行数":       strings.Count(doc.Content, "\n") + 1,
			"格式":       doc.Format,
			"元数据字段": len(doc.Metadata),
		}
		fmt.Printf("  文档统计: %v\n", docStats)
	}

	// 演示批量文档处理
	fmt.Printf("\n2. 批量文档处理:\n")

	type batchDoc struct {
		Content    string
		DocumentID string
		Metadata   map[string]any
	}

	var batchDocuments []batchDoc
	for i := range 3 {
		batchDocuments = append(batchDocuments, batchDoc{
			Content: fmt.Sprintf(`# 批量文档 %d

这是第 %d 个批量处理的文档。它包含了关于人工智能发展的重要信息。

## 主要内容
- AI技术趋势分析
- 行业应用案例
- 未来发展预测

## 详细描述
人工智能技术在过去几年中取得了显著进展，特别是在深度学习、自然语言处理和计算机视觉领域。
`, i+1, i+1),
			DocumentID: fmt.Sprintf("batch_doc_%d", i+1),
			Metadata: map[string]any{
				"batch_id": "batch_001",
				"sequence": i + 1,
				"topic":    "artificial_intelligence",
			},
		})
	}

	// 批量处理
	startTime := time.Now()
	for _, doc := range batchDocuments {
		params := map[string]any{
			"action":      "add_text",
			"text":        doc.Content,
			"document_id": doc.DocumentID,
		}
		for k, v := range doc.Metadata {
			params[k] = v
		}
		result := d.RAGTool.Run(params)
		fmt.Printf("  批量处理 %s: %s\n", doc.DocumentID, result)
	}

	batchTime := time.Since(startTime)
	fmt.Printf("  批量处理耗时: %.3f秒\n", batchTime.Seconds())

	// 获取摄取统计
	stats := d.RAGTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n📊 文档摄取统计: %s\n", stats)
}

// DemonstrateChunkingStrategies 演示文档分块策略
func (d *RAGPipelineComplete) DemonstrateChunkingStrategies() {
	fmt.Println("\n✂️ 文档分块策略演示")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 分块策略特点:")
	fmt.Println("• 📏 基于语义的智能分块")
	fmt.Println("• 🔗 保持上下文连贯性")
	fmt.Println("• ⚖️ 平衡块大小和信息完整性")
	fmt.Println("• 🎯 优化检索效果")

	// 演示不同分块策略
	fmt.Printf("\n1. 分块策略对比:\n")

	// 长文档示例
	longDocument := `# 人工智能发展史

## 引言
人工智能（Artificial Intelligence, AI）的发展历程可以追溯到20世纪50年代。

## 第一阶段：符号主义时代（1950s-1980s）
1950年，阿兰·图灵发表了著名的论文《计算机器与智能》，提出了"图灵测试"的概念。

## 第二阶段：连接主义复兴（1980s-2000s）
1986年，Rumelhart等人重新发现了反向传播算法，使得多层神经网络的训练成为可能。

## 第三阶段：深度学习革命（2000s-至今）
2006年，Geoffrey Hinton等人提出了深度信念网络，开启了深度学习的新时代。

## 第四阶段：通用人工智能探索（2020s-未来）
人工智能正朝着更加通用、可解释和安全的方向发展。

## 结论
人工智能的发展是一个螺旋上升的过程，每个阶段都有其独特的贡献和局限性。
`

	// 添加长文档并观察分块效果
	chunkingResult := d.RAGTool.Run(map[string]any{
		"action":            "add_text",
		"text":              longDocument,
		"document_id":       "ai_history_long",
		"title":             "人工智能发展史",
		"type":              "historical_overview",
		"chunking_strategy": "semantic",
	})
	fmt.Printf("长文档分块结果: %s\n", chunkingResult)

	// 演示不同分块大小的影响
	fmt.Printf("\n2. 分块大小影响分析:\n")

	testQueries := []string{
		"图灵测试是什么？",
		"深度学习的关键技术突破",
		"AlphaGo的意义",
		"通用人工智能的未来",
	}

	for _, query := range testQueries {
		startTime := time.Now()
		results := d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  3,
		})
		searchTime := time.Since(startTime)
		fmt.Printf("  查询: '%s' (%.4f秒)\n", query, searchTime.Seconds())
		if len(results) > 120 {
			fmt.Printf("    结果: %s...\n", results[:120])
		} else {
			fmt.Printf("    结果: %s\n", results)
		}
	}

	// 演示结构化文档的分块
	fmt.Printf("\n3. 结构化文档分块:\n")

	structuredDoc := `# 机器学习算法手册

## 监督学习算法

### 线性回归
**定义**: 线性回归是一种用于预测连续数值的算法。
**公式**: y = wx + b
**优点**: 简单易懂，计算效率高
**缺点**: 只能处理线性关系

### 逻辑回归
**定义**: 逻辑回归用于二分类问题。
**优点**: 输出概率值，可解释性强
**缺点**: 对特征工程要求高

### 决策树
**定义**: 基于特征进行分层决策的树形结构。
**优点**: 可解释性强，处理非线性关系
**缺点**: 容易过拟合

## 无监督学习算法

### K-means聚类
**定义**: 将数据分为K个簇的聚类算法。
**优点**: 简单高效
**缺点**: 需要预设簇数

### 主成分分析(PCA)
**定义**: 降维算法，保留主要信息。
**优点**: 降低维度，去除噪声
**缺点**: 损失部分信息
`

	structuredResult := d.RAGTool.Run(map[string]any{
		"action":      "add_text",
		"text":        structuredDoc,
		"document_id": "ml_algorithms_handbook",
		"title":       "机器学习算法手册",
		"type":        "reference_manual",
		"structure":   "hierarchical",
	})
	fmt.Printf("结构化文档分块: %s\n", structuredResult)

	// 测试结构化检索
	structuredQueries := []string{
		"线性回归的优缺点",
		"K-means聚类算法",
		"PCA降维原理",
	}

	for _, query := range structuredQueries {
		results := d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  2,
		})
		if len(results) > 100 {
			fmt.Printf("  结构化查询 '%s': %s...\n", query, results[:100])
		} else {
			fmt.Printf("  结构化查询 '%s': %s\n", query, results)
		}
	}
}

// DemonstrateAdvancedRetrieval 演示高级检索策略
func (d *RAGPipelineComplete) DemonstrateAdvancedRetrieval() {
	fmt.Println("\n🔍 高级检索策略演示")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 高级检索特点:")
	fmt.Println("• 🎯 多查询扩展（MQE）")
	fmt.Println("• 💭 假设文档嵌入（HyDE）")
	fmt.Println("• 🔄 混合检索策略")
	fmt.Println("• 📊 相关性重排序")

	// 演示多查询扩展
	fmt.Printf("\n1. 多查询扩展（MQE）演示:\n")

	baseQuery := "如何提高机器学习模型的性能？"
	fmt.Printf("原始查询: %s\n", baseQuery)

	// 模拟查询扩展
	expandedQueries := []string{
		"机器学习模型性能优化方法",
		"提升ML模型准确率的技巧",
		"模型调优和超参数优化",
		"机器学习模型评估指标",
	}

	fmt.Println("扩展查询:")
	for i, query := range expandedQueries {
		fmt.Printf("  %d. %s\n", i+1, query)
	}

	// 执行多查询检索
	allQueries := append([]string{baseQuery}, expandedQueries...)
	for _, query := range allQueries {
		results := d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  3,
		})
		if len(results) > 80 {
			fmt.Printf("  查询结果 '%s...': %s...\n", query[:20], results[:80])
		} else {
			fmt.Printf("  查询结果 '%s': %s\n", query, results)
		}
	}

	// 演示假设文档嵌入（HyDE）
	fmt.Printf("\n2. 假设文档嵌入（HyDE）演示:\n")

	userQuestion := "什么是深度学习？"
	fmt.Printf("用户问题: %s\n", userQuestion)

	hypotheticalAnswer := "深度学习是机器学习的一个子领域，它使用多层神经网络来学习数据的复杂模式。深度学习模型通过多个隐藏层来提取数据的层次化特征表示。常见的深度学习架构包括卷积神经网络（CNN）、循环神经网络（RNN）和Transformer。深度学习在图像识别、自然语言处理、语音识别等领域取得了突破性进展。"

	if len(hypotheticalAnswer) > 100 {
		fmt.Printf("假设答案: %s...\n", hypotheticalAnswer[:100])
	} else {
		fmt.Printf("假设答案: %s\n", hypotheticalAnswer)
	}

	// 使用假设答案进行检索
	hydeResults := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  hypotheticalAnswer,
		"limit":  5,
	})
	if len(hydeResults) > 120 {
		fmt.Printf("HyDE检索结果: %s...\n", hydeResults[:120])
	} else {
		fmt.Printf("HyDE检索结果: %s\n", hydeResults)
	}

	// 对比直接查询结果
	directResults := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  userQuestion,
		"limit":  5,
	})
	if len(directResults) > 120 {
		fmt.Printf("直接查询结果: %s...\n", directResults[:120])
	} else {
		fmt.Printf("直接查询结果: %s\n", directResults)
	}

	// 演示混合检索策略
	fmt.Printf("\n3. 混合检索策略演示:\n")

	complexQuery := "比较监督学习和无监督学习的区别，并给出具体应用例子"
	fmt.Printf("复杂查询: %s\n", complexQuery)

	subQueries := []string{
		"监督学习的定义和特点",
		"无监督学习的定义和特点",
		"监督学习的应用例子",
		"无监督学习的应用例子",
		"监督学习和无监督学习的区别",
	}

	fmt.Println("查询分解:")
	for _, subQuery := range subQueries {
		results := d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  subQuery,
			"limit":  2,
		})
		fmt.Printf("  子查询: %s\n", subQuery)
		if len(results) > 80 {
			fmt.Printf("    结果: %s...\n", results[:80])
		} else {
			fmt.Printf("    结果: %s\n", results)
		}
	}

	// 演示相关性重排序
	fmt.Printf("\n4. 相关性重排序演示:\n")

	rankingQuery := "神经网络训练过程"
	fmt.Printf("排序查询: %s\n", rankingQuery)

	// 获取初始结果
	initialResults := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  rankingQuery,
		"limit":  8,
	})
	if len(initialResults) > 150 {
		fmt.Printf("初始检索结果: %s...\n", initialResults[:150])
	} else {
		fmt.Printf("初始检索结果: %s\n", initialResults)
	}

	// 模拟重排序过程（基于多个因素）
	fmt.Println("重排序因素:")
	fmt.Println("  • 语义相似度权重: 0.6")
	fmt.Println("  • 文档新鲜度权重: 0.2")
	fmt.Println("  • 文档权威性权重: 0.2")

	// 最终排序结果
	finalResults := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  rankingQuery,
		"limit":  5,
	})
	if len(finalResults) > 150 {
		fmt.Printf("重排序后结果: %s...\n", finalResults[:150])
	} else {
		fmt.Printf("重排序后结果: %s\n", finalResults)
	}
}

// DemonstrateIntelligentQA 演示智能问答生成
func (d *RAGPipelineComplete) DemonstrateIntelligentQA() {
	fmt.Println("\n🤖 智能问答生成演示")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 智能问答特点:")
	fmt.Println("• 🎯 问题理解和分类")
	fmt.Println("• 📚 上下文构建")
	fmt.Println("• 💡 答案生成和优化")
	fmt.Println("• 🔗 引用和溯源")

	// 演示不同类型问题的处理
	fmt.Printf("\n1. 不同类型问题处理:\n")

	type qaExample struct {
		Question         string
		Type             string
		ExpectedApproach string
	}

	qaExamples := []qaExample{
		{"什么是机器学习？", "定义类问题", "提供清晰定义和基本概念"},
		{"如何选择合适的机器学习算法？", "方法类问题", "提供步骤和决策框架"},
		{"深度学习和传统机器学习有什么区别？", "比较类问题", "对比分析优缺点"},
		{"为什么神经网络需要激活函数？", "原理类问题", "解释技术原理和必要性"},
		{"在图像分类项目中应该使用哪种算法？", "应用类问题", "结合场景给出具体建议"},
	}

	for _, example := range qaExamples {
		fmt.Printf("\n问题类型: %s\n", example.Type)
		fmt.Printf("问题: %s\n", example.Question)
		fmt.Printf("处理策略: %s\n", example.ExpectedApproach)

		startTime := time.Now()
		answer := d.RAGTool.Run(map[string]any{
			"action":   "ask",
			"question": example.Question,
			"limit":    4,
		})
		qaTime := time.Since(startTime)

		if len(answer) > 200 {
			fmt.Printf("回答 (%.3f秒): %s...\n", qaTime.Seconds(), answer[:200])
		} else {
			fmt.Printf("回答 (%.3f秒): %s\n", qaTime.Seconds(), answer)
		}
	}

	// 演示上下文构建过程
	fmt.Printf("\n2. 上下文构建过程演示:\n")

	contextQuestion := "如何防止神经网络过拟合？"
	fmt.Printf("问题: %s\n", contextQuestion)

	fmt.Println("上下文构建步骤:")
	fmt.Println("  1. 问题分析 - 识别关键概念：过拟合、神经网络、防止方法")
	fmt.Println("  2. 相关文档检索 - 搜索相关技术文档")
	fmt.Println("  3. 上下文筛选 - 选择最相关的信息片段")
	fmt.Println("  4. 上下文排序 - 按相关性和重要性排序")

	contextSearch := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  "神经网络过拟合防止方法",
		"limit":  6,
	})
	if len(contextSearch) > 180 {
		fmt.Printf("  检索到的上下文: %s...\n", contextSearch[:180])
	} else {
		fmt.Printf("  检索到的上下文: %s\n", contextSearch)
	}

	finalAnswer := d.RAGTool.Run(map[string]any{
		"action":   "ask",
		"question": contextQuestion,
		"limit":    5,
	})
	if len(finalAnswer) > 250 {
		fmt.Printf("  最终答案: %s...\n", finalAnswer[:250])
	} else {
		fmt.Printf("  最终答案: %s\n", finalAnswer)
	}

	// 演示多轮对话支持
	fmt.Printf("\n3. 多轮对话支持:\n")

	conversation := []string{
		"什么是卷积神经网络？",
		"它主要用于什么任务？",
		"相比传统方法有什么优势？",
		"在实际项目中如何使用？",
	}

	fmt.Println("模拟对话场景:")
	for i, question := range conversation {
		fmt.Printf("\n  轮次 %d: %s\n", i+1, question)

		// 在多轮对话中，后续问题可能需要前面的上下文
		contextQuery := question
		if i > 0 {
			contextQuery = fmt.Sprintf("卷积神经网络 %s", question)
		}

		answer := d.RAGTool.Run(map[string]any{
			"action":   "ask",
			"question": contextQuery,
			"limit":    3,
		})
		if len(answer) > 150 {
			fmt.Printf("  回答: %s...\n", answer[:150])
		} else {
			fmt.Printf("  回答: %s\n", answer)
		}
	}

	// 演示答案质量评估
	fmt.Printf("\n4. 答案质量评估:\n")

	qualityQuestion := "解释反向传播算法的工作原理"
	fmt.Printf("评估问题: %s\n", qualityQuestion)

	answer := d.RAGTool.Run(map[string]any{
		"action":   "ask",
		"question": qualityQuestion,
		"limit":    5,
	})

	if len(answer) > 300 {
		fmt.Printf("生成答案: %s...\n", answer[:300])
	} else {
		fmt.Printf("生成答案: %s\n", answer)
	}

	// 模拟质量评估指标
	qualityMetrics := map[string]string{
		"相关性":   "高 - 答案直接回应了问题",
		"准确性":   "高 - 技术描述准确",
		"完整性":   "中 - 涵盖了主要概念",
		"可读性":   "高 - 结构清晰易懂",
		"引用质量": "中 - 基于可靠来源",
	}

	fmt.Println("质量评估:")
	for metric, score := range qualityMetrics {
		fmt.Printf("  %s: %s\n", metric, score)
	}
}

// DemonstratePerformanceOptimization 演示性能优化
func (d *RAGPipelineComplete) DemonstratePerformanceOptimization() {
	fmt.Println("\n⚡ 性能优化演示")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("🔍 性能优化特点:")
	fmt.Println("• 🚀 检索速度优化")
	fmt.Println("• 💾 内存使用优化")
	fmt.Println("• 🎯 结果质量提升")
	fmt.Println("• 📊 系统监控")

	// 演示检索性能测试
	fmt.Printf("\n1. 检索性能测试:\n")

	performanceQueries := []string{
		"机器学习基础概念",
		"深度学习应用场景",
		"神经网络训练技巧",
		"数据预处理方法",
		"模型评估指标",
	}

	totalTime := 0.0
	totalQueries := len(performanceQueries)

	fmt.Printf("执行 %d 个查询的性能测试:\n", totalQueries)

	for i, query := range performanceQueries {
		startTime := time.Now()
		d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  5,
		})
		queryTime := time.Since(startTime)
		totalTime += queryTime.Seconds()

		fmt.Printf("  查询 %d: '%s' - %.4f秒\n", i+1, query, queryTime.Seconds())
	}

	avgTime := totalTime / float64(totalQueries)
	fmt.Printf("\n性能统计:\n")
	fmt.Printf("  总耗时: %.4f秒\n", totalTime)
	fmt.Printf("  平均查询时间: %.4f秒\n", avgTime)
	if avgTime > 0 {
		fmt.Printf("  查询吞吐量: %.2f 查询/秒\n", 1/avgTime)
	}

	// 演示批量处理优化
	fmt.Printf("\n2. 批量处理优化:\n")

	batchQueries := []string{
		"什么是监督学习？",
		"什么是无监督学习？",
		"什么是强化学习？",
		"什么是深度学习？",
		"什么是神经网络？",
	}

	// 单个处理
	startTime := time.Now()
	for _, query := range batchQueries {
		d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  2,
		})
	}
	individualTime := time.Since(startTime)

	fmt.Printf("  单个处理耗时: %.4f秒\n", individualTime.Seconds())

	// 模拟批量处理
	startTime = time.Now()
	for _, query := range batchQueries {
		d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  query,
			"limit":  2,
		})
	}
	batchTime := time.Since(startTime)

	fmt.Printf("  批量处理耗时: %.4f秒\n", batchTime.Seconds())
	if individualTime.Seconds() > 0 {
		fmt.Printf("  性能提升: %.1f%%\n", (individualTime.Seconds()-batchTime.Seconds())/individualTime.Seconds()*100)
	}

	// 演示缓存机制
	fmt.Printf("\n3. 缓存机制演示:\n")

	cacheQuery := "机器学习算法分类"

	// 第一次查询（无缓存）
	startTime = time.Now()
	d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  cacheQuery,
		"limit":  3,
	})
	firstTime := time.Since(startTime)
	fmt.Printf("  首次查询: %.4f秒\n", firstTime.Seconds())

	// 第二次查询（可能有缓存）
	startTime = time.Now()
	d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  cacheQuery,
		"limit":  3,
	})
	secondTime := time.Since(startTime)
	fmt.Printf("  重复查询: %.4f秒\n", secondTime.Seconds())

	if secondTime.Seconds() < firstTime.Seconds() && firstTime.Seconds() > 0 {
		speedup := (firstTime.Seconds() - secondTime.Seconds()) / firstTime.Seconds() * 100
		fmt.Printf("  缓存加速: %.1f%%\n", speedup)
	}

	// 演示系统监控
	fmt.Printf("\n4. 系统监控:\n")

	systemStats := d.RAGTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("  系统统计: %s\n", systemStats)

	// 模拟资源使用监控
	fmt.Println("  资源使用情况:")
	fmt.Println("    文档数量: 15个")
	fmt.Println("    索引大小: 约2.5MB")
	fmt.Println("    内存使用: 约128MB")
	fmt.Printf("    平均响应时间: %.4f秒\n", avgTime)
	fmt.Println("    成功率: 100%")
}

func main() {
	fmt.Println("📚 RAG完整处理管道演示")
	fmt.Println("展示从文档处理到智能问答的完整RAG流程")
	fmt.Println(strings.Repeat("=", 80))

	demo := NewRAGPipelineComplete()

	// 1. 文档摄取演示
	demo.DemonstrateDocumentIngestion()

	// 2. 分块策略演示
	demo.DemonstrateChunkingStrategies()

	// 3. 高级检索演示
	demo.DemonstrateAdvancedRetrieval()

	// 4. 智能问答演示
	demo.DemonstrateIntelligentQA()

	// 5. 性能优化演示
	demo.DemonstratePerformanceOptimization()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎉 RAG完整处理管道演示完成！")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n✨ RAG管道核心特性:")
	fmt.Println("1. 📥 多格式文档摄取 - 支持PDF、DOCX、TXT、MD等")
	fmt.Println("2. ✂️ 智能文档分块 - 基于语义的分块策略")
	fmt.Println("3. 🔍 高级检索策略 - MQE、HyDE、混合检索")
	fmt.Println("4. 🤖 智能问答生成 - 上下文构建和答案优化")
	fmt.Println("5. ⚡ 性能优化 - 缓存、批量处理、监控")

	fmt.Println("\n🎯 技术优势:")
	fmt.Println("• 端到端处理流程")
	fmt.Println("• 多策略检索优化")
	fmt.Println("• 智能上下文构建")
	fmt.Println("• 高质量答案生成")
	fmt.Println("• 全面性能监控")

	fmt.Println("\n💡 应用场景:")
	fmt.Println("• 企业知识库问答")
	fmt.Println("• 技术文档助手")
	fmt.Println("• 学习辅导系统")
	fmt.Println("• 智能客服系统")
}
