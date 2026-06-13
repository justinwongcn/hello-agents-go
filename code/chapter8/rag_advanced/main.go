// 代码示例 05: RAGTool高级检索策略
// 展示MQE、HyDE等先进检索技术的实现和应用

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
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关文档片段并按相关性排序", query)
	case "ask":
		question, _ := params["question"].(string)
		return fmt.Sprintf("💡 问答 '%s': 基于检索上下文生成的答案", question)
	case "stats":
		return "📊 统计: 文档数4, 分块数32, 向量维度384"
	default:
		return "操作完成"
	}
}

// AdvancedSearchDemo 高级检索演示类
type AdvancedSearchDemo struct {
	RAGTool *RAGTool
}

// NewAdvancedSearchDemo 创建演示实例
func NewAdvancedSearchDemo() *AdvancedSearchDemo {
	demo := &AdvancedSearchDemo{
		RAGTool: NewRAGTool("./advanced_search_kb", "advanced_search_demo"),
	}
	demo.setupKnowledgeBase()
	return demo
}

// setupKnowledgeBase 设置知识库内容
func (d *AdvancedSearchDemo) setupKnowledgeBase() {
	fmt.Println("📚 设置知识库内容")
	fmt.Println(strings.Repeat("=", 50))

	// 添加技术文档
	techDocuments := []struct {
		ID      string
		Content string
	}{
		{
			ID: "transformer_architecture",
			Content: `# Transformer架构详解

## 注意力机制
Transformer的核心是自注意力机制（Self-Attention），它允许模型在处理序列时关注到序列中的不同位置。

### 多头注意力
多头注意力机制将输入投影到多个不同的子空间，每个头关注不同的表示子空间。

### 位置编码
由于Transformer没有循环结构，需要位置编码来提供序列中位置信息。

## 编码器-解码器结构
- 编码器：将输入序列编码为表示
- 解码器：基于编码器输出生成目标序列

## 应用领域
- 机器翻译
- 文本摘要
- 问答系统
- 代码生成
`,
		},
		{
			ID: "deep_learning_optimization",
			Content: `# 深度学习优化技术

## 梯度下降算法
梯度下降是深度学习中最基础的优化算法。

### 随机梯度下降（SGD）
- 每次使用单个样本更新参数
- 计算效率高，但收敛不稳定

### 批量梯度下降
- 使用全部训练数据计算梯度
- 收敛稳定，但计算成本高

### 小批量梯度下降
- 平衡了SGD和批量梯度下降的优缺点
- 是实际应用中最常用的方法

## 自适应学习率算法
- Adam：结合动量和自适应学习率
- AdaGrad：根据历史梯度调整学习率
- RMSprop：解决AdaGrad学习率衰减过快的问题

## 正则化技术
- Dropout：随机丢弃神经元防止过拟合
- Batch Normalization：标准化层输入
- Weight Decay：权重衰减正则化
`,
		},
		{
			ID: "nlp_applications",
			Content: `# 自然语言处理应用

## 文本分类
文本分类是NLP中的基础任务，包括情感分析、主题分类、垃圾邮件检测等。

### 传统方法
- 词袋模型（Bag of Words）
- TF-IDF特征
- 朴素贝叶斯分类器

### 深度学习方法
- CNN用于文本分类
- RNN和LSTM处理序列信息
- BERT等预训练模型

## 命名实体识别（NER）
识别文本中的人名、地名、组织名等实体。

## 机器翻译
将一种语言的文本翻译成另一种语言。

### 神经机器翻译
- Seq2Seq模型
- 注意力机制
- Transformer架构
`,
		},
		{
			ID: "computer_vision",
			Content: `# 计算机视觉技术

## 图像分类
图像分类是计算机视觉的基础任务，目标是将图像分配到预定义的类别中。

### 卷积神经网络（CNN）
- 卷积层：提取局部特征
- 池化层：降低维度和计算量
- 全连接层：进行最终分类

### 经典架构
- LeNet：最早的CNN架构
- AlexNet：深度学习在图像识别的突破
- VGG：使用小卷积核的深层网络
- ResNet：残差连接解决梯度消失

## 目标检测
在图像中定位和识别多个对象。

### 两阶段方法
- R-CNN：区域提议+CNN分类
- Faster R-CNN：RPN网络生成提议

### 单阶段方法
- YOLO：将检测作为回归问题
- SSD：多尺度特征检测

## 图像分割
- FCN：全卷积网络
- U-Net：编码器-解码器结构
- Mask R-CNN：在Faster R-CNN基础上添加分割分支
`,
		},
	}

	// 批量添加文档
	for _, doc := range techDocuments {
		d.RAGTool.Run(map[string]any{
			"action":      "add_text",
			"text":        doc.Content,
			"document_id": doc.ID,
		})
		fmt.Printf("✅ 添加文档: %s\n", doc.ID)
	}

	fmt.Printf("📊 知识库设置完成，共添加 %d 个文档\n", len(techDocuments))
}

