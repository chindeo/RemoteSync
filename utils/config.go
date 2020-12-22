package utils

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var Config = struct {
	Host         string `default:"" env:"host"`
	Appid        string `default:"" env:"appid"`
	Appsecret    string `default:"" env:"appsecret"`
	DB           string `default:"" env:"db"`
	Timetype     string `default:"" env:"timetype"`
	Timeduration int64  `default:"" env:"timeduration"`
	DevType      string `default:"2" env:"devtype"`
}{}

func init() {
	if err := configor.Load(&Config, ConfigFile()); err != nil {
		panic(fmt.Sprintf("Config Path:%s ,Error:%+v\n", ConfigFile(), err))
	}
}
