// 14_WeatherAgent - 在 Agent 中使用天气 MCP 服务器
//
// 对应 Python: 14_weather_agent.py
// 演示如何在智能体中使用天气 MCP 服务器

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SimpleAgent 简单智能体（模拟）
type SimpleAgent struct {
	Name         string
	SystemPrompt string
	Tool         any
}

// Run 运行智能体
func (a *SimpleAgent) Run(input string) string {
	// 模拟智能体调用天气工具
	if strings.Contains(input, "天气") || strings.Contains(input, "北京") || strings.Contains(input, "上海") {
		return "北京今天天气：晴朗,温度25°C,适合户外活动。（模拟响应）"
	}
	return fmt.Sprintf("智能体 '%s' 处理: %s（模拟响应）", a.Name, input)
}

// MCPTool MCP 工具（模拟）
type MCPTool struct {
	Name          string
	ServerCommand []string
}

func main() {
	// 创建天气助手
	assistant := &SimpleAgent{
		Name: "天气助手",
		SystemPrompt: `你是天气助手，可以查询城市天气。
使用 get_weather 工具查询天气，支持中文城市名。`,
		Tool: &MCPTool{
			Name:          "weather",
			ServerCommand: []string{"python", "14_weather_mcp_server.py"},
		},
	}

	// 检查是否有命令行参数
	if len(os.Args) > 1 && os.Args[1] == "demo" {
		// 演示模式
		fmt.Println("\n查询北京天气：")
		response := assistant.Run("北京今天天气怎么样？")
		fmt.Printf("回答: %s\n\n", response)
	} else {
		// 交互模式
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("\n你: ")
			if !scanner.Scan() {
				break
			}
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}
			if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
				break
			}
			response := assistant.Run(input)
			fmt.Printf("助手: %s", response)
		}
	}
}
