// NoteTool 与 ContextBuilder 集成示例
//
// 展示如何将 NoteTool 与 ContextBuilder 集成，实现：
// 1. 长期项目追踪
// 2. 笔记检索与上下文注入
// 3. 基于历史笔记的连贯建议

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ContextConfig 上下文配置
type ContextConfig struct {
	MaxTokens         int
	ReserveRatio      float64
	MinRelevance      float64
	EnableCompression bool
}

// ContextPacket 上下文包
type ContextPacket struct {
	Content        string
	Timestamp      time.Time
	TokenCount     int
	RelevanceScore float64
	Metadata       map[string]any
}

// Message 消息
type Message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	MemoryTool any
	RAGTool    any
	Config     *ContextConfig
}

// NewContextBuilder 创建 ContextBuilder
func NewContextBuilder(memoryTool, ragTool any, config *ContextConfig) *ContextBuilder {
	return &ContextBuilder{
		MemoryTool: memoryTool,
		RAGTool:    ragTool,
		Config:     config,
	}
}

// Build 构建上下文
func (cb *ContextBuilder) Build(userQuery string, conversationHistory []Message, systemInstructions string, additionalPackets []ContextPacket) string {
	var sb strings.Builder
	sb.WriteString("=== System Instructions ===\n")
	sb.WriteString(systemInstructions)
	sb.WriteString("\n\n=== Conversation History ===\n")
	if len(conversationHistory) == 0 {
		sb.WriteString("(无历史记录)\n")
	} else {
		for _, msg := range conversationHistory {
			sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Role, msg.Content))
		}
	}

	if len(additionalPackets) > 0 {
		sb.WriteString("\n=== Additional Context (Notes) ===\n")
		for _, pkt := range additionalPackets {
			sb.WriteString(pkt.Content + "\n---\n")
		}
	}

	sb.WriteString("\n=== User Query ===\n")
	sb.WriteString(userQuery)
	sb.WriteString("\n")
	return sb.String()
}

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
	return &NoteTool{
		Workspace: workspace,
		Notes:     make(map[string]*Note),
		Counter:   0,
	}
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

// Invoke 调用 LLM
func (llm *HelloAgentsLLM) Invoke(messages []map[string]string) string {
	for _, msg := range messages {
		if msg["role"] == "user" {
			content := msg["content"]
			if strings.Contains(content, "完成") || strings.Contains(content, "重构") {
				return "很好！数据模型层重构完成且测试覆盖率达到85%是一个很好的里程碑。建议下一步: 1) 解决依赖冲突 2) 制定业务逻辑层重构计划"
			}
			if strings.Contains(content, "依赖") || strings.Contains(content, "冲突") {
				return "依赖版本冲突的解决方案: 1) 使用virtualenv隔离环境 2) 使用pip-tools锁定版本 3) 逐步升级依赖"
			}
		}
	}
	return "请问有什么可以帮助您的?"
}

// ProjectAssistant 长期项目助手,集成 NoteTool 和 ContextBuilder
type ProjectAssistant struct {
	Name              string
	ProjectName       string
	LLM               *HelloAgentsLLM
	NoteTool          *NoteTool
	ContextBuilder    *ContextBuilder
	ConversationHistory []Message
}

// NewProjectAssistant 创建 ProjectAssistant
func NewProjectAssistant(name, projectName string) *ProjectAssistant {
	llm := &HelloAgentsLLM{}
	noteTool := NewNoteTool(fmt.Sprintf("./%s_notes", projectName))
	config := &ContextConfig{MaxTokens: 4000}
	builder := NewContextBuilder(nil, nil, config)

	return &ProjectAssistant{
		Name:              name,
		ProjectName:       projectName,
		LLM:               llm,
		NoteTool:          noteTool,
		ContextBuilder:    builder,
		ConversationHistory: []Message{},
	}
}

