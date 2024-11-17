package service

import (
	"ai-qa-service/internal/api"
	"ai-qa-service/internal/coze"
	"ai-qa-service/pkg/consts"
	"ai-qa-service/pkg/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type QuestionService struct{}

var QuestionSrv = &QuestionService{}

func (q *QuestionService) GetQuestionList(userID int64) (questionListID int64, questions []models.TbQuestionBank, err error) {
	// 1. 随机抽取 10 道题
	questions, err = models.QuestionBankGetRandom(consts.QUESTION_BATCH_SIZE)
	if err != nil {
		return
	}

	// 2. 创建题单
	questionListID, err = models.QuestionListCreate(userID)
	if err != nil {
		return
	}

	// 3. 将题单和题目关联
	for _, question := range questions {
		if err = models.UserAnswerCreate(questionListID, question.ID); err != nil {
			return
		}
	}

	return
}

func (q *QuestionService) SaveAnswer(input *api.AnswerSaveRequest) error {
	userAnswer := models.TbUserAnswer{
		UserAnswer: input.Answer,
	}

	return models.UserAnswerUpdate(input.QuestionListID, input.QuestionID, userAnswer)
}

func (q *QuestionService) SubmitAnswer(input *api.AnswerSubmitRequest) error {
	// 1. 校验是否已经提交过
	questionList, err := models.QuestionListGetByID(input.QuestionListID)
	if err != nil {
		return err
	}
	if questionList.Status == consts.QUESTION_LIST_FINISHED {
		return nil
	}

	// 2. 获取用户答题记录
	userAnswers, err := models.UserAnswerGetByQuestionListID(input.QuestionListID)
	if err != nil {
		return err
	}
	questions := make([]models.TbQuestionBank, 0, len(userAnswers))
	for _, answer := range userAnswers {
		question, err := models.QuestionBankGetByID(answer.QuestionID)
		if err != nil {
			return err
		}
		questions = append(questions, question)
	}

	var (
		totalAnswer   = len(userAnswers)
		correctAnswer = 0
	)

	// 3. 使用大语言模型判题，得到答案以及建议
	for i := 0; i < len(questions); i++ {
		cozeAnswer, err := coze.GetCorrectAnswerAndSuggestion(questions[i].Title, questions[i].Content, userAnswers[i].UserAnswer)
		if err != nil {
			return err
		}
		userAnswers[i].IsCorrect = cozeAnswer.IsCorrect
		userAnswers[i].Suggestion = cozeAnswer.Suggestion

		if cozeAnswer.IsCorrect {
			correctAnswer++
		}
	}

	// 4. 更新用户答案表
	for _, userAnswer := range userAnswers {
		if err = models.UserAnswerUpdate(userAnswer.QuestionListID, userAnswer.QuestionID, userAnswer); err != nil {
			return err
		}
	}

	// 5. 更新题单状态
	if err := models.QuestionListUpdateStatusByID(input.QuestionListID, consts.QUESTION_LIST_FINISHED); err != nil {
		return err
	}

	// 6. 更新用户答题情况
	statistics, err := models.UserAnswerStatisticsGetByUserID(questionList.UserID)
	if err != nil {
		// 如果记录不存在，则创建一条新的记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statistics = &models.TbUserAnswerStatistics{
				UserID:         questionList.UserID,
				TotalAnswers:   totalAnswer,
				CorrectAnswers: correctAnswer,
				AccuracyRate:   float64(correctAnswer) / float64(totalAnswer),
			}
			return models.UserAnswerStatisticsCreate(statistics)
		} else {
			return err
		}
	}
	statistics.TotalAnswers += totalAnswer
	statistics.CorrectAnswers += correctAnswer
	statistics.AccuracyRate = float64(statistics.CorrectAnswers) / float64(statistics.TotalAnswers)

	return models.UserAnswerStatisticsUpdate(statistics)
}

func (q *QuestionService) GetQuestionListResult(input *api.AnswerGetResultRequest) ([]models.TbUserAnswer, error) {
	return models.UserAnswerGetByQuestionListID(input.QuestionListID)
}
