# Eino 智能助手演示

基于 [CloudWeGo Eino](https://github.com/cloudwego/eino) 框架的智能助手演示项目。

## 项目简介

本项目是《Hello Agents》第六章的 Go 语言实现，展示了如何使用 Eino 框架构建一个支持工具调用的智能助手。

### 功能特性

- 🌤️ **天气查询** - 查询指定城市的天气信息
- 🔢 **数学计算** - 执行简单的数学计算
- 🔍 **网络搜索** - 搜索互联网上的最新信息
- 🤖 **智能对话** - 基于 LLM 的自然语言理解（需要配置 API 密钥）
- 🎭 **模拟模式** - 无需 API 密钥即可体验基本功能

## 技术架构

```
┌─────────────────────────────────────────────┐
│                 用户界面                      │
│            (命令行交互)                       │
└─────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────┐
│              ADK Runner                      │
│         (Agent 执行引擎)                     │
└─────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────┐
│          ChatModelAgent                      │
│      (ReAct 循环 + 工具调用)                 │
└─────────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│ 天气工具 │ │ 计算工具 │ │ 搜索工具 │
└──────────┘ └──────────┘ └──────────┘
```

## 快速开始

### 前置条件

- Go 1.21+
- OpenAI 兼容的 API 密钥（可选，不配置则使用模拟模式）

### 安装依赖

```bash
cd code/chapter6/eino-demo
go mod tidy
```

### 运行演示

```bash
# 模拟模式（无需 API 密钥）
go run main.go

# 使用真实 LLM
export LLM_API_KEY=*** LLM_BASE_URL="https://api.openai.com/v1"
export LLM_MODEL_ID="gpt-4o-mini"
go run main.go
```

### 交互示例

```
🤖 Eino 智能助手 - 基于 CloudWeGo Eino 框架
支持天气查询、数学计算、网络搜索等功能
(输入 'quit' 退出)

⚠️ 未配置有效的 LLM_API_KEY，将使用模拟模式运行
   设置环境变量 LLM_API_KEY 以启用完整功能

🤔 您的问题: 北京天气怎么样？

============================================================
💡 回答: 收到您的问题：北京天气怎么样？

这是一个使用 Eino 框架的智能助手演示。
在实际应用中，这里会调用 LLM 并使用工具来回答您的问题。

可用工具：
- get_weather: 查询天气
- calculator: 数学计算
- web_search: 网络搜索

请配置 LLM_API_KEY 环境变量以启用完整功能。
============================================================
```

## 代码结构

```
eino-demo/
├── main.go          # 主程序，包含工具定义和 Agent 配置
├── go.mod           # Go 模块定义
├── go.sum           # 依赖校验
├── README.md        # 项目说明
└── eino-demo        # 编译后的可执行文件
```

### 核心组件

1. **WeatherTool** - 天气查询工具
   - 实现 `tool.InvokableTool` 接口
   - 支持中文城市名称查询

2. **CalculatorTool** - 计算器工具
   - 支持基本数学运算
   - 可扩展为支持复杂表达式

3. **SearchTool** - 搜索工具
   - 模拟网络搜索功能
   - 可集成真实的搜索 API

4. **MockAgent** - 模拟 Agent
   - 当没有配置 API 密钥时使用
   - 展示 Agent 的基本结构

## Eino 框架核心概念

### 1. Components（组件）

Eino 提供了可复用的组件抽象：

- **ChatModel** - 聊天模型接口
- **Tool** - 工具接口
- **Retriever** - 检索器接口
- **Embedding** - 嵌入模型接口

### 2. ADK（Agent Development Kit）

ADK 是 Eino 的 Agent 开发套件：

- **ChatModelAgent** - 基础 Agent 模式
- **DeepAgent** - 深度推理 Agent
- **Runner** - Agent 执行引擎

### 3. Orchestration（编排）

Eino 提供了强大的编排能力：

- **Chain** - 链式编排
- **Graph** - 图式编排
- **Workflow** - 工作流编排

## API 参考

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `LLM_API_KEY` | OpenAI 兼容的 API 密钥 | 空（使用模拟模式） |
| `LLM_BASE_URL` | API 基础 URL | `https://api.openai.com/v1` |
| `LLM_MODEL_ID` | 模型 ID | `gpt-4o-mini` |

### 工具列表

| 工具名 | 说明 | 参数 |
|--------|------|------|
| `get_weather` | 查询天气 | `city` - 城市名称 |
| `calculator` | 数学计算 | `expression` - 数学表达式 |
| `web_search` | 网络搜索 | `query` - 搜索关键词 |

## 与其他框架对比

| 特性 | Eino (Go) | LangChain (Python) | AutoGen (Python) |
|------|-----------|-------------------|------------------|
| 语言 | Go | Python | Python |
| 性能 | 高 | 中 | 中 |
| 类型安全 | 强 | 弱 | 弱 |
| 工具调用 | ✅ | ✅ | ✅ |
| 多 Agent | ✅ | ✅ | ✅ |
| 流式处理 | ✅ | ✅ | ✅ |
| 生产就绪 | ✅ | ✅ | ⚠️ |

## 扩展阅读

- [Eino 官方文档](https://www.cloudwego.cn/docs/eino/)
- [Eino GitHub 仓库](https://github.com/cloudwego/eino)
- [Eino 示例代码](https://github.com/cloudwego/eino-examples)
- [ADK 文档](https://www.cloudwego.cn/docs/eino/core_modules/adk/)

## 许可证

本项目基于 Apache 2.0 许可证开源。
