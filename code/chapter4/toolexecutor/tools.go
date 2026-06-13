package toolexecutor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// Search 一个基于SerpApi的实战网页搜索引擎工具。
// 它会智能地解析搜索结果，优先返回直接答案或知识图谱信息。
func Search(query string) string {
	fmt.Printf("🔍 正在执行 [SerpApi] 网页搜索: %s\n", query)

	apiKey := os.Getenv("SERPAPI_API_KEY")
	if apiKey == "" {
		return "错误：SERPAPI_API_KEY 未在 .env 文件中配置。"
	}

	// 构建请求URL
	reqURL := fmt.Sprintf("https://serpapi.com/search?engine=google&q=%s&api_key=%s&gl=cn&hl=zh-cn",
		url.QueryEscape(query), apiKey)

	resp, err := httpClient.Get(reqURL)
	if err != nil {
		return fmt.Sprintf("搜索时发生错误: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("搜索时发生错误: %v", err)
	}

	var results map[string]any
	if err := json.Unmarshal(body, &results); err != nil {
		return fmt.Sprintf("搜索时发生错误: %v", err)
	}

	// 智能解析：优先寻找最直接的答案
	if answerBoxList, ok := results["answer_box_list"].([]any); ok && len(answerBoxList) > 0 {
		var parts []string
		for _, item := range answerBoxList {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, "\n")
	}

	if answerBox, ok := results["answer_box"].(map[string]any); ok {
		if answer, ok := answerBox["answer"].(string); ok && answer != "" {
			return answer
		}
	}

	if kg, ok := results["knowledge_graph"].(map[string]any); ok {
		if desc, ok := kg["description"].(string); ok && desc != "" {
			return desc
		}
	}

	if organicResults, ok := results["organic_results"].([]any); ok && len(organicResults) > 0 {
		var snippets []string
		for i, res := range organicResults {
			if i >= 3 {
				break
			}
			resMap, ok := res.(map[string]any)
			if !ok {
				continue
			}
			title, _ := resMap["title"].(string)
			snippet, _ := resMap["snippet"].(string)
			snippets = append(snippets, fmt.Sprintf("[%d] %s\n%s", i+1, title, snippet))
		}
		return strings.Join(snippets, "\n\n")
	}

	return fmt.Sprintf("对不起，没有找到关于 '%s' 的信息。", query)
}

// ToolInfo 工具信息
type ToolInfo struct {
	Description string
	Func        func(string) string
}

// ToolExecutor 一个工具执行器，负责管理和执行工具。
type ToolExecutor struct {
	tools map[string]ToolInfo
}

// NewToolExecutor 创建新的工具执行器
func NewToolExecutor() *ToolExecutor {
	return &ToolExecutor{
		tools: make(map[string]ToolInfo),
	}
}

// RegisterTool 向工具箱中注册一个新工具。
func (te *ToolExecutor) RegisterTool(name, description string, fn func(string) string) {
	if _, exists := te.tools[name]; exists {
		fmt.Printf("警告：工具 '%s' 已存在，将被覆盖。\n", name)
	}
	te.tools[name] = ToolInfo{Description: description, Func: fn}
	fmt.Printf("工具 '%s' 已注册。\n", name)
}

// GetTool 根据名称获取一个工具的执行函数。
func (te *ToolExecutor) GetTool(name string) func(string) string {
	if info, ok := te.tools[name]; ok {
		return info.Func
	}
	return nil
}

// GetAvailableTools 获取所有可用工具的格式化描述字符串。
func (te *ToolExecutor) GetAvailableTools() string {
	var lines []string
	for name, info := range te.tools {
		lines = append(lines, fmt.Sprintf("- %s: %s", name, info.Description))
	}
	return strings.Join(lines, "\n")
}
