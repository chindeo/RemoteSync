package models

import (
	"fmt"
	"github.com/snowlyg/RemoteSync/utils"
	my "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var sqliteDB *gorm.DB
var mysql *gorm.DB

func GetSqlite() *gorm.DB {
	var single sync.Once
	var err error
	single.Do(func() {
		file := utils.DBFile()
		sqliteDB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("database sqlite error:%+v", err))
		}

		err = sqliteDB.AutoMigrate(&RemoteDev{}, &Loc{}, &UserType{})
		if err != nil {
			fmt.Println(fmt.Sprintf("database model init error:%+v", err))
		}
	})
	return sqliteDB
}

func GetMysql() *gorm.DB {
	var single sync.Once
	var err error
	single.Do(func() {
		mysql, err = gorm.Open(my.Open(utils.Config.DB), &gorm.Config{})
		if err != nil {
			fmt.Println(fmt.Sprintf("database mysql init error:%+v", err))
		}
	})

	mysql, err = gorm.Open(my.Open(utils.Config.DB), &gorm.Config{})
	if err != nil {
		fmt.Println(fmt.Sprintf("database mysql init error:%+v", err))
	}
	return mysql

}

func Close() {
}
