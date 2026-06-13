package main

import (
	"fmt"
	"strings"
)

func main() {
	// 示例语料库，与上方案例讲解中的语料库保持一致
	corpus := "datawhale agent learns datawhale agent works"
	tokens := strings.Fields(corpus)
	totalTokens := len(tokens)

	// --- 第一步：计算 P(datawhale) ---
	countDatawhale := 0
	for _, t := range tokens {
		if t == "datawhale" {
			countDatawhale++
		}
	}
	pDatawhale := float64(countDatawhale) / float64(totalTokens)
	fmt.Printf("第一步: P(datawhale) = %d/%d = %.3f\n", countDatawhale, totalTokens, pDatawhale)

	// --- 第二步：计算 P(agent|datawhale) ---
	// 先计算 bigrams 用于后续步骤
	bigramCounts := make(map[[2]string]int)
	for i := range len(tokens) - 1 {
		bigram := [2]string{tokens[i], tokens[i+1]}
		bigramCounts[bigram]++
	}
	countDatawhaleAgent := bigramCounts[[2]string{"datawhale", "agent"}]
	// countDatawhale 已在第一步计算
	pAgentGivenDatawhale := float64(countDatawhaleAgent) / float64(countDatawhale)
	fmt.Printf("第二步: P(agent|datawhale) = %d/%d = %.3f\n", countDatawhaleAgent, countDatawhale, pAgentGivenDatawhale)

	// --- 第三步：计算 P(learns|agent) ---
	countAgentLearns := bigramCounts[[2]string{"agent", "learns"}]
	countAgent := 0
	for _, t := range tokens {
		if t == "agent" {
			countAgent++
		}
	}
	pLearnsGivenAgent := float64(countAgentLearns) / float64(countAgent)
	fmt.Printf("第三步: P(learns|agent) = %d/%d = %.3f\n", countAgentLearns, countAgent, pLearnsGivenAgent)

	// --- 最后：将概率连乘 ---
	pSentence := pDatawhale * pAgentGivenDatawhale * pLearnsGivenAgent
	fmt.Printf("最后: P('datawhale agent learns') ≈ %.3f * %.3f * %.3f = %.3f\n", pDatawhale, pAgentGivenDatawhale, pLearnsGivenAgent, pSentence)
}
