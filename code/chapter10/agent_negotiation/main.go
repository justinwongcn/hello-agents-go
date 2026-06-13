// 10_AgentNegotiation - Agent 间协商
//
// 对应 Python: 10_AgentNegotiation.py
// 高级用法：展示两个 Agent 之间的协商流程

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

func main() {
	// 创建两个需要协商的Agent
	agent1 := &A2AServer{
		Name:        "agent1",
		Description: "Agent 1",
		Skills:      make(map[string]func(string) string),
	}

	agent1.Skills["propose"] = func(text string) string {
		// 处理协商提案
		proposalStr := text
		if idx := strings.Index(strings.ToLower(text), "propose "); idx >= 0 {
			proposalStr = strings.TrimSpace(text[idx+len("propose "):])
		}

		// 评估提案（模拟）
		if strings.Contains(proposalStr, "deadline") {
			// 模拟评估：假设需要至少7天
			return `{"accepted":true,"message":"接受提案"}`
		}
		return `{"accepted":false,"message":"无效的提案格式"}`
	}

	agent2 := &A2AServer{
		Name:        "agent2",
		Description: "Agent 2",
		Skills:      make(map[string]func(string) string),
	}

	agent2.Skills["negotiate"] = func(text string) string {
		// 发起协商
		task := "未知任务"
		deadline := 0
		// 解析任务和截止日期
		if idx := strings.Index(strings.ToLower(text), "negotiate "); idx >= 0 {
			parts := strings.TrimSpace(text[idx+len("negotiate "):])
			if tIdx := strings.Index(parts, "task:"); tIdx >= 0 {
				taskPart := parts[tIdx+len("task:"):]
				if dIdx := strings.Index(taskPart, " deadline:"); dIdx >= 0 {
					task = strings.TrimSpace(taskPart[:dIdx])
					fmt.Sscanf(taskPart[dIdx+len(" deadline:"):], "%d", &deadline)
				}
			}
		}
		proposal := fmt.Sprintf(`{"task":"%s","deadline":%d}`, task, deadline)
		return fmt.Sprintf(`{"status":"negotiating","proposal":%s}`, proposal)
	}

	// 启动服务
	go agent1.Run("localhost", 7000)
	go agent2.Run("localhost", 7001)
	time.Sleep(2 * time.Second)

	// 测试协商流程
	client1 := &A2AClient{Endpoint: "http://localhost:7000"}
	client2 := &A2AClient{Endpoint: "http://localhost:7001"}

	// Agent2发起协商
	negotiation := client2.ExecuteSkill("negotiate", "negotiate task:开发新功能 deadline:5")
	fmt.Printf("协商请求：%v\n", negotiation["result"])

	// Agent1评估提案
	proposal := client1.ExecuteSkill("propose", "propose {'task': '开发新功能', 'deadline': 5}")
	fmt.Printf("提案评估：%v\n", proposal["result"])

	// 保持服务运行
	fmt.Println("服务已停止")
}
