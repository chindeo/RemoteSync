package utils

import (
	"encoding/json"
	"fmt"

	"github.com/chindeo/pkg/net"
)

type AppInfo struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Tel           string `json:"tel"`
	Addr          string `json:"addr"`
	Describe      string `json:"describle"`
	BusinessHours string `json:"business_hours"`
}

type ResponseAppInfo struct {
	net.ResponseInfo
	Data *AppInfo `json:"data"`
}

func GetAppInfo() (*AppInfo, error) {
	serviceResponse := &net.ServerResponse{
		FullPath:     Config.Host + "/api/v1/profile/hospital",
		Auth:         true,
		ResponseInfo: &net.ResponseInfo{},
	}
	result, err := net.NetClient.GetNet(serviceResponse)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(result))

	responseAppInfo := &ResponseAppInfo{}
	err = json.Unmarshal(result, responseAppInfo)
	if err != nil {
		return nil, err
	}

	if responseAppInfo.Data == nil {
		return nil, fmt.Errorf("未获取到数据")
	}

	return responseAppInfo.Data, nil
}
