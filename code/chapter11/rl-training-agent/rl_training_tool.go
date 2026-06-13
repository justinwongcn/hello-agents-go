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

// ============================================================================
// RLTrainingTool - 强化学习训练工具
// ============================================================================

// RLTrainingToolConfig 配置
type RLTrainingToolConfig struct {
	// Python脚本目录（包含01_dataset_loading.py等脚本）
	ScriptDir string
	// Python可执行文件路径
	PythonPath string
}

// RLTrainingTool 强化学习训练工具
type RLTrainingTool struct {
	config RLTrainingToolConfig
}

// NewRLTrainingTool 创建工具实例
func NewRLTrainingTool(config RLTrainingToolConfig) *RLTrainingTool {
	if config.PythonPath == "" {
		config.PythonPath = "python3"
	}
	return &RLTrainingTool{config: config}
}

// Info 返回工具元信息
func (t *RLTrainingTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "rl_training",
		Desc: "强化学习训练工具：支持加载数据集(load_dataset)、训练模型(train)、评估模型(evaluate)、创建奖励函数(create_reward)。" +
			"支持SFT和GRPO算法，支持LoRA配置。",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Type:     "string",
				Desc:     "操作类型: load_dataset(加载数据集), train(训练模型), evaluate(评估模型), create_reward(创建奖励函数)",
				Required: true,
				Enum:     []string{"load_dataset", "train", "evaluate", "create_reward"},
			},
			"model_name": {
				Type: "string",
				Desc: "模型名称或路径，如 'Qwen/Qwen3-0.6B' 或本地路径 './output/sft_model'",
			},
			"algorithm": {
				Type: "string",
				Desc: "训练算法: sft(监督微调) 或 grpo(强化学习)",
				Enum: []string{"sft", "grpo"},
			},
			"max_samples": {
				Type: "integer",
				Desc: "最大样本数，0表示使用全部数据",
			},
			"num_epochs": {
				Type: "integer",
				Desc: "训练轮数，默认3",
			},
			"batch_size": {
				Type: "integer",
				Desc: "批量大小，默认4",
			},
			"use_lora": {
				Type: "boolean",
				Desc: "是否使用LoRA参数高效微调",
			},
			"lora_r": {
				Type: "integer",
				Desc: "LoRA秩，默认16",
			},
			"lora_alpha": {
				Type: "integer",
				Desc: "LoRA缩放因子，默认32",
			},
			"output_dir": {
				Type: "string",
				Desc: "输出目录路径",
			},
			"learning_rate": {
				Type: "number",
				Desc: "学习率",
			},
			"reward_type": {
				Type: "string",
				Desc: "奖励函数类型: accuracy, length_penalty, step",
				Enum: []string{"accuracy", "length_penalty", "step"},
			},
			"model_path": {
				Type: "string",
				Desc: "评估时使用的模型路径",
			},
			"format": {
				Type: "string",
				Desc: "数据集格式: sft 或 rl",
				Enum: []string{"sft", "rl"},
			},
			"split": {
				Type: "string",
				Desc: "数据集划分: train 或 test",
				Enum: []string{"train", "test"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (t *RLTrainingTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析JSON参数
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	action, _ := params["action"].(string)
	if action == "" {
		return "", fmt.Errorf("必须指定 action 参数")
	}

	switch action {
	case "load_dataset":
		return t.loadDataset(ctx, params)
	case "train":
		return t.train(ctx, params)
	case "evaluate":
		return t.evaluate(ctx, params)
	case "create_reward":
		return t.createReward(ctx, params)
	default:
		return "", fmt.Errorf("不支持的操作: %s", action)
	}
}

// loadDataset 加载数据集
func (t *RLTrainingTool) loadDataset(ctx context.Context, params map[string]interface{}) (string, error) {
	scriptPath := t.resolveScript("01_dataset_loading.py")

	configJSON, err := buildConfigJSON(params)
	if err != nil {
		return "", fmt.Errorf("构建配置失败: %w", err)
	}

	script := `
import sys, json, os
sys.path.insert(0, sys.argv[1])
from hello_agents.tools import RLTrainingTool
tool = RLTrainingTool()
config = json.loads(os.environ['RL_CONFIG_JSON'])
result = tool.run(config)
print(result)
`

	result, err := t.runPythonWithConfig(ctx, script, filepath.Dir(scriptPath), configJSON)
	if err != nil {
		return "", fmt.Errorf("加载数据集失败: %w", err)
	}

	return result, nil
}

// train 训练模型
func (t *RLTrainingTool) train(ctx context.Context, params map[string]interface{}) (string, error) {
	algorithm, _ := params["algorithm"].(string)
	if algorithm == "" {
		algorithm = "grpo"
	}

	scriptPath := t.resolveScript("05_grpo_training.py")
	if algorithm == "sft" {
		scriptPath = t.resolveScript("04_sft_training.py")
	}

	configJSON, err := buildConfigJSON(params)
	if err != nil {
		return "", fmt.Errorf("构建配置失败: %w", err)
	}

	script := `
import sys, json, os
sys.path.insert(0, sys.argv[1])
from hello_agents.tools import RLTrainingTool
tool = RLTrainingTool()
config = json.loads(os.environ['RL_CONFIG_JSON'])
result = tool.run(config)
print(result)
`

	result, err := t.runPythonWithConfig(ctx, script, filepath.Dir(scriptPath), configJSON)
	if err != nil {
		return "", fmt.Errorf("训练失败: %w", err)
	}

	return result, nil
}

// evaluate 评估模型
func (t *RLTrainingTool) evaluate(ctx context.Context, params map[string]interface{}) (string, error) {
	scriptPath := t.resolveScript("07_model_evaluation.py")
	configJSON, err := buildConfigJSON(params)
	if err != nil {
		return "", fmt.Errorf("构建配置失败: %w", err)
	}

	script := `
import sys, json, os
sys.path.insert(0, sys.argv[1])
from hello_agents.tools import RLTrainingTool
tool = RLTrainingTool()
config = json.loads(os.environ['RL_CONFIG_JSON'])
result = tool.run(config)
print(result)
`

	result, err := t.runPythonWithConfig(ctx, script, filepath.Dir(scriptPath), configJSON)
	if err != nil {
		return "", fmt.Errorf("评估失败: %w", err)
	}

	return result, nil
}

// createReward 创建奖励函数
func (t *RLTrainingTool) createReward(ctx context.Context, params map[string]interface{}) (string, error) {
	rewardType, _ := params["reward_type"].(string)
	if rewardType == "" {
		rewardType = "accuracy"
	}

	resultMap := map[string]interface{}{
		"reward_type": rewardType,
		"status":      "created",
	}

	switch rewardType {
	case "accuracy":
		resultMap["description"] = "准确性奖励: 答案正确奖励1.0, 答案错误奖励0.0"
	case "length_penalty":
		penaltyWeight, ok := params["penalty_weight"].(float64)
		if !ok {
			penaltyWeight = 0.001
		}
		resultMap["description"] = "长度惩罚奖励: 在准确性基础上惩罚过长回答"
		resultMap["penalty_weight"] = penaltyWeight
	case "step":
		stepBonus, ok := params["step_bonus"].(float64)
		if !ok {
			stepBonus = 0.1
		}
		resultMap["description"] = "步骤奖励: 在准确性基础上奖励详细推理步骤"
		resultMap["step_bonus"] = stepBonus
	default:
		return "", fmt.Errorf("不支持的奖励类型: %s", rewardType)
	}

	result, _ := json.MarshalIndent(resultMap, "", "  ")
	return string(result), nil
}

// resolveScript 解析脚本路径
func (t *RLTrainingTool) resolveScript(name string) string {
	if t.config.ScriptDir != "" {
		return filepath.Join(t.config.ScriptDir, name)
	}
	return name
}

// buildConfigJSON 构建配置JSON
func buildConfigJSON(params map[string]interface{}) (string, error) {
	// 转换参数名称以匹配Python RLTrainingTool接口
	config := make(map[string]interface{})
	for k, v := range params {
		config[k] = v
	}
	b, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %w", err)
	}
	return string(b), nil
}

// runPythonWithConfig 运行Python脚本，通过环境变量安全传递配置JSON
func (t *RLTrainingTool) runPythonWithConfig(ctx context.Context, script string, scriptDir string, configJSON string) (string, error) {
	cmd := exec.CommandContext(ctx, t.config.PythonPath, "-c", script, scriptDir)
	cmd.Dir = t.config.ScriptDir
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8", "RL_CONFIG_JSON="+configJSON)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Python执行失败: %w\n输出: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
