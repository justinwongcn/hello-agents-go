package main

import (
	"fmt"
	"math"
)

// cosineSimilarity 计算两个向量的余弦相似度
func cosineSimilarity(vec1, vec2 []float64) float64 {
	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0
	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}
	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func main() {
	// 假设我们已经学习到了简化的二维词向量
	embeddings := map[string][]float64{
		"king":  {0.9, 0.8},
		"queen": {0.9, 0.2},
		"man":   {0.7, 0.9},
		"woman": {0.7, 0.3},
	}

	// king - man + woman
	resultVec := make([]float64, 2)
	for i := range resultVec {
		resultVec[i] = embeddings["king"][i] - embeddings["man"][i] + embeddings["woman"][i]
	}

	// 计算结果向量与 "queen" 的相似度
	sim := cosineSimilarity(resultVec, embeddings["queen"])

	fmt.Printf("king - man + woman 的结果向量: %v\n", resultVec)
	fmt.Printf("该结果与 'queen' 的相似度: %.4f\n", sim)
}
