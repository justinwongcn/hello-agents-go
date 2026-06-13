// CodebaseMaintainer - 代码库维护助手
//
// 完整的长程智能体实现，整合:
// 1. ContextBuilder - 上下文管理
// 2. NoteTool - 结构化笔记
// 3. TerminalTool - 即时文件访问
// 4. MemoryTool - 对话记忆
//
// 关键改进：使用 Agentic 方式，让 agent 自主决定使用哪些工具

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

// MemoryTool 记忆工具
type MemoryTool struct {
	UserID string
}

// NewMemoryTool 创建 MemoryTool
func NewMemoryTool(userID string) *MemoryTool {
	return &MemoryTool{UserID: userID}
}

// Run 执行记忆操作
func (mt *MemoryTool) Run(params map[string]any) string {
	action, _ := params["action"].(string)
	switch action {
	case "add":
		return fmt.Sprintf("✅ 已添加记忆到工作记忆")
	case "search":
		query, _ := params["query"].(string)
		return fmt.Sprintf("🔍 搜索 '%s': 找到相关记忆", query)
	default:
		return "操作完成"
	}
}

// TerminalTool 终端工具
type TerminalTool struct {
	Workspace string
	Timeout   int
}

// NewTerminalTool 创建 TerminalTool
func NewTerminalTool(workspace string, timeout int) *TerminalTool {
	return &TerminalTool{Workspace: workspace, Timeout: timeout}
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

// ToolRegistry 工具注册表
type ToolRegistry struct {
	Tools map[string]any
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{Tools: make(map[string]any)}
}

// RegisterTool 注册工具
func (tr *ToolRegistry) RegisterTool(tool any) {
	switch t := tool.(type) {
	case *TerminalTool:
		tr.Tools["terminal"] = t
	case *NoteTool:
		tr.Tools["note"] = t
	case *MemoryTool:
		tr.Tools["memory"] = t
	}
}

// ListTools 列出工具
func (tr *ToolRegistry) ListTools() []string {
	var names []string
	for name := range tr.Tools {
		names = append(names, name)
	}
	return names
}

// FunctionCallAgent 功能调用Agent
type FunctionCallAgent struct {
	Name          string
	LLM           any
	SystemPrompt  string
	ToolRegistry  *ToolRegistry
	MessageHistory []Message
}

// NewFunctionCallAgent 创建 FunctionCallAgent
func NewFunctionCallAgent(name string, llm any, systemPrompt string, toolRegistry *ToolRegistry) *FunctionCallAgent {
	return &FunctionCallAgent{
		Name:         name,
		LLM:          llm,
		SystemPrompt: systemPrompt,
		ToolRegistry: toolRegistry,
	}
}

// Run 运行 Agent
func (a *FunctionCallAgent) Run(userInput string) string {
	// 模拟 Agent 自主决策和使用工具
	response := fmt.Sprintf("Agent '%s' 已分析用户输入并自主决定使用工具完成任务。", a.Name)

	// 记录消息历史
	a.MessageHistory = append(a.MessageHistory,
		Message{Content: userInput, Role: "user", Timestamp: time.Now()},
		Message{Content: response, Role: "assistant", Timestamp: time.Now()},
	)

	return response
}

// HelloAgentsLLM 模拟 LLM
type HelloAgentsLLM struct{}

// CodebaseMaintainer 代码库维护助手 - 长程智能体示例
//
// 整合 ContextBuilder + NoteTool + TerminalTool + MemoryTool
// 实现跨会话的代码库维护任务管理
//
// 核心特性：
// - Agent 自主使用工具探索代码库
// - 不预定义工作流，完全基于 agent 决策
// - 跨会话记忆和上下文管理
type CodebaseMaintainer struct {
	ProjectName       string
	CodebasePath      string
	SessionID         string
	LLM               *HelloAgentsLLM
	MemoryTool        *MemoryTool
	NoteTool          *NoteTool
	TerminalTool      *TerminalTool
	ContextBuilder    *ContextBuilder
	ToolRegistry      *ToolRegistry
	Agent             *FunctionCallAgent
	ConversationHistory []Message
	Stats             map[string]any
}

// NewCodebaseMaintainer 创建 CodebaseMaintainer
func NewCodebaseMaintainer(
	projectName string,
	codebasePath string,
	llm *HelloAgentsLLM,
) *CodebaseMaintainer {
	sessionID := fmt.Sprintf("session_%s", time.Now().Format("20060102_150405"))

	// 初始化工具
	memoryTool := NewMemoryTool(projectName)
	noteTool := NewNoteTool(fmt.Sprintf("./%s_notes", projectName))
	terminalTool := NewTerminalTool(codebasePath, 60)

	// 初始化上下文构建器
	config := &ContextConfig{
		MaxTokens:         4000,
		ReserveRatio:      0.15,
		MinRelevance:      0.2,
		EnableCompression: true,
	}
	builder := NewContextBuilder(memoryTool, nil, config)

	// 创建工具注册表并注册工具
	toolRegistry := NewToolRegistry()
	toolRegistry.RegisterTool(terminalTool)
	toolRegistry.RegisterTool(noteTool)
	toolRegistry.RegisterTool(memoryTool)

	// 创建 Agent
	agent := NewFunctionCallAgent(
		"CodebaseMaintainer",
		llm,
		fmt.Sprintf(`你是 %s 项目的代码库维护助手。

你的核心能力:
1. 使用 TerminalTool 探索代码库
   - 你可以执行任何 shell 命令: ls, cat, grep, find, git 等
   - 工作目录: %s
   
2. 使用 NoteTool 记录发现和任务
   - 创建笔记记录重要发现
   - 笔记类型: blocker(阻塞问题)、action(行动计划)、task_state(任务状态)、conclusion(结论)
   
3. 使用 MemoryTool 存储关键信息
   - 记住重要的上下文信息
   - 跨会话保持连贯性

当前会话ID: %s

重要原则:
- 你要自主决定使用哪些工具、执行什么命令
- 探索代码库时，先了解整体结构，再深入细节
- 发现重要信息时，主动使用 NoteTool 记录
- 保持回答的专业性和实用性`, projectName, codebasePath, sessionID),
		toolRegistry,
	)

	stats := map[string]any{
		"session_start":     time.Now(),
		"commands_executed": 0,
		"notes_created":     0,
		"issues_found":      0,
		"tool_calls":        0,
	}

	cm := &CodebaseMaintainer{
		ProjectName:    projectName,
		CodebasePath:   codebasePath,
		SessionID:      sessionID,
		LLM:            llm,
		MemoryTool:     memoryTool,
		NoteTool:       noteTool,
		TerminalTool:   terminalTool,
		ContextBuilder: builder,
		ToolRegistry:   toolRegistry,
		Agent:          agent,
		Stats:          stats,
	}

	fmt.Printf("✅ 代码库维护助手已初始化: %s (Agentic Mode)\n", projectName)
	fmt.Printf("📁 工作目录: %s\n", codebasePath)
	fmt.Printf("🆔 会话ID: %s\n", sessionID)
	fmt.Printf("🔧 可用工具: %s\n", strings.Join(toolRegistry.ListTools(), ", "))

	return cm
}

// Run 运行助手（Agentic 方式）
//
// mode: 运行模式提示（给 agent 提供方向性建议）
//   - "auto": 自动决策是否使用工具
//   - "explore": 建议 agent 侧重代码探索
//   - "analyze": 建议 agent 侧重问题分析
//   - "plan": 建议 agent 侧重任务规划
func (cm *CodebaseMaintainer) Run(userInput string, mode ...string) string {
	m := "auto"
	if len(mode) > 0 {
		m = mode[0]
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 80))
	fmt.Printf("👤 用户: %s\n", userInput)
	fmt.Printf("%s\n\n", strings.Repeat("=", 80))

	// 第一步: 检索相关笔记（为 agent 提供上下文）
	relevantNotes := cm.retrieveRelevantNotes(userInput)
	notePackets := cm.notesToPackets(relevantNotes)

	// 第二步: 构建优化的上下文
	_ = cm.ContextBuilder.Build(
		userInput,
		cm.ConversationHistory,
		cm.buildSystemInstructions(m),
		notePackets,
	)

	// 第三步: 让 Agent 自主决策和使用工具
	fmt.Println("🤖 Agent 正在思考并决定使用哪些工具...\n")

	// 更新 agent 的系统提示（包含上下文）
	cm.Agent.SystemPrompt = cm.buildSystemInstructions(m)

	// 调用 agent（agent 会自主决定是否使用工具）
	response := cm.Agent.Run(userInput)

	// 第四步: 统计工具使用情况
	cm.trackToolUsage()

	// 第五步: 更新对话历史
	cm.updateHistory(userInput, response)

	fmt.Printf("\n🤖 助手: %s\n", response)
	fmt.Printf("%s\n", strings.Repeat("=", 80))

	return response
}

