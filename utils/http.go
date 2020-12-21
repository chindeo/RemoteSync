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
	XToken string `json:"X-Token"`
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
		fmt.Sprintf("appid=%s&appsecret=%s&apptype=%s", Config.Appid, Config.Appsecret, "hospital"),
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
		SetCacheToken(re.Data.XToken)
		return nil
	} else {
		return errors.New(re.Message)
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
	result := Request("GET", "/api/v1/application", "", false)
	if len(result) == 0 {
		return errors.New("请求没有返回数据")
	}

	err := json.Unmarshal(result, &air)
	if err != nil {
		return err
	}

	if air.Code == 200 {
		SetAppInfoCache(air.Data)
		return nil
	} else {
		return errors.New(air.Message)
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
		fullUrl := fmt.Sprintf("http://%s/%s", Config.Host, url)
		if strings.Contains(url, "http") {
			fullUrl = url
		}
		req, _ := http.NewRequest(method, fullUrl, strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		if auth {
			req.Header.Set("X-Token", GetCacheToken())
		}
		resp, err := Client.Do(req)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s: %+v", url, err))
			return
		}
		defer resp.Body.Close()
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
