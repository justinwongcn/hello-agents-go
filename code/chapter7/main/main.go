package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"hello-agents-go/chapter4/llmclient"
	"hello-agents-go/chapter7/common"
)

func main() {
	// 加载环境变量
	_ = godotenv.Load()

	// 实例化我们重写的客户端，并指定provider
	llm, err := common.NewMyLLM("", "", "", "modelscope", 0, 0, 0)
	if err != nil {
		fmt.Printf("❌ 创建LLM客户端失败: %v\n", err)
		return
	}

	// 准备消息
	messages := []llmclient.Message{{Role: "user", Content: "你好，请介绍一下你自己。"}}

	// 发起流式调用，think等方法都已从父类继承，无需重写
	fmt.Println("ModelScope Response:")
	for chunk := range llm.ThinkStream(messages, llm.Temperature) {
		// chunk在流式输出中已经打印过一遍，这里只需要忽略即可
		_ = chunk
	}
}