// Run 运行助手,自动集成笔记
func (pa *ProjectAssistant) Run(userInput string, noteAsAction ...bool) string {
	saveNote := len(noteAsAction) > 0 && noteAsAction[0]

	// 1. 从 NoteTool 检索相关笔记
	relevantNotes := pa.retrieveRelevantNotes(userInput)

	// 2. 将笔记转换为 ContextPacket
	notePackets := pa.notesToPackets(relevantNotes)

	// 3. 构建优化的上下文
	optimizedContext := pa.ContextBuilder.Build(
		userInput,
		pa.ConversationHistory,
		pa.buildSystemInstructions(),
		notePackets,
	)

	// 4. 调用 LLM (以 messages 数组形式传入)
	messages := []map[string]string{
		{"role": "system", "content": optimizedContext},
		{"role": "user", "content": userInput},
	}
	response := pa.LLM.Invoke(messages)

	// 5. 如果需要,将交互记录为笔记
	if saveNote {
		pa.saveAsNote(userInput, response)
	}

	// 6. 更新对话历史
	pa.updateHistory(userInput, response)

	return response
}

// retrieveRelevantNotes 检索相关笔记
func (pa *ProjectAssistant) retrieveRelevantNotes(query string) []map[string]any {
	// 优先检索 blocker 和 action 类型的笔记
	blockersRaw := pa.NoteTool.Run(map[string]any{
		"action":    "list",
		"note_type": "blocker",
		"limit":     2,
	})

	// 通用搜索
	searchResultsRaw := pa.NoteTool.Run(map[string]any{
		"action": "search",
		"query":  query,
		"limit":  3,
	})

	blockers := ensureListOfDicts(blockersRaw)
	searchResults := ensureListOfDicts(searchResultsRaw)

	// 合并并去重
	allNotes := make(map[string]map[string]any)
	for _, note := range append(blockers, searchResults...) {
		noteID := getNoteField(note, "note_id", "id", "uuid", "title")
		allNotes[noteID] = note
	}

	var result []map[string]any
	count := 0
	for _, note := range allNotes {
		if count >= 3 {
			break
		}
		result = append(result, note)
		count++
	}
	return result
}

// ensureListOfDicts 将 NoteTool 返回规范化为字典列表
func ensureListOfDicts(data any) []map[string]any {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case string:
		var parsed any
		if err := json.Unmarshal([]byte(v), &parsed); err != nil {
			return nil
		}
		return ensureListOfDicts(parsed)
	case map[string]any:
		if items, ok := v["items"].([]any); ok {
			var result []map[string]any
			for _, item := range items {
				if m, ok := item.(map[string]any); ok {
					result = append(result, m)
				}
			}
			return result
		}
		return []map[string]any{v}
	case []any:
		var result []map[string]any
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				result = append(result, m)
			}
		}
		return result
	}
	return nil
}

// getNoteField 从笔记中获取第一个存在的字段值
func getNoteField(note map[string]any, keys ...string) string {
	for _, key := range keys {
		if val, ok := note[key]; ok {
			return fmt.Sprintf("%v", val)
		}
	}
	return fmt.Sprintf("%v", note)
}

