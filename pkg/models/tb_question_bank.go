package models

import (
	"ai-qa-service/pkg/consts"

	"github.com/pkg/errors"
)

type TbQuestionBank struct {
	ID        int64  `gorm:"type:bigint;primaryKey;auto_increment" json:"id,string"` // 题目 ID
	Title     string `gorm:"type:varchar(512);not null;" json:"title"`               // 题目标题
	Content   string `gorm:"type:varchar(4096);not null" json:"content"`             // 题目内容
	CreatedAt Time   `gorm:"type:timestamp default CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt Time   `gorm:"type:timestamp default CURRENT_TIMESTAMP" json:"updated_at"`
}

// 随机抽取 count 道题
func QuestionBankGetRandom(count int) ([]TbQuestionBank, error) {
	var questions []TbQuestionBank

	err := GetOrmDB().Order("RAND()").Limit(count).Find(&questions).Error
	if err != nil {
		return nil, errors.Wrapf(err, "get random question failed, count: %d", count)
	}

	return questions, nil
}

func QuestionBankGetByID(id int64) (TbQuestionBank, error) {
	var question TbQuestionBank

	err := GetOrmDB().Where("id = ?", id).First(&question).Error
	if err != nil {
		return question, errors.Wrapf(err, "get question by id failed, id: %d", id)
	}

	return question, nil
}

func QuestionBankCreate(questionList []TbQuestionBank) error {
	// 1. 获取题库中，题目总数量
	var count int64
	err := GetOrmDB().Model(&TbQuestionBank{}).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "get question count failed")
	}

	// 2. 如果总数量超过最大数量，则不插入数据
	if count >= consts.QUESTION_BANK_MAX_SIZE {
		return nil
	}
	
	// 3. 插入数据
	err = GetOrmDB().Create(&questionList).Error
	if err != nil {
		return errors.Wrap(err, "create question failed")
	}

	return nil
}