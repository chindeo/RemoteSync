package utils

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var CC *cache.Cache

func init() {
	CC = cache.New(1*time.Hour, 2*time.Hour)
}

func SetCacheToken(token string) {
	CC.Set("XToken", token, cache.DefaultExpiration)
}

func GetCacheToken() string {
	foo, found := CC.Get("XToken")
	if found {
		return foo.(string)
	}
	return ""
}

type AppInfo struct {
	CtHospitalId int64 `json:"ct_hospital_id"`
}

const APPINFPINDEX = "APPINFO"

func SetAppInfoCache(appInfo *AppInfo) {
	CC.Set(APPINFPINDEX, appInfo, cache.DefaultExpiration)
}

func GetAppInfoCache() *AppInfo {
	foo, found := CC.Get(APPINFPINDEX)
	if found {
		return foo.(*AppInfo)
	}
	return nil
}
