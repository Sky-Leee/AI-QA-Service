package errno

import "github.com/pkg/errors"

var (
	// user
	ErrUserExist     = errors.New("用户已经存在")
	ErrUserNotExist  = errors.New("用户不存在")
	ErrWrongPassword = errors.New("密码错误")

	// common
	ErrGenToken     = errors.New("生成 Token 失败")
	ErrInvalidToken = errors.New("无效的 Token")
	ErrExpiredToken = errors.New("过期的 Token")
	ErrNotFound     = errors.New("未找到")
	ErrInternal     = errors.New("内部错误")
	ErrTimeout      = errors.New("请求超时")

	// params
	ErrInvalidParam = errors.New("无效参数")

	// permissions
	ErrForbidden = errors.New("禁止访问")

	// email
	ErrEmailExist = errors.New("邮箱已存在")
)
