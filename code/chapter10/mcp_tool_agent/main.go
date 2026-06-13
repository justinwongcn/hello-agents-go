// 05_UseMCPToolInAgent - 在智能体中使用 MCP 工具
//
// 对应 Python: 05_UseMCPToolInAgent.py
// 演示如何将 MCP 工具集成到智能体中

package main

import (
	"fmt"
	"strings"
)

// SimpleAgent 简单智能体
type SimpleAgent struct {
	Name         string
	Tools        map[string]*MCPTool
	SystemPrompt string
}

// AddTool 添加工具
func (a *SimpleAgent) AddTool(tool *MCPTool) {
	a.Tools[tool.Name] = tool
}

// Run 运行智能体
func (a *SimpleAgent) Run(input string) string {
	// 模拟智能体根据输入选择合适的工具
	for _, tool := range a.Tools {
		if tool.Name == "mcp" {
			result := tool.Run(map[string]any{
				"action":    "call_tool",
				"tool_name": "add",
				"arguments": map[string]any{"a": 123, "b": 456},
			})
			return fmt.Sprintf("计算结果: %v", result)
		}
	}
	return "未找到合适的工具"
}

// MCPTool MCP 工具
type MCPTool struct {
	Name          string
	Description   string
	ServerCommand []string
}

// Run 执行 MCP 工具操作
func (t *MCPTool) Run(params map[string]any) any {
	action, _ := params["action"].(string)
	switch action {
	case "call_tool":
		toolName, _ := params["tool_name"].(string)
		args, _ := params["arguments"].(map[string]any)
		if toolName == "add" {
			a, _ := args["a"].(float64)
			b, _ := args["b"].(float64)
			return a + b
		}
		return nil
	default:
		return nil
	}
}

func main() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("方式1：使用内置演示服务器")
	fmt.Println(strings.Repeat("=", 70))

	agent := &SimpleAgent{
		Name:  "助手",
		Tools: make(map[string]*MCPTool),
	}

	// 无需任何配置，自动使用内置演示服务器
	// 内置服务器提供：add, subtract, multiply, divide, greet, get_system_info
	mcpTool := &MCPTool{Name: "mcp", Description: "MCP工具"} // 默认name="mcp"
	agent.AddTool(mcpTool)

	// 智能体可以使用内置工具
	response := agent.Run("计算 123 + 456")
	fmt.Println(response) // 智能体会自动调用add工具

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("方式2：连接外部MCP服务器（使用多个服务器）")
	fmt.Println(strings.Repeat("=", 70))

	// 重要：为每个MCP服务器指定不同的name，避免工具名称冲突

	// 示例1：连接到社区提供的文件系统服务器
	fsTool := &MCPTool{
		Name:          "filesystem", // 指定唯一名称
		Description:   "访问本地文件系统",
		ServerCommand: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."},
	}
	agent.AddTool(fsTool)

	// 示例2：连接到自定义的 Python MCP 服务器
	// 关于如何编写自定义MCP服务器，请参考10.5章节
	customTool := &MCPTool{
		Name:          "custom_server", // 使用不同的名称
		Description:   "自定义业务逻辑服务器",
		ServerCommand: []string{"python", "my_mcp_server.py"},
	}
	agent.AddTool(customTool)

	fmt.Println("\n当前Agent拥有的工具：")
	fmt.Printf("- %s: %s\n", mcpTool.Name, mcpTool.Description)
	fmt.Printf("- %s: %s\n", fsTool.Name, fsTool.Description)
	fmt.Printf("- %s: %s\n", customTool.Name, customTool.Description)

	// Agent现在可以自动使用这些工具！
	response = agent.Run("请读取my_README.md文件，并总结其中的主要内容")
	fmt.Println(response)
}
