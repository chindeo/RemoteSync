package models

import (
	"fmt"
	"github.com/snowlyg/RemoteSync/utils"
	my "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Sqlite *gorm.DB
var Mysql *gorm.DB
var err error

func init() {
	file := utils.DBFile()
	Sqlite, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("database sqlite error:%+v", err))
	}

	err = Sqlite.AutoMigrate(&RemoteDev{}, &Loc{}, &UserType{})
	if err != nil {
		panic(fmt.Sprintf("database model init error:%+v", err))
	}

	Mysql, err = gorm.Open(my.Open(utils.Config.DB), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("database mysql init error:%+v", err))
	}
}

func Close() {
}
