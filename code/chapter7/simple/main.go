package main

import (
	"fmt"
	"regexp"
	"strings"

	"hello-agents-go/chapter4/llmclient"
	"hello-agents-go/chapter7/common"
)

// 工具调用正则表达式（预编译为包级变量）
var reToolCall = regexp.MustCompile(`\[TOOL_CALL:([^:]+):([^\]]+)\]`)

// Message 表示一条对话消息
type Message struct {
	Role    string
	Content string
}

// MySimpleAgent 重写的简单对话Agent
// 展示如何基于框架基类构建自定义Agent
type MySimpleAgent struct {
	Name               string
	LLM                *llmclient.HelloAgentsLLM
	SystemPrompt       string
	ToolRegistry       *common.ToolRegistry
	EnableToolCalling  bool
	History            []Message
}

// NewMySimpleAgent 创建简单对话Agent
func NewMySimpleAgent(name string, llm *llmclient.HelloAgentsLLM, systemPrompt string, toolRegistry *common.ToolRegistry, enableToolCalling bool) *MySimpleAgent {
	enableToolCalling = enableToolCalling && toolRegistry != nil
	status := "禁用"
	if enableToolCalling {
		status = "启用"
	}
	fmt.Printf("✅ %s 初始化完成，工具调用: %s\n", name, status)

	return &MySimpleAgent{
		Name:              name,
		LLM:               llm,
		SystemPrompt:      systemPrompt,
		ToolRegistry:      toolRegistry,
		EnableToolCalling: enableToolCalling,
	}
}

// Run 运行Agent - 实现简单对话逻辑，支持可选工具调用
func (a *MySimpleAgent) Run(inputText string, maxToolIterations int) string {
	if maxToolIterations == 0 {
		maxToolIterations = 3
	}

	fmt.Printf("🤖 %s 正在处理: %s\n", a.Name, inputText)

	// 构建消息列表
	var messages []llmclient.Message

	// 添加系统消息（可能包含工具信息）
	enhancedSystemPrompt := a.getEnhancedSystemPrompt()
	messages = append(messages, llmclient.Message{Role: "system", Content: enhancedSystemPrompt})

	// 添加历史消息
	for _, msg := range a.History {
		messages = append(messages, llmclient.Message{Role: msg.Role, Content: msg.Content})
	}

	// 添加当前用户消息
	messages = append(messages, llmclient.Message{Role: "user", Content: inputText})

	// 如果没有启用工具调用，使用简单对话逻辑
	if !a.EnableToolCalling {
		response := a.LLM.Think(messages, 0)
		a.addMessage(inputText, "user")
		a.addMessage(response, "assistant")
		fmt.Printf("✅ %s 响应完成\n", a.Name)
		return response
	}

	// 支持多轮工具调用的逻辑
	return a.runWithTools(messages, inputText, maxToolIterations)
}

// getEnhancedSystemPrompt 构建增强的系统提示词，包含工具信息
func (a *MySimpleAgent) getEnhancedSystemPrompt() string {
	basePrompt := a.SystemPrompt
	if basePrompt == "" {
		basePrompt = "你是一个有用的AI助手。"
	}

	if !a.EnableToolCalling || a.ToolRegistry == nil {
		return basePrompt
	}

	// 获取工具描述
	toolsDescription := a.ToolRegistry.GetToolsDescription()
	if toolsDescription == "" || toolsDescription == "暂无可用工具" {
		return basePrompt
	}

	toolsSection := "\n\n## 可用工具\n"
	toolsSection += "你可以使用以下工具来帮助回答问题：\n"
	toolsSection += toolsDescription + "\n"

	toolsSection += "\n## 工具调用格式\n"
	toolsSection += "当需要使用工具时，请使用以下格式：\n"
	toolsSection += "`[TOOL_CALL:{tool_name}:{parameters}]`\n"
	toolsSection += "例如：`[TOOL_CALL:search:Python编程]` 或 `[TOOL_CALL:memory:recall=用户信息]`\n\n"
	toolsSection += "工具调用结果会自动插入到对话中，然后你可以基于结果继续回答。\n"

	return basePrompt + toolsSection
}

