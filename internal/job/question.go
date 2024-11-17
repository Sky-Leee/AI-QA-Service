package job

import (
	"ai-qa-service/internal/coze"
	"ai-qa-service/internal/logger"
	"ai-qa-service/pkg/models"
	"time"

	"github.com/spf13/viper"
)

var (
	maxRetryTime time.Duration
	sleepTime    time.Duration
)

func GenQuestions() {
	maxRetryTime = time.Duration(viper.GetInt("job.question.max_retry_time")) * time.Second
	sleepTime = time.Duration(viper.GetInt("job.question.sleep_time")) * time.Second

	go func() {
		var err error
		retryTime := time.Second
		for {
			time.Sleep(sleepTime)
			if checkIfExit() {
				return				
			}

			// 上一次执行失败
			if err != nil {
				markAsExit()
				time.Sleep(retryTime)
				retryTime = max(retryTime*2, maxRetryTime)
			}
	
			var questionList []models.TbQuestionBank
			
			// 大语言模型生成多个问题
			logger.Infof("Trying to get question list from coze")
			questionList, err = coze.GetQuestionList()
			if err != nil {
				logger.Errorf("Get question list from coze failed: %v", err)
				markAsExit()
				continue
			}
			logger.Infof("Get question list from coze success")
	
			// 插入到题库表
			err = models.QuestionBankCreate(questionList)
			if err != nil {
				logger.Errorf("Insert question list to question bank failed: %v")
				markAsExit()
				continue
			}
			logger.Infof("Insert question list to question bank success")
	
			// 正常执行，清空上一次执行状态
			retryTime = time.Second
			err = nil
			markAsExit()
		}
	}()
}
