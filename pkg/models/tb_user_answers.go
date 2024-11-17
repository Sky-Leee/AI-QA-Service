package models

import "github.com/pkg/errors"

type TbUserAnswer struct {
	QuestionListID int64  `gorm:"type:bigint" json:"question_list_id,string"`     // 题单 ID
	QuestionID     int64  `gorm:"type:bigint" json:"question_id,string"`          // 题目 ID
	UserAnswer     string `gorm:"type:varchar(1024);not null" json:"user_answer"` // 用户的答案
	IsCorrect      bool   `gorm:"type:tinyint(1);not null"`                       // 是否答对
	Suggestion     string `gorm:"type:varchar(1024);not null" json:"suggestion"`  // 大语言模型的建议
}

func UserAnswerCreate(questionListID, questionID int64) error {
	err := GetOrmDB().Create(&TbUserAnswer{
		QuestionListID: questionListID,
		QuestionID:     questionID,
		UserAnswer:     "",
		IsCorrect:      false,
	}).Error

	if err != nil {
		return errors.Wrapf(err, "create user answer failed, questionListID: %d, questionID: %d", questionListID, questionID)
	}

	return nil
}

func UserAnswerUpdate(questionListID, questionID int64, userAnswer TbUserAnswer) error {
	err := GetOrmDB().Model(&TbUserAnswer{}).
		Where("question_list_id = ? and question_id = ?", questionListID, questionID).
		Updates(userAnswer).Error

	if err != nil {
		return errors.Wrapf(err, "update user answer failed, questionListID: %d, questionID: %d", questionListID, questionID)
	}

	return nil
}

func UserAnswerGetByQuestionListID(questionListID int64) ([]TbUserAnswer, error) {
	var userAnswers []TbUserAnswer
	err := GetOrmDB().Where("question_list_id = ?", questionListID).Find(&userAnswers).Error

	if err != nil {
		return nil, errors.Wrapf(err, "get user answer failed, questionListID: %d", questionListID)
	}

	return userAnswers, nil

}
