// 02_Connect2MCP - 连接到 MCP 服务器
//
// 对应 Python: 02_Connect2MCP.py
// 演示如何连接到 MCP 服务器，发现工具，调用工具，以及错误处理

package main

import (
	"fmt"
	"strings"
)

// MCPClient MCP 客户端（模拟）
type MCPClient struct {
	Command []string
}

// ListTools 列出可用工具
func (c *MCPClient) ListTools() []map[string]any {
	// 模拟返回工具列表
	return []map[string]any{
		{"name": "read_file", "description": "读取文件内容", "inputSchema": map[string]any{
			"properties": map[string]any{
				"path": map[string]any{"type": "string", "description": "文件路径"},
			},
		}},
		{"name": "write_file", "description": "写入文件内容", "inputSchema": map[string]any{
			"properties": map[string]any{
				"path":    map[string]any{"type": "string", "description": "文件路径"},
				"content": map[string]any{"type": "string", "description": "文件内容"},
			},
		}},
		{"name": "list_directory", "description": "列出目录内容"},
		{"name": "create_directory", "description": "创建目录"},
		{"name": "move_file", "description": "移动文件"},
	}
}

// CallTool 调用工具
func (c *MCPClient) CallTool(toolName string, args map[string]any) (string, error) {
	// 模拟工具调用
	switch toolName {
	case "read_file":
		path, _ := args["path"].(string)
		if path == "nonexistent.txt" {
			return "", fmt.Errorf("文件不存在: %s", path)
		}
		return fmt.Sprintf("文件 '%s' 的内容（模拟）", path), nil
	case "write_file":
		path, _ := args["path"].(string)
		return fmt.Sprintf("已写入文件: %s", path), nil
	case "list_directory":
		return "main.go\nREADME.md\ngo.mod", nil
	default:
		return "", fmt.Errorf("未知工具: %s", toolName)
	}
}

// Close 关闭客户端
func (c *MCPClient) Close() {
	// 模拟关闭连接
}

// connectToServer 连接到服务器
func connectToServer() {
	// 方式1：连接到社区提供的文件系统服务器
	// npx会自动下载并运行@modelcontextprotocol/server-filesystem包
	client := &MCPClient{
		Command: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."},
	}
	defer client.Close()

	// 获取可用工具
	tools := client.ListTools()
	toolNames := make([]string, len(tools))
	for i, t := range tools {
		toolNames[i] = t["name"].(string)
	}
	fmt.Printf("可用工具: %v\n", toolNames)

	// 方式2：连接到自定义的Python MCP服务器
	client2 := &MCPClient{Command: []string{"python", "my_mcp_server.py"}}
	defer client2.Close()
	// 使用client2...
	_ = client2
}

// discoverTools 发现工具
func discoverTools() {
	client := &MCPClient{
		Command: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."},
	}
	defer client.Close()

	// 获取所有可用工具
	tools := client.ListTools()

	fmt.Printf("服务器提供了 %d 个工具：\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("\n工具名称: %s\n", tool["name"])
		desc, ok := tool["description"].(string)
		if !ok || desc == "" {
			desc = "无描述"
		}
		fmt.Printf("描述: %s\n", desc)

		// 打印参数信息
		if schema, ok := tool["inputSchema"].(map[string]any); ok {
			if props, ok := schema["properties"].(map[string]any); ok {
				fmt.Println("参数:")
				for paramName, paramInfo := range props {
					info := paramInfo.(map[string]any)
					paramType, _ := info["type"].(string)
					if paramType == "" {
						paramType = "any"
					}
					paramDesc, _ := info["description"].(string)
					fmt.Printf("  - %s (%s): %s\n", paramName, paramType, paramDesc)
				}
			}
		}
	}

	// 输出示例：
	// 服务器提供了 5 个工具：
	//
	// 工具名称: read_file
	// 描述: 读取文件内容
	// 参数:
	//   - path (string): 文件路径
	//
	// 工具名称: write_file
	// 描述: 写入文件内容
	// 参数:
	//   - path (string): 文件路径
	//   - content (string): 文件内容
}

// useTools 使用工具
func useTools() {
	client := &MCPClient{
		Command: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."},
	}
	defer client.Close()

	// 读取文件
	result, _ := client.CallTool("read_file", map[string]any{"path": "my_README.md"})
	fmt.Printf("文件内容：\n%s\n", result)

	// 列出目录
	result, _ = client.CallTool("list_directory", map[string]any{"path": "."})
	fmt.Printf("当前目录文件：%s\n", result)

	// 写入文件
	result, _ = client.CallTool("write_file", map[string]any{
		"path":    "output.txt",
		"content": "Hello from MCP!",
	})
	fmt.Printf("写入结果：%s\n", result)
}

// safeToolCall 安全的工具调用（带错误处理）
func safeToolCall() {
	client := &MCPClient{
		Command: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."},
	}
	defer client.Close()

	// 尝试读取可能不存在的文件
	result, err := client.CallTool("read_file", map[string]any{"path": "nonexistent.txt"})
	if err != nil {
		fmt.Printf("工具调用失败: %v\n", err)
		// 可以选择重试、使用默认值或向用户报告错误
		return
	}
	fmt.Println(result)
}

func main() {
	fmt.Println("=== 连接到服务器 ===")
	connectToServer()

	fmt.Println("\n=== 发现工具 ===")
	discoverTools()

	fmt.Println("\n=== 使用工具 ===")
	useTools()

	fmt.Println("\n=== 安全工具调用 ===")
	safeToolCall()

	_ = strings.TrimSpace("") // 避免未使用导入警告
}
