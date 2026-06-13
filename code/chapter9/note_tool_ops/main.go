// NoteTool 基本操作示例
//
// 展示 NoteTool 的核心操作：
// 1. 创建笔记 (create)
// 2. 读取笔记 (read)
// 3. 更新笔记 (update)
// 4. 搜索笔记 (search)
// 5. 列出笔记 (list)
// 6. 笔记摘要 (summary)
// 7. 删除笔记 (delete)

package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
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
	mu        sync.Mutex
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
func (nt *NoteTool) Run(params map[string]any) string {
	action, _ := params["action"].(string)
	switch action {
	case "create":
		return nt.create(params)
	case "read":
		return nt.read(params)
	case "update":
		return nt.update(params)
	case "search":
		return nt.search(params)
	case "list":
		return nt.list(params)
	case "summary":
		return nt.summary()
	case "delete":
		return nt.delete(params)
	default:
		return "未知操作"
	}
}

func (nt *NoteTool) create(params map[string]any) string {
	nt.mu.Lock()
	defer nt.mu.Unlock()

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
	} else if tagSlice, ok := params["tags"].([]any); ok {
		for _, t := range tagSlice {
			if s, ok := t.(string); ok {
				tags = append(tags, s)
			}
		}
	}

	now := time.Now()
	note := &Note{
		NoteID:    noteID,
		Title:     title,
		Content:   content,
		NoteType:  noteType,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	nt.Notes[noteID] = note

	return fmt.Sprintf("✅ 笔记已创建\nID: %s\n标题: %s\n类型: %s", noteID, title, noteType)
}

func (nt *NoteTool) read(params map[string]any) string {
	noteID, _ := params["note_id"].(string)
	note, ok := nt.Notes[noteID]
	if !ok {
		return fmt.Sprintf("❌ 笔记未找到: %s", noteID)
	}
	return fmt.Sprintf("📝 笔记详情\nID: %s\n标题: %s\n类型: %s\n标签: %v\n创建时间: %s\n更新时间: %s\n\n%s",
		note.NoteID, note.Title, note.NoteType, note.Tags,
		note.CreatedAt.Format("2006-01-02 15:04:05"),
		note.UpdatedAt.Format("2006-01-02 15:04:05"),
		note.Content)
}

func (nt *NoteTool) update(params map[string]any) string {
	nt.mu.Lock()
	defer nt.mu.Unlock()

	noteID, _ := params["note_id"].(string)
	note, ok := nt.Notes[noteID]
	if !ok {
		return fmt.Sprintf("❌ 笔记未找到: %s", noteID)
	}

	if content, ok := params["content"].(string); ok {
		note.Content = content
	}
	note.UpdatedAt = time.Now()

	return fmt.Sprintf("✅ 笔记已更新: %s", noteID)
}

func (nt *NoteTool) search(params map[string]any) string {
	query, _ := params["query"].(string)
	limit := 10
	if l, ok := params["limit"].(int); ok {
		limit = l
	}

	var results []string
	count := 0
	for _, note := range nt.Notes {
		if count >= limit {
			break
		}
		if strings.Contains(note.Title, query) || strings.Contains(note.Content, query) {
			results = append(results, fmt.Sprintf("- [%s] %s (%s)", note.NoteID, note.Title, note.NoteType))
			count++
		}
	}

	if len(results) == 0 {
		return fmt.Sprintf("🔍 搜索 '%s': 未找到相关笔记", query)
	}
	return fmt.Sprintf("🔍 搜索 '%s' 的结果 (%d条):\n%s", query, len(results), strings.Join(results, "\n"))
}

func (nt *NoteTool) list(params map[string]any) string {
	noteType, _ := params["note_type"].(string)
	limit := 10
	if l, ok := params["limit"].(int); ok {
		limit = l
	}

	var results []string
	count := 0
	for _, note := range nt.Notes {
		if count >= limit {
			break
		}
		if noteType == "" || note.NoteType == noteType {
			results = append(results, fmt.Sprintf("- [%s] %s (%s)", note.NoteID, note.Title, note.NoteType))
			count++
		}
	}

	if len(results) == 0 {
		return fmt.Sprintf("📋 笔记列表 (%s): 无笔记", noteType)
	}
	return fmt.Sprintf("📋 笔记列表 (%s) (%d条):\n%s", noteType, len(results), strings.Join(results, "\n"))
}

