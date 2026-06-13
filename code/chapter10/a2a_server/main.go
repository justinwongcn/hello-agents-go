// 09_A2A_Server - A2A Agent 服务端
//
// 对应 Python: 09_A2A_Server.py
// 创建 A2A 研究员 Agent 服务，提供 research 技能

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// A2AServer A2A 服务器（模拟）
type A2AServer struct {
	Name        string
	Description string
	Version     string
	Skills      map[string]func(string) string
}

// Skill 添加技能
func (s *A2AServer) Skill(name string) func(func(string) string) {
	return func(fn func(string) string) {
		s.Skills[name] = fn
	}
}

// Run 运行服务器
func (s *A2AServer) Run(host string, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"name":"%s","description":"%s","version":"%s"}`,
			s.Name, s.Description, s.Version)
	})
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("A2A 服务器 %s 启动在 http://%s\n", s.Name, addr)
	http.ListenAndServe(addr, mux)
}

func main() {
	// 创建研究员Agent服务
	researcher := &A2AServer{
		Name:        "researcher",
		Description: "负责搜索和分析资料的Agent",
		Version:     "1.0.0",
		Skills:      make(map[string]func(string) string),
	}

	// 定义技能
	researcher.Skills["research"] = func(text string) string {
		// 处理研究请求
		topic := text
		if idx := strings.Index(strings.ToLower(text), "research "); idx >= 0 {
			topic = strings.TrimSpace(text[idx+len("research "):])
		}
		// 实际的研究逻辑（这里简化）
		return fmt.Sprintf(`{"topic":"%s","findings":"关于%s的研究结果...","sources":["来源1","来源2","来源3"]}`,
			topic, topic)
	}

	// 在后台启动服务
	go func() {
		researcher.Run("localhost", 5000)
	}()

	fmt.Println("✅ 研究员Agent服务已启动在 http://localhost:5000")

	// 保持程序运行
	time.Sleep(2 * time.Second)
	fmt.Println("服务已停止")
}
