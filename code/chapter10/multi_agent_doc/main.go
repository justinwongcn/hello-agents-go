// 06_MultiAgentDocumentAssist - 多Agent协作的智能文档助手
//
// 对应 Python: 06_MultiAgentDocumentAssist.py
// 使用两个SimpleAgent分工协作：
//   - Agent1：GitHub搜索专家
//   - Agent2：文档生成专家

package main

import (
	"fmt"
	"os"
	"strings"
)

// SimpleAgent 简单智能体
type SimpleAgent struct {
	Name         string
	SystemPrompt string
}

// Run 运行智能体
func (a *SimpleAgent) Run(input string) string {
	// 模拟智能体根据系统提示处理输入
	if strings.Contains(a.Name, "搜索") {
		return fmt.Sprintf("搜索结果（模拟）：\n"+
			"1. langchain - LLM应用开发框架\n"+
			"2. auto-gpt - 自主AI Agent\n"+
			"3. agentgpt - 浏览器端AI Agent\n"+
			"4. metagpt - 多Agent框架\n"+
			"5. babyagi - 任务驱动的AI Agent")
	}
	// 文档生成专家
	return fmt.Sprintf("# AI Agent框架研究报告\n\n"+
		"## 简介\n"+
		"这是关于AI Agent的GitHub项目调研报告。\n\n"+
		"## 主要发现\n"+
		"1. **langchain** - LLM应用开发框架\n"+
		"2. **auto-gpt** - 自主AI Agent\n"+
		"3. **agentgpt** - 浏览器端AI Agent\n"+
		"4. **metaGPT** - 多Agent框架\n"+
		"5. **babyagi** - 任务驱动的AI Agent\n\n"+
		"## 总结\n"+
		"这些项目的共同特点是利用LLM构建自主Agent系统。")
}

// MCPTool MCP 工具
type MCPTool struct {
	Name          string
	Description   string
	ServerCommand []string
}

func main() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("多Agent协作的智能文档助手")
	fmt.Println(strings.Repeat("=", 70))

	// ============================================================
	// Agent 1: GitHub搜索专家
	// ============================================================
	fmt.Println("\n【步骤1】创建GitHub搜索专家...")

	githubSearcher := &SimpleAgent{
		Name: "GitHub搜索专家",
		SystemPrompt: `你是一个GitHub搜索专家。
你的任务是搜索GitHub仓库并返回结果。
请返回清晰、结构化的搜索结果，包括：
- 仓库名称
- 简短描述

保持简洁，不要添加额外的解释。`,
	}

	// 添加GitHub工具
	githubTool := &MCPTool{
		Name:          "gh",
		ServerCommand: []string{"npx", "-y", "@modelcontextprotocol/server-github"},
	}
	_ = githubTool

	// ============================================================
	// Agent 2: 文档生成专家
	// ============================================================
	fmt.Println("\n【步骤2】创建文档生成专家...")

	documentWriter := &SimpleAgent{
		Name: "文档生成专家",
		SystemPrompt: `你是一个文档生成专家。
你的任务是根据提供的信息生成结构化的Markdown报告。

报告应该包括：
- 标题
- 简介
- 主要内容（分点列出，包括项目名称、描述等）
- 总结

请直接输出完整的Markdown格式报告内容，不要使用工具保存。`,
	}

	// 添加文件系统工具
	fsTool := &MCPTool{
		Name:          "fs",
		ServerCommand: []string{"npx", "-y", "@modelcontextprotocol/server-filesystem", "."},
	}
	_ = fsTool

	// ============================================================
	// 执行任务
	// ============================================================
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("开始执行任务...")
	fmt.Println(strings.Repeat("=", 70))

	// 步骤1：GitHub搜索
	fmt.Println("\n【步骤3】Agent1 搜索GitHub...")
	searchTask := "搜索关于'AI agent'的GitHub仓库，返回前5个最相关的结果"

	searchResults := githubSearcher.Run(searchTask)

	fmt.Println("\n搜索结果:")
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println(searchResults)
	fmt.Println(strings.Repeat("-", 70))

	// 步骤2：生成报告
	fmt.Println("\n【步骤4】Agent2 生成报告...")
	reportTask := fmt.Sprintf(`
根据以下GitHub搜索结果，生成一份Markdown格式的研究报告：

%s

报告要求：
1. 标题：# AI Agent框架研究报告
2. 简介：说明这是关于AI Agent的GitHub项目调研
3. 主要发现：列出找到的项目及其特点（包括名称、描述等）
4. 总结：总结这些项目的共同特点

请直接输出完整的Markdown格式报告。
`, searchResults)

	reportContent := documentWriter.Run(reportTask)

	fmt.Println("\n报告内容:")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(reportContent)
	fmt.Println(strings.Repeat("=", 70))

	// 步骤3：保存报告
	fmt.Println("\n【步骤5】保存报告到文件...")
	err := os.WriteFile("report.md", []byte(reportContent), 0644)
	if err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
	} else {
		fmt.Println("✅ 报告已保存到 report.md")

		// 验证文件
		info, err := os.Stat("report.md")
		if err == nil {
			fmt.Printf("✅ 文件大小: %d 字节\n", info.Size())
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("任务完成！")
	fmt.Println(strings.Repeat("=", 70))
}
