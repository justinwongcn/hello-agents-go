package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// getStats 统计词元对频率
func getStats(vocab map[string]int) map[[2]string]int {
	pairs := make(map[[2]string]int)
	for word, freq := range vocab {
		symbols := strings.Split(word, " ")
		for i := range len(symbols) - 1 {
			pair := [2]string{symbols[i], symbols[i+1]}
			pairs[pair] += freq
		}
	}
	return pairs
}

// mergeVocab 合并词元对
func mergeVocab(pair [2]string, vIn map[string]int) map[string]int {
	vOut := make(map[string]int)
	bigram := regexp.QuoteMeta(strings.Join(pair[:], " "))
	p := regexp.MustCompile(`(?<!\S)` + bigram + `(?!\S)`)
	for word, freq := range vIn {
		wOut := p.ReplaceAllString(word, strings.Join(pair[:], ""))
		vOut[wOut] = freq
	}
	return vOut
}

func main() {
	// 准备语料库，每个词末尾加上</w>表示结束，并切分好字符
	vocab := map[string]int{
		"h u g </w>": 1,
		"p u g </w>": 1,
		"p u n </w>": 1,
		"b u n </w>": 1,
	}
	numMerges := 4 // 设置合并次数

	for i := range numMerges {
		pairs := getStats(vocab)
		if len(pairs) == 0 {
			break
		}
		// 找到频率最高的 pair
		best := maxPair(pairs)
		vocab = mergeVocab(best, vocab)
		fmt.Printf("第%d次合并: %v -> %s\n", i+1, best, strings.Join(best[:], ""))
		fmt.Printf("新词表（部分）: %v\n", mapKeys(vocab))
		fmt.Println(strings.Repeat("-", 20))
	}
}

// maxPair 找到频率最高的 pair
func maxPair(pairs map[[2]string]int) [2]string {
	var best [2]string
	maxFreq := -1
	for pair, freq := range pairs {
		if freq > maxFreq {
			maxFreq = freq
			best = pair
		}
	}
	return best
}

// mapKeys 返回 map 的 keys（排序后）
func mapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
