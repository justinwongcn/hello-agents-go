# 第十二章：智能体性能评估 - Go/Eino 实现

本目录包含使用 CloudWeGo Eino 框架实现的智能体评估工具。

## 📁 文件结构

```
eino-evaluation/
├── main.go                 # 主程序入口
├── evaluation_tools.go     # 评估工具实现
├── go.mod                  # Go 模块定义
└── README.md               # 本文档
```

## 🛠️ 评估工具

### 1. BFCLEvaluationTool
- **名称**: `bfcl_evaluation`
- **功能**: 评估智能体的函数调用能力
- **参数**:
  - `category`: 评估类别（simple_python, multiple, parallel 等）
  - `max_samples`: 最大样本数
  - `model_name`: 模型名称

### 2. GAIAEvaluationTool
- **名称**: `gaia_evaluation`
- **功能**: 评估智能体解决真实世界问题的能力
- **参数**:
  - `level`: 评估级别（1=简单, 2=中等, 3=困难）
  - `max_samples`: 最大样本数

### 3. LLMJudgeTool
- **名称**: `llm_judge`
- **功能**: 使用 LLM 评估生成内容质量
- **参数**:
  - `input_file`: 待评估的题目 JSON 文件
  - `output_file`: 输出文件路径（可选）
- **评估维度**: 正确性、清晰度、难度匹配、完整性（1-5 分）

### 4. WinRateTool
- **名称**: `win_rate`
- **功能**: 通过对比评估生成质量
- **参数**:
  - `input_file`: 待评估的题目 JSON 文件
  - `num_comparisons`: 对比次数（默认 20）

## 🚀 快速开始

### 环境准备

```bash
# 设置环境变量
export LLM_API_KEY="your-api-key"
export LLM_BASE_URL="https://api.openai.com/v1"
export LLM_MODEL_ID="gpt-4o-mini"
```

### 运行

```bash
cd /Users/john/GolandProjects/hello-agents-go/code/chapter12/eino-evaluation
go mod tidy
go run .
```

### 使用示例

启动后，可以通过自然语言与评估助手交互：

```
🤔 您的问题: 帮我运行 BFCL 评估，使用 simple_python 类别，评估 5 个样本
🤔 您的问题: 用 LLM Judge 评估 ./data/problems.json 的质量
🤔 您的问题: 对 ./data/generated.json 进行 Win Rate 评估
```

## 📝 与 Python 实现的关系

本 Go 实现通过调用 Chapter 12 目录下的 Python 脚本来执行实际评估：

- `04_run_bfcl_evaluation.py` → BFCL 评估
- `05_gaia_quick_start.py` → GAIA 评估
- `08_data_generation_llm_judge.py` → LLM Judge
- `09_data_generation_win_rate.py` → Win Rate

Go 端提供 Eino 框架集成和统一的 Agent 交互界面。
