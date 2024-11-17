package coze

import (
	"ai-qa-service/internal/logger"
	"ai-qa-service/pkg/consts"
	"ai-qa-service/pkg/models"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func GetQuestionList() ([]models.TbQuestionBank, error) {
	args := getArgs(genBotID, consts.COZE_GEN_KEYWORD)
	reply := &CozeResponseBody{}

	if err := SendCozeRequest(args, reply); err != nil {
		return nil, errors.Wrap(err, "get question list from coze failed")
	}

	content, err := parseCozeResponseToAnswer(reply)
	if err != nil {
		return nil, err
	}

	return parseGenQuestionOutput(content)
}

func GetCorrectAnswerAndSuggestion(title, content, answer string) (*models.TbUserAnswer, error) {
	query := fmt.Sprintf(consts.COZE_JUDGE_KEYWORD, title, content, answer)
	args := getArgs(judgeBotID, query)
	reply := &CozeResponseBody{}

	if err := SendCozeRequest(args, reply); err != nil {
		return nil, errors.Wrap(err, "get correct answer from coze failed")
	}

	logger.Debugf("reply from coze: %+v", reply)
	content, err := parseCozeResponseToAnswer(reply)
	if err != nil {
		return nil, err
	}

	return parseJudgeOutput(content)
}

// parseGenQuestionOutput 将内容解析为 TbQuestionBank 结构体的切片
func parseGenQuestionOutput(content string) ([]models.TbQuestionBank, error) {
	// 去除开头和结尾的标识符
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "<-START->") || !strings.HasSuffix(content, "<-END->") {
		return nil, errors.New("invalid content format")
	}

	// 去除标识符，提取中间的内容
	content = strings.TrimPrefix(content, "<-START->")
	content = strings.TrimSuffix(content, "<-END->")

	// 分割每一题
	questions := []models.TbQuestionBank{}
	questionBlocks := strings.Split(content, "---")

	// 使用正则提取题目和选项
	re := regexp.MustCompile(`^(\d+\..+)`)

	for _, block := range questionBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		// 提取题目标题，题目是以数字和句点开始的
		match := re.FindStringSubmatch(block)
		if match == nil {
			return nil, fmt.Errorf("invalid question format: %s", block)
		}

		// 获取题目标题
		title := match[1]

		// 剩余部分是选项
		content := strings.TrimSpace(block[len(title):])

		questions = append(questions, models.TbQuestionBank{
			Title:   title,
			Content: content,
		})
	}

	return questions, nil
}

func parseJudgeOutput(content string) (*models.TbUserAnswer, error) {
	logger.Infof("judge content: %s", content)
	// 正则表达式匹配 "Result: " 和 "Explanation: " 后面的内容
	resultPattern := `(?i)Result:\s*(\w)`
	explanationPattern := `(?i)Explanation:\s*(.*)`

	// 使用正则表达式提取结果（Correct/Incorrect）
	resultRe := regexp.MustCompile(resultPattern)
	resultMatch := resultRe.FindStringSubmatch(content)
	if len(resultMatch) < 2 {
		return nil, errors.New("unable to parse Result from content")
	}

	// 判断答案是否正确
	isCorrect := resultMatch[1] == "T"

	// 使用正则表达式提取解释建议
	explanationRe := regexp.MustCompile(explanationPattern)
	explanationMatch := explanationRe.FindStringSubmatch(content)
	if len(explanationMatch) < 2 {
		return nil, errors.New("unable to parse Explanation from content")
	}
	suggestion := explanationMatch[1]

	// 去掉多余的空格
	suggestion = strings.TrimSpace(suggestion)

	return &models.TbUserAnswer{
		IsCorrect:  isCorrect,
		Suggestion: suggestion,
	}, nil
}
