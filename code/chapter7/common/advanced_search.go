package common

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

// MyAdvancedSearchTool 自定义高级搜索工具类
// 展示多源整合和智能选择的设计模式
type MyAdvancedSearchTool struct {
	Name           string
	Description    string
	SearchSources  []string
	tavilyClient   *http.Client
	tavilyAPIKey   string
	serpAPIKey     string
}

// NewMyAdvancedSearchTool 创建高级搜索工具实例
func NewMyAdvancedSearchTool() *MyAdvancedSearchTool {
	t := &MyAdvancedSearchTool{
		Name:        "my_advanced_search",
		Description: "智能搜索工具，支持多个搜索源，自动选择最佳结果",
	}
	t.setupSearchSources()
	return t
}

// setupSearchSources 设置可用的搜索源
func (t *MyAdvancedSearchTool) setupSearchSources() {
	// 检查Tavily可用性
	if apiKey := os.Getenv("TAVILY_API_KEY"); apiKey != "" {
		t.tavilyAPIKey = apiKey
		t.tavilyClient = &http.Client{Timeout: 30 * time.Second}
		t.SearchSources = append(t.SearchSources, "tavily")
		fmt.Println("✅ Tavily搜索源已启用")
	}

	// 检查SerpApi可用性
	if apiKey := os.Getenv("SERPAPI_API_KEY"); apiKey != "" {
		t.serpAPIKey = apiKey
		t.SearchSources = append(t.SearchSources, "serpapi")
		fmt.Println("✅ SerpApi搜索源已启用")
	}

	if len(t.SearchSources) > 0 {
		fmt.Printf("🔧 可用搜索源: %s\n", strings.Join(t.SearchSources, ", "))
	} else {
		fmt.Println("⚠️ 没有可用的搜索源，请配置API密钥")
	}
}

// Run 执行智能搜索（满足 Tool 接口）
func (t *MyAdvancedSearchTool) Run(input string) string {
	return t.Search(input)
}

// Search 执行智能搜索
func (t *MyAdvancedSearchTool) Search(query string) string {
	if strings.TrimSpace(query) == "" {
		return "❌ 错误：搜索查询不能为空"
	}

	// 检查是否有可用的搜索源
	if len(t.SearchSources) == 0 {
		return `❌ 没有可用的搜索源，请配置以下API密钥之一：

1. Tavily API: 设置环境变量 TAVILY_API_KEY
   获取地址: https://tavily.com/

2. SerpAPI: 设置环境变量 SERPAPI_API_KEY
   获取地址: https://serpapi.com/

配置后重新运行程序。`
	}

	fmt.Printf("🔍 开始智能搜索: %s\n", query)

	// 尝试多个搜索源，返回最佳结果
	for _, source := range t.SearchSources {
		var result string
		var err error

		if source == "tavily" {
			result, err = t.searchWithTavily(query)
			if err == nil && result != "" && !strings.Contains(result, "未找到") {
				return fmt.Sprintf("📊 Tavily AI搜索结果：\n\n%s", result)
			}
		} else if source == "serpapi" {
			result, err = t.searchWithSerpAPI(query)
			if err == nil && result != "" && !strings.Contains(result, "未找到") {
				return fmt.Sprintf("🌐 SerpApi Google搜索结果：\n\n%s", result)
			}
		}

		if err != nil {
			fmt.Printf("⚠️ %s 搜索失败: %v\n", source, err)
		}
	}

	return "❌ 所有搜索源都失败了，请检查网络连接和API密钥配置"
}

// tavilyResponse Tavily API 响应结构
type tavilyResponse struct {
	Answer  string         `json:"answer"`
	Results []tavilyResult `json:"results"`
}

type tavilyResult struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// searchWithTavily 使用Tavily搜索
func (t *MyAdvancedSearchTool) searchWithTavily(query string) (string, error) {
	if t.tavilyClient == nil || t.tavilyAPIKey == "" {
		return "", fmt.Errorf("Tavily未配置")
	}

	reqBody := map[string]any{
		"query":       query,
		"max_results": 3,
		"api_key":     t.tavilyAPIKey,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.tavily.com/search", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.tavilyClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tavilyResp tavilyResponse
	if err := json.Unmarshal(body, &tavilyResp); err != nil {
		return "", err
	}

	var result string
	if tavilyResp.Answer != "" {
		result = fmt.Sprintf("💡 AI直接答案：%s\n\n", tavilyResp.Answer)
	}

	result += "🔗 相关结果：\n"
	for i, item := range tavilyResp.Results {
		if i >= 3 {
			break
		}
		content := item.Content
		if len(content) > 150 {
			content = content[:150] + "..."
		}
		result += fmt.Sprintf("[%d] %s\n", i+1, item.Title)
		result += fmt.Sprintf("    %s\n\n", content)
	}

	return result, nil
}

// searchWithSerpAPI 使用SerpApi搜索
func (t *MyAdvancedSearchTool) searchWithSerpAPI(query string) (string, error) {
	if t.serpAPIKey == "" {
		return "", fmt.Errorf("SerpApi未配置")
	}

	reqURL := fmt.Sprintf("https://serpapi.com/search?engine=google&q=%s&api_key=%s&num=3",
		url.QueryEscape(query), t.serpAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var results map[string]any
	if err := json.Unmarshal(body, &results); err != nil {
		return "", err
	}

	result := "🔗 Google搜索结果：\n"
	if organicResults, ok := results["organic_results"].([]any); ok {
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
			result += fmt.Sprintf("[%d] %s\n", i+1, title)
			result += fmt.Sprintf("    %s\n\n", snippet)
		}
	}

	return result, nil
}

// CreateAdvancedSearchRegistry 创建包含高级搜索工具的注册表
func CreateAdvancedSearchRegistry() *ToolRegistry {
	registry := NewToolRegistry()

	// 创建搜索工具实例
	searchTool := NewMyAdvancedSearchTool()

	// 注册搜索工具的方法作为函数
	registry.RegisterFunction(
		"advanced_search",
		"高级搜索工具，整合Tavily和SerpAPI多个搜索源，提供更全面的搜索结果",
		searchTool.Run,
	)

	return registry
}
