// 07_SimpleA2AAgent - 简单 A2A 智能体
//
// 对应 Python: 07_SimpleA2AAgent.py
// 创建一个简单的 A2A 计算器智能体，支持基本数学运算

package main

import (
	"fmt"
	"strings"
)

// A2AServer A2A 服务器（模拟）
type A2AServer struct {
	Name         string
	Description  string
	Version      string
	Capabilities map[string][]string
	Skills       map[string]func(string) string
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
}

// createCalculatorAgent 创建一个计算器智能体
func createCalculatorAgent() *A2AServer {
	fmt.Println("🧮 创建计算器智能体")

	// 创建 A2A 服务器
	calculator := &A2AServer{
		Name:        "calculator-agent",
		Description: "专业的数学计算智能体",
		Version:     "1.0.0",
		Capabilities: map[string][]string{
			"math":     {"addition", "subtraction", "multiplication", "division"},
			"advanced": {"power", "sqrt", "factorial"},
		},
		Skills: make(map[string]func(string) string),
	}

	// 添加基础计算技能
	calculator.Skills["add"] = func(query string) string {
		// 简单解析 "计算 5 + 3" 格式
		parts := query
		parts = strings.ReplaceAll(parts, "计算", "")
		parts = strings.ReplaceAll(parts, "加", "+")
		parts = strings.ReplaceAll(parts, "加上", "+")
		if strings.Contains(parts, "+") {
			nums := strings.Split(parts, "+")
			var sum float64
			var numStrs []string
			for _, n := range nums {
				var val float64
				fmt.Sscanf(strings.TrimSpace(n), "%f", &val)
				sum += val
				numStrs = append(numStrs, strings.TrimSpace(n))
			}
			return fmt.Sprintf("计算结果: %s = %g", strings.Join(numStrs, " + "), sum)
		}
		return "请使用格式: 计算 5 + 3"
	}

	calculator.Skills["multiply"] = func(query string) string {
		parts := query
		parts = strings.ReplaceAll(parts, "计算", "")
		parts = strings.ReplaceAll(parts, "乘以", "*")
		parts = strings.ReplaceAll(parts, "×", "*")
		if strings.Contains(parts, "*") {
			nums := strings.Split(parts, "*")
			result := 1.0
			var numStrs []string
			for _, n := range nums {
				var val float64
				fmt.Sscanf(strings.TrimSpace(n), "%f", &val)
				result *= val
				numStrs = append(numStrs, strings.TrimSpace(n))
			}
			return fmt.Sprintf("计算结果: %s = %g", strings.Join(numStrs, " × "), result)
		}
		return "请使用格式: 计算 5 * 3"
	}

	calculator.Skills["info"] = func(query string) string {
		skillNames := make([]string, 0, len(calculator.Skills))
		for k := range calculator.Skills {
			skillNames = append(skillNames, k)
		}
		return fmt.Sprintf("我是 %s，可以进行基础数学计算。支持的技能: %v", calculator.Name, skillNames)
	}

	skillNames := make([]string, 0, len(calculator.Skills))
	for k := range calculator.Skills {
		skillNames = append(skillNames, k)
	}
	fmt.Printf("✅ 计算器智能体创建成功，支持技能: %v\n", skillNames)
	return calculator
}

func main() {
	calcAgent := createCalculatorAgent()
	if calcAgent != nil {
		// 测试技能
		fmt.Println("\n🧪 测试智能体技能:")
		testQueries := []string{
			"获取信息",
			"计算 10 + 5",
			"计算 6 * 7",
		}

		for _, query := range testQueries {
			var result string
			if strings.Contains(query, "信息") {
				result = calcAgent.Skills["info"](query)
			} else if strings.Contains(query, "+") {
				result = calcAgent.Skills["add"](query)
			} else if strings.Contains(query, "*") || strings.Contains(query, "×") {
				result = calcAgent.Skills["multiply"](query)
			} else {
				result = "未知查询类型"
			}

			fmt.Printf("  📝 查询: %s\n", query)
			fmt.Printf("  🤖 回复: %s\n", result)
			fmt.Println()
		}
	}
}
