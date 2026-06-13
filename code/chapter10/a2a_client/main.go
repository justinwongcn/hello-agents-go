// 09_A2A_Client - A2A Agent 客户端
//
// 对应 Python: 09_A2A_Client.py
// 创建 A2A 客户端，连接到研究员 Agent 并发送请求

package main

import (
	"fmt"
	"time"
)

// A2AClient A2A 客户端（模拟）
type A2AClient struct {
	Endpoint string
}

// ExecuteSkill 执行技能
func (c *A2AClient) ExecuteSkill(skillName string, text string) map[string]any {
	// 模拟远程调用
	return map[string]any{
		"status": "success",
		"result": map[string]any{
			"topic":    "AI在医疗领域的应用",
			"findings": "关于AI在医疗领域的应用的研究结果...",
			"sources":  []string{"来源1", "来源2", "来源3"},
		},
	}
}

func main() {
	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// 创建客户端连接到研究员Agent
	client := &A2AClient{Endpoint: "http://localhost:5000"}

	// 发送研究请求
	response := client.ExecuteSkill("research", "research AI在医疗领域的应用")
	fmt.Printf("收到响应：%v\n", response["result"])

	// 输出：
	// 收到响应：{'topic': 'AI在医疗领域的应用', 'findings': '关于AI在医疗领域的应用的研究结果...', 'sources': ['来源1', '来源2', '来源3']}
}
