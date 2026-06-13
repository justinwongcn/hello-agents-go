// 09_A2A_Network - 创建 Agent 网络
//
// 对应 Python: 09_A2A_Network.py
// 创建多个 Agent 服务并协作完成任务

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
	// 模拟服务器运行
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
	// 1. 创建多个Agent服务
	researcher := &A2AServer{
		Name:        "researcher",
		Description: "研究员",
		Skills:      make(map[string]func(string) string),
	}

	researcher.Skills["research"] = func(text string) string {
		topic := text
		if idx := strings.Index(strings.ToLower(text), "research "); idx >= 0 {
			topic = strings.TrimSpace(text[idx+len("research "):])
		}
		return fmt.Sprintf(`{"topic":"%s","findings":"%s的研究结果"}`, topic, topic)
	}

	writer := &A2AServer{
		Name:        "writer",
		Description: "撰写员",
		Skills:      make(map[string]func(string) string),
	}

	writer.Skills["write"] = func(text string) string {
		content := text
		if idx := strings.Index(strings.ToLower(text), "write "); idx >= 0 {
			content = strings.TrimSpace(text[idx+len("write "):])
		}
		// 尝试解析研究数据
		topic := "未知主题"
		findings := content
		if strings.Contains(content, "topic") {
			topic = "AI在医疗领域的应用"
			findings = "AI在医疗领域的应用的研究结果"
		}
		return fmt.Sprintf("# %s\n\n基于研究：%s\n\n文章内容...", topic, findings)
	}

	editor := &A2AServer{
		Name:        "editor",
		Description: "编辑",
		Skills:      make(map[string]func(string) string),
	}

	editor.Skills["edit"] = func(text string) string {
		article := text
		if idx := strings.Index(strings.ToLower(text), "edit "); idx >= 0 {
			article = strings.TrimSpace(text[idx+len("edit "):])
		}
		return fmt.Sprintf(`{"article":"%s\n\n[已编辑优化]","feedback":"文章质量良好","approved":true}`, article)
	}

	// 2. 启动所有服务
	go researcher.Run("localhost", 5000)
	go writer.Run("localhost", 5001)
	go editor.Run("localhost", 5002)
	time.Sleep(2 * time.Second) // 等待服务启动

	// 3. 创建客户端连接到各个Agent
	researcherClient := &A2AClient{Endpoint: "http://localhost:5000"}
	writerClient := &A2AClient{Endpoint: "http://localhost:5001"}
	editorClient := &A2AClient{Endpoint: "http://localhost:5002"}

	// 4. 协作流程
	createContent := func(topic string) string {
		// 步骤1：研究
		research := researcherClient.ExecuteSkill("research", fmt.Sprintf("research %s", topic))
		researchData, _ := research["result"].(string)

		// 步骤2：撰写
		article := writerClient.ExecuteSkill("write", fmt.Sprintf("write %s", researchData))
		articleContent, _ := article["result"].(string)

		// 步骤3：编辑
		final := editorClient.ExecuteSkill("edit", fmt.Sprintf("edit %s", articleContent))
		result, _ := final["result"].(string)
		return result
	}

	// 使用
	result := createContent("AI在医疗领域的应用")
	fmt.Printf("\n最终结果：\n%s\n", result)
}
