package api

import "ai-qa-service/pkg/models"

type UserDTO struct {
	UserID   int64  `json:"user_id,string"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type UserRegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserRegisterResponse struct {
	models.TbUser `json:"user_info"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	models.TbUser `json:"user_info"`
}
