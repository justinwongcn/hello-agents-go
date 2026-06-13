// 智能文档问答助手 - 基于HelloAgents的智能文档问答系统
//
// 这是一个完整的PDF学习助手应用，支持：
// - 加载PDF文档并构建知识库
// - 智能问答（基于RAG）
// - 学习历程记录（基于Memory）
// - 学习回顾和报告生成

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// MemoryTool 模拟MemoryTool
type MemoryTool struct {
	UserID string
}

// NewMemoryTool 创建MemoryTool
func NewMemoryTool(userID string) *MemoryTool {
	return &MemoryTool{UserID: userID}
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
	RAGNamespace string
}

// NewRAGTool 创建RAGTool
func NewRAGTool(namespace string) *RAGTool {
	return &RAGTool{RAGNamespace: namespace}
}

// Run 执行RAG操作
func (r *RAGTool) Run(params map[string]any) string {
	action, _ := params["action"].(string)
	switch action {
	case "add_document":
		return "✅ 文档已处理并添加到知识库"
	case "add_text":
		return "✅ 文本已添加到知识库"
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关文档片段", query)
	case "ask":
		question, _ := params["question"].(string)
		return fmt.Sprintf("💡 问答 '%s': 基于RAG检索的答案", question)
	case "stats":
		return "📊 统计: 文档数1, 分块数12, 向量维度384"
	default:
		return "操作完成"
	}
}

// PDFLearningAssistant 智能文档问答助手
type PDFLearningAssistant struct {
	UserID       string
	SessionID    string
	MemoryTool   *MemoryTool
	RAGTool      *RAGTool
	Stats        map[string]any
	CurrentDoc   string
	StartTime    time.Time
}

// NewPDFLearningAssistant 初始化学习助手
func NewPDFLearningAssistant(userID string) *PDFLearningAssistant {
	now := time.Now()
	return &PDFLearningAssistant{
		UserID:    userID,
		SessionID: fmt.Sprintf("session_%s", now.Format("20060102_150405")),
		MemoryTool: NewMemoryTool(userID),
		RAGTool:    NewRAGTool(fmt.Sprintf("pdf_%s", userID)),
		Stats: map[string]any{
			"documents_loaded": 0,
			"questions_asked":  0,
			"concepts_learned": 0,
		},
		StartTime: now,
	}
}

// LoadDocument 加载PDF文档到知识库
func (a *PDFLearningAssistant) LoadDocument(pdfPath string) map[string]any {
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return map[string]any{"success": false, "message": fmt.Sprintf("文件不存在: %s", pdfPath)}
	}

	startTime := time.Now()

	// 使用RAG工具处理PDF
	a.RAGTool.Run(map[string]any{
		"action":        "add_document",
		"file_path":     pdfPath,
		"chunk_size":    1000,
		"chunk_overlap": 200,
	})

	processTime := time.Since(startTime)

	// 提取文件名
	parts := strings.Split(pdfPath, "/")
	a.CurrentDoc = parts[len(parts)-1]
	a.Stats["documents_loaded"] = a.Stats["documents_loaded"].(int) + 1

	// 记录到学习记忆
	a.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     fmt.Sprintf("加载了文档《%s》", a.CurrentDoc),
		"memory_type": "episodic",
		"importance":  0.9,
		"event_type":  "document_loaded",
		"session_id":  a.SessionID,
	})

	return map[string]any{
		"success":  true,
		"message":  fmt.Sprintf("加载成功！(耗时: %.1f秒)", processTime.Seconds()),
		"document": a.CurrentDoc,
	}
}

// Ask 向文档提问
func (a *PDFLearningAssistant) Ask(question string, useAdvancedSearch bool) string {
	if a.CurrentDoc == "" {
		return "⚠️ 请先加载文档！使用 LoadDocument() 方法加载PDF文档。"
	}

	// 记录问题到工作记忆
	a.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     fmt.Sprintf("提问: %s", question),
		"memory_type": "working",
		"importance":  0.6,
		"session_id":  a.SessionID,
	})

	// 使用RAG检索答案
	answer := a.RAGTool.Run(map[string]any{
		"action":                 "ask",
		"question":               question,
		"limit":                  5,
		"enable_advanced_search": useAdvancedSearch,
		"enable_mqe":             useAdvancedSearch,
		"enable_hyde":            useAdvancedSearch,
	})

	// 记录到情景记忆
	a.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     fmt.Sprintf("关于'%s'的学习", question),
		"memory_type": "episodic",
		"importance":  0.7,
		"event_type":  "qa_interaction",
		"session_id":  a.SessionID,
	})

	a.Stats["questions_asked"] = a.Stats["questions_asked"].(int) + 1

	return answer
}

// AddNote 添加学习笔记
func (a *PDFLearningAssistant) AddNote(content string, concept string) {
	if concept == "" {
		concept = "general"
	}

	a.MemoryTool.Run(map[string]any{
		"action":      "add",
		"content":     content,
		"memory_type": "semantic",
		"importance":  0.8,
		"concept":     concept,
		"session_id":  a.SessionID,
	})

	a.Stats["concepts_learned"] = a.Stats["concepts_learned"].(int) + 1
}

// Recall 回顾学习历程
func (a *PDFLearningAssistant) Recall(query string, limit int) string {
	return a.MemoryTool.Run(map[string]any{
		"action": "search",
		"query":  query,
		"limit":  limit,
	})
}

