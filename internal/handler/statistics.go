package handler

import (
	"ai-qa-service/internal/api"
	"ai-qa-service/internal/logger"
	"ai-qa-service/internal/service"
	"ai-qa-service/pkg/consts"

	"github.com/gin-gonic/gin"
)

func RegistStatisticsHandler(statGrp *gin.RouterGroup) {
	{
		// 获取用户答题情况
		statGrp.GET("/get", StatisticsGetHandler)
	}
}

// StatisticsGetHandler 获取用户答题情况
//
//	@Summary		获取用户答题情况
//	@Description	获取用户答题情况
//	@Tags			统计相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			Authorization	header		string						false	"Bearer 用户令牌"
//	@Param			user_info		query		api.StatisticsGetRequest	false	"分页参数"
//	@Success		200				{object}	handler.Response{data=api.StatisticsGetResponse}
//	@Router			/statistics/get [get]
func StatisticsGetHandler(ctx *gin.Context) {
	input := api.StatisticsGetRequest{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		ResponseError(ctx, consts.CodeInvalidParam)
		return
	}

	userID := ctx.GetInt64("user_id")

	statistics, err := service.StatSrv.GetStatistics(userID, &input)
	if err != nil {
		logger.ErrorWithStack(err)
		ResponseError(ctx, consts.CodeInternalErr)
		return
	}

	ResponseSuccess(ctx, statistics)
}