// runWithTools 支持工具调用的运行逻辑
func (a *MySimpleAgent) runWithTools(messages []llmclient.Message, inputText string, maxToolIterations int) string {
	currentIteration := 0
	finalResponse := ""

	for currentIteration < maxToolIterations {
		// 调用LLM
		response := a.LLM.Think(messages, 0)

		// 检查是否有工具调用
		toolCalls := a.parseToolCalls(response)

		if len(toolCalls) > 0 {
			fmt.Printf("🔧 检测到 %d 个工具调用\n", len(toolCalls))
			// 执行所有工具调用并收集结果
			var toolResults []string
			cleanResponse := response

			for _, call := range toolCalls {
				result := a.executeToolCall(call["tool_name"], call["parameters"])
				toolResults = append(toolResults, result)
				// 从响应中移除工具调用标记
				cleanResponse = strings.ReplaceAll(cleanResponse, call["original"], "")
			}

			// 构建包含工具结果的消息
			messages = append(messages, llmclient.Message{Role: "assistant", Content: cleanResponse})

			// 添加工具结果
			toolResultsText := strings.Join(toolResults, "\n\n")
			messages = append(messages, llmclient.Message{
				Role:    "user",
				Content: fmt.Sprintf("工具执行结果：\n%s\n\n请基于这些结果给出完整的回答。", toolResultsText),
			})

			currentIteration++
			continue
		}

		// 没有工具调用，这是最终回答
		finalResponse = response
		break
	}

	// 如果超过最大迭代次数，获取最后一次回答
	if currentIteration >= maxToolIterations && finalResponse == "" {
		finalResponse = a.LLM.Think(messages, 0)
	}

	// 保存到历史记录
	a.addMessage(inputText, "user")
	a.addMessage(finalResponse, "assistant")
	fmt.Printf("✅ %s 响应完成\n", a.Name)

	return finalResponse
}

// parseToolCalls 解析文本中的工具调用
func (a *MySimpleAgent) parseToolCalls(text string) []map[string]string {
	matches := reToolCall.FindAllStringSubmatch(text, -1)

	var toolCalls []map[string]string
	for _, match := range matches {
		if len(match) >= 3 {
			toolCalls = append(toolCalls, map[string]string{
				"tool_name":  strings.TrimSpace(match[1]),
				"parameters": strings.TrimSpace(match[2]),
				"original":   match[0],
			})
		}
	}

	return toolCalls
}

// executeToolCall 执行工具调用
func (a *MySimpleAgent) executeToolCall(toolName, parameters string) string {
	if a.ToolRegistry == nil {
		return "❌ 错误：未配置工具注册表"
	}

	// 智能参数解析
	if toolName == "calculator" {
		// 计算器工具直接传入表达式
		result := a.ToolRegistry.ExecuteTool(toolName, parameters)
		return fmt.Sprintf("🔧 工具 %s 执行结果：\n%s", toolName, result)
	}

	// 其他工具使用智能参数解析
	paramDict := a.parseToolParameters(toolName, parameters)
	tool := a.ToolRegistry.GetTool(toolName)
	if tool == nil {
		return fmt.Sprintf("❌ 错误：未找到工具 '%s'", toolName)
	}

	// 将参数字典转换为输入字符串
	input := parameters
	if query, ok := paramDict["query"]; ok {
		input = query
	} else if inputVal, ok := paramDict["input"]; ok {
		input = inputVal
	}

	result := tool.Run(input)
	return fmt.Sprintf("🔧 工具 %s 执行结果：\n%s", toolName, result)
}

