// 代码示例 04: RAGTool的MarkItDown处理管道
// 展示Any格式→Markdown→分块→向量化的完整流程

package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	case "add_document":
		return "✅ 文档已处理并添加到知识库"
	case "add_text":
		docID, _ := params["document_id"].(string)
		return fmt.Sprintf("✅ 文本 '%s' 已分块并添加到知识库", docID)
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关文档片段", query)
	case "ask":
		question, _ := params["question"].(string)
		return fmt.Sprintf("💡 基于知识库回答 '%s': [相关答案内容]", question)
	case "stats":
		return "📊 统计: 文档数10, 分块数45, 向量维度384"
	default:
		return "操作完成"
	}
}

// BatchAddTexts 批量添加文本
func (r *RAGTool) BatchAddTexts(texts []string, documentIDs []string) string {
	return fmt.Sprintf("✅ 批量添加完成: %d个文档", len(texts))
}

// MarkItDownPipelineDemo 处理管道演示类
type MarkItDownPipelineDemo struct {
	RAGTool *RAGTool
	TempDir string
}

// NewMarkItDownPipelineDemo 创建演示实例
func NewMarkItDownPipelineDemo() *MarkItDownPipelineDemo {
	tempDir, _ := os.MkdirTemp("", "rag_demo_")
	return &MarkItDownPipelineDemo{
		RAGTool: NewRAGTool("./demo_rag_kb", "markitdown_demo"),
		TempDir: tempDir,
	}
}

