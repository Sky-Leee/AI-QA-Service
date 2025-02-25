package models

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitMySQL() {
	dbHost := viper.Get("mysql.host")
	dbPort := viper.GetInt("mysql.port")
	userName := viper.Get("mysql.username")
	password := viper.Get("mysql.password")
	database := viper.Get("mysql.database")
	charset := viper.Get("mysql.charset")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", userName, password, dbHost, dbPort, database, charset)
	debug := viper.GetBool("mysql.debug")
	var err error
	if debug {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}
	if err != nil {
		panic(fmt.Sprintf("mysql: %s", err.Error()))
	}
	initTables()
	initIndices()
}

func initTables() {
	db.AutoMigrate(&TbUser{})
	db.AutoMigrate(&TbUserAnswer{})
	db.AutoMigrate(&TbQuestionBank{})
	db.AutoMigrate(&TbQuestionList{})
	db.AutoMigrate(&TbUserAnswerStatistics{})
}

func initIndices() {
	// createUnionIndexIfNotExists("idx_cid_uid", "comment_user_like_mappings", "comment_id, user_id", false)
	// createUnionIndexIfNotExists("idx_cid_uid", "comment_user_hate_mappings", "comment_id, user_id", false)
	// createUnionIndexIfNotExists("idx_uid_oid_otype", "comment_user_like_mappings", "user_id, obj_id, obj_type", false)
	// createUnionIndexIfNotExists("idx_uid_oid_otype", "comment_user_like_mappings", "user_id, obj_id, obj_type", false)
	// createUnionIndexIfNotExists("idx_oid_otype", "comment_subjects", "obj_id, obj_type", true)
}

func createUnionIndexIfNotExists(indexName, tableName, columns string, unique bool) {
	var indexCount int64
	db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME = ? AND INDEX_NAME = ?", tableName, indexName).Count(&indexCount)

	if indexCount == 0 {
		// 不存在则创建索引
		sql := fmt.Sprintf("CREATE INDEX %s ON %s(%s);", indexName, tableName, columns)
		if unique {
			sql = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s(%s);", indexName, tableName, columns)
		}
		db.Exec(sql)
	}
}

func GetOrmDB() *gorm.DB {
	return db
}
