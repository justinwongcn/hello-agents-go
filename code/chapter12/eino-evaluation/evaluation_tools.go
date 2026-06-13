package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// BFCLEvaluationTool BFCL 评估工具
// 调用 Python 脚本进行 BFCL（Berkeley Function Calling Leaderboard）评估
type BFCLEvaluationTool struct {
	// Python 脚本路径，默认使用同目录下的 Python 脚本
	ScriptDir string
	// 评估类别，如 simple_python, multiple, parallel 等
	Category string
	// 最大样本数，0 表示全部
	MaxSamples int
	// 模型名称
	ModelName string
}

func NewBFCLEvaluationTool() *BFCLEvaluationTool {
	return &BFCLEvaluationTool{
		ScriptDir:  defaultScriptDir(),
		Category:   "simple_python",
		MaxSamples: 5,
		ModelName:  "HelloAgents",
	}
}

func (t *BFCLEvaluationTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "bfcl_evaluation",
		Desc: "运行 BFCL (Berkeley Function Calling Leaderboard) 评估，评估智能体的函数调用能力。支持多种评估类别：simple_python, multiple, parallel 等",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"category": {
				Type:     "string",
				Desc:     "评估类别，可选值: simple_python, multiple, parallel, parallel_multiple, java, javascript, rest, sql, relevance, irrelevance",
				Required: false,
			},
			"max_samples": {
				Type:     "integer",
				Desc:     "最大评估样本数，0 表示评估全部样本",
				Required: false,
			},
			"model_name": {
				Type:     "string",
				Desc:     "被评估的模型或智能体名称",
				Required: false,
			},
		}),
	}, nil
}

func (t *BFCLEvaluationTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	category := extractStringArg(argumentsInJSON, "category")
	if category == "" {
		category = t.Category
	}
	maxSamples := extractIntArg(argumentsInJSON, "max_samples")
	if maxSamples == 0 {
		maxSamples = t.MaxSamples
	}
	modelName := extractStringArg(argumentsInJSON, "model_name")
	if modelName == "" {
		modelName = t.ModelName
	}

	// 构建 Python 命令
	scriptPath := filepath.Join(t.ScriptDir, "04_run_bfcl_evaluation.py")
	args := []string{
		scriptPath,
		"--category", category,
		"--samples", fmt.Sprintf("%d", maxSamples),
		"--model-name", modelName,
	}

	// 执行 Python 脚本
	output, err := executePythonScript(ctx, "python3", args)
	if err != nil {
		return fmt.Sprintf(`{"error": "BFCL 评估失败: %s", "stderr": "%s"}`, err.Error(), escapeJSON(output)), nil
	}

	// 解析评估结果
	result := parseEvaluationOutput(output)
	result["evaluation_type"] = "bfcl"
	result["category"] = category
	result["model_name"] = modelName
	result["max_samples"] = maxSamples

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON 序列化失败: %s"}`, err.Error()), nil
	}
	return string(resultJSON), nil
}

// GAIAEvaluationTool GAIA 评估工具
// 调用 Python 脚本进行 GAIA（General AI Assistants）评估
type GAIAEvaluationTool struct {
	ScriptDir  string
	Level      int
	MaxSamples int
}

func NewGAIAEvaluationTool() *GAIAEvaluationTool {
	return &GAIAEvaluationTool{
		ScriptDir:  defaultScriptDir(),
		Level:      1,
		MaxSamples: 2,
	}
}

func (t *GAIAEvaluationTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "gaia_evaluation",
		Desc: "运行 GAIA (General AI Assistants) 评估，评估智能体解决真实世界问题的能力。GAIA 是受限数据集，需要 HuggingFace 访问权限",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"level": {
				Type:     "integer",
				Desc:     "评估级别: 1=简单, 2=中等, 3=困难",
				Required: false,
			},
			"max_samples": {
				Type:     "integer",
				Desc:     "最大评估样本数，0 表示全部",
				Required: false,
			},
		}),
	}, nil
}

func (t *GAIAEvaluationTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	level := extractIntArg(argumentsInJSON, "level")
	if level == 0 {
		level = t.Level
	}
	maxSamples := extractIntArg(argumentsInJSON, "max_samples")
	if maxSamples == 0 {
		maxSamples = t.MaxSamples
	}

	scriptPath := filepath.Join(t.ScriptDir, "05_gaia_quick_start.py")
	args := []string{
		scriptPath,
		"--level", fmt.Sprintf("%d", level),
		"--max-samples", fmt.Sprintf("%d", maxSamples),
	}

	output, err := executePythonScript(ctx, "python3", args)
	if err != nil {
		return fmt.Sprintf(`{"error": "GAIA 评估失败: %s", "stderr": "%s"}`, err.Error(), escapeJSON(output)), nil
	}

	result := parseEvaluationOutput(output)
	result["evaluation_type"] = "gaia"
	result["level"] = level

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON 序列化失败: %s"}`, err.Error()), nil
	}
	return string(resultJSON), nil
}

// LLMJudgeTool LLM Judge 评估工具
// 使用 LLM 对生成内容进行质量评估
type LLMJudgeTool struct {
	ScriptDir string
}

