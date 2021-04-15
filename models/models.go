package models

import (
	"fmt"

	"github.com/chindeo/RemoteSync/utils"
	my "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMysql() (*gorm.DB, error) {
	mysql, err := gorm.Open(my.Open(utils.Config.DB), &gorm.Config{})
	if err != nil {
		fmt.Println(fmt.Sprintf("database mysql init error:%+v", err))
		return nil, err
	}

	db, err := mysql.DB()
	if err != nil {
		fmt.Println(fmt.Sprintf("database mysql error:%+v", err))
		return nil, err
	}
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return mysql, nil
}

func CloseMysql(mysql *gorm.DB) {
	db, _ := mysql.DB()
	db.Close()
}
