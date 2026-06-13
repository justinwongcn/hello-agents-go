package common

import (
	"fmt"
	"strings"
)

// Tool 是自定义工具的接口
type Tool interface {
	Run(input string) string
}

// ToolFunc 将普通函数适配为 Tool 接口
type ToolFunc struct {
	Name        string
	Description string
	Fn          func(string) string
}

// Run 执行工具函数
func (t *ToolFunc) Run(input string) string {
	return t.Fn(input)
}

// toolEntry 工具注册表中的条目
type toolEntry struct {
	tool Tool
}

// ToolRegistry 工具注册表，用于管理和执行工具
type ToolRegistry struct {
	tools map[string]*toolEntry
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*toolEntry),
	}
}

// RegisterFunction 注册一个函数作为工具
func (r *ToolRegistry) RegisterFunction(name, description string, fn func(string) string) {
	r.tools[name] = &toolEntry{
		tool: &ToolFunc{Name: name, Description: description, Fn: fn},
	}
	fmt.Printf("🔧 工具 '%s' 已注册。\n", name)
}

// RegisterTool 注册一个工具实例
func (r *ToolRegistry) RegisterTool(t Tool) {
	name := ""
	switch v := t.(type) {
	case *ToolFunc:
		name = v.Name
	case *Calculator:
		name = v.Name
	case *MyAdvancedSearchTool:
		name = v.Name
	default:
		name = fmt.Sprintf("tool_%d", len(r.tools))
	}
	r.tools[name] = &toolEntry{tool: t}
	fmt.Printf("🔧 工具 '%s' 已注册。\n", name)
}

// GetTool 根据名称获取工具
func (r *ToolRegistry) GetTool(name string) Tool {
	if entry, ok := r.tools[name]; ok {
		return entry.tool
	}
	return nil
}

// ExecuteTool 执行指定名称的工具
func (r *ToolRegistry) ExecuteTool(name, input string) string {
	entry, ok := r.tools[name]
	if !ok {
		return fmt.Sprintf("❌ 错误：未找到工具 '%s'", name)
	}
	return entry.tool.Run(input)
}

// GetToolsDescription 获取所有工具的格式化描述
func (r *ToolRegistry) GetToolsDescription() string {
	if len(r.tools) == 0 {
		return "暂无可用工具"
	}

	var lines []string
	for name, entry := range r.tools {
		desc := ""
		switch v := entry.tool.(type) {
		case *ToolFunc:
			desc = v.Description
		case *Calculator:
			desc = v.Description
		case *MyAdvancedSearchTool:
			desc = v.Description
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", name, desc))
	}
	return strings.Join(lines, "\n")
}

// Unregister 移除指定工具
func (r *ToolRegistry) Unregister(name string) {
	delete(r.tools, name)
}

// ListTools 列出所有工具名称
func (r *ToolRegistry) ListTools() []string {
	var names []string
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// CreateCalculatorRegistry 创建包含计算器的工具注册表
func CreateCalculatorRegistry() *ToolRegistry {
	registry := NewToolRegistry()

	calculator := NewCalculator()
	registry.RegisterFunction(
		"my_calculator",
		"简单的数学计算工具，支持基本运算(+,-,*,/)和sqrt函数",
		calculator.Run,
	)

	return registry
}