// CreateSampleDocuments 创建多格式示例文档
func (d *MarkItDownPipelineDemo) CreateSampleDocuments() map[string]string {
	fmt.Println("📄 创建多格式示例文档")
	fmt.Println(strings.Repeat("=", 50))

	// 创建Markdown文档
	markdownContent := `# Python编程指南

## 基础语法
Python是一种解释型、高级编程语言。

### 变量和数据类型
- 整数：` + "`42`" + `
- 字符串：` + "`\"Hello World\"`" + `
- 列表：` + "`[1, 2, 3]`" + `

### 函数定义
` + "```python" + `
def greet(name):
    return f"Hello, {name}!"
` + "```" + `

## 面向对象编程
Python支持面向对象编程范式。

### 类定义
` + "```python" + `
class Person:
    def __init__(self, name):
        self.name = name
    
    def say_hello(self):
        return f"Hello, I'm {self.name}"
` + "```" + `
`

	// 创建HTML文档
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Web开发基础</title>
</head>
<body>
    <h1>HTML基础</h1>
    <p>HTML是超文本标记语言，用于创建网页结构。</p>
    
    <h2>常用标签</h2>
    <ul>
        <li>h1-h6: 标题标签</li>
        <li>p: 段落标签</li>
        <li>div: 容器标签</li>
        <li>span: 行内标签</li>
    </ul>
    
    <h2>CSS样式</h2>
    <p>CSS用于控制网页的样式和布局。</p>
</body>
</html>`

	// 创建JSON文档
	jsonContent := `{
    "project": "HelloAgents",
    "version": "1.0.0",
    "description": "AI Agent开发框架",
    "features": [
        "记忆系统",
        "RAG检索",
        "工具集成",
        "多模态支持"
    ]
}`

	// 创建CSV文档
	csvContent := `名称,类型,重要性,描述
工作记忆,临时存储,0.7,存储当前会话的临时信息
情景记忆,事件记录,0.8,记录具体的事件和经历
语义记忆,知识存储,0.9,存储概念性知识和规则
感知记忆,多模态,0.6,处理图像音频等感知数据
向量检索,技术组件,0.8,基于语义相似度的检索
知识图谱,技术组件,0.9,实体关系的结构化表示`

	// 保存文档到临时目录
	documents := map[string]string{
		"python_guide.md":  markdownContent,
		"web_basics.html":  htmlContent,
		"project_info.json": jsonContent,
		"memory_types.csv":  csvContent,
	}

	filePaths := make(map[string]string)
	for filename, content := range documents {
		filePath := filepath.Join(d.TempDir, filename)
		os.WriteFile(filePath, []byte(content), 0644)
		filePaths[filename] = filePath
		fmt.Printf("✅ 创建文档: %s\n", filename)
	}

	return filePaths
}

// DemonstrateMarkItDownConversion 演示MarkItDown转换过程
func (d *MarkItDownPipelineDemo) DemonstrateMarkItDownConversion(filePaths map[string]string) map[string]map[string]any {
	fmt.Println("\n🔄 MarkItDown转换过程演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("MarkItDown处理流程:")
	fmt.Println("1. 📄 检测文档格式")
	fmt.Println("2. 🔄 转换为Markdown")
	fmt.Println("3. 📝 保持结构信息")
	fmt.Println("4. ✨ 统一格式输出")

	conversionResults := make(map[string]map[string]any)

	for filename, filePath := range filePaths {
		fmt.Printf("\n处理文档: %s\n", filename)
		fmt.Printf("原始格式: %s\n", filepath.Ext(filename))

		startTime := time.Now()

		// 使用RAGTool添加文档，内部会调用MarkItDown
		result := d.RAGTool.Run(map[string]any{
			"action":    "add_document",
			"file_path": filePath,
		})

		processTime := time.Since(startTime)

		fmt.Printf("处理结果: %s\n", result)
		fmt.Printf("处理时间: %.3f秒\n", processTime.Seconds())
		fmt.Printf("✅ %s → Markdown → 分块 → 向量化\n", filename)

		conversionResults[filename] = map[string]any{
			"result": result,
			"time":   processTime.Seconds(),
		}
	}

	return conversionResults
}

// DemonstrateMarkdownChunking 演示基于Markdown的智能分块
func (d *MarkItDownPipelineDemo) DemonstrateMarkdownChunking() {
	fmt.Println("\n📊 基于Markdown的智能分块演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("Markdown分块策略:")
	fmt.Println("• 🏷️ 标题层次感知 - 利用#、##、###结构")
	fmt.Println("• 📝 段落语义保持 - 保持内容完整性")
	fmt.Println("• 🔢 Token精确控制 - 适配嵌入模型")
	fmt.Println("• 🔗 智能重叠策略 - 避免信息丢失")

	// 添加一个复杂的Markdown文档来演示分块
	complexMarkdown := `# 人工智能技术栈

## 机器学习基础

### 监督学习
监督学习使用标注数据训练模型，包括分类和回归任务。

#### 分类算法
- 逻辑回归：用于二分类和多分类问题
- 决策树：基于特征分割的树形结构
- 随机森林：多个决策树的集成方法
- 支持向量机：寻找最优分离超平面

#### 回归算法
- 线性回归：建立特征与目标的线性关系
- 多项式回归：处理非线性关系
- 岭回归：添加L2正则化的线性回归

### 无监督学习
无监督学习从无标注数据中发现模式和结构。

#### 聚类算法
- K-means：基于距离的聚类方法
- 层次聚类：构建聚类树状结构
- DBSCAN：基于密度的聚类算法

#### 降维算法
- PCA：主成分分析，线性降维
- t-SNE：非线性降维，适合可视化
- UMAP：保持局部和全局结构的降维

## 深度学习

### 神经网络基础
神经网络是深度学习的基础，模拟人脑神经元结构。

### 常见架构
- CNN：卷积神经网络，适合图像处理
- RNN：循环神经网络，处理序列数据
- LSTM：长短期记忆网络，解决梯度消失
- Transformer：注意力机制，处理长序列

## 自然语言处理

### 文本预处理
- 分词：将文本分割为词汇单元
- 词性标注：识别词汇的语法角色
- 命名实体识别：提取人名、地名等实体
- 情感分析：判断文本的情感倾向

### 语言模型
- N-gram：基于统计的语言模型
- Word2Vec：词向量表示学习
- BERT：双向编码器表示
- GPT：生成式预训练模型
`

	fmt.Printf("\n📝 添加复杂Markdown文档进行分块测试...\n")
	result := d.RAGTool.Run(map[string]any{
		"action":        "add_text",
		"text":          complexMarkdown,
		"document_id":   "ai_tech_stack",
		"chunk_size":    800,
		"chunk_overlap": 100,
	})

	fmt.Printf("分块结果: %s\n", result)

	// 测试基于结构的检索
	fmt.Printf("\n🔍 测试基于Markdown结构的检索:\n")

	type searchQuery struct {
		Query       string
		Description string
	}

	searchQueries := []searchQuery{
		{"监督学习算法", "测试二级标题内容检索"},
		{"神经网络基础", "测试跨层级内容检索"},
		{"BERT GPT", "测试具体技术检索"},
		{"聚类降维", "测试相关概念检索"},
	}

	for _, sq := range searchQueries {
		fmt.Printf("\n查询: '%s' (%s)\n", sq.Query, sq.Description)
		searchResult := d.RAGTool.Run(map[string]any{
			"action": "search",
			"query":  sq.Query,
			"limit":  2,
		})
		if len(searchResult) > 200 {
			fmt.Printf("检索结果: %s...\n", searchResult[:200])
		} else {
			fmt.Printf("检索结果: %s\n", searchResult)
		}
	}
}

// DemonstrateEmbeddingOptimization 演示面向嵌入的Markdown预处理
func (d *MarkItDownPipelineDemo) DemonstrateEmbeddingOptimization() {
	fmt.Println("\n🎯 面向嵌入的Markdown预处理演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("Markdown预处理优化:")
	fmt.Println("• 🏷️ 移除格式标记，保留语义内容")
	fmt.Println("• 🔗 处理链接格式，保留链接文本")
	fmt.Println("• 💻 清理代码块，保留代码内容")
	fmt.Println("• 🧹 清理多余空白，优化向量表示")

	// 演示预处理前后的对比
	rawMarkdown := `## 代码示例

这是一个**重要的**Python函数：

` + "```python" + `
def process_data(data):
    """处理数据的函数"""
    return [item.strip() for item in data if item]
` + "```" + `

更多信息请参考[官方文档](https://docs.python.org)。

*注意*：这个函数会` + "`自动过滤`" + `空值。
`

	fmt.Printf("\n📝 原始Markdown内容:\n")
	fmt.Println(rawMarkdown)

	// 添加到RAG系统，内部会进行预处理
	result := d.RAGTool.Run(map[string]any{
		"action":      "add_text",
		"text":        rawMarkdown,
		"document_id": "preprocessing_demo",
	})

	fmt.Printf("\n✅ 预处理并添加完成: %s\n", result)

	// 测试预处理后的检索效果
	fmt.Printf("\n🔍 测试预处理后的检索效果:\n")
	searchResult := d.RAGTool.Run(map[string]any{
		"action": "search",
		"query":  "Python函数处理数据",
		"limit":  1,
	})
	fmt.Printf("检索结果: %s\n", searchResult)
}

// DemonstratePipelinePerformance 演示处理管道性能
func (d *MarkItDownPipelineDemo) DemonstratePipelinePerformance() {
	fmt.Println("\n⚡ 处理管道性能演示")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("性能测试指标:")
	fmt.Println("• 📄 文档转换速度")
	fmt.Println("• 📊 分块处理效率")
	fmt.Println("• 🎯 向量化时间")
	fmt.Println("• 💾 存储操作耗时")

	// 批量处理性能测试
	batchTexts := make([]string, 10)
	batchDocIDs := make([]string, 10)
	for i := range 10 {
		batchTexts[i] = fmt.Sprintf("批量处理测试文档 %d：这是一个用于测试MarkItDown处理管道性能的示例文档。"+
			"文档包含了多种格式的内容，包括标题、段落、列表等结构化信息。"+
			"通过批量处理可以评估系统的整体性能表现。", i+1)
		batchDocIDs[i] = fmt.Sprintf("perf_test_%d", i+1)
	}

	fmt.Printf("\n⏱️ 批量处理性能测试 (10个文档):\n")
	startTime := time.Now()

	batchResult := d.RAGTool.BatchAddTexts(batchTexts, batchDocIDs)

	batchTime := time.Since(startTime)

	fmt.Printf("批量处理结果: %s\n", batchResult)
	fmt.Printf("总耗时: %.3f秒\n", batchTime.Seconds())
	fmt.Printf("平均每文档: %.3f秒\n", batchTime.Seconds()/10)

	// 获取最终统计
	stats := d.RAGTool.Run(map[string]any{"action": "stats"})
	fmt.Printf("\n📊 最终统计: %s\n", stats)
}

func main() {
	fmt.Println("🔄 RAGTool的MarkItDown处理管道演示")
	fmt.Println("展示Any格式→Markdown→分块→向量化的完整流程")
	fmt.Println(strings.Repeat("=", 70))

	demo := NewMarkItDownPipelineDemo()

	// 1. 创建多格式示例文档
	filePaths := demo.CreateSampleDocuments()

	// 2. 演示MarkItDown转换过程
	demo.DemonstrateMarkItDownConversion(filePaths)

	// 3. 演示基于Markdown的智能分块
	demo.DemonstrateMarkdownChunking()

	// 4. 演示面向嵌入的预处理优化
	demo.DemonstrateEmbeddingOptimization()

	// 5. 演示处理管道性能
	demo.DemonstratePipelinePerformance()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎉 MarkItDown处理管道演示完成！")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n✨ 处理管道核心特性:")
	fmt.Println("1. 🔄 格式统一 - Any格式→Markdown标准化")
	fmt.Println("2. 🏗️ 结构保持 - 保留文档逻辑结构")
	fmt.Println("3. 📊 智能分块 - 基于Markdown结构的语义分割")
	fmt.Println("4. 🎯 嵌入优化 - 针对向量化的预处理")
	fmt.Println("5. ⚡ 高效处理 - 批量处理和性能优化")

	fmt.Println("\n🎯 技术优势:")
	fmt.Println("• 统一处理 - 一套流程处理所有格式")
	fmt.Println("• 结构感知 - 充分利用Markdown结构信息")
	fmt.Println("• 语义保持 - 在格式转换中保持语义完整性")
	fmt.Println("• 检索优化 - 为向量检索优化的文本表示")

	// 清理临时文件
	os.RemoveAll(demo.TempDir)
	fmt.Printf("\n🧹 清理临时文件: %s\n", demo.TempDir)
}
