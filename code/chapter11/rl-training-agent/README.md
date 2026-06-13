# RL Training Agent - 基于 Eino 框架

使用 CloudWeGo Eino 框架构建的强化学习训练助手，通过 Go 调用 Python 脚本实现 RL 训练。

## 功能

- **加载数据集** (`load_dataset`): 加载 GSM8K 等数据集，支持 SFT/RL 格式
- **训练模型** (`train`): 支持 SFT 监督微调和 GRPO 强化学习训练
- **评估模型** (`evaluate`): 评估训练后的模型性能
- **创建奖励函数** (`create_reward`): 创建准确性/长度惩罚/步骤奖励

## 环境变量

```bash
# LLM 配置（可选，未配置则使用模拟模式）
export LLM_API_KEY="your-api-key"
export LLM_BASE_URL="https://api.openai.com/v1"
export LLM_MODEL_ID="gpt-4o-mini"

# Python 配置（可选）
export SCRIPT_DIR="../"        # Python脚本目录
export PYTHON_PATH="python3"   # Python可执行文件路径
```

## 运行

```bash
# 下载依赖
go mod tidy

# 运行
go run .
```

## 参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| action | string | 操作类型: load_dataset, train, evaluate, create_reward |
| model_name | string | 模型名称或路径 |
| algorithm | string | 训练算法: sft 或 grpo |
| max_samples | int | 最大样本数 |
| num_epochs | int | 训练轮数 |
| batch_size | int | 批量大小 |
| use_lora | bool | 是否使用 LoRA |
| lora_r | int | LoRA 秩，默认 16 |
| lora_alpha | int | LoRA 缩放因子，默认 32 |
| output_dir | string | 输出目录 |
| learning_rate | float | 学习率 |
| reward_type | string | 奖励函数类型: accuracy, length_penalty, step |
| model_path | string | 评估时使用的模型路径 |
| format | string | 数据集格式: sft 或 rl |
| split | string | 数据集划分: train 或 test |

## 文件结构

```
rl-training-agent/
├── go.mod                  # Go 模块定义
├── main.go                 # 入口文件，Eino Agent 配置
├── rl_training_tool.go     # RLTrainingTool 实现
└── README.md               # 说明文档
```

## 依赖的 Python 脚本

此 Go 工具通过 `os/exec` 调用上层目录的 Python 脚本：

- `01_dataset_loading.py` - 数据集加载
- `04_sft_training.py` - SFT 训练
- `05_grpo_training.py` - GRPO 训练
- `07_model_evaluation.py` - 模型评估

需要确保 Python 环境已安装 `hello_agents` 包。
