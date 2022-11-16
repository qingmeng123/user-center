/*******
* @Author:qingmeng
* @Description:
* @File:init
* @Date2022/7/14
 */

package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"user-center/model"
)

var db *gorm.DB

func InitDatabase(dsn string) {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true}, //禁用表名加s

		Logger:                                   logger.Default.LogMode(logger.Info), //打印sql语句
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true, //禁止创建外键约束
	})
	if err != nil {
		panic("Connecting database failed:" + err.Error())
	}

	//迁移
	db.AutoMigrate(&model.User{})
}
