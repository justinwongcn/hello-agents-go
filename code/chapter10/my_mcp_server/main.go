// my_mcp_server - 自定义 MCP 服务器示例
//
// 对应 Python: my_mcp_server.py
// 这是一个简单的 MCP 服务器,提供基础的数学计算和文本处理工具。
// 用于演示如何创建自己的 MCP 服务器。

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

// MCPServer MCP 服务器（模拟）
type MCPServer struct {
	Name     string
	Tools    map[string]any
	Resources map[string]func() string
	Prompts  map[string]func() string
}

// NewMCPServer 创建 MCP 服务器
func NewMCPServer(name string) *MCPServer {
	return &MCPServer{
		Name:      name,
		Tools:     make(map[string]any),
		Resources: make(map[string]func() string),
		Prompts:   make(map[string]func() string),
	}
}

// ==================== 数学工具 ====================

// add 加法计算器
func add(a, b float64) float64 {
	return a + b
}

// subtract 减法计算器
func subtract(a, b float64) float64 {
	return a - b
}

// multiply 乘法计算器
func multiply(a, b float64) float64 {
	return a * b
}

// divide 除法计算器
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// ==================== 文本处理工具 ====================

// reverseText 反转文本
func reverseText(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// countWords 统计文本中的单词数量
func countWords(text string) int {
	return len(strings.Fields(text))
}

// toUppercase 将文本转换为大写
func toUppercase(text string) string {
	return strings.ToUpper(text)
}

// toLowercase 将文本转换为小写
func toLowercase(text string) string {
	return strings.ToLower(text)
}

// ==================== 资源定义 ====================

// getServerConfig 获取服务器配置信息
func getServerConfig() string {
	config := map[string]any{
		"name":        "MyCustomServer",
		"version":     "1.0.0",
		"tools_count": 8,
		"description": "自定义MCP服务器示例",
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	return string(data)
}

// getCapabilities 获取服务器能力列表
func getCapabilities() string {
	return `服务器能力列表：

数学计算：
- add: 加法计算
- subtract: 减法计算
- multiply: 乘法计算
- divide: 除法计算

文本处理：
- reverse_text: 反转文本
- count_words: 统计单词数
- to_uppercase: 转换为大写
- to_lowercase: 转换为小写

资源：
- config://server: 服务器配置
- info://capabilities: 能力列表（本资源）`
}

// ==================== 提示词模板 ====================

// mathHelper 数学计算助手提示词
func mathHelper() string {
	return `你是一个数学计算助手。你可以使用以下工具：
- add(a, b): 计算两数之和
- subtract(a, b): 计算两数之差
- multiply(a, b): 计算两数之积
- divide(a, b): 计算两数之商

请根据用户的问题选择合适的工具进行计算。`
}

// textProcessor 文本处理助手提示词
func textProcessor() string {
	return `你是一个文本处理助手。你可以使用以下工具：
- reverse_text(text): 反转文本
- count_words(text): 统计单词数
- to_uppercase(text): 转换为大写
- to_lowercase(text): 转换为小写

请根据用户的需求选择合适的工具处理文本。`
}

// ==================== 主程序 ====================

func main() {
	// 创建MCP服务器实例
	mcp := NewMCPServer("MyCustomServer")

	// 注册工具
	mcp.Tools["add"] = add
	mcp.Tools["subtract"] = subtract
	mcp.Tools["multiply"] = multiply
	mcp.Tools["divide"] = divide
	mcp.Tools["reverse_text"] = reverseText
	mcp.Tools["count_words"] = countWords
	mcp.Tools["to_uppercase"] = toUppercase
	mcp.Tools["to_lowercase"] = toLowercase

	// 注册资源
	mcp.Resources["config://server"] = getServerConfig
	mcp.Resources["info://capabilities"] = getCapabilities

	// 注册提示词
	mcp.Prompts["math_helper"] = mathHelper
	mcp.Prompts["text_processor"] = textProcessor

	// 运行MCP服务器
	// 演示工具调用
	fmt.Println("=== MCP 服务器演示 ===")
	fmt.Printf("服务器名称: %s\n", mcp.Name)
	fmt.Printf("工具数量: %d\n", len(mcp.Tools))
	fmt.Printf("资源数量: %d\n", len(mcp.Resources))
	fmt.Printf("提示词数量: %d\n", len(mcp.Prompts))

	fmt.Println("\n--- 数学工具测试 ---")
	fmt.Printf("add(10, 20) = %g\n", add(10, 20))
	fmt.Printf("subtract(10, 20) = %g\n", subtract(10, 20))
	fmt.Printf("multiply(10, 20) = %g\n", multiply(10, 20))
	result, err := divide(10, 20)
	if err != nil {
		fmt.Printf("divide(10, 20) = 错误: %v\n", err)
	} else {
		fmt.Printf("divide(10, 20) = %g\n", result)
	}

	fmt.Println("\n--- 文本处理工具测试 ---")
	fmt.Printf("reverse_text('Hello') = %s\n", reverseText("Hello"))
	fmt.Printf("count_words('Hello World') = %d\n", countWords("Hello World"))
	fmt.Printf("to_uppercase('hello') = %s\n", toUppercase("hello"))
	fmt.Printf("to_lowercase('HELLO') = %s\n", toLowercase("HELLO"))

	_ = unicode.IsLetter('a') // 避免未使用导入警告
}
