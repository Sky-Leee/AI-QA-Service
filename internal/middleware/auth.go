package middleware

import (
	"ai-qa-service/internal/handler"
	"ai-qa-service/internal/utils"
	"ai-qa-service/pkg/consts"
	"ai-qa-service/pkg/errno"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 认证中间件，基于 JWT
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Authorization 头
		header := ctx.Request.Header.Get("Authorization")
		if len(header) == 0 {
			handler.ResponseError(ctx, consts.CodeNeedLogin)
			ctx.Abort()
			return
		}

		// 获取协议 和 access_token
		// 这里使用 Bearer 作为 协议
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.ResponseError(ctx, consts.CodeUnsupportedAuthProtocol)
			ctx.Abort()
			return
		}
		if parts[1] == "null" {
			handler.ResponseError(ctx, consts.CodeInvalidToken)
			ctx.Abort()
			return
		}

		// 检验 token
		UserID, err := utils.ParseToken(parts[1])
		if err != nil {
			if errors.Is(err, errno.ErrInvalidToken) {
				handler.ResponseError(ctx, consts.CodeInvalidToken)
			} else if errors.Is(err, errno.ErrExpiredToken) {
				handler.ResponseError(ctx, consts.CodeExpiredToken)
			} else {
				handler.ResponseErrorWithMsg(ctx, consts.CodeInternalErr, "解析 token 失败")
			}
			ctx.Abort()
			return
		}

		ctx.Set("user_id", UserID)
		ctx.Set("access_token", parts[1]) // 用于后续限制一个用户登录
		ctx.Next()
	}
}
