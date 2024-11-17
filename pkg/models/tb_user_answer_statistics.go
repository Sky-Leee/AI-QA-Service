package models

import "github.com/pkg/errors"

type TbUserAnswerStatistics struct {
	UserID         int64   `gorm:"type:bigint;primaryKey;auto_increment" json:"user_id,string"`
	TotalAnswers   int     `gorm:"type:int(11);not null" json:"total_answers"`
	CorrectAnswers int     `gorm:"type:int(11);not null" json:"correct_answers"`
	AccuracyRate   float64 `gorm:"type:float;not null" json:"accuracy_rate"`
}

func UserAnswerStatisticsCreate(userAnswerStatistics *TbUserAnswerStatistics) error {
	err := GetOrmDB().Create(userAnswerStatistics).Error
	if err != nil {
		return errors.Wrapf(err, "create user answer statistics failed, userID: %d", userAnswerStatistics.UserID)
	}

	return nil
}

func UserAnswerStatisticsGetByUserID(userID int64) (*TbUserAnswerStatistics, error) {
	var userAnswertatistics TbUserAnswerStatistics

	err := db.Where("user_id = ?", userID).First(&userAnswertatistics).Error
	if err != nil {
		return nil, errors.Wrapf(err, "get user answer statistics failed, userID: %d", userID)
	}

	return &userAnswertatistics, nil
}

func UserAnswerStatisticsUpdate(userAnswertatistics *TbUserAnswerStatistics) error {
	err := GetOrmDB().Model(&TbUserAnswerStatistics{}).
		Where("user_id = ?", userAnswertatistics.UserID).
		Updates(userAnswertatistics).Error

	if err != nil {
		return errors.Wrapf(err, "update user answer statistics failed, userID: %d", userAnswertatistics.UserID)
	}

	return nil
}
