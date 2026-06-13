// TerminalTool 使用示例
//
// 展示 TerminalTool 的典型使用模式：
// 1. 探索式导航
// 2. 数据文件分析
// 3. 日志文件分析
// 4. 代码库分析

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ScriptDir 脚本所在目录
var ScriptDir string

func init() {
	ex, err := os.Executable()
	if err != nil {
		ScriptDir, _ = os.Getwd()
	} else {
		ScriptDir = filepath.Dir(ex)
	}
	// 如果是从 go run 运行，使用当前工作目录
	if cwd, err := os.Getwd(); err == nil {
		ScriptDir = cwd
	}
}

// TerminalTool 终端工具
type TerminalTool struct {
	Workspace string
	Timeout   int
}

// NewTerminalTool 创建 TerminalTool
func NewTerminalTool(workspace string, timeout ...int) *TerminalTool {
	t := 60
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &TerminalTool{Workspace: workspace, Timeout: t}
}

// Run 执行终端命令
func (tt *TerminalTool) Run(params map[string]any) string {
	command, _ := params["command"].(string)

	// 安全检查: 拒绝危险命令
	if tt.isDangerousCommand(command) {
		return fmt.Sprintf("🚫 安全拒绝: 命令 '%s' 被安全策略阻止", command)
	}

	// 安全检查: 检查是否尝试逃逸工作目录
	if tt.isPathEscape(command) {
		return "🚫 安全拒绝: 尝试访问工作目录外的文件"
	}

	// 执行命令
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = tt.Workspace
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s\n[命令执行完成，退出码: %d]", string(output), cmd.ProcessState.ExitCode())
	}
	return string(output)
}

// isDangerousCommand 检查危险命令
func (tt *TerminalTool) isDangerousCommand(command string) bool {
	dangerous := []string{"rm -rf /", "rm -rf /*", "mkfs", "dd if=", ":(){ :|:& };:"}
	for _, d := range dangerous {
		if strings.Contains(command, d) {
			return true
		}
	}
	return false
}

// isPathEscape 检查路径逃逸
func (tt *TerminalTool) isPathEscape(command string) bool {
	// 简单检查是否包含大量 ../ 尝试逃逸
	if strings.Contains(command, "cd ../../../") {
		return true
	}
	return false
}

func demoExploratoryNavigation() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("场景1: 探索式导航")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	terminal := NewTerminalTool(ScriptDir)

	// 第一步:查看当前目录
	fmt.Println("1. 查看当前目录:")
	result := terminal.Run(map[string]any{"command": "ls -la"})
	fmt.Println(result)

	// 第二步:查看Go文件
	fmt.Println("\n2. 查看Go文件:")
	result = terminal.Run(map[string]any{"command": "ls -la */main.go 2>/dev/null || echo '未找到main.go文件'"})
	fmt.Println(result)

	// 第三步:查找特定文件
	fmt.Println("\n3. 查找特定模式的文件:")
	result = terminal.Run(map[string]any{"command": "find . -name 'codebase_maintainer' -type d"})
	fmt.Println(result)

	// 第四步:查看文件内容
	fmt.Println("\n4. 查看文件内容:")
	result = terminal.Run(map[string]any{"command": "head -n 20 codebase_maintainer/main.go 2>/dev/null || echo '文件不存在'"})
	fmt.Println(result)
}

func demoDataFileAnalysis() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("场景2: 数据文件分析")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	dataDir := filepath.Join(ScriptDir, "data")
	terminal := NewTerminalTool(dataDir)

	// 查看 CSV 文件的前几行
	fmt.Println("1. 查看 CSV 文件前5行:")
	result := terminal.Run(map[string]any{"command": "head -n 5 sales_2024.csv 2>/dev/null || echo '数据文件不存在（示例演示）'"})
	fmt.Println(result)

	// 统计总行数
	fmt.Println("\n2. 统计文件行数:")
	result = terminal.Run(map[string]any{"command": "wc -l *.csv 2>/dev/null || echo '无CSV文件（示例演示）'"})
	fmt.Println(result)

	// 提取和统计产品类别
	fmt.Println("\n3. 统计产品类别分布:")
	result = terminal.Run(map[string]any{"command": "tail -n +2 sales_2024.csv 2>/dev/null | cut -d',' -f3 | sort | uniq -c 2>/dev/null || echo '数据文件不存在（示例演示）'"})
	fmt.Println(result)
}

