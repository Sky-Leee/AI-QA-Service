package models

import "github.com/pkg/errors"

type TbUser struct {
	ID          int64  `gorm:"type:bigint;primaryKey;auto_increment" json:"id,string"`
	UserName    string `gorm:"type:varchar(64);not null;unique" json:"username"`
	Password    string `gorm:"type:varchar(64);not null" json:"password"`
	Email       string `gorm:"type:varchar(64);not null;unique" json:"email"`
	AccessToken string `gorm:"type:varchar(256);not null" json:"access_token"`
	CreatedAt   Time   `gorm:"type:timestamp default CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   Time   `gorm:"type:timestamp default CURRENT_TIMESTAMP" json:"updated_at"`
}

func UserGetAccessToken(userID int64) (string, error) {
	accessToken := ""
	err := GetOrmDB().Model(&TbUser{}).
		Select("access_token").
		Where("id = ?", userID).
		First(&accessToken).Error

	if err != nil {
		return "", errors.Wrapf(err, "UserGetAccessToken failed, userID: %d", userID)
	}
	return accessToken, nil
}

func UserUpdateAccessToken(userID int64, accessToken string) error {
	err := GetOrmDB().Model(&TbUser{}).
		Where("id = ?", userID).
		Update("access_token", accessToken).Error

	if err != nil {
		return errors.Wrapf(err, "UserUpdateAccessToken failed, userID: %d, accessToken: %s", userID, accessToken)
	}

	return nil
}

func UserCreate(user *TbUser) error {
	err := GetOrmDB().Create(user).Error
	if err != nil {
		return errors.Wrapf(err, "UserCreate failed, user: %+v", user)
	}
	return nil
}

func UserGetByName(username string) (*TbUser, error) {
	user := &TbUser{}

	err := GetOrmDB().Where("user_name = ?", username).First(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, "UserGetByName failed, username: %s", username)
	}

	return user, nil
}

func UserGetByEmail(email string) (*TbUser, error) {
	user := &TbUser{}

	err := GetOrmDB().Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, "UserGetByEmail failed, email: %s", email)
	}

	return user, nil
}

func UserGetByID(userID int64) (*TbUser, error) {
	user := &TbUser{}

	err := GetOrmDB().Where("id = ?", userID).First(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, "UserGetByID failed, userID: %d", userID)
	}

	return user, nil
}
