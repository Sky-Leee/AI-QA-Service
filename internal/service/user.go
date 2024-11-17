package service

import (
	"ai-qa-service/internal/api"
	"ai-qa-service/internal/utils"
	"ai-qa-service/pkg/errno"
	"ai-qa-service/pkg/models"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

var UserSrv *UserService

func (u *UserService) Regist(input *api.UserRegisterRequest) (*models.TbUser, error) {
	// 1. 校验邮箱是否存在
	if err := u.checkEmailIfExist(input.Email); err != nil {
		return nil, err
	}
	// 2. 校验用户名是否存在
	if err := u.checkUserIfExist(input.UserName); err != nil {
		return nil, err
	}

	// 3. 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户
	err = models.UserCreate(&models.TbUser{
		UserName: input.UserName,
		Password: string(hashedPassword),
		Email:    input.Email,
	})
	if err != nil {
		return nil, err
	}

	// 5. 获取 userID，以生成 accessToken
	user, err := models.UserGetByName(input.UserName)
	if err != nil {
		return nil, err
	}

	// 6. 生成 accessToken
	accessToken, err := utils.GenToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 7. 更新 accessToken
	if err := models.UserUpdateAccessToken(user.ID, accessToken); err != nil {
		return nil, err
	}

	user.Password = input.Password
	user.AccessToken = accessToken
	return user, nil
}

func (u *UserService) Login(input *api.UserLoginRequest) (*models.TbUser, error) {
	// 1. 获取 user
	user, err := models.UserGetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotExist
		}
		return nil, err
	}

	// 2. 校验密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errno.ErrWrongPassword
	}

	// 3. 生成 accessToken
	accessToken, err := utils.GenToken(user.ID)

	// 4. 更新 accessToken
	if err := models.UserUpdateAccessToken(user.ID, accessToken); err != nil {
		return nil, err
	}
	user.AccessToken = accessToken

	return user, nil
}

func (u *UserService) checkUserIfExist(username string) error {
	_, err := models.UserGetByName(username)
	if err == nil {
		return errno.ErrUserExist
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return nil
}
func (u *UserService) checkEmailIfExist(email string) error {
	_, err := models.UserGetByEmail(email)
	if err == nil {
		return errno.ErrEmailExist
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return nil
}