package handler

import (
	"ai-qa-service/internal/api"
	"ai-qa-service/internal/logger"
	"ai-qa-service/internal/service"
	"ai-qa-service/pkg/consts"

	"github.com/gin-gonic/gin"
)

func RegistQuestionHandler(quesGrp *gin.RouterGroup) {
	{
		quesGrp.GET("/get", QuestionGetListHandler)      // 获取题单
		quesGrp.POST("/save", QuestionSaveHandler)       // 保存某一道题的答案
		quesGrp.POST("/submit", QuestionSubmitHandler)   // 提交题单所有答案
		quesGrp.GET("/result", QuestionGetResultHandler) // 获取某一个题单的答题结果
	}
}

// QuestionGetListHandler 获取题单接口
//
//	@Summary		获取题单接口
//	@Description	用户获取题单接口
//	@Tags			题单相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@Success		200				{object}	handler.Response{data=api.QuestionsGetResponse}
//	@Router			/question/get [get]
func QuestionGetListHandler(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	questionListID, questions, err := service.QuestionSrv.GetQuestionList(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		ResponseError(ctx, consts.CodeInternalErr)
		return
	}

	resp := api.QuestionsGetResponse{
		Total:          len(questions),
		QuestionListID: questionListID,
		QuestionList:   questions,
	}

	ResponseSuccess(ctx, resp)
}

// QuestionSaveHandler 保存答案接口
//
//	@Summary		保存答案接口
//	@Description	保存某一道题的答案
//	@Tags			题单相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			answer_info		body		api.AnswerSaveRequest	false	"保存的答案信息"
//	@Param			Authorization	header		string					false	"Bearer 用户令牌"
//	@Success		200				{object}	handler.Response{data=api.AnswerSaveResponse}
//	@Router			/question/save [post]
func QuestionSaveHandler(ctx *gin.Context) {
	input := api.AnswerSaveRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ResponseError(ctx, consts.CodeInvalidParam)
		return
	}

	if err := service.QuestionSrv.SaveAnswer(&input); err != nil {
		logger.ErrorWithStack(err)
		ResponseError(ctx, consts.CodeInternalErr)
		return
	}

	ResponseSuccess(ctx, api.AnswerSaveResponse{})
}

// QuestionSubmitHandler 提交题单接口
//
//	@Summary		提交题单接口
//	@Description	用户完成一个题单的所有题目后，提交题单（注意，调用此接口前，需要保证所有答案已经保存过）
//	@Tags			题单相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			answer_info		body		api.AnswerSubmitRequest	false	"题单号"
//	@Param			Authorization	header		string					false	"Bearer 用户令牌"
//	@Success		200				{object}	handler.Response{data=api.AnswerSubmitResponse}
//	@Router			/question/submit [post]
func QuestionSubmitHandler(ctx *gin.Context) {
	input := api.AnswerSubmitRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ResponseError(ctx, consts.CodeInvalidParam)
		return
	}

	if err := service.QuestionSrv.SubmitAnswer(&input); err != nil {
		logger.ErrorWithStack(err)
		ResponseError(ctx, consts.CodeInternalErr)
		return
	}

	ResponseSuccess(ctx, api.AnswerSubmitResponse{})
}

// QuestionGetResultHandler 获取答题结果接口
//
//	@Summary		获取答题结果接口
//	@Description	获取一个题单的答题结果
//	@Tags			题单相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			answer_info		query		api.AnswerGetResultRequest	false	"题单号"
//	@Param			Authorization	header		string						false	"Bearer 用户令牌"
//	@Success		200				{object}	handler.Response{data=api.AnswerGetResultResponse}
//	@Router			/question/result [get]
func QuestionGetResultHandler(ctx *gin.Context) {
	input := api.AnswerGetResultRequest{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		ResponseError(ctx, consts.CodeInvalidParam)
		return
	}

	answerList, err := service.QuestionSrv.GetQuestionListResult(&input)
	if err != nil {
		logger.ErrorWithStack(err)
		ResponseError(ctx, consts.CodeInternalErr)
		return
	}

	ResponseSuccess(ctx, api.AnswerGetResultResponse{
		Total:      len(answerList),
		AnswerList: answerList,
	})
}
