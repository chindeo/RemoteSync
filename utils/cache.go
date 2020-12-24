package utils

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

var CC *cache.Cache
var ai *AppInfo
var phpsess *http.Cookie
var token string

func init() {
	CC = cache.New(1*time.Hour, 2*time.Hour)
}

func SetCacheToken(token string) {
	CC.Set(fmt.Sprintf("XToken_%s", Config.Appid), token, cache.DefaultExpiration)
}

func GetCacheToken() string {
	if token != "" {
		return token
	}
	foo, found := CC.Get(fmt.Sprintf("XToken_%s", Config.Appid))
	if found {
		token = foo.(string)
	}
	return token
}
func SetSessionId(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			CC.Set(fmt.Sprintf("PHPSESSIONID_%s", Config.Appid), cookie, cache.DefaultExpiration)
		}
	}
}

func GetSessionId() *http.Cookie {
	if phpsess != nil {
		return phpsess
	}
	foo, found := CC.Get(fmt.Sprintf("PHPSESSIONID_%s", Config.Appid))
	if found {
		phpsess = foo.(*http.Cookie)
	}
	return phpsess
}

type AppInfo struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Tel           string `json:"tel"`
	Addr          string `json:"addr"`
	Describe      string `json:"describle"`
	BusinessHours string `json:"business_hours"`
}

func SetAppInfoCache(appInfo *AppInfo) {
	CC.Set(fmt.Sprintf("APPINFO_%s", Config.Appid), appInfo, cache.DefaultExpiration)
}

func GetAppInfoCache() *AppInfo {
	if ai != nil {
		return ai
	}

	foo, found := CC.Get(fmt.Sprintf("APPINFO_%s", Config.Appid))
	if found {
		ai = foo.(*AppInfo)
	}

	return ai
}

func GetAppID() int64 {
	if ai == nil {
		return 0
	}
	return ai.Id
}

func GetAppName() string {
	if ai == nil {
		return ""
	}
	return ai.Name
}