func NewLLMJudgeTool() *LLMJudgeTool {
	return &LLMJudgeTool{
		ScriptDir: defaultScriptDir(),
	}
}

func (t *LLMJudgeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "llm_judge",
		Desc: "使用 LLM Judge 评估生成的题目质量，从正确性、清晰度、难度匹配、完整性四个维度评分（1-5分）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"input_file": {
				Type:     "string",
				Desc:     "待评估的题目 JSON 文件路径",
				Required: true,
			},
			"output_file": {
				Type:     "string",
				Desc:     "评估结果输出文件路径",
				Required: false,
			},
		}),
	}, nil
}

func (t *LLMJudgeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	inputFile := extractStringArg(argumentsInJSON, "input_file")
	if inputFile == "" {
		return `{"error": "缺少必要参数 input_file"}`, nil
	}
	outputFile := extractStringArg(argumentsInJSON, "output_file")

	scriptPath := filepath.Join(t.ScriptDir, "08_data_generation_llm_judge.py")
	args := []string{scriptPath, "--input", inputFile}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}

	output, err := executePythonScript(ctx, "python3", args)
	if err != nil {
		return fmt.Sprintf(`{"error": "LLM Judge 评估失败: %s"}`, err.Error()), nil
	}

	result := parseEvaluationOutput(output)
	result["evaluation_type"] = "llm_judge"

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON 序列化失败: %s"}`, err.Error()), nil
	}
	return string(resultJSON), nil
}

// WinRateTool Win Rate 评估工具
// 通过对比生成题目和真题评估生成质量
type WinRateTool struct {
	ScriptDir      string
	NumComparisons int
}

func NewWinRateTool() *WinRateTool {
	return &WinRateTool{
		ScriptDir:      defaultScriptDir(),
		NumComparisons: 20,
	}
}

func (t *WinRateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "win_rate",
		Desc: "使用 Win Rate 评估生成题目的质量，通过与参考真题对比，计算胜率、平局率和败率。Win Rate ≈ 50% 表示生成质量与真题相当",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"input_file": {
				Type:     "string",
				Desc:     "待评估的题目 JSON 文件路径",
				Required: true,
			},
			"num_comparisons": {
				Type:     "integer",
				Desc:     "对比次数，默认 20",
				Required: false,
			},
		}),
	}, nil
}

func (t *WinRateTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	inputFile := extractStringArg(argumentsInJSON, "input_file")
	if inputFile == "" {
		return `{"error": "缺少必要参数 input_file"}`, nil
	}
	numComparisons := extractIntArg(argumentsInJSON, "num_comparisons")
	if numComparisons == 0 {
		numComparisons = t.NumComparisons
	}

	scriptPath := filepath.Join(t.ScriptDir, "09_data_generation_win_rate.py")
	args := []string{
		scriptPath,
		"--input", inputFile,
		"--comparisons", fmt.Sprintf("%d", numComparisons),
	}

	output, err := executePythonScript(ctx, "python3", args)
	if err != nil {
		return fmt.Sprintf(`{"error": "Win Rate 评估失败: %s"}`, err.Error()), nil
	}

	result := parseEvaluationOutput(output)
	result["evaluation_type"] = "win_rate"
	result["num_comparisons"] = numComparisons

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON 序列化失败: %s"}`, err.Error()), nil
	}
	return string(resultJSON), nil
}

// ==================== 辅助函数 ====================

// defaultScriptDir 获取默认脚本目录（使用绝对路径）
func defaultScriptDir() string {
	exePath, err := os.Executable()
	if err != nil {
		return "../"
	}
	return filepath.Join(filepath.Dir(exePath), "..")
}

// executePythonScript 执行 Python 脚本并返回输出
func executePythonScript(ctx context.Context, pythonCmd string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, pythonCmd, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// parseEvaluationOutput 尝试从输出中提取 JSON 结果
// 使用括号匹配算法，支持嵌套 JSON 对象
func parseEvaluationOutput(output string) map[string]interface{} {
	result := make(map[string]interface{})

	// 查找第一个 '{' 的位置，然后用括号匹配找到完整的 JSON 对象
	start := strings.Index(output, "{")
	if start == -1 {
		result["raw_output"] = output
		return result
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(output); i++ {
		ch := output[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inString {
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				jsonStr := output[start : i+1]
				if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
					return result
				}
				break
			}
		}
	}

	// 如果没有找到合法的 JSON，返回原始输出
	result["raw_output"] = output
	return result
}

// extractStringArg 从 JSON 参数中提取字符串值
func extractStringArg(argumentsInJSON, key string) string {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
		return ""
	}
	val, ok := args[key]
	if !ok {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

// extractIntArg 从 JSON 参数中提取整数值
func extractIntArg(argumentsInJSON, key string) int {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
		return 0
	}
	val, ok := args[key]
	if !ok {
		return 0
	}
	// JSON 数字默认解析为 float64
	f, ok := val.(float64)
	if !ok {
		return 0
	}
	return int(f)
}

// escapeJSON 转义字符串中的特殊字符用于 JSON
func escapeJSON(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return s
	}
	// json.Marshal 返回带引号的 JSON 字符串，去掉首尾引号
	return string(b[1 : len(b)-1])
}
