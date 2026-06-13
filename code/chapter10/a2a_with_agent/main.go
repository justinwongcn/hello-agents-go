// 09_A2A_WithAgent - A2A 协议 + SimpleAgent 集成案例
//
// 对应 Python: 09_A2A_WithAgent.py
// 展示如何将 A2A 协议的 Agent 作为工具集成到 SimpleAgent 中

package main

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================
// 1. 创建专业 A2A Agent 服务
// ============================================================

// A2AServer A2A 服务器（模拟）
type A2AServer struct {
	Name        string
	Description string
	Version     string
	Skills      map[string]func(string) string
}

// Skill 添加技能
func (s *A2AServer) Skill(name string) func(func(string) string) {
	return func(fn func(string) string) {
		s.Skills[name] = fn
	}
}

// Run 运行服务器
func (s *A2AServer) Run(host string, port int) {
	fmt.Printf("A2A 服务器 %s 启动在 %s:%d\n", s.Name, host, port)
	select {}
}

// A2AClient A2A 客户端（模拟）
type A2AClient struct {
	Endpoint string
}

// ExecuteSkill 执行技能
func (c *A2AClient) ExecuteSkill(skillName string, text string) map[string]any {
	return map[string]any{
		"status": "success",
		"result": fmt.Sprintf("模拟执行 %s: %s", skillName, text),
	}
}

// ============================================================
// 3. 创建 A2A 工具（封装 A2A Agent 为 Tool）
// ============================================================

// A2ATool 将 A2A Agent 封装为 Tool
type A2ATool struct {
	Name        string
	Description string
	AgentURL    string
	SkillName   string
	Client      *A2AClient
	Parameters  []map[string]any
}

// NewA2ATool 创建 A2ATool
func NewA2ATool(name, description, agentURL, skillName string) *A2ATool {
	return &A2ATool{
		Name:        name,
		Description: description,
		AgentURL:    agentURL,
		SkillName:   skillName,
		Client:      &A2AClient{Endpoint: agentURL},
		Parameters: []map[string]any{
			{"name": "question", "type": "string", "description": "要问的问题", "required": true},
		},
	}
}

// Run 执行工具
func (t *A2ATool) Run(question string) string {
	result := t.Client.ExecuteSkill(t.SkillName, fmt.Sprintf("answer %s", question))
	if status, ok := result["status"].(string); ok && status == "success" {
		if r, ok := result["result"].(string); ok {
			return r
		}
		return "No response"
	}
	return fmt.Sprintf("Error: %v", result["error"])
}

// SimpleAgent 简单智能体（模拟）
type SimpleAgent struct {
	Name         string
	SystemPrompt string
	Tools        map[string]*A2ATool
}

// AddTool 添加工具
func (a *SimpleAgent) AddTool(tool *A2ATool) {
	a.Tools[tool.Name] = tool
}

// Run 运行智能体
func (a *SimpleAgent) Run(input string) string {
	// 模拟智能体根据输入选择工具
	for _, tool := range a.Tools {
		response := tool.Run(input)
		return response
	}
	return "未找到合适的工具"
}

func main() {
	// ============================================================
	// 1. 创建专业 A2A Agent 服务
	// ============================================================

	// 技术专家 Agent
	techExpert := &A2AServer{
		Name:        "tech_expert",
		Description: "技术专家,回答技术相关问题",
		Version:     "1.0.0",
		Skills:      make(map[string]func(string) string),
	}

	techExpert.Skills["answer"] = func(text string) string {
		question := text
		if idx := strings.Index(strings.ToLower(text), "answer "); idx >= 0 {
			question = strings.TrimSpace(text[idx+len("answer "):])
		}
		fmt.Printf("  [技术专家] 回答问题: %s\n", question)
		return fmt.Sprintf("技术回答：关于'%s',这是一个技术问题的专业解答...", question)
	}

	// 销售顾问 Agent
	salesAdvisor := &A2AServer{
		Name:        "sales_advisor",
		Description: "销售顾问,回答销售问题",
		Version:     "1.0.0",
		Skills:      make(map[string]func(string) string),
	}

	salesAdvisor.Skills["answer"] = func(text string) string {
		question := text
		if idx := strings.Index(strings.ToLower(text), "answer "); idx >= 0 {
			question = strings.TrimSpace(text[idx+len("answer "):])
		}
		fmt.Printf("  [销售顾问] 回答问题: %s\n", question)
		return fmt.Sprintf("销售回答：关于'%s',我们有特别优惠...", question)
	}

	// ============================================================
	// 2. 启动 A2A Agent 服务
	// ============================================================

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🚀 启动专业 Agent 服务")
	fmt.Println(strings.Repeat("=", 60))

	go func() {
		techExpert.Run("localhost", 6000)
	}()
	go func() {
		salesAdvisor.Run("localhost", 6001)
	}()

	fmt.Println("✓ 技术专家 Agent 启动在 http://localhost:6000")
	fmt.Println("✓ 销售顾问 Agent 启动在 http://localhost:6001")

	fmt.Println("\n⏳ 等待服务启动...")
	time.Sleep(3 * time.Second)

	// ============================================================
	// 3. 创建 A2A 工具（封装 A2A Agent 为 Tool）
	// ============================================================

	techTool := NewA2ATool(
		"tech_expert",
		"技术专家,回答技术相关问题",
		"http://localhost:6000",
		"answer",
	)

	salesTool := NewA2ATool(
		"sales_advisor",
		"销售顾问,回答销售相关问题",
		"http://localhost:6001",
		"answer",
	)

	// ============================================================
	// 4. 创建 SimpleAgent（使用 A2A 工具）
	// ============================================================

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🤖 创建接待员 SimpleAgent")
	fmt.Println(strings.Repeat("=", 60))

	receptionist := &SimpleAgent{
		Name: "接待员",
		SystemPrompt: `你是客服接待员,负责：
1. 分析客户问题类型（技术问题 or 销售问题）
2. 使用合适的工具（tech_expert 或 sales_advisor）获取答案
3. 整理答案并返回给客户

可用工具：
- tech_expert: 回答技术问题
- sales_advisor: 回答销售问题

请保持礼貌和专业。`,
		Tools: make(map[string]*A2ATool),
	}

	// 添加 A2A 工具
	receptionist.AddTool(techTool)
	receptionist.AddTool(salesTool)

	fmt.Println("✓ 接待员 Agent 创建完成")
	fmt.Println("✓ 已集成 A2A 工具: tech_expert, sales_advisor")

	// ============================================================
	// 5. 测试集成系统
	// ============================================================

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 测试 A2A + SimpleAgent 集成")
	fmt.Println(strings.Repeat("=", 60))

	testQuestions := []string{
		"你们的产品有什么优惠活动吗？",
		"如何配置服务器的SSL证书？",
		"我想了解一下价格方案",
	}

	for i, question := range testQuestions {
		fmt.Printf("\n问题 %d: %s\n", i+1, question)
		fmt.Println(strings.Repeat("-", 60))
		response := receptionist.Run(question)
		fmt.Printf("回答: %s\n", response)
		fmt.Println()
	}

	// ============================================================
	// 6. 保持服务运行
	// ============================================================

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("💡 系统仍在运行")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("你可以继续测试或按 Ctrl+C 停止\n")
}
