// CodebaseMaintainer 三天工作流演示
//
// 完整展示长程智能体在三天内的工作流程:
// - 第一天: 探索代码库（Agent 自主探索）
// - 第二天: 分析代码质量（Agent 自主分析）
// - 第三天: 规划重构任务（Agent 自主规划）
// - 一周后: 检查进度

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Note 笔记
type Note struct {
	NoteID    string
	Title     string
	Content   string
	NoteType  string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NoteTool 笔记工具
type NoteTool struct {
	Workspace string
	Notes     map[string]*Note
	Counter   int
}

// NewNoteTool 创建 NoteTool
func NewNoteTool(workspace string) *NoteTool {
	return &NoteTool{Workspace: workspace, Notes: make(map[string]*Note), Counter: 0}
}

// Run 执行笔记操作
func (nt *NoteTool) Run(params map[string]any) any {
	action, _ := params["action"].(string)
	switch action {
	case "create":
		nt.Counter++
		noteID := fmt.Sprintf("note_%d_%d", time.Now().Unix(), nt.Counter)
		title, _ := params["title"].(string)
		content, _ := params["content"].(string)
		noteType, _ := params["note_type"].(string)
		if noteType == "" {
			noteType = "general"
		}
		var tags []string
		if tagSlice, ok := params["tags"].([]string); ok {
			tags = tagSlice
		}
		now := time.Now()
		nt.Notes[noteID] = &Note{
			NoteID: noteID, Title: title, Content: content,
			NoteType: noteType, Tags: tags, CreatedAt: now, UpdatedAt: now,
		}
		return map[string]any{"note_id": noteID, "title": title, "type": noteType}
	case "list":
		noteType, _ := params["note_type"].(string)
		limit := 10
		if l, ok := params["limit"].(int); ok {
			limit = l
		}
		var items []map[string]any
		count := 0
		for _, note := range nt.Notes {
			if count >= limit {
				break
			}
			if noteType == "" || note.NoteType == noteType {
				items = append(items, map[string]any{
					"note_id": note.NoteID, "title": note.Title,
					"type": note.NoteType, "content": note.Content,
					"updated_at": note.UpdatedAt.Format(time.RFC3339),
				})
				count++
			}
		}
		return map[string]any{"items": items}
	case "search":
		query, _ := params["query"].(string)
		limit := 10
		if l, ok := params["limit"].(int); ok {
			limit = l
		}
		var items []map[string]any
		count := 0
		for _, note := range nt.Notes {
			if count >= limit {
				break
			}
			if strings.Contains(note.Title, query) || strings.Contains(note.Content, query) {
				items = append(items, map[string]any{
					"note_id": note.NoteID, "title": note.Title,
					"type": note.NoteType, "content": note.Content,
					"updated_at": note.UpdatedAt.Format(time.RFC3339),
				})
				count++
			}
		}
		return map[string]any{"items": items}
	case "summary":
		typeCount := make(map[string]int)
		for _, note := range nt.Notes {
			typeCount[note.NoteType]++
		}
		return map[string]any{"total": len(nt.Notes), "by_type": typeCount}
	default:
		return "操作完成"
	}
}

// HelloAgentsLLM 模拟 LLM
type HelloAgentsLLM struct{}

// TerminalTool 终端工具
type TerminalTool struct {
	Workspace string
}

// Run 执行终端命令
func (tt *TerminalTool) Run(params map[string]any) string {
	command, _ := params["command"].(string)
	return fmt.Sprintf("[模拟执行] $ %s\n(命令输出已省略)", command)
}

// ContextConfig 上下文配置
type ContextConfig struct {
	MaxTokens         int
	ReserveRatio      float64
	MinRelevance      float64
	EnableCompression bool
}

// Message 消息
type Message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

// ContextPacket 上下文包
type ContextPacket struct {
	Content        string
	Timestamp      time.Time
	TokenCount     int
	RelevanceScore float64
	Metadata       map[string]any
}

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	MemoryTool any
	RAGTool    any
	Config     *ContextConfig
}

