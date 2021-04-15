package main

import (
	"flag"
	"fmt"
	"os"

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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  -action= <install remove start stop restart version remote_sync loc_sync user_type_sync cache_clear>\n")
		fmt.Fprintf(os.Stderr, "    程序操作指令 \n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()
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
		// Host:        "127.0.0.1:6379",
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
