// 04_MCPTransport - MCP 传输方式示例
//
// 对应 Python: 04_MCPTransport.py
// 演示 MCP 的各种传输方式：Memory、Stdio、HTTP/SSE

package main

import "fmt"

// MCPTool MCP 工具（模拟）
type MCPTool struct {
	Name           string
	Description    string
	ServerCommand  []string
	transport      string
}

// Run 执行 MCP 工具操作
func (t *MCPTool) Run(params map[string]any) any {
	action, _ := params["action"].(string)
	switch action {
	case "list_tools":
		return []map[string]any{
			{"name": "add", "description": "加法计算"},
			{"name": "subtract", "description": "减法计算"},
			{"name": "multiply", "description": "乘法计算"},
			{"name": "divide", "description": "除法计算"},
			{"name": "greet", "description": "问候"},
			{"name": "get_system_info", "description": "获取系统信息"},
		}
	case "call_tool":
		toolName, _ := params["tool_name"].(string)
		args, _ := params["arguments"].(map[string]any)
		if toolName == "add" {
			a, _ := args["a"].(float64)
			b, _ := args["b"].(float64)
			return map[string]any{"result": a + b}
		}
		return nil
	default:
		return nil
	}
}

// MCPClient MCP 客户端（用于 HTTP/SSE 传输）
type MCPClient struct {
	Endpoint string
}

// ListTools 列出工具
func (c *MCPClient) ListTools() ([]map[string]any, error) {
	// 模拟远程调用
	return []map[string]any{
		{"name": "process_data", "description": "处理数据"},
	}, nil
}

// CallTool 调用工具
func (c *MCPClient) CallTool(toolName string, args map[string]any) (string, error) {
	return fmt.Sprintf("远程处理结果: %v", args), nil
}

func main() {
	// 1. Memory Transport - 内存传输（用于测试）
	// 不指定任何参数，使用内置演示服务器
	fmt.Println("=== 1. Memory Transport（内存传输）===")
	mcpTool := &MCPTool{} // 使用默认配置

	// 2. Stdio Transport - 标准输入输出传输（本地开发）
	// 使用命令列表启动本地服务器
	fmt.Println("\n=== 2. Stdio Transport（标准输入输出传输）===")
	mcpTool = &MCPTool{ServerCommand: []string{"python", "examples/mcp_example_server.py"}}
	_ = mcpTool

	// 3. Stdio Transport with Args - 带参数的命令传输
	// 可以传递额外参数
	fmt.Println("\n=== 3. Stdio Transport with Args（带参数的命令传输）===")
	mcpTool = &MCPTool{ServerCommand: []string{"python", "examples/mcp_example_server.py", "--debug"}}
	_ = mcpTool

	// 4. Stdio Transport - 社区服务器（npx方式）
	// 使用npx启动社区MCP服务器
	fmt.Println("\n=== 4. Stdio Transport - 社区服务器（npx方式）===")
	mcpTool = &MCPTool{ServerCommand: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."}}
	_ = mcpTool

	// 5. HTTP/SSE/StreamableHTTP Transport
	// 注意：MCPTool主要用于Stdio和Memory传输
	// 对于HTTP/SSE等远程传输，建议直接使用MCPClient
	fmt.Println("\n=== 5. HTTP/SSE Transport（远程传输）===")
	fmt.Println("注意：MCPTool 主要用于 Stdio 和 Memory 传输")
	fmt.Println("对于 HTTP/SSE 等远程传输，建议使用底层的 MCPClient")

	// 使用内置演示服务器（Memory传输）
	fmt.Println("\n=== 使用内置演示服务器（Memory传输）===")
	mcpTool = &MCPTool{}

	// 列出可用工具
	result := mcpTool.Run(map[string]any{"action": "list_tools"})
	fmt.Println(result)

	// 调用工具
	result = mcpTool.Run(map[string]any{
		"action":    "call_tool",
		"tool_name": "add",
		"arguments": map[string]any{"a": 10, "b": 20},
	})
	fmt.Println(result)

	// 使用自定义Python服务器
	fmt.Println("\n=== 使用自定义Python服务器 ===")
	mcpTool = &MCPTool{ServerCommand: []string{"python", "my_mcp_server.py"}}

	// 使用社区服务器（文件系统）
	mcpTool = &MCPTool{ServerCommand: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."}}

	// 列出工具
	result = mcpTool.Run(map[string]any{"action": "list_tools"})
	fmt.Println(result)

	// 调用工具
	result = mcpTool.Run(map[string]any{
		"action":    "call_tool",
		"tool_name": "read_file",
		"arguments": map[string]any{"path": "my_README.md"},
	})
	fmt.Println(result)

	// HTTP 传输示例
	fmt.Println("\n=== HTTP 传输示例 ===")
	fmt.Println("注意：需要实际的 HTTP MCP 服务器")
	client := &MCPClient{Endpoint: "http://api.example.com/mcp"}
	tools, _ := client.ListTools()
	fmt.Printf("远程服务器工具: %d 个\n", len(tools))
	resultStr, _ := client.CallTool("process_data", map[string]any{
		"data":      "Hello, World!",
		"operation": "uppercase",
	})
	fmt.Printf("远程处理结果: %s\n", resultStr)
}