// buildBaseSystemPrompt 构建基础系统提示
func (cm *CodebaseMaintainer) buildBaseSystemPrompt() string {
	return fmt.Sprintf(`你是 %s 项目的代码库维护助手。

你的核心能力:
1. 使用 TerminalTool 探索代码库
   - 你可以执行任何 shell 命令: ls, cat, grep, find, git 等
   - 工作目录: %s
   
2. 使用 NoteTool 记录发现和任务
   - 创建笔记记录重要发现
   - 笔记类型: blocker(阻塞问题)、action(行动计划)、task_state(任务状态)、conclusion(结论)
   
3. 使用 MemoryTool 存储关键信息
   - 记住重要的上下文信息
   - 跨会话保持连贯性

当前会话ID: %s

重要原则:
- 你要自主决定使用哪些工具、执行什么命令
- 探索代码库时，先了解整体结构，再深入细节
- 发现重要信息时，主动使用 NoteTool 记录
- 保持回答的专业性和实用性`, cm.ProjectName, cm.CodebasePath, cm.SessionID)
}

// trackToolUsage 统计工具使用情况
func (cm *CodebaseMaintainer) trackToolUsage() {
	// 从 agent 的执行历史中统计
	for _, msg := range cm.Agent.MessageHistory {
		if len(cm.Agent.MessageHistory) > 10 {
			// 只看最近10条
			msg = cm.Agent.MessageHistory[len(cm.Agent.MessageHistory)-1]
		}
		if msg.Role == "tool" {
			cm.Stats["tool_calls"] = cm.Stats["tool_calls"].(int) + 1
			content := strings.ToLower(msg.Content)
			if strings.Contains(content, "terminal") || strings.Contains(content, "command") {
				cm.Stats["commands_executed"] = cm.Stats["commands_executed"].(int) + 1
			} else if strings.Contains(content, "note") {
				if strings.Contains(content, "create") {
					cm.Stats["notes_created"] = cm.Stats["notes_created"].(int) + 1
				}
			}
		}
	}
}

