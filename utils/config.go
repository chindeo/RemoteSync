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
	DevType      int64  `default:"2" env:"devtype"` // device type 设备类型
	AuthType     string `default:"2" env:"authtype"`
	Outdir       string `default:"" env:"outdir"`
	Timeout      int64  `default:"10" env:"timeout"`
	Timeover     int64  `default:"5" env:"timeover"`
	Loginuri     string `default:"/api/v1/get_access_token" env:"Loginuri"`
	Refreshuri   string `default:"/api/v1/refresh_token" env:"Refreshuri"`
	IsDev        string `default:"" env:"Isdev"`
}{}

func InitConfig() {
	if err := configor.Load(&Config, ConfigFile()); err != nil {
		panic(fmt.Sprintf("Config Path:%s ,Error:%+v\n", ConfigFile(), err))
	}
	fmt.Printf("config: %+v\n", Config)
}
