package handler

import (
	"ai-qa-service/internal/api"
	"ai-qa-service/internal/logger"
	"ai-qa-service/internal/service"
	"ai-qa-service/pkg/consts"
	"ai-qa-service/pkg/errno"
	"errors"

	"github.com/gin-gonic/gin"
)

func RegistUserHandler(usrGrp *gin.RouterGroup) {
	{
		// 用户注册
		usrGrp.POST("/regist", UserRegistHandler)
		// 用户登录
		usrGrp.POST("/login", UserLoginHandler)
	}
}

// UserRegistHandler 用户注册接口
//
//	@Summary		用户注册接口
//	@Description	用户注册接口
//	@Tags			用户相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			user_info	body		api.UserRegisterRequest	false	"用户注册信息"
//	@Success		200			{object}	handler.Response{data=api.UserRegisterResponse}
//	@Router			/user/regist [post]
func UserRegistHandler(ctx *gin.Context) {
	input := api.UserRegisterRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ResponseError(ctx, consts.CodeInvalidParam)
		return
	}

	user, err := service.UserSrv.Regist(&input)
	if err != nil {
		if errors.Is(err, errno.ErrEmailExist) || errors.Is(err, errno.ErrUserExist) {
			ResponseError(ctx, consts.CodeUserExist)
		} else {
			logger.ErrorWithStack(err)
			ResponseError(ctx, consts.CodeInternalErr)
		}
		return
	}

	ResponseSuccess(ctx, api.UserRegisterResponse{
		TbUser: *user,
	})
}

// UserLoginHandler 用户登录接口
//
//	@Summary		用户登录接口
//	@Description	用户登录接口
//	@Tags			用户相关接口
//	@Accept			application/json
//	@Produce		application/json
//	@Param			user_info	body		api.UserLoginRequest	false	"用户登录信息"
//	@Success		200			{object}	handler.Response{data=api.UserLoginResponse}
//	@Router			/user/login [post]
func UserLoginHandler(ctx *gin.Context) {
	input := api.UserLoginRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ResponseError(ctx, consts.CodeInvalidParam)
		return
	}

	user, err := service.UserSrv.Login(&input)
	if err != nil {
		if errors.Is(err, errno.ErrUserNotExist) {
			ResponseError(ctx, consts.CodeUserNotExist)
		} else if errors.Is(err, errno.ErrWrongPassword) {
			ResponseError(ctx, consts.CodeWrongPassword) 
		} else {
			logger.ErrorWithStack(err)
			ResponseError(ctx, consts.CodeInternalErr)
		}
	}

	ResponseSuccess(ctx, api.UserLoginResponse{
		TbUser: *user,
	})
}
