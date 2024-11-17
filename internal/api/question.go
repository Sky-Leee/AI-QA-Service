package api

import "ai-qa-service/pkg/models"

type QuestionsGetRequest struct {
	// empty
}

type QuestionsGetResponse struct {
	Total          int                     `json:"total"`
	QuestionListID int64                   `json:"question_list_id,string"`
	QuestionList   []models.TbQuestionBank `json:"question_list"`
}
