// 13_ANPLoadBalancing - ANP 负载均衡
//
// 对应 Python: 13_ANPLoadBalancing.py
// 演示使用 ANP 进行负载均衡

package main

import (
	"fmt"
	"math/rand"
	"sort"
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

// getBestServer 选择负载最低的服务器
func getBestServer(discovery *ANPDiscovery) *ServiceInfo {
	servers := discovery.DiscoverServices("api")
	if len(servers) == 0 {
		return nil
	}

	best := servers[0]
	bestLoad := 1.0
	if load, ok := best.Metadata["load"].(float64); ok {
		bestLoad = load
	}

	for _, s := range servers[1:] {
		load := 1.0
		if l, ok := s.Metadata["load"].(float64); ok {
			load = l
		}
		if load < bestLoad {
			bestLoad = load
			best = s
		}
	}
	return best
}

func main() {
	// 创建服务发现中心
	discovery := NewANPDiscovery()

	// 注册多个相同类型的服务
	for i := range 5 {
		discovery.RegisterService(
			fmt.Sprintf("api_server_%d", i),
			fmt.Sprintf("API服务器%d", i),
			"api",
			[]string{"rest_api"},
			fmt.Sprintf("http://api%d:8000", i),
			map[string]any{"load": rand.Float64()*0.8 + 0.1},
		)
	}

	// 模拟请求分配
	for i := range 10 {
		server := getBestServer(discovery)
		if server != nil {
			load := server.Metadata["load"].(float64)
			fmt.Printf("请求 %d -> %s (负载: %.2f)\n", i+1, server.ServiceName, load)

			// 更新负载（模拟）
			server.Metadata["load"] = load + 0.1
		}
	}
}
