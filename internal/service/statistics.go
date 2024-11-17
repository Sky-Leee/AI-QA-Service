package service

import (
	"ai-qa-service/internal/api"
	"ai-qa-service/pkg/models"
)

type StatisticsService struct{}

var StatSrv *StatisticsService

func (s *StatisticsService) GetStatistics(userID int64, input *api.StatisticsGetRequest) (*api.StatisticsGetResponse, error) {
	statistics, err := models.UserAnswerStatisticsGetByUserID(userID)
	if err != nil {
		return nil, err
	}

	questionList, err := models.QuestionListGetByUserID(userID, input.Page, input.Size)
	if err != nil {
		return nil, err
	}

	return &api.StatisticsGetResponse{
		Total:          statistics.TotalAnswers,
		CorrectAnswers: statistics.CorrectAnswers,
		AccuracyRate:   statistics.AccuracyRate,
		QuestionList:   questionList,
	}, nil
}