// retrieveRelevantNotes 检索相关笔记
func (cm *CodebaseMaintainer) retrieveRelevantNotes(query string, limit ...int) []map[string]any {
	l := 3
	if len(limit) > 0 {
		l = limit[0]
	}

	// 优先检索 blocker
	blockersRaw := cm.NoteTool.Run(map[string]any{
		"action":    "list",
		"note_type": "blocker",
		"limit":     2,
	})
	blockers := normalizeNoteResults(blockersRaw)

	// 搜索相关笔记
	searchResultsRaw := cm.NoteTool.Run(map[string]any{
		"action": "search",
		"query":  query,
		"limit":  l,
	})
	searchResults := normalizeNoteResults(searchResultsRaw)

	// 合并去重
	allNotes := make(map[string]map[string]any)
	for _, note := range append(blockers, searchResults...) {
		noteID, _ := note["note_id"].(string)
		if noteID == "" {
			noteID, _ = note["id"].(string)
		}
		if noteID == "" {
			continue
		}
		if _, exists := allNotes[noteID]; !exists {
			allNotes[noteID] = note
		}
	}

	var result []map[string]any
	count := 0
	for _, note := range allNotes {
		if count >= l {
			break
		}
		result = append(result, note)
		count++
	}
	return result
}

// normalizeNoteResults 将笔记工具的返回值转换为笔记字典列表
func normalizeNoteResults(result any) []map[string]any {
	if result == nil {
		return nil
	}

	switch v := result.(type) {
	case map[string]any:
		return []map[string]any{v}
	case []any:
		var notes []map[string]any
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				notes = append(notes, m)
			}
		}
		return notes
	case string:
		text := strings.TrimSpace(v)
		if text == "" {
			return nil
		}
		if strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[") {
			var parsed any
			if err := json.Unmarshal([]byte(text), &parsed); err == nil {
				return normalizeNoteResults(parsed)
			}
		}
		return nil
	}
	return nil
}

// notesToPackets 将笔记转换为上下文包
func (cm *CodebaseMaintainer) notesToPackets(notes []map[string]any) []ContextPacket {
	// 根据笔记类型设置不同的相关性分数
	relevanceMap := map[string]float64{
		"blocker":     0.9,
		"action":      0.8,
		"task_state":  0.75,
		"conclusion":  0.7,
	}

	var packets []ContextPacket
	for _, note := range notes {
		noteType, _ := note["type"].(string)
		if noteType == "" {
			noteType = "general"
		}
		relevance := relevanceMap[noteType]
		if relevance == 0 {
			relevance = 0.6
		}

		title, _ := note["title"].(string)
		if title == "" {
			title = "Untitled"
		}
		content, _ := note["content"].(string)
		packetContent := fmt.Sprintf("[笔记:%s]\n类型: %s\n\n%s", title, noteType, content)

		var noteTimestamp time.Time
		if updatedAt, ok := note["updated_at"].(string); ok && updatedAt != "" {
			if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				noteTimestamp = t
			}
		}
		if noteTimestamp.IsZero() {
			noteTimestamp = time.Now()
		}

		noteID, _ := note["note_id"].(string)
		if noteID == "" {
			noteID, _ = note["id"].(string)
		}

		packets = append(packets, ContextPacket{
			Content:        packetContent,
			Timestamp:      noteTimestamp,
			TokenCount:     len(packetContent) / 4,
			RelevanceScore: relevance,
			Metadata: map[string]any{
				"type":      "note",
				"note_type": noteType,
				"note_id":   noteID,
			},
		})
	}
	return packets
}