func (nt *NoteTool) summary() string {
	typeCount := make(map[string]int)
	for _, note := range nt.Notes {
		typeCount[note.NoteType]++
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 笔记摘要 (总计: %d条)\n", len(nt.Notes)))
	for t, c := range typeCount {
		sb.WriteString(fmt.Sprintf("  - %s: %d条\n", t, c))
	}
	return sb.String()
}

func (nt *NoteTool) delete(params map[string]any) string {
	nt.mu.Lock()
	defer nt.mu.Unlock()

	noteID, _ := params["note_id"].(string)
	if _, ok := nt.Notes[noteID]; !ok {
		return fmt.Sprintf("❌ 笔记未找到: %s", noteID)
	}
	delete(nt.Notes, noteID)
	return fmt.Sprintf("🗑️ 笔记已删除: %s", noteID)
}

// extractNoteId 从 NoteTool 的输出文本中提取 note_id
func extractNoteId(output string) (string, error) {
	re := regexp.MustCompile(`ID:\s*(note_[0-9_]+)`)
	matches := re.FindStringSubmatch(output)
	if matches == nil {
		return "", fmt.Errorf("无法从输出解析 note_id:\n%s", output)
	}
	return matches[1], nil
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("NoteTool 基本操作示例")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// 初始化 NoteTool
	notes := NewNoteTool("./project_notes")

	// 1. 创建笔记
	fmt.Println("1. 创建笔记...")
	createOutput1 := notes.Run(map[string]any{
		"action":    "create",
		"title":     "重构项目 - 第一阶段",
		"content":   "## 完成情况\n已完成数据模型层的重构,测试覆盖率达到85%。\n\n## 下一步\n重构业务逻辑层",
		"note_type": "task_state",
		"tags":      []string{"refactoring", "phase1"},
	})
	fmt.Println(createOutput1 + "\n")
	noteID1, _ := extractNoteId(createOutput1)

	// 创建第二个笔记
	createOutput2 := notes.Run(map[string]any{
		"action":    "create",
		"title":     "依赖冲突问题",
		"content":   "## 问题描述\n发现某些第三方库版本不兼容,需要解决。\n\n## 影响范围\n业务逻辑层的3个模块\n\n## 下一步\n1. 使用虚拟环境隔离\n2. 锁定版本\n3. 使用 pipdeptree 分析依赖树",
		"note_type": "blocker",
		"tags":      []string{"dependency", "urgent"},
	})
	fmt.Println(createOutput2 + "\n")
	noteID2, _ := extractNoteId(createOutput2)

	// 2. 读取笔记
	fmt.Println("2. 读取笔记...")
	noteDetail := notes.Run(map[string]any{
		"action":  "read",
		"note_id": noteID1,
	})
	fmt.Println(noteDetail + "\n")

	// 3. 更新笔记
	fmt.Println("3. 更新笔记...")
	updateResult := notes.Run(map[string]any{
		"action":  "update",
		"note_id": noteID1,
		"content": "## 完成情况\n已完成数据模型层的重构,测试覆盖率达到85%。\n\n## 问题\n遇到依赖版本冲突,已记录到单独笔记。\n\n## 下一步\n先解决依赖冲突,再继续重构业务逻辑层",
	})
	fmt.Println(updateResult + "\n")

	// 4. 搜索笔记
	fmt.Println("4. 搜索笔记...")
	searchResults := notes.Run(map[string]any{
		"action": "search",
		"query":  "依赖",
		"limit":  5,
	})
	fmt.Println(searchResults + "\n")

	// 5. 列出笔记
	fmt.Println("5. 列出所有 blocker 类型的笔记...")
	blockers := notes.Run(map[string]any{
		"action":    "list",
		"note_type": "blocker",
		"limit":     10,
	})
	fmt.Println(blockers + "\n")

	// 6. 笔记摘要
	fmt.Println("6. 生成笔记摘要...")
	summaryOutput := notes.Run(map[string]any{
		"action": "summary",
	})
	fmt.Println(summaryOutput + "\n")

	// 7. 删除笔记 (演示，实际使用时谨慎)
	fmt.Println("7. 删除笔记 (演示)...")
	// delete_result = notes.Run({
	//     "action": "delete",
	//     "note_id": note_id_2
	// })
	// print(delete_result + "\n")
	fmt.Printf("(已跳过实际删除操作, 笔记ID: %s)\n\n", noteID2)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("NoteTool 操作演示完成!")
	fmt.Println(strings.Repeat("=", 80))
}