func demoLogAnalysis() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("场景3: 日志文件分析")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	logsDir := filepath.Join(ScriptDir, "logs")
	terminal := NewTerminalTool(logsDir)

	// 查看最新的错误日志
	fmt.Println("1. 查看最新的错误日志:")
	result := terminal.Run(map[string]any{"command": "tail -n 50 app.log 2>/dev/null | grep ERROR || echo '日志文件不存在（示例演示）'"})
	fmt.Println(result)

	// 统计错误类型分布
	fmt.Println("\n2. 统计错误类型分布:")
	result = terminal.Run(map[string]any{"command": "grep ERROR app.log 2>/dev/null | awk '{print $4}' | sort | uniq -c | sort -rn || echo '日志文件不存在（示例演示）'"})
	fmt.Println(result)

	// 查找特定时间段的日志
	fmt.Println("\n3. 查找特定时间段的日志:")
	result = terminal.Run(map[string]any{"command": "grep '2024-01-19 15:' app.log 2>/dev/null | tail -n 20 || echo '日志文件不存在（示例演示）'"})
	fmt.Println(result)
}

func demoCodebaseAnalysis() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("场景4: 代码库分析")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	codebaseDir := filepath.Join(ScriptDir, "codebase")
	terminal := NewTerminalTool(codebaseDir)

	// 统计代码行数
	fmt.Println("1. 统计代码行数:")
	result := terminal.Run(map[string]any{"command": "find . -name '*.py' -exec wc -l {} + 2>/dev/null | tail -n 1 || echo '代码库目录不存在（示例演示）'"})
	fmt.Println(result)

	// 查找所有 TODO 注释
	fmt.Println("\n2. 查找所有 TODO 注释:")
	result = terminal.Run(map[string]any{"command": "grep -rn 'TODO' --include='*.py' 2>/dev/null || echo '未找到TODO（示例演示）'"})
	fmt.Println(result)

	// 查找特定函数的定义
	fmt.Println("\n3. 查找特定函数的定义:")
	result = terminal.Run(map[string]any{"command": "grep -rn 'def process_data' --include='*.py' 2>/dev/null || echo '未找到函数定义（示例演示）'"})
	fmt.Println(result)
}

func demoSecurityFeatures() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("安全特性演示")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	projectDir := filepath.Join(ScriptDir, "project")
	terminal := NewTerminalTool(projectDir)

	// 尝试执行不允许的命令
	fmt.Println("1. 尝试执行危险命令 (rm):")
	result := terminal.Run(map[string]any{"command": "rm -rf /"})
	fmt.Println(result)

	// 尝试访问工作目录外的文件
	fmt.Println("\n2. 尝试访问工作目录外的文件:")
	result = terminal.Run(map[string]any{"command": "cat /etc/passwd"})
	fmt.Println(result)

	// 尝试逃逸工作目录
	fmt.Println("\n3. 尝试通过 .. 逃逸工作目录:")
	result = terminal.Run(map[string]any{"command": "cd ../../../etc"})
	fmt.Println(result)
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("TerminalTool 使用示例")
	fmt.Println(strings.Repeat("=", 80))

	// 演示各种使用场景
	demoExploratoryNavigation()
	demoDataFileAnalysis()
	demoLogAnalysis()
	demoCodebaseAnalysis()
	demoSecurityFeatures()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("演示完成!")
	fmt.Println(strings.Repeat("=", 80))
}
