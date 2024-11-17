package models

import (
	"ai-qa-service/pkg/consts"

	"github.com/pkg/errors"
)

type TbQuestionList struct {
	ID        int64 `gorm:"type:bigint;primaryKey;auto_increment" json:"id,string"` // 题单ID，主键，自增
	UserID    int64 `gorm:"type:int;not null" json:"user_id"`                       // 用户 ID
	Status    int   `gorm:"type:int;not null" json:"status"`                        // 完成状态，0: 未完成，1: 已完成
	CreatedAt Time  `gorm:"type:timestamp default CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt Time  `gorm:"type:timestamp default CURRENT_TIMESTAMP" json:"updated_at"`
}

func QuestionListCreate(userID int64) (int64, error) {
	questionList := TbQuestionList{
		UserID: userID,
		Status: consts.QUESTION_LIST_NOT_FINISHED,
	}

	res := GetOrmDB().Create(&questionList)
	if res.Error != nil {
		return 0, errors.Wrapf(res.Error, "create question list failed, user_id: %d", userID)
	}

	return questionList.ID, res.Error
}

func QuestionListGetByID(id int64) (*TbQuestionList, error) {
	var questionList TbQuestionList
	res := GetOrmDB().Where("id = ?", id).First(&questionList)
	if res.Error != nil {
		return nil, errors.Wrapf(res.Error, "get question list by id failed, id: %d", id)
	}

	return &questionList, nil
}

func QuestionListGetByUserID(userID int64, page, size int) ([]TbQuestionList, error) {
	var questionLists []TbQuestionList

	res := GetOrmDB().Where("user_id = ?", userID).Offset((page - 1) * size).Limit(size).Find(&questionLists)
	if res.Error != nil {
		return nil, errors.Wrapf(res.Error, "get question list by user_id failed, user_id: %d", userID)
	}

	return questionLists, nil
}

func QuestionListUpdateStatusByID(id int64, status int) error {
	res := GetOrmDB().Model(&TbQuestionList{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return errors.Wrapf(res.Error, "update question list status failed, id: %d", id)
	}

	return nil
}
