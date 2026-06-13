// 12_ANPTaskDistribution - ANP 任务分配
//
// 对应 Python: 12_ANPTaskDistribution.py
// 演示使用 ANP 进行智能任务分配

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceID    string
	ServiceName  string
	ServiceType  string
	Capabilities []string
	Endpoint     string
	Metadata     map[string]any
}

// ANPDiscovery ANP 服务发现中心
type ANPDiscovery struct {
	services map[string]*ServiceInfo
}

// NewANPDiscovery 创建 ANP 服务发现中心
func NewANPDiscovery() *ANPDiscovery {
	return &ANPDiscovery{services: make(map[string]*ServiceInfo)}
}

// RegisterService 注册服务
func (d *ANPDiscovery) RegisterService(serviceID, serviceName, serviceType string, capabilities []string, endpoint string, metadata map[string]any) {
	d.services[serviceID] = &ServiceInfo{
		ServiceID:    serviceID,
		ServiceName:  serviceName,
		ServiceType:  serviceType,
		Capabilities: capabilities,
		Endpoint:     endpoint,
		Metadata:     metadata,
	}
}

// DiscoverServices 发现服务
func (d *ANPDiscovery) DiscoverServices(serviceType string) []*ServiceInfo {
	var result []*ServiceInfo
	for _, s := range d.services {
		if serviceType == "" || s.ServiceType == serviceType {
			result = append(result, s)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ServiceID < result[j].ServiceID
	})
	return result
}

// ListAllServices 列出所有服务
func (d *ANPDiscovery) ListAllServices() []*ServiceInfo {
	return d.DiscoverServices("")
}

// SimpleAgent 简单智能体（模拟）
type SimpleAgent struct {
	Name         string
	SystemPrompt string
	Tool         any
}

// Run 运行智能体
func (a *SimpleAgent) Run(input string) string {
	return fmt.Sprintf("智能体 '%s' 分析任务并选择最优节点（模拟）", a.Name)
}

// ANPTool ANP 工具（模拟）
type ANPTool struct {
	Name        string
	Description string
	Discovery   *ANPDiscovery
}

// Run 执行工具
func (t *ANPTool) Run(params map[string]any) any {
	action, _ := params["action"].(string)
	switch action {
	case "discover_services":
		serviceType, _ := params["service_type"].(string)
		return t.Discovery.DiscoverServices(serviceType)
	case "get_stats":
		return map[string]any{"total_services": len(t.Discovery.ListAllServices())}
	default:
		return nil
	}
}

func main() {
	// 1. 创建服务发现中心
	discovery := NewANPDiscovery()

	// 2. 注册多个计算节点
	for i := range 10 {
		discovery.RegisterService(
			fmt.Sprintf("compute_node_%d", i),
			fmt.Sprintf("计算节点%d", i),
			"compute",
			[]string{"data_processing", "ml_training"},
			fmt.Sprintf("http://node%d:8000", i),
			map[string]any{
				"load":      rand.Float64()*0.8 + 0.1,
				"cpu_cores": []int{4, 8, 16}[rand.Intn(3)],
				"memory_gb": []int{16, 32, 64}[rand.Intn(3)],
				"gpu":       rand.Intn(2) == 1,
			},
		)
	}

	fmt.Printf("✅ 注册了 %d 个计算节点\n", len(discovery.ListAllServices()))

	// 3. 创建任务调度Agent
	scheduler := &SimpleAgent{
		Name: "任务调度器",
		SystemPrompt: `你是一个智能任务调度器,负责：
1. 分析任务需求
2. 选择最合适的计算节点
3. 分配任务

选择节点时考虑：负载、CPU核心数、内存、GPU等因素。

使用 service_discovery 工具时,必须提供 action 参数：
- 查看所有节点：{"action": "discover_services", "service_type": "compute"}
- 获取网络统计：{"action": "get_stats"}`,
		Tool: &ANPTool{
			Name:        "service_discovery",
			Description: "服务发现工具,可以查找和选择计算节点",
			Discovery:   discovery,
		},
	}

	_ = scheduler

	// 4. 智能任务分配
	assignTask := func(taskDescription string) {
		fmt.Printf("\n任务：%s\n", taskDescription)
		fmt.Println(strings.Repeat("=", 50))

		// 让Agent智能选择节点
		response := scheduler.Run(fmt.Sprintf(`
请为以下任务选择最合适的计算节点：
%s

步骤：
1. 使用 service_discovery 工具查看所有可用的计算节点（service_type="compute"）
2. 分析每个节点的特点（负载、CPU核心数、内存、GPU等）
3. 根据任务需求选择最合适的节点
4. 说明选择理由

请直接给出最终选择的节点ID和理由。
`, taskDescription))

		fmt.Println(response)
		fmt.Println(strings.Repeat("=", 50))
	}

	// 测试不同类型的任务
	assignTask("训练一个大型深度学习模型,需要GPU支持")
	assignTask("处理大量文本数据,需要高内存")
	assignTask("运行轻量级数据分析任务")
}