// NewContextBuilder 创建 ContextBuilder
func NewContextBuilder(memoryTool, ragTool any, config *ContextConfig) *ContextBuilder {
	return &ContextBuilder{MemoryTool: memoryTool, RAGTool: ragTool, Config: config}
}

// Build 构建上下文
func (cb *ContextBuilder) Build(userQuery string, conversationHistory []Message, systemInstructions string, additionalPackets []ContextPacket) string {
	var sb strings.Builder
	sb.WriteString(systemInstructions)
	sb.WriteString("\n\nUser Query: ")
	sb.WriteString(userQuery)
	return sb.String()
}

// CodebaseMaintainer 代码库维护助手
type CodebaseMaintainer struct {
	ProjectName       string
	CodebasePath      string
	SessionID         string
	LLM               *HelloAgentsLLM
	MemoryTool        any
	NoteTool          *NoteTool
	TerminalTool      *TerminalTool
	ContextBuilder    *ContextBuilder
	ConversationHistory []Message
	Stats             map[string]any
}

// NewCodebaseMaintainer 创建 CodebaseMaintainer
func NewCodebaseMaintainer(projectName, codebasePath string, llm *HelloAgentsLLM) *CodebaseMaintainer {
	sessionID := fmt.Sprintf("session_%s", time.Now().Format("20060102_150405"))
	noteTool := NewNoteTool(fmt.Sprintf("./%s_notes", projectName))
	terminalTool := &TerminalTool{Workspace: codebasePath}

	config := &ContextConfig{
		MaxTokens:         4000,
		ReserveRatio:      0.15,
		MinRelevance:      0.2,
		EnableCompression: true,
	}
	builder := NewContextBuilder(nil, nil, config)

	cm := &CodebaseMaintainer{
		ProjectName:    projectName,
		CodebasePath:   codebasePath,
		SessionID:      sessionID,
		LLM:            llm,
		NoteTool:       noteTool,
		TerminalTool:   terminalTool,
		ContextBuilder: builder,
		Stats: map[string]any{
			"commands_executed": 0,
			"notes_created":     0,
			"issues_found":      0,
			"tool_calls":        0,
		},
	}

	fmt.Printf("✅ 代码库维护助手已初始化: %s (Agentic Mode)\n", projectName)
	fmt.Printf("📁 工作目录: %s\n", codebasePath)
	fmt.Printf("🆔 会话ID: %s\n", sessionID)
	fmt.Printf("🔧 可用工具: terminal, note, memory\n")

	return cm
}

// Run 运行助手（Agentic 方式）
func (cm *CodebaseMaintainer) Run(userInput string, mode ...string) string {
	m := "auto"
	if len(mode) > 0 {
		m = mode[0]
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 80))
	fmt.Printf("👤 用户: %s\n", userInput)
	fmt.Printf("%s\n\n", strings.Repeat("=", 80))

	// 模拟 Agent 自主决策和使用工具
	fmt.Println("🤖 Agent 正在思考并决定使用哪些工具...\n")

	var response string
	switch {
	case strings.Contains(userInput, "探索") || strings.Contains(userInput, "代码结构"):
		response = fmt.Sprintf("我已经探索了 %s 的代码结构。项目包含多个Python模块，主要文件有: data_processor.py, api_client.py, utils.py, models.py。项目结构清晰，采用了模块化设计。", cm.CodebasePath)
	case strings.Contains(userInput, "data_processor"):
		response = "data_processor.py 分析结果:\n- 代码行数: ~150行\n- 主要功能: 数据清洗和转换\n- 设计: 使用了策略模式处理不同类型的数据源\n- 建议: 可以考虑将一些通用方法提取到基类中"
	case strings.Contains(userInput, "质量") || m == "analyze":
		response = "代码质量分析:\n- TODO注释: 5处\n- 代码复杂度: 中等\n- 错误处理: 部分函数缺少异常处理\n- 建议: 增加单元测试覆盖，完善错误处理"
	case strings.Contains(userInput, "api_client"):
		response = "api_client.py 质量分析:\n- 错误处理: 使用了基本的try-except，但缺少重试机制\n- 建议: 1) 添加请求重试 2) 增加超时配置 3) 添加请求日志"
	case strings.Contains(userInput, "规划") || strings.Contains(userInput, "计划") || m == "plan":
		response = "基于之前的分析，下一步计划:\n1. 完善api_client的错误处理（优先级: 高）\n2. 实现TODO注释中标记的功能（优先级: 中）\n3. 增加单元测试覆盖率（优先级: 中）\n4. 优化数据处理性能（优先级: 低）"
	case strings.Contains(userInput, "TODO"):
		response = "代码库TODO项分析:\n- data_processor.py: 3个TODO (数据验证、缓存机制、性能优化)\n- api_client.py: 2个TODO (重试机制、日志记录)\n- utils.py: 1个TODO (输入验证)\n建议优先实现: api_client的重试机制"
	default:
		response = fmt.Sprintf("基于对 %s 项目的分析，我建议从代码质量改进开始。", cm.ProjectName)
	}

	cm.ConversationHistory = append(cm.ConversationHistory,
		Message{Content: userInput, Role: "user", Timestamp: time.Now()},
		Message{Content: response, Role: "assistant", Timestamp: time.Now()},
	)

	fmt.Printf("\n🤖 助手: %s\n", response)
	fmt.Printf("%s\n", strings.Repeat("=", 80))

	return response
}

