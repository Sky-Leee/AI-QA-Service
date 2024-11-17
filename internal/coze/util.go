package coze

import "fmt"

// 解析响应结果，获取 assistant 的 answer
func parseCozeResponseToAnswer(resp *CozeResponseBody) (string, error) {
	for _, msg := range resp.Messages {
		if msg.Role == "assistant" && msg.Type == "answer" {
			return msg.Content, nil
		}
	}

	return "", fmt.Errorf("no assistant answer found in response")
}
