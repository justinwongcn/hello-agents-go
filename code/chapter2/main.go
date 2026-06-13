package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------
// 1. 规则库：模式(正则表达式) -> 响应模板列表（有序 slice，保证匹配优先级）
// ---------------------------------------------------------------------------

type rule struct {
	pattern   *regexp.Regexp
	responses []string
}

var rules = []rule{
	{regexp.MustCompile(`(?i)I need (.*)`), []string{
		"Why do you need {0}?",
		"Would it really help you to get {0}?",
		"Are you sure you need {0}?",
	}},
	{regexp.MustCompile(`(?i)Why don't you (.*)\?`), []string{
		"Do you really think I don't {0}?",
		"Perhaps eventually I will {0}.",
		"Do you really want me to {0}?",
	}},
	{regexp.MustCompile(`(?i)Why can't I (.*)\?`), []string{
		"Do you think you should be able to {0}?",
		"If you could {0}, what would you do?",
		"I don't know -- why can't you {0}?",
	}},
	{regexp.MustCompile(`(?i)I am (.*)`), []string{
		"Did you come to me because you are {0}?",
		"How long have you been {0}?",
		"How do you feel about being {0}?",
	}},
	{regexp.MustCompile(`(?i).* mother .*`), []string{
		"Tell me more about your mother.",
		"What was your relationship with your mother like?",
		"How do you feel about your mother?",
	}},
	{regexp.MustCompile(`(?i).* father .*`), []string{
		"Tell me more about your father.",
		"How did your father make you feel?",
		"What has your father taught you?",
	}},
	{regexp.MustCompile(`(?i).*`), []string{
		"Please tell me more.",
		"Let's change focus a bit... Tell me about your family.",
		"Can you elaborate on that?",
	}},
}

// ---------------------------------------------------------------------------
// 2. 代词转换规则
// ---------------------------------------------------------------------------

var pronounSwap = map[string]string{
	"i": "you", "you": "i", "me": "you", "my": "your",
	"am": "are", "are": "am", "was": "were", "i'd": "you would",
	"i've": "you have", "i'll": "you will", "yours": "mine",
	"mine": "yours",
}

// swapPronouns 对输入短语中的代词进行第一/第二人称转换
func swapPronouns(phrase string) string {
	words := strings.Fields(strings.ToLower(phrase))
	for i, word := range words {
		if swapped, ok := pronounSwap[word]; ok {
			words[i] = swapped
		}
	}
	return strings.Join(words, " ")
}

// ---------------------------------------------------------------------------
// 3. 响应生成
// ---------------------------------------------------------------------------

// respond 根据规则库生成响应
func respond(userInput string) string {
	for _, r := range rules {
		if match := r.pattern.FindStringSubmatch(userInput); match != nil {
			// 捕获匹配到的部分
			capturedGroup := ""
			if len(match) > 1 {
				capturedGroup = match[1]
			}
			// 进行代词转换
			swappedGroup := swapPronouns(capturedGroup)
			// 从模板中随机选择一个并格式化
			response := r.responses[rand.N(len(r.responses))]
			return strings.ReplaceAll(response, "{0}", swappedGroup)
		}
	}
	// 如果没有匹配任何特定规则，使用最后的通配符规则
	return rules[len(rules)-1].responses[rand.N(len(rules[len(rules)-1].responses))]
}

// ---------------------------------------------------------------------------
// 4. 主聊天循环
// ---------------------------------------------------------------------------

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Therapist: Hello! How can I help you today?")

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		lower := strings.ToLower(userInput)
		if lower == "quit" || lower == "exit" || lower == "bye" {
			fmt.Println("Therapist: Goodbye. It was nice talking to you.")
			break
		}

		response := respond(userInput)
		fmt.Printf("Therapist: %s\n", response)
	}
}
