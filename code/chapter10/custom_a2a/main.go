// 08_CustomA2AAgent - 自定义 A2A 智能体
//
// 对应 Python: 08_CustomA2AAgent.py
// 创建自定义 A2A 智能体，支持问候和计算技能

package main

import "fmt"

// A2AServer A2A 服务器（模拟）
type A2AServer struct {
	Name         string
	Description  string
	Capabilities map[string][]string
	Skills       map[string]func(string) string
}

// Skill 添加技能（装饰器风格）
func (s *A2AServer) Skill(name string) func(func(string) string) {
	return func(fn func(string) string) {
		s.Skills[name] = fn
	}
}

// createCustomAgent 创建自定义智能体
func createCustomAgent() *A2AServer {
	// 创建智能体
	agent := &A2AServer{
		Name:        "my-custom-agent",
		Description: "我的自定义智能体",
		Capabilities: map[string][]string{
			"custom": {"skill1", "skill2"},
		},
		Skills: make(map[string]func(string) string),
	}

	// 添加技能
	agent.Skills["greet"] = func(name string) string {
		// 问候用户
		return fmt.Sprintf("你好，%s！我是自定义智能体。", name)
	}

	agent.Skills["calculate"] = func(expression string) string {
		// 简单计算（仅支持基本运算）
		allowedChars := "0123456789+-*/(). "
		for _, c := range expression {
			found := false
			for _, a := range allowedChars {
				if c == a {
					found = true
					break
				}
			}
			if !found {
				return "错误: 只支持基本数学运算"
			}
		}
		// 简单表达式计算（模拟）
		return fmt.Sprintf("计算结果: %s = (模拟计算结果)", expression)
	}

	return agent
}

func main() {
	customAgent := createCustomAgent()
	if customAgent != nil {
		// 测试技能
		fmt.Println("测试问候技能:")
		result1 := customAgent.Skills["greet"]("张三")
		fmt.Println(result1)

		fmt.Println("\n测试计算技能:")
		result2 := customAgent.Skills["calculate"]("10 + 5 * 2")
		fmt.Println(result2)
	}
}