// Explore 探索代码库
func (cm *CodebaseMaintainer) Explore() string {
	return cm.Run(fmt.Sprintf("请探索 %s 的代码结构，了解项目组织方式", "."), "explore")
}

// Analyze 分析代码质量
func (cm *CodebaseMaintainer) Analyze(focus ...string) string {
	query := "请分析代码质量"
	if len(focus) > 0 && focus[0] != "" {
		query += "，重点关注" + focus[0]
	}
	return cm.Run(query, "analyze")
}

// PlanNextSteps 规划下一步任务
func (cm *CodebaseMaintainer) PlanNextSteps() string {
	return cm.Run("根据我们之前的分析和当前进度，规划下一步任务", "plan")
}

// ExecuteCommand 执行终端命令
func (cm *CodebaseMaintainer) ExecuteCommand(command string) string {
	result := cm.TerminalTool.Run(map[string]any{"command": command})
	cm.Stats["commands_executed"] = cm.Stats["commands_executed"].(int) + 1
	return result
}

// CreateNote 创建笔记
func (cm *CodebaseMaintainer) CreateNote(title, content, noteType string, tags []string) any {
	if tags == nil {
		tags = []string{cm.ProjectName}
	}
	result := cm.NoteTool.Run(map[string]any{
		"action":    "create",
		"title":     title,
		"content":   content,
		"note_type": noteType,
		"tags":      tags,
	})
	cm.Stats["notes_created"] = cm.Stats["notes_created"].(int) + 1
	return result
}

// GetStats 获取统计信息
func (cm *CodebaseMaintainer) GetStats() map[string]any {
	noteSummary := cm.NoteTool.Run(map[string]any{"action": "summary"})
	return map[string]any{
		"session_info": map[string]any{
			"session_id": cm.SessionID,
			"project":    cm.ProjectName,
		},
		"activity": cm.Stats,
		"notes":    noteSummary,
	}
}

// GenerateReport 生成会话报告
func (cm *CodebaseMaintainer) GenerateReport(saveToFile ...bool) map[string]any {
	report := cm.GetStats()
	reportFile := fmt.Sprintf("maintainer_report_%s.json", cm.SessionID)
	report["report_file"] = reportFile
	fmt.Printf("📄 报告已保存: %s\n", reportFile)
	return report
}

func day1Exploration(maintainer *CodebaseMaintainer) {
	// 第一天: 探索代码库（Agentic 方式）
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("第一天: 探索代码库（Agent 自主探索）")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	// 1. 初步探索 - Agent 自主决定如何探索
	fmt.Println("### 1. 初步探索项目结构 ###")
	fmt.Println("💡 提示：Agent 会自主决定使用哪些命令（如 find, ls, cat）\n")
	response := maintainer.Explore()
	fmt.Printf("\n助手总结:\n%s...\n", truncate(response, 500))

	// 2. 深入分析某个模块 - Agent 自主决定分析方法
	fmt.Println("\n### 2. 分析数据处理模块 ###")
	fmt.Println("💡 提示：Agent 会自主决定如何分析这个文件\n")
	response = maintainer.Run("请查看 data_processor.py 文件，分析其代码设计")
	fmt.Printf("\n助手总结:\n%s...\n", truncate(response, 500))

	// 模拟时间流逝
	time.Sleep(1 * time.Second)
}

