// 11_ANPInit - ANP 服务发现初始化
//
// 对应 Python: 11_ANPInit.py
// 演示 ANP 服务发现、注册和网络构建

package main

import (
	"fmt"
	"sort"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceID   string
	ServiceName string
	ServiceType string
	Capabilities []string
	Endpoint    string
	Metadata    map[string]any
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

// ANPNetwork ANP 网络
type ANPNetwork struct {
	NetworkID string
	Nodes     map[string]string
	Edges     map[string][]string
}

// NewANPNetwork 创建 ANP 网络
func NewANPNetwork(networkID string) *ANPNetwork {
	return &ANPNetwork{
		NetworkID: networkID,
		Nodes:     make(map[string]string),
		Edges:     make(map[string][]string),
	}
}

// AddNode 添加节点
func (n *ANPNetwork) AddNode(nodeID, endpoint string) {
	n.Nodes[nodeID] = endpoint
}

// ConnectNodes 连接节点
func (n *ANPNetwork) ConnectNodes(node1, node2 string) {
	n.Edges[node1] = append(n.Edges[node1], node2)
	n.Edges[node2] = append(n.Edges[node2], node1)
}

// GetNetworkStats 获取网络统计
func (n *ANPNetwork) GetNetworkStats() map[string]any {
	return map[string]any{
		"total_nodes": len(n.Nodes),
		"total_edges": len(n.Edges),
		"network_id":  n.NetworkID,
	}
}

func main() {
	// 创建服务发现中心
	discovery := NewANPDiscovery()

	// 注册Agent服务
	discovery.RegisterService(
		"nlp_agent_1", "NLP处理专家A", "nlp",
		[]string{"text_analysis", "sentiment_analysis", "ner"},
		"http://localhost:8001",
		map[string]any{"load": 0.3, "price": 0.01, "version": "1.0.0"},
	)

	discovery.RegisterService(
		"nlp_agent_2", "NLP处理专家B", "nlp",
		[]string{"text_analysis", "translation"},
		"http://localhost:8002",
		map[string]any{"load": 0.7, "price": 0.02, "version": "1.1.0"},
	)

	fmt.Println("✅ 服务注册完成")

	// 按类型查找
	nlpServices := discovery.DiscoverServices("nlp")
	fmt.Printf("找到 %d 个NLP服务\n", len(nlpServices))

	// 选择负载最低的服务
	var bestService *ServiceInfo
	bestLoad := 1.0
	for _, s := range nlpServices {
		load, ok := s.Metadata["load"].(float64)
		if !ok {
			load = 1.0
		}
		if load < bestLoad {
			bestLoad = load
			bestService = s
		}
	}
	if bestService != nil {
		fmt.Printf("最佳服务：%s (负载: %v)\n", bestService.ServiceName, bestService.Metadata["load"])
	}

	// 创建网络
	network := NewANPNetwork("ai_cluster")

	// 添加节点
	for _, service := range discovery.ListAllServices() {
		network.AddNode(service.ServiceID, service.Endpoint)
	}

	// 建立连接（根据能力匹配）
	network.ConnectNodes("nlp_agent_1", "nlp_agent_2")

	stats := network.GetNetworkStats()
	fmt.Printf("✅ 网络构建完成，共 %v 个节点\n", stats["total_nodes"])
}
