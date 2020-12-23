package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type getToken struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *Token `json:"data"`
}

type Req struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	AccessToken string `json:"AccessToken"`
}

func SyncServices(path, data string) (interface{}, error) {
	var re Req
	result := Request("POST", path, data, true)
	if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("SyncServices 同步数据请求没有返回数据"))
	}
	err := json.Unmarshal(result, &re)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SyncServices dopost: %s json.Unmarshal error：%v ,with result: %v", path, err, string(result)))
	}
	return re, nil
}

func GetToken() error {
	token := GetCacheToken()
	if token != "" {
		return nil
	}

	var re getToken
	result := Request(
		"POST",
		"/api/v1/get_access_token",
		fmt.Sprintf("app_id=%s&app_secret=%s", Config.Appid, Config.Appsecret),
		false,
	)
	if len(result) == 0 {
		return errors.New("请求没有返回数据")
	}

	err := json.Unmarshal(result, &re)
	if err != nil {
		return err
	}

	if re.Code == 200 {
		if re.Data.AccessToken == "" {
			return errors.New(fmt.Sprintf("get token return response %+v", re))
		}
		fmt.Println(fmt.Sprintf("get token return response %+v", re.Data))
		SetCacheToken(re.Data.AccessToken)
		if err := GetAppInfo(); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("get token return response %+v", re))
	}
}

type AppInfoRequest struct {
	Code    int64    `json:"code"`
	Message string   `json:"message"`
	Data    *AppInfo `json:"data"`
}

func GetAppInfo() error {
	appInfo := GetAppInfoCache()
	if appInfo != nil {
		return nil
	}

	var air AppInfoRequest
	result := Request("GET", "/api/v1/application", "", true)
	if len(result) == 0 {
		return errors.New("请求没有返回数据")
	}

	err := json.Unmarshal(result, &air)
	if err != nil {
		return err
	}

	if air.Code == 200 {
		if air.Data == nil {
			return errors.New(fmt.Sprintf("get appinfo return response %+v", air))
		}
		fmt.Println(fmt.Sprintf("get appinfo return response %+v", air.Data))
		SetAppInfoCache(air.Data)
		return nil
	} else {
		return errors.New(fmt.Sprintf("get appinfo return response %+v", air))
	}
}

func Request(method, url, data string, auth bool) []byte {
	timeout := 3
	timeover := 3
	T := time.Tick(time.Duration(timeover) * time.Second)
	var result = make(chan []byte, 10)
	t := time.Duration(timeout) * time.Second
	Client := http.Client{Timeout: t}
	go func() {
		fullUrl := fmt.Sprintf("%s/%s", Config.Host, url)
		if strings.Contains(url, "http") || strings.Contains(url, "https") {
			fullUrl = url
		}
		req, _ := http.NewRequest(method, fullUrl, strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		req.Header.Set("AuthType", "4")
		if auth {
			req.Header.Set("X-Token", GetCacheToken())
			phpSessionId := GetSessionId()
			if phpSessionId != nil {
				req.AddCookie(phpSessionId)
			}
		}
		resp, err := Client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("%s: %+v", url, err))
			return
		}

		if !auth {
			SetSessionId(resp.Cookies())
		}

		b, _ := ioutil.ReadAll(resp.Body)
		result <- b
	}()

	for {
		select {
		case x := <-result:
			return x
		case <-T:
			return nil
		}
	}

}
