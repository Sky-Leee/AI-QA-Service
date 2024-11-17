package api

import "ai-qa-service/pkg/models"

type AnswerSaveRequest struct {
	QuestionListID int64  `json:"question_list_id,string"`
	QuestionID     int64  `json:"question_id,string"`
	Answer         string `json:"answer"`
}

type AnswerSaveResponse struct {
	// Empty
}

type AnswerSubmitRequest struct {
	QuestionListID int64 `json:"question_list_id,string"`
}

type AnswerSubmitResponse struct {
	// Empty
}

type AnswerGetResultRequest struct {
	QuestionListID int64 `form:"question_list_id,string"`
}

type AnswerGetResultResponse struct {
	Total      int                   `json:"total"`
	AnswerList []models.TbUserAnswer `json:"answer_list"`
}