// parseToolParameters 智能解析工具参数
func (a *MySimpleAgent) parseToolParameters(toolName, parameters string) map[string]string {
	paramDict := make(map[string]string)

	if strings.Contains(parameters, "=") {
		// 格式: key=value 或 action=search,query=Python
		if strings.Contains(parameters, ",") {
			// 多个参数：action=search,query=Python,limit=3
			pairs := strings.Split(parameters, ",")
			for _, pair := range pairs {
				if strings.Contains(pair, "=") {
					parts := strings.SplitN(pair, "=", 2)
					paramDict[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		} else {
			// 单个参数：key=value
			parts := strings.SplitN(parameters, "=", 2)
			paramDict[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	} else {
		// 直接传入参数，根据工具类型智能推断
		if toolName == "search" {
			paramDict["query"] = parameters
		} else if toolName == "memory" {
			paramDict["action"] = "search"
			paramDict["query"] = parameters
		} else {
			paramDict["input"] = parameters
		}
	}

	return paramDict
}

// StreamRun 自定义的流式运行方法
func (a *MySimpleAgent) StreamRun(inputText string) string {
	fmt.Printf("🌊 %s 开始流式处理: %s\n", a.Name, inputText)

	var messages []llmclient.Message

	if a.SystemPrompt != "" {
		messages = append(messages, llmclient.Message{Role: "system", Content: a.SystemPrompt})
	}

	for _, msg := range a.History {
		messages = append(messages, llmclient.Message{Role: msg.Role, Content: msg.Content})
	}

	messages = append(messages, llmclient.Message{Role: "user", Content: inputText})

	// 流式调用LLM
	var fullResponse strings.Builder
	fmt.Print("📝 实时响应: ")
	for chunk := range a.LLM.ThinkStream(messages, 0) {
		fullResponse.WriteString(chunk)
		fmt.Print(chunk)
	}
	fmt.Println() // 换行

	// 保存完整对话到历史记录
	a.addMessage(inputText, "user")
	a.addMessage(fullResponse.String(), "assistant")
	fmt.Printf("✅ %s 流式响应完成\n", a.Name)

	return fullResponse.String()
}

// AddTool 添加工具到Agent（便利方法）
func (a *MySimpleAgent) AddTool(t common.Tool) {
	if a.ToolRegistry == nil {
		a.ToolRegistry = common.NewToolRegistry()
		a.EnableToolCalling = true
	}
	a.ToolRegistry.RegisterTool(t)
}

// HasTools 检查是否有可用工具
func (a *MySimpleAgent) HasTools() bool {
	return a.EnableToolCalling && a.ToolRegistry != nil
}

// RemoveTool 移除工具（便利方法）
func (a *MySimpleAgent) RemoveTool(toolName string) bool {
	if a.ToolRegistry != nil {
		a.ToolRegistry.Unregister(toolName)
		return true
	}
	return false
}

// ListTools 列出所有可用工具
func (a *MySimpleAgent) ListTools() []string {
	if a.ToolRegistry != nil {
		return a.ToolRegistry.ListTools()
	}
	return nil
}

// addMessage 添加消息到历史记录
func (a *MySimpleAgent) addMessage(content, role string) {
	a.History = append(a.History, Message{Role: role, Content: content})
}

func main() {
	// 使用示例（需要配置 .env 文件和 API 密钥才能实际运行）
	fmt.Println("=== MySimpleAgent 示例 ===")
	fmt.Println("请确保已配置 .env 文件中的 LLM_API_KEY 等环境变量。")

	// 创建LLM客户端
	llm, err := llmclient.NewHelloAgentsLLM("", "", "", 0)
	if err != nil {
		fmt.Printf("❌ 创建LLM客户端失败: %v\n", err)
		fmt.Println("请配置 .env 文件后重试。")
		return
	}

	// 创建工具注册表
	registry := common.CreateCalculatorRegistry()

	// 创建Agent
	agent := NewMySimpleAgent("MySimpleAgent", llm, "你是一个有用的AI助手，擅长数学计算。", registry, true)

	// 运行Agent
	response := agent.Run("请帮我计算 sqrt(144) + 3 * 7", 3)
	fmt.Printf("\n📋 最终回答: %s\n", response)
}