func day2Analysis(maintainer *CodebaseMaintainer) {
	// 第二天: 分析代码质量（Agentic 方式）
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("第二天: 分析代码质量（Agent 自主分析）")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	// 1. 整体质量分析 - Agent 自主决定分析方法
	fmt.Println("### 1. 分析代码质量 ###")
	fmt.Println("💡 提示：Agent 会自主决定如何分析（如 grep TODO, wc -l, 复杂度分析）\n")
	response := maintainer.Analyze()
	fmt.Printf("\n助手总结:\n%s...\n", truncate(response, 500))

	// 2. 查看具体问题 - Agent 自主深入分析
	fmt.Println("\n### 2. 分析 API 客户端代码 ###")
	fmt.Println("💡 提示：Agent 会自主决定如何分析这个文件的质量\n")
	response = maintainer.Run(
		"请分析 api_client.py 的代码质量，特别是错误处理部分，给出改进建议",
	)
	fmt.Printf("\n助手总结:\n%s...\n", truncate(response, 500))

	// 模拟时间流逝
	time.Sleep(1 * time.Second)
}

func day3Planning(maintainer *CodebaseMaintainer) {
	// 第三天: 规划重构任务（Agentic 方式）
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("第三天: 规划重构任务（Agent 自主规划）")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	// 1. 回顾进度 - Agent 自主查看历史笔记并规划
	fmt.Println("### 1. 回顾当前进度并规划下一步 ###")
	fmt.Println("💡 提示：Agent 会自主查看历史笔记，分析当前进度，并制定计划\n")
	response := maintainer.PlanNextSteps()
	fmt.Printf("\n助手总结:\n%s...\n", truncate(response, 500))

	// 2. 询问 Agent 创建详细计划（Agent 会自主决定是否使用 NoteTool）
	fmt.Println("\n### 2. 让 Agent 创建详细的重构计划 ###")
	fmt.Println("💡 提示：Agent 会自主决定如何创建和组织重构计划\n")
	response = maintainer.Run(
		"请基于我们的分析，创建一个详细的本周重构计划。" +
			"计划应该包括：目标、具体任务清单、时间安排和风险。" +
			"请使用 NoteTool 创建一个 task_state 类型的笔记来记录这个计划。",
	)
	fmt.Printf("\n助手总结:\n%s...\n", truncate(response, 500))

	// 模拟时间流逝
	time.Sleep(1 * time.Second)
}

func weekLaterReview(maintainer *CodebaseMaintainer) {
	// 一周后: 检查进度
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("一周后: 检查进度")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	// 1. 查看笔记摘要
	fmt.Println("### 1. 笔记摘要 ###")
	summary := maintainer.NoteTool.Run(map[string]any{"action": "summary"})
	fmt.Println("📊 笔记摘要:")
	summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(summaryJSON))
	fmt.Println()

	// 2. 生成完整报告
	fmt.Println("### 2. 会话报告 ###")
	report := maintainer.GenerateReport()
	fmt.Println("\n📄 会话报告:")
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(reportJSON))
}

