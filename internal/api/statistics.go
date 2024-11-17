package api

import "ai-qa-service/pkg/models"

type StatisticsGetRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type StatisticsGetResponse struct {
	Total          int                     `json:"total"`
	CorrectAnswers int                     `json:"correct_answers"`
	AccuracyRate   float64                 `json:"accuracy_rate"`
	QuestionList   []models.TbQuestionList `json:"question_list"`
}
