// weather_mcp_server - 天气查询 MCP 服务器
//
// 对应 Python: 14_weather_mcp_server.py
// 天气查询 MCP 服务器,提供天气查询工具

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// CITY_MAP 城市映射
var CITY_MAP = map[string]string{
	"北京": "Beijing", "上海": "Shanghai", "广州": "Guangzhou",
	"深圳": "Shenzhen", "杭州": "Hangzhou", "成都": "Chengdu",
	"重庆": "Chongqing", "武汉": "Wuhan", "西安": "Xi'an",
	"南京": "Nanjing", "天津": "Tianjin", "苏州": "Suzhou",
}

// WeatherData 天气数据
type WeatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	Humidity    int     `json:"humidity"`
	Condition   string  `json:"condition"`
	WindSpeed   float64 `json:"wind_speed"`
	Visibility  float64 `json:"visibility"`
	Timestamp   string  `json:"timestamp"`
}

// MCPServer MCP 服务器（模拟）
type MCPServer struct {
	Name        string
	Description string
	Tools       map[string]func(map[string]any) string
}

// AddTool 添加工具
func (s *MCPServer) AddTool(name string, fn func(map[string]any) string) {
	s.Tools[name] = fn
}

// Run 运行服务器
func (s *MCPServer) Run() {
	fmt.Printf("🌤️  Starting %s...\n", s.Name)
	fmt.Println("📡 Transport: stdio")
	fmt.Println("✨ Ready to serve weather data!")
}

// getWeatherData 从 wttr.in 获取天气数据
func getWeatherData(city string) (*WeatherData, error) {
	cityEn, ok := CITY_MAP[city]
	if !ok {
		cityEn = city
	}
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", cityEn)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	currentCondition, ok := data["current_condition"].([]any)
	if !ok || len(currentCondition) == 0 {
		return nil, fmt.Errorf("无法获取当前天气数据")
	}

	current, ok := currentCondition[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("天气数据格式错误")
	}

	tempStr, _ := current["temp_C"].(string)
	var temp float64
	fmt.Sscanf(tempStr, "%f", &temp)

	feelsLikeStr, _ := current["FeelsLikeC"].(string)
	var feelsLike float64
	fmt.Sscanf(feelsLikeStr, "%f", &feelsLike)

	humidityStr, _ := current["humidity"].(string)
	var humidity int
	fmt.Sscanf(humidityStr, "%d", &humidity)

	weatherDesc, _ := current["weatherDesc"].([]any)
	condition := ""
	if len(weatherDesc) > 0 {
		if desc, ok := weatherDesc[0].(map[string]any); ok {
			condition, _ = desc["value"].(string)
		}
	}

	windSpeedStr, _ := current["windspeedKmph"].(string)
	var windSpeed float64
	fmt.Sscanf(windSpeedStr, "%f", &windSpeed)

	visibilityStr, _ := current["visibility"].(string)
	var visibility float64
	fmt.Sscanf(visibilityStr, "%f", &visibility)

	return &WeatherData{
		City:        city,
		Temperature: temp,
		FeelsLike:   feelsLike,
		Humidity:    humidity,
		Condition:   condition,
		WindSpeed:   windSpeed / 3.6,
		Visibility:  visibility,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// getWeather 获取指定城市的当前天气
func getWeather(args map[string]any) string {
	city, _ := args["city"].(string)
	weatherData, err := getWeatherData(city)
	if err != nil {
		errJSON, _ := json.Marshal(map[string]any{"error": err.Error(), "city": city})
		return string(errJSON)
	}
	result, _ := json.MarshalIndent(weatherData, "", "  ")
	return string(result)
}

// listSupportedCities 列出所有支持的中文城市
func listSupportedCities(args map[string]any) string {
	cities := make([]string, 0, len(CITY_MAP))
	for city := range CITY_MAP {
		cities = append(cities, city)
	}
	sort.Strings(cities)
	result := map[string]any{
		"cities": cities,
		"count":  len(cities),
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data)
}

// getServerInfo 获取服务器信息
func getServerInfo(args map[string]any) string {
	info := map[string]any{
		"name":    "Weather MCP Server",
		"version": "1.0.0",
		"tools":   []string{"get_weather", "list_supported_cities", "get_server_info"},
	}
	data, _ := json.MarshalIndent(info, "", "  ")
	return string(data)
}

func main() {
	// 创建 MCP 服务器
	weatherServer := &MCPServer{
		Name:        "weather-server",
		Description: "真实天气查询服务",
		Tools:       make(map[string]func(map[string]any) string),
	}

	// 注册工具到服务器
	weatherServer.AddTool("get_weather", getWeather)
	weatherServer.AddTool("list_supported_cities", listSupportedCities)
	weatherServer.AddTool("get_server_info", getServerInfo)

	// 运行服务器
	weatherServer.Run()
}
