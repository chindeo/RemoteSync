package models

import (
	"encoding/json"
	"fmt"
	"github.com/snowlyg/RemoteSync/logging"
	"github.com/snowlyg/RemoteSync/utils"
	"time"
)

type UserType struct {
	ID          uint   `gorm:"primarykey"`
	UtpId       int64  `json:"utp_id"`
	UtpCode     string `json:"utp_code"`
	UtpDesc     string `json:"utp_desc"`
	UtpType     string `json:"utp_type"`
	UtpActive   string `json:"utp_active"`
	UtpContrast string `json:"utp_contrast"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RequestUserType struct {
	UtpId           int64  `json:"utp_id"`
	UtpCode         string `json:"utp_code"`
	UtpDesc         string `json:"utp_desc"`
	UtpType         string `json:"utp_type"`
	UtpActive       string `json:"utp_active"`
	UtpContrast     string `json:"utp_contrast"`
	ApplicationName string `json:"application_name"`
	ApplicationId   int64  `json:"application_id"`
}

func UserTypeSync() {
	logger := logging.GetMyLogger("user_type")
	var userTypes []UserType
	var requestUserTypesJson []byte
	mysql, err := GetMysql()
	if err != nil {
		fmt.Println(fmt.Sprintf("get mysql error %+v", err))
		return
	}
	defer CloseMysql(mysql)
	appId := utils.GetAppID()
	appName := utils.GetAppName()

	query := "select utp_id,utp_code,utp_desc,utp_type,utp_active,utp_contrast from ct_user_type where utp_active = '1'"
	rows, err := mysql.Raw(query).Rows()
	if err != nil {
		logger.Error("mysql raw error :", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userType UserType
		mysql.ScanRows(rows, &userType)
		userTypes = append(userTypes, userType)
	}

	if len(userTypes) == 0 {
		return
	}

	var requestUserTypes []*RequestUserType

	// 没有旧数据
	path := "common/v1/data_sync/remote"

	for _, re := range userTypes {
		requestUserType := &RequestUserType{
			UtpId:           re.UtpId,
			UtpCode:         re.UtpCode,
			UtpDesc:         re.UtpDesc,
			UtpType:         re.UtpType,
			UtpActive:       re.UtpActive,
			UtpContrast:     re.UtpContrast,
			ApplicationId:   appId,
			ApplicationName: appName,
		}
		requestUserTypes = append(requestUserTypes, requestUserType)
	}

	requestUserTypesJson, err = json.Marshal(&requestUserTypes)
	if len(requestUserTypesJson) > 0 {
		var res interface{}
		res, err = utils.SyncServices(path, fmt.Sprintf("delUserTypeIds=%s&requestUserTypes=%s", "", string(requestUserTypesJson)))
		if err != nil {
			logger.Error(err)
		}
		logger.Infof("职位数据同步提交返回信息:", res)
	}

	return
}
