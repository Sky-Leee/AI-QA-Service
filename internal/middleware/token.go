package middleware

import (
	"ai-qa-service/internal/handler"
	"ai-qa-service/internal/logger"
	"ai-qa-service/pkg/consts"
	"ai-qa-service/pkg/models"

	"github.com/gin-gonic/gin"
)

// 校验上下文的 access_token 是否与 db 中的一致
func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetInt64("user_id")
		access_token := ctx.GetString("access_token")

		r_access_token, err := models.UserGetAccessToken(userID)
		if err != nil {
			logger.ErrorWithStack(err)
			handler.ResponseError(ctx, consts.CodeInternalErr)
			ctx.Abort()
		} else if access_token != r_access_token {
			handler.ResponseError(ctx, consts.CodeNeedLogin)
			ctx.Abort()
		}
		ctx.Next()
	}
}