// GetStats 获取学习统计
func (a *PDFLearningAssistant) GetStats() map[string]any {
	duration := time.Since(a.StartTime).Seconds()

	currentDoc := "未加载"
	if a.CurrentDoc != "" {
		currentDoc = a.CurrentDoc
	}

	return map[string]any{
		"会话时长": fmt.Sprintf("%.0f秒", duration),
		"加载文档": a.Stats["documents_loaded"],
		"提问次数": a.Stats["questions_asked"],
		"学习笔记": a.Stats["concepts_learned"],
		"当前文档": currentDoc,
	}
}

// GenerateReport 生成学习报告
func (a *PDFLearningAssistant) GenerateReport(saveToFile bool) map[string]any {
	// 获取记忆摘要
	memorySummary := a.MemoryTool.Run(map[string]any{"action": "summary", "limit": 10})

	// 获取RAG统计
	ragStats := a.RAGTool.Run(map[string]any{"action": "stats"})

	// 生成报告
	duration := time.Since(a.StartTime).Seconds()
	report := map[string]any{
		"session_info": map[string]any{
			"session_id":      a.SessionID,
			"user_id":         a.UserID,
			"start_time":      a.StartTime.Format(time.RFC3339),
			"duration_seconds": duration,
		},
		"learning_metrics": map[string]any{
			"documents_loaded": a.Stats["documents_loaded"],
			"questions_asked":  a.Stats["questions_asked"],
			"concepts_learned": a.Stats["concepts_learned"],
		},
		"memory_summary": memorySummary,
		"rag_status":     ragStats,
	}

	// 保存到文件
	if saveToFile {
		reportFile := fmt.Sprintf("learning_report_%s.json", a.SessionID)
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err == nil {
			err = os.WriteFile(reportFile, jsonData, 0644)
			if err == nil {
				report["report_file"] = reportFile
			} else {
				report["save_error"] = err.Error()
			}
		}
	}

	return report
}

func main() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("智能文档问答助手")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("基于HelloAgents的智能文档问答系统")
	fmt.Println("支持: PDF加载 | 智能问答 | 学习笔记 | 学习回顾 | 报告生成")
	fmt.Println()

	// 创建学习助手
	assistant := NewPDFLearningAssistant("demo_user")
	fmt.Printf("✅ 助手已初始化 (用户: %s, 会话: %s)\n", assistant.UserID, assistant.SessionID)

	// 模拟加载文档
	fmt.Println("\n📄 加载文档演示:")
	result := assistant.LoadDocument("./example.pdf")
	if success, ok := result["success"].(bool); ok && success {
		fmt.Printf("✅ %s\n", result["message"])
		fmt.Printf("📄 文档: %s\n", result["document"])
	} else {
		fmt.Printf("❌ %s\n", result["message"])
	}

	// 模拟智能问答
	fmt.Println("\n💬 智能问答演示:")
	questions := []string{
		"什么是大语言模型？",
		"Transformer架构有哪些核心组件？",
		"如何训练大语言模型？",
	}

	for _, question := range questions {
		fmt.Printf("\n❓ 问题: %s\n", question)
		answer := assistant.Ask(question, true)
		fmt.Printf("💡 回答: %s\n", answer)
	}

	// 添加学习笔记
	fmt.Println("\n📝 学习笔记演示:")
	assistant.AddNote("Transformer的核心是自注意力机制，它允许模型关注输入序列的不同部分", "transformer")
	assistant.AddNote("大语言模型通过预训练和微调两个阶段进行训练", "llm_training")
	fmt.Println("✅ 笔记已保存")

	// 回顾学习历程
	fmt.Println("\n🧠 学习回顾演示:")
	recallResult := assistant.Recall("Transformer", 5)
	fmt.Printf("回顾结果: %s\n", recallResult)

	// 获取统计信息
	fmt.Println("\n📊 学习统计:")
	stats := assistant.GetStats()
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// 生成学习报告
	fmt.Println("\n📋 生成学习报告:")
	report := assistant.GenerateReport(false)

	fmt.Println("✅ 学习报告已生成")
	fmt.Printf("\n**会话信息**\n")
	if sessionInfo, ok := report["session_info"].(map[string]any); ok {
		fmt.Printf("  - 会话时长: %.0f秒\n", sessionInfo["duration_seconds"])
	}
	if metrics, ok := report["learning_metrics"].(map[string]any); ok {
		fmt.Printf("  - 加载文档: %v\n", metrics["documents_loaded"])
		fmt.Printf("  - 提问次数: %v\n", metrics["questions_asked"])
		fmt.Printf("  - 学习笔记: %v\n", metrics["concepts_learned"])
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 智能文档问答助手演示完成！")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n✨ 核心功能:")
	fmt.Println("1. 📄 PDF文档加载 - 支持PDF文件的自动解析和索引")
	fmt.Println("2. 💬 智能问答 - 基于RAG的精准问答")
	fmt.Println("3. 📝 学习笔记 - 记录学习心得和重要概念")
	fmt.Println("4. 🧠 学习回顾 - 回顾学习历程和已学内容")
	fmt.Println("5. 📊 学习报告 - 生成详细的学习统计报告")

	fmt.Println("\n🎯 应用场景:")
	fmt.Println("• 学术论文阅读助手")
	fmt.Println("• 技术文档学习工具")
	fmt.Println("• 课程资料问答系统")
	fmt.Println("• 企业知识库查询")
}
