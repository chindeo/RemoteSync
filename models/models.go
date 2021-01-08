package models

import (
	"fmt"
	"github.com/snowlyg/RemoteSync/utils"
	my "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var sqliteDB *gorm.DB

func GetSqlite() *gorm.DB {
	var err error
	//if sqliteDB != nil {
	//	return sqliteDB
	//}
	file := utils.DBFile()
	sqliteDB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("database sqlite error:%+v", err))
	}
	db, _ := sqliteDB.DB()
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	return sqliteDB
}

func GetMysql() *gorm.DB {
	mysql, err := gorm.Open(my.Open(utils.Config.DB), &gorm.Config{})
	if err != nil {
		fmt.Println(fmt.Sprintf("database mysql init error:%+v", err))
		return nil
	}

	db, _ := mysql.DB()
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return mysql
}

func CloseMysql(mysql *gorm.DB) {
	db, _ := mysql.DB()
	db.Close()
}