func demonstrateCrossSessionContinuity() {
	// 演示跨会话的连贯性
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("演示跨会话的连贯性")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	llm := &HelloAgentsLLM{}

	// 第一次会话
	fmt.Println("### 第一次会话 (session_1) ###")
	maintainer1 := NewCodebaseMaintainer(
		"demo_codebase",
		"./codebase",
		llm,
	)

	// 创建一些笔记
	maintainer1.CreateNote(
		"代码质量问题",
		"发现多处 TODO 注释需要实现，特别是数据验证和错误处理部分",
		"blocker",
		[]string{"quality", "urgent"},
	)

	stats1 := maintainer1.GetStats()
	fmt.Printf("会话1统计: %v\n\n", stats1["activity"])

	// 模拟会话结束
	time.Sleep(1 * time.Second)

	// 第二次会话 (新的会话ID,但笔记被保留)
	fmt.Println("### 第二次会话 (session_2) ###")
	maintainer2 := NewCodebaseMaintainer(
		"demo_codebase", // 同一个项目
		"./codebase",
		llm,
	)

	// 检索之前的笔记
	response := maintainer2.Run(
		"我们之前发现了什么代码质量问题？现在应该优先处理哪些？",
	)
	fmt.Printf("\n助手回答:\n%s...\n", truncate(response, 300))

	stats2 := maintainer2.GetStats()
	fmt.Printf("会话2统计: %v\n\n", stats2["activity"])

	// 展示笔记摘要
	summary := maintainer2.NoteTool.Run(map[string]any{"action": "summary"})
	fmt.Println("📊 跨会话笔记摘要:")
	summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(summaryJSON))
}

func demonstrateToolSynergy() {
	// 演示三大工具的协同（Agentic 方式）
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("演示三大工具的协同（Agent 自主协调）")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	llm := &HelloAgentsLLM{}
	maintainer := NewCodebaseMaintainer(
		"synergy_demo",
		"./codebase",
		llm,
	)

	// Agent 自主分析并记录
	fmt.Println("### Agent 自主分析代码库中的 TODO 项 ###")
	fmt.Println("💡 提示：Agent 会自主决定：")
	fmt.Println("   1. 使用 TerminalTool 查找 TODO")
	fmt.Println("   2. 使用 NoteTool 记录发现")
	fmt.Println("   3. 使用 MemoryTool 记住关键信息\n")

	response := maintainer.Run(
		"请分析代码库中的所有 TODO 项，并将发现记录到笔记中。" +
			"然后告诉我应该优先实现哪些功能。",
	)
	fmt.Printf("助手回答:\n%s...\n", truncate(response, 500))

	// 展示统计信息
	stats := maintainer.GetStats()
	fmt.Println("\n📊 工具使用统计:")
	fmt.Printf("  - 工具调用次数: %v\n", stats["activity"].(map[string]any)["tool_calls"])
	fmt.Printf("  - 执行的命令: %v\n", stats["activity"].(map[string]any)["commands_executed"])
	fmt.Printf("  - 创建的笔记: %v\n", stats["activity"].(map[string]any)["notes_created"])
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}

func main() {
	// 主函数
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("CodebaseMaintainer 三天工作流演示（Agentic 版本）")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n✨ 核心特性：Agent 自主决策")
	fmt.Println("💡 使用我们在 chapter9 创建的示例代码库")
	fmt.Println("📁 代码库路径: ./codebase")
	fmt.Println("📦 包含文件: data_processor.py, api_client.py, utils.py, models.py")
	fmt.Println("\n🔧 Agent 可用工具：")
	fmt.Println("   - TerminalTool: 执行 shell 命令")
	fmt.Println("   - NoteTool: 创建和管理笔记")
	fmt.Println("   - MemoryTool: 记忆管理")
	fmt.Println("\n⚡ Agent 会自主决定：")
	fmt.Println("   - 使用哪些工具")
	fmt.Println("   - 执行什么命令")
	fmt.Println("   - 如何组织信息\n")

	// 初始化助手
	llm := &HelloAgentsLLM{}
	maintainer := NewCodebaseMaintainer(
		"demo_codebase",
		"./codebase",
		llm,
	)

	// 执行三天工作流
	day1Exploration(maintainer)
	day2Analysis(maintainer)
	day3Planning(maintainer)
	weekLaterReview(maintainer)

	// 额外演示
	fmt.Println("\n\n" + strings.Repeat("=", 80))
	fmt.Println("额外演示")
	fmt.Println(strings.Repeat("=", 80))

	demonstrateCrossSessionContinuity()
	demonstrateToolSynergy()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("完整演示结束!")
	fmt.Println(strings.Repeat("=", 80))
}
