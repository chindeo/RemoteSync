package utils

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

var CC *cache.Cache

func init() {
	CC = cache.New(1*time.Hour, 2*time.Hour)
}

func SetCacheToken(token string) {
	CC.Set(fmt.Sprintf("XToken_%s", Config.Appid), token, cache.DefaultExpiration)
}

func GetCacheToken() string {
	foo, found := CC.Get(fmt.Sprintf("XToken_%s", Config.Appid))
	if found {
		return foo.(string)
	}
	return ""
}
func SetSessionId(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			CC.Set(fmt.Sprintf("PHPSESSIONID_%s", Config.Appid), cookie, cache.DefaultExpiration)
		}
	}
}

func GetSessionId() *http.Cookie {
	foo, found := CC.Get(fmt.Sprintf("PHPSESSIONID_%s", Config.Appid))
	if found {
		return foo.(*http.Cookie)
	}
	return nil
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
	foo, found := CC.Get(fmt.Sprintf("APPINFO_%s", Config.Appid))
	if found {
		return foo.(*AppInfo)
	}
	return nil
}
