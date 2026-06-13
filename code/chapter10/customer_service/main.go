// 10_CustomerService - 智能客服系统
//
// 对应 Python: 10_CustomerService.py
// 实战案例：使用 A2A 工具构建智能客服系统

package main

import (
	"fmt"
	"strings"
	"time"
)

// A2AServer A2A 服务器（模拟）
type A2AServer struct {
	Name        string
	Description string
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

// A2ATool A2A 工具（模拟）
type A2ATool struct {
	Name        string
	Description string
	AgentURL    string
	Client      *A2AClient
}

// NewA2ATool 创建 A2ATool
func NewA2ATool(name, description, agentURL string) *A2ATool {
	return &A2ATool{
		Name:        name,
		Description: description,
		AgentURL:    agentURL,
		Client:      &A2AClient{Endpoint: agentURL},
	}
}

// Run 执行工具
func (t *A2ATool) Run(question string) string {
	result := t.Client.ExecuteSkill("answer", fmt.Sprintf("answer %s", question))
	if r, ok := result["result"].(string); ok {
		return r
	}
	return "No response"
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

// handleCustomerQuery 处理客户咨询
func handleCustomerQuery(agent *SimpleAgent, query string) {
	fmt.Printf("\n客户咨询：%s\n", query)
	fmt.Println(strings.Repeat("=", 50))
	response := agent.Run(query)
	fmt.Printf("\n客服回复：%s\n", response)
	fmt.Println(strings.Repeat("=", 50))
}

func main() {
	// 1. 创建技术专家Agent服务
	techExpert := &A2AServer{
		Name:        "tech_expert",
		Description: "技术专家,回答技术问题",
		Skills:      make(map[string]func(string) string),
	}

	techExpert.Skills["answer"] = func(text string) string {
		question := text
		if idx := strings.Index(strings.ToLower(text), "answer "); idx >= 0 {
			question = strings.TrimSpace(text[idx+len("answer "):])
		}
		// 实际应用中,这里会调用LLM或知识库
		return fmt.Sprintf("技术回答：关于'%s',我建议您查看我们的技术文档...", question)
	}

	// 2. 创建销售顾问Agent服务
	salesAdvisor := &A2AServer{
		Name:        "sales_advisor",
		Description: "销售顾问,回答销售问题",
		Skills:      make(map[string]func(string) string),
	}

	salesAdvisor.Skills["answer"] = func(text string) string {
		question := text
		if idx := strings.Index(strings.ToLower(text), "answer "); idx >= 0 {
			question = strings.TrimSpace(text[idx+len("answer "):])
		}
		return fmt.Sprintf("销售回答：关于'%s',我们有特别优惠...", question)
	}

	// 3. 启动服务
	go techExpert.Run("localhost", 6000)
	go salesAdvisor.Run("localhost", 6001)
	time.Sleep(2 * time.Second)

	// 4. 创建接待员Agent（使用HelloAgents的SimpleAgent）
	receptionist := &SimpleAgent{
		Name: "接待员",
		SystemPrompt: `你是客服接待员,负责：
1. 分析客户问题类型（技术问题 or 销售问题）
2. 将问题转发给相应的专家
3. 整理专家的回答并返回给客户

请保持礼貌和专业。`,
		Tools: make(map[string]*A2ATool),
	}

	// 添加技术专家工具
	techTool := NewA2ATool(
		"tech_expert",
		"技术专家,回答技术相关问题",
		"http://localhost:6000",
	)
	receptionist.AddTool(techTool)

	// 添加销售顾问工具
	salesTool := NewA2ATool(
		"sales_advisor",
		"销售顾问,回答价格、购买相关问题",
		"http://localhost:6001",
	)
	receptionist.AddTool(salesTool)

	// 5. 处理客户咨询
	handleCustomerQuery(receptionist, "你们的API如何调用？")
	handleCustomerQuery(receptionist, "企业版的价格是多少？")
	handleCustomerQuery(receptionist, "如何集成到我的Python项目中？")
}
