// 10_A2ATool_Simple - 使用 A2ATool 包装器
//
// 对应 Python: 10_A2ATool_Simple.py
// 展示如何使用 A2ATool 包装器将 A2A Agent 集成到协调者 Agent

package main

import "fmt"

// A2AClient A2A 客户端（模拟）
type A2AClient struct {
	Endpoint string
}

// ExecuteSkill 执行技能
func (c *A2AClient) ExecuteSkill(skillName string, text string) map[string]any {
	return map[string]any{
		"status": "success",
		"result": fmt.Sprintf("研究结果: AI在教育领域的应用（模拟）"),
	}
}

// A2ATool A2A 工具（模拟）
type A2ATool struct {
	AgentURL string
	Client   *A2AClient
}

// Run 执行工具
func (t *A2ATool) Run(action string, text string) map[string]any {
	return t.Client.ExecuteSkill("research", text)
}

// SimpleAgent 简单智能体（模拟）
type SimpleAgent struct {
	Name  string
	Tools []any
}

// AddTool 添加工具
func (a *SimpleAgent) AddTool(tool any) {
	a.Tools = append(a.Tools, tool)
}

// Run 运行智能体
func (a *SimpleAgent) Run(input string) string {
	return fmt.Sprintf("智能体 '%s' 处理请求: %s（模拟响应）", a.Name, input)
}

func main() {
	// 假设已经有一个研究员Agent服务运行在 http://localhost:5000

	// 创建协调者Agent
	coordinator := &SimpleAgent{Name: "协调者"}

	// 添加A2A工具，连接到研究员Agent
	researcherTool := &A2ATool{
		AgentURL: "http://localhost:5000",
		Client:   &A2AClient{Endpoint: "http://localhost:5000"},
	}
	coordinator.AddTool(researcherTool)

	// 协调者可以调用研究员Agent
	// 使用 action="ask" 向 Agent 提问
	response := coordinator.Run("使用a2a工具，向Agent提问：请研究AI在教育领域的应用")
	fmt.Println(response)
}