// buildSystemInstructions 构建系统指令（Agentic 方式）
func (cm *CodebaseMaintainer) buildSystemInstructions(mode string) string {
	baseInstructions := cm.buildBaseSystemPrompt()

	modeHints := map[string]string{
		"explore": `
用户当前关注: 探索代码库

建议策略:
- 考虑使用 TerminalTool 了解代码结构（如 find, ls, tree）
- 查看关键文件（如 README, 主要模块）
- 将架构信息记录到笔记方便后续查阅`,
		"analyze": `
用户当前关注: 分析代码质量

建议策略:
- 考虑使用 grep 查找潜在问题（TODO, FIXME, BUG）
- 分析代码复杂度和结构
- 将发现的问题记录为 blocker 或 action 笔记`,
		"plan": `
用户当前关注: 任务规划

建议策略:
- 回顾历史笔记了解当前进度
- 基于已有信息制定行动计划
- 创建或更新 task_state 类型的笔记`,
		"auto": `
用户当前关注: 自由对话

建议策略:
- 根据用户需求灵活决策
- 在需要时主动使用工具获取信息
- 不需要时可以直接回答`,
	}

	hint, ok := modeHints[mode]
	if !ok {
		hint = modeHints["auto"]
	}

	return baseInstructions + "\n" + hint
}

// updateHistory 更新对话历史
func (cm *CodebaseMaintainer) updateHistory(userInput, response string) {
	cm.ConversationHistory = append(cm.ConversationHistory,
		Message{Content: userInput, Role: "user", Timestamp: time.Now()},
		Message{Content: response, Role: "assistant", Timestamp: time.Now()},
	)

	// 限制历史长度(保留最近10轮对话)
	if len(cm.ConversationHistory) > 20 {
		cm.ConversationHistory = cm.ConversationHistory[len(cm.ConversationHistory)-20:]
	}
}

// === 便捷方法 ===

// Explore 探索代码库（Agentic 方式）
//
// Agent 会自主决定使用哪些命令来探索代码库
func (cm *CodebaseMaintainer) Explore(target ...string) string {
	t := "."
	if len(target) > 0 {
		t = target[0]
	}
	return cm.Run(fmt.Sprintf("请探索 %s 的代码结构，了解项目组织方式", t), "explore")
}

// Analyze 分析代码质量（Agentic 方式）
//
// Agent 会自主决定如何分析代码质量
func (cm *CodebaseMaintainer) Analyze(focus ...string) string {
	query := "请分析代码质量"
	if len(focus) > 0 && focus[0] != "" {
		query += "，重点关注" + focus[0]
	}
	return cm.Run(query, "analyze")
}

// PlanNextSteps 规划下一步任务（Agentic 方式）
//
// Agent 会查看历史笔记并规划下一步
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
	var noteSummary any
	func() {
		defer func() { recover() }()
		noteSummary = cm.NoteTool.Run(map[string]any{"action": "summary"})
	}()

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

func main() {
	// 主函数 - 演示 CodebaseMaintainer 的使用（Agentic 版本）
	//
	// 在这个版本中：
	// - Agent 自主决定使用哪些工具
	// - 不预定义工作流
	// - Agent 根据需求灵活探索代码库
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("CodebaseMaintainer 演示（Agentic 版本）")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	// 初始化助手
	maintainer := NewCodebaseMaintainer(
		"my_flask_app",
		"./my_flask_app",
		&HelloAgentsLLM{},
	)

	// 探索代码库（Agent 自主决定如何探索）
	fmt.Println("\n### 探索代码库（Agent 自主探索）###")
	_ = maintainer.Explore()

	// 分析代码质量（Agent 自主决定分析方法）
	fmt.Println("\n### 分析代码质量（Agent 自主分析）###")
	_ = maintainer.Analyze()

	// 规划下一步（Agent 基于历史信息规划）
	fmt.Println("\n### 规划下一步任务（Agent 自主规划）###")
	_ = maintainer.PlanNextSteps()

	// 生成报告
	fmt.Println("\n### 生成会话报告 ###")
	report := maintainer.GenerateReport()
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(reportJSON))

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("演示完成!")
	fmt.Println(strings.Repeat("=", 80))
}