// notesToPackets 将笔记转换为上下文包
func (pa *ProjectAssistant) notesToPackets(notes []map[string]any) []ContextPacket {
	var packets []ContextPacket

	for _, note := range notes {
		title, _ := note["title"].(string)
		body, _ := note["content"].(string)
		content := fmt.Sprintf("[笔记:%s]\n%s", title, body)

		// 安全解析时间戳
		var parsedTs time.Time
		for _, key := range []string{"updated_at", "updatedAt", "time", "timestamp"} {
			if ts, ok := note[key]; ok {
				switch v := ts.(type) {
				case string:
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						parsedTs = t
					}
				case float64:
					parsedTs = time.Unix(int64(v), 0)
				case int64:
					parsedTs = time.Unix(v, 0)
				}
				if !parsedTs.IsZero() {
					break
				}
			}
		}
		if parsedTs.IsZero() {
			parsedTs = time.Now()
		}

		noteType := "note"
		if t, ok := note["type"].(string); ok && t != "" {
			noteType = t
		} else if t, ok := note["note_type"].(string); ok && t != "" {
			noteType = t
		}

		noteID := getNoteField(note, "note_id", "id", "uuid", "title")

		packets = append(packets, ContextPacket{
			Content:        content,
			Timestamp:      parsedTs,
			TokenCount:     len(content) / 4, // 简单估算
			RelevanceScore: 0.75,             // 笔记具有较高相关性
			Metadata: map[string]any{
				"type":      "note",
				"note_type": noteType,
				"note_id":   noteID,
			},
		})
	}

	return packets
}

// saveAsNote 将交互保存为笔记
func (pa *ProjectAssistant) saveAsNote(userInput, response string) {
	// 判断应该保存为什么类型的笔记
	noteType := "conclusion"
	if strings.Contains(userInput, "问题") || strings.Contains(userInput, "阻塞") {
		noteType = "blocker"
	} else if strings.Contains(userInput, "计划") || strings.Contains(userInput, "下一步") {
		noteType = "action"
	}

	title := userInput
	if len([]rune(title)) > 30 {
		title = string([]rune(title)[:30]) + "..."
	}

	pa.NoteTool.Run(map[string]any{
		"action":    "create",
		"title":     title,
		"content":   fmt.Sprintf("## 问题\n%s\n\n## 分析\n%s", userInput, response),
		"note_type": noteType,
		"tags":      []string{pa.ProjectName, "auto_generated"},
	})
}

// buildSystemInstructions 构建系统指令
func (pa *ProjectAssistant) buildSystemInstructions() string {
	return fmt.Sprintf(`你是 %s 项目的长期助手。

你的职责:
1. 基于历史笔记提供连贯的建议
2. 追踪项目进展和待解决问题
3. 在回答时引用相关的历史笔记
4. 提供具体、可操作的下一步建议

注意:
- 优先关注标记为 blocker 的问题
- 在建议中说明依据来源(笔记、记忆或知识库)
- 保持对项目整体进度的认识`, pa.ProjectName)
}

// updateHistory 更新对话历史
func (pa *ProjectAssistant) updateHistory(userInput, response string) {
	now := time.Now()
	pa.ConversationHistory = append(pa.ConversationHistory, Message{Content: userInput, Role: "user", Timestamp: now})
	pa.ConversationHistory = append(pa.ConversationHistory, Message{Content: response, Role: "assistant", Timestamp: now})

	// 限制历史长度
	if len(pa.ConversationHistory) > 10 {
		pa.ConversationHistory = pa.ConversationHistory[len(pa.ConversationHistory)-10:]
	}
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("NoteTool 与 ContextBuilder 集成示例")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// 使用示例
	assistant := NewProjectAssistant("项目助手", "data_pipeline_refactoring")

	// 第一次交互:记录项目状态
	fmt.Println("第一次交互:记录项目状态")
	response := assistant.Run(
		"我们已经完成了数据模型层的重构,测试覆盖率达到85%。下一步计划重构业务逻辑层。",
		true,
	)
	fmt.Printf("助手回答: %s\n\n", response)

	// 第二次交互:提出问题
	fmt.Println("第二次交互:提出问题")
	response = assistant.Run(
		"在重构业务逻辑层时,我遇到了依赖版本冲突的问题,该如何解决?",
	)
	fmt.Printf("助手回答: %s\n\n", response)

	// 查看笔记摘要
	fmt.Println("查看笔记摘要:")
	summary := assistant.NoteTool.Run(map[string]any{"action": "summary"})
	summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(summaryJSON))

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("演示完成!")
	fmt.Println(strings.Repeat("=", 80))
}