// DemonstrateBasicSearch 演示基础搜索功能
func (d *AdvancedSearchDemo) DemonstrateBasicSearch() {
	fmt.Println("\n🔍 基础搜索功能演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("基础搜索特点:")
	fmt.Println("• 向量相似度匹配")
	fmt.Println("• 基于嵌入的语义理解")
	fmt.Println("• 相关性排序")
	fmt.Println("• 快速响应")

	type basicQuery struct {
		Query       string
		Description string
	}

	basicQueries := []basicQuery{
		{"注意力机制", "测试精确概念匹配"},
		{"深度学习优化", "测试主题匹配"},
		{"图像分类CNN", "测试多词匹配"},
		{"机器翻译模型", "测试跨文档匹配"},
	}

	fmt.Printf("\n🔍 基础搜索测试:\n")
	for _, bq := range basicQueries {
		fmt.Printf("\n查询: '%s' (%s)\n", bq.Query, bq.Description)

		startTime := time.Now()
		result := d.RAGTool.Run(map[string]any{
			"action":                   "search",
			"query":                    bq.Query,
			"limit":                    2,
			"enable_advanced_search":   false,
		})
		searchTime := time.Since(startTime)

		fmt.Printf("耗时: %.3f秒\n", searchTime.Seconds())
		if len(result) > 200 {
			fmt.Printf("结果: %s...\n", result[:200])
		} else {
			fmt.Printf("结果: %s\n", result)
		}
	}
}

// DemonstrateMQESearch 演示多查询扩展（MQE）搜索
func (d *AdvancedSearchDemo) DemonstrateMQESearch() {
	fmt.Println("\n🔄 多查询扩展（MQE）搜索演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("MQE搜索原理:")
	fmt.Println("• 🤖 使用LLM生成语义等价查询")
	fmt.Println("• 🔍 并行执行多个查询")
	fmt.Println("• 📊 合并和去重结果")
	fmt.Println("• 🎯 提高召回率和覆盖面")

	type mqeQuery struct {
		Query       string
		Description string
	}

	mqeQueries := []mqeQuery{
		{"深度学习", "测试概念扩展"},
		{"优化算法", "测试技术扩展"},
		{"神经网络", "测试架构扩展"},
	}

	fmt.Printf("\n🔄 MQE搜索测试:\n")
	for _, mq := range mqeQueries {
		fmt.Printf("\n查询: '%s' (%s)\n", mq.Query, mq.Description)

		// 基础搜索对比
		startTime := time.Now()
		basicResult := d.RAGTool.Run(map[string]any{
			"action":                   "search",
			"query":                    mq.Query,
			"limit":                    3,
			"enable_advanced_search":   false,
		})
		basicTime := time.Since(startTime)

		// MQE搜索
		startTime = time.Now()
		mqeResult := d.RAGTool.Run(map[string]any{
			"action":                   "search",
			"query":                    mq.Query,
			"limit":                    3,
			"enable_advanced_search":   true,
		})
		mqeTime := time.Since(startTime)

		fmt.Printf("基础搜索耗时: %.3f秒\n", basicTime.Seconds())
		fmt.Printf("MQE搜索耗时: %.3f秒\n", mqeTime.Seconds())
		if len(basicResult) > 150 {
			fmt.Printf("基础结果: %s...\n", basicResult[:150])
		} else {
			fmt.Printf("基础结果: %s\n", basicResult)
		}
		if len(mqeResult) > 150 {
			fmt.Printf("MQE结果: %s...\n", mqeResult[:150])
		} else {
			fmt.Printf("MQE结果: %s\n", mqeResult)
		}
		if basicTime.Seconds() > 0 {
			fmt.Printf("性能对比: MQE搜索耗时是基础搜索的 %.1f 倍\n", mqeTime.Seconds()/basicTime.Seconds())
		}
	}
}

// DemonstrateHyDESearch 演示假设文档嵌入（HyDE）搜索
func (d *AdvancedSearchDemo) DemonstrateHyDESearch() {
	fmt.Println("\n📝 假设文档嵌入（HyDE）搜索演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("HyDE搜索原理:")
	fmt.Println("• 🤖 LLM生成假设性答案文档")
	fmt.Println("• 📄 将假设文档作为查询向量")
	fmt.Println("• 🎯 改善查询-文档匹配效果")
	fmt.Println("• 🔍 特别适合复杂问题检索")

	type hydeQuery struct {
		Query       string
		Description string
	}

	hydeQueries := []hydeQuery{
		{"如何提高深度学习模型的性能？", "测试方法性问题"},
		{"Transformer相比RNN有什么优势？", "测试对比性问题"},
		{"什么是计算机视觉中的目标检测？", "测试定义性问题"},
	}

	fmt.Printf("\n📝 HyDE搜索测试:\n")
	for _, hq := range hydeQueries {
		fmt.Printf("\n查询: '%s' (%s)\n", hq.Query, hq.Description)

		// 使用智能问答（内部使用HyDE）
		startTime := time.Now()
		hydeResult := d.RAGTool.Run(map[string]any{
			"action":                   "ask",
			"question":                 hq.Query,
			"limit":                    3,
			"enable_advanced_search":   true,
		})
		hydeTime := time.Since(startTime)

		fmt.Printf("HyDE问答耗时: %.3f秒\n", hydeTime.Seconds())
		if len(hydeResult) > 300 {
			fmt.Printf("HyDE结果: %s...\n", hydeResult[:300])
		} else {
			fmt.Printf("HyDE结果: %s\n", hydeResult)
		}
	}
}

// DemonstrateCombinedAdvancedSearch 演示组合高级搜索
func (d *AdvancedSearchDemo) DemonstrateCombinedAdvancedSearch() {
	fmt.Println("\n🚀 组合高级搜索演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("组合搜索策略:")
	fmt.Println("• 🔄 MQE + HyDE 双重扩展")
	fmt.Println("• 📊 多策略结果融合")
	fmt.Println("• 🎯 最大化检索效果")
	fmt.Println("• ⚡ 智能缓存优化")

	type complexQuery struct {
		Query       string
		Description string
	}

	complexQueries := []complexQuery{
		{"深度学习中的注意力机制是如何工作的？", "复杂技术问题"},
		{"比较不同的梯度下降优化算法", "对比分析问题"},
		{"计算机视觉和自然语言处理的共同技术", "跨领域问题"},
	}

	fmt.Printf("\n🚀 组合高级搜索测试:\n")
	for _, cq := range complexQueries {
		fmt.Printf("\n查询: '%s' (%s)\n", cq.Query, cq.Description)

		// 组合高级搜索
		startTime := time.Now()

		// 先进行高级搜索获取相关片段
		searchResult := d.RAGTool.Run(map[string]any{
			"action":                   "search",
			"query":                    cq.Query,
			"limit":                    4,
			"enable_advanced_search":   true,
		})

		// 再进行智能问答生成完整答案
		qaResult := d.RAGTool.Run(map[string]any{
			"action":                   "ask",
			"question":                 cq.Query,
			"limit":                    4,
			"enable_advanced_search":   true,
			"include_citations":        true,
		})

		combinedTime := time.Since(startTime)

		fmt.Printf("组合搜索耗时: %.3f秒\n", combinedTime.Seconds())
		if len(searchResult) > 200 {
			fmt.Printf("搜索片段: %s...\n", searchResult[:200])
		} else {
			fmt.Printf("搜索片段: %s\n", searchResult)
		}
		if len(qaResult) > 400 {
			fmt.Printf("智能问答: %s...\n", qaResult[:400])
		} else {
			fmt.Printf("智能问答: %s\n", qaResult)
		}
	}
}

// DemonstrateSearchPerformanceAnalysis 演示搜索性能分析
func (d *AdvancedSearchDemo) DemonstrateSearchPerformanceAnalysis() {
	fmt.Println("\n📊 搜索性能分析")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("性能分析指标:")
	fmt.Println("• ⏱️ 响应时间对比")
	fmt.Println("• 🎯 检索质量评估")
	fmt.Println("• 💾 资源使用情况")
	fmt.Println("• 📈 扩展性分析")

	// 性能测试查询
	performanceQueries := []string{
		"机器学习",
		"深度学习优化算法",
		"Transformer注意力机制原理",
		"计算机视觉目标检测方法比较",
	}

	fmt.Printf("\n📊 性能对比测试:\n")

	// 测试不同搜索策略的性能
	type strategy struct {
		Name   string
		Params map[string]any
	}

	strategies := []strategy{
		{"基础搜索", map[string]any{"enable_advanced_search": false}},
		{"高级搜索", map[string]any{"enable_advanced_search": true}},
	}

	type performanceResult struct {
		Times   []float64
		Average float64
	}

	performanceResults := make(map[string]*performanceResult)

	for _, s := range strategies {
		fmt.Printf("\n%s性能测试:\n", s.Name)
		var strategyTimes []float64

		for _, query := range performanceQueries {
			startTime := time.Now()

			params := map[string]any{
				"action": "search",
				"query":  query,
				"limit":  3,
			}
			for k, v := range s.Params {
				params[k] = v
			}

			d.RAGTool.Run(params)

			queryTime := time.Since(startTime)
			strategyTimes = append(strategyTimes, queryTime.Seconds())

			displayQuery := query
			if len(displayQuery) > 20 {
				displayQuery = displayQuery[:20]
			}
			fmt.Printf("  查询: '%s...' 耗时: %.3f秒\n", displayQuery, queryTime.Seconds())
		}

		var totalTime float64
		for _, t := range strategyTimes {
			totalTime += t
		}
		avgTime := totalTime / float64(len(strategyTimes))
		performanceResults[s.Name] = &performanceResult{
			Times:   strategyTimes,
			Average: avgTime,
		}

		fmt.Printf("  平均耗时: %.3f秒\n", avgTime)
	}

	// 性能对比分析
	fmt.Printf("\n📈 性能对比分析:\n")
	basicAvg := performanceResults["基础搜索"].Average
	advancedAvg := performanceResults["高级搜索"].Average

	fmt.Printf("基础搜索平均耗时: %.3f秒\n", basicAvg)
	fmt.Printf("高级搜索平均耗时: %.3f秒\n", advancedAvg)
	if basicAvg > 0 {
		ratio := advancedAvg / basicAvg
		fmt.Printf("性能比值: %.1fx\n", ratio)
		fmt.Printf("分析: 高级搜索通过多策略提升检索质量，耗时增加 %.0f%%\n", (ratio-1)*100)
	}

	// 获取系统统计
	stats := d.RAGTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n📊 系统统计: %s\n", stats)
}

func main() {
	fmt.Println("🚀 RAGTool高级检索策略演示")
	fmt.Println("展示MQE、HyDE等先进检索技术的实现和应用")
	fmt.Println(strings.Repeat("=", 70))

	demo := NewAdvancedSearchDemo()

	// 1. 基础搜索演示
	demo.DemonstrateBasicSearch()

	// 2. MQE搜索演示
	demo.DemonstrateMQESearch()

	// 3. HyDE搜索演示
	demo.DemonstrateHyDESearch()

	// 4. 组合高级搜索演示
	demo.DemonstrateCombinedAdvancedSearch()

	// 5. 搜索性能分析
	demo.DemonstrateSearchPerformanceAnalysis()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎉 高级检索策略演示完成！")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n✨ 高级检索核心技术:")
	fmt.Println("1. 🔄 MQE多查询扩展 - 提高召回率和覆盖面")
	fmt.Println("2. 📝 HyDE假设文档嵌入 - 改善查询匹配效果")
	fmt.Println("3. 🚀 组合搜索策略 - 多技术融合优化")
	fmt.Println("4. 📊 智能结果排序 - 多因素评分机制")
	fmt.Println("5. ⚡ 性能优化 - 缓存和批量处理")

	fmt.Println("\n🎯 技术优势:")
	fmt.Println("• 语义理解 - 超越关键词匹配的语义检索")
	fmt.Println("• 查询扩展 - 自动生成相关查询提升召回")
	fmt.Println("• 上下文感知 - 理解查询意图和上下文")
	fmt.Println("• 质量优化 - 多策略融合提升检索质量")

	fmt.Println("\n💡 应用场景:")
	fmt.Println("• 技术文档问答 - 复杂技术问题的精准回答")
	fmt.Println("• 知识发现 - 从大量文档中发现相关知识")
	fmt.Println("• 智能搜索 - 理解用户意图的智能搜索")
	fmt.Println("• 内容推荐 - 基于语义相似度的内容推荐")
}
