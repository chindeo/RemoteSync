package main

import (
	"flag"
	"fmt"

	"github.com/chindeo/RemoteSync/models"
	"github.com/chindeo/RemoteSync/utils"
	"github.com/chindeo/pkg/net"

	// _ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

var Version string

var Action = flag.String("action", "", "程序操作指令")

func main() {
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	utils.InitConfig()
	loginData := fmt.Sprintf("app_id=%s&app_secret=%s", utils.Config.Appid, utils.Config.Appsecret)
	err := net.NewNetClient(&net.Config{
		Appid:      utils.Config.Appid,
		AppSecret:  utils.Config.Appsecret,
		LoginUrl:   utils.Config.Host + utils.Config.Loginuri,
		RefreshUrl: utils.Config.Host + utils.Config.Refreshuri,
		LoginData:  loginData,
		TimeOver:   utils.Config.Timeover,
		TimeOut:    utils.Config.Timeout,
		// TokenDriver: "redis",
		// Host:        "10.0.0.26:6379",
		// Pwd:         "Chindeo",
		Headers: map[string]string{
			"AuthType": utils.Config.AuthType,
			"MAC":      utils.Config.Appid,
		},
	})
	if err != nil {
		fmt.Printf("new net client %v \n", err)
		return
	}
	token, err := net.NetClient.GetToken()
	if err != nil {
		fmt.Printf("get token  %v \n", err)
		return
	}
	fmt.Printf("token: %s\n", token)
	if token == "" {
		return
	}

	go models.RemoteSync()
	go models.LocSync()
	go models.UserTypeSync()

	select {}
}
