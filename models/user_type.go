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
	appId := utils.GetAppID()
	appName := utils.GetAppName()

	query := "select utp_id,utp_code,utp_desc,utp_type,utp_active,utp_contrast from ct_user_type where utp_active = '1'"
	rows, err := GetMysql().Raw(query).Rows()
	if err != nil {
		logging.GetUserTypeLogger().Error("mysql raw error :", err)
	}
	defer rows.Close()

	var userTypes []UserType
	for rows.Next() {
		var userType UserType
		GetSqlite().ScanRows(rows, &userType)
		userTypes = append(userTypes, userType)
	}

	if len(userTypes) == 0 {
		return
	}

	var oldUserTypes []UserType
	GetSqlite().Find(&oldUserTypes)
	var delUserTypeIds []int64
	var newUserTypes []UserType
	var requestUserTypes []*RequestUserType

	// 没有旧数据
	path := "common/v1/data_sync/remote"
	if len(oldUserTypes) == 0 {
		newUserTypes = userTypes
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
		GetSqlite().Create(&newUserTypes)
		requestUserTypesJson, _ := json.Marshal(&requestUserTypes)
		if len(requestUserTypesJson) > 0 {
			var res interface{}
			res, err = utils.SyncServices(path, fmt.Sprintf("delUserTypeIds=%s&requestUserTypes=%s", "", string(requestUserTypesJson)))
			if err != nil {
				logging.GetUserTypeLogger().Error(err)
			}
			logging.GetUserTypeLogger().Infof("职位数据同步提交返回信息:", res)
		}

		return

	}

	// not in new
	for _, ore := range oldUserTypes {
		in := false
		for _, re := range userTypes {
			if ore.UtpId == re.UtpId {
				in = true
			}
		}
		if !in {
			delUserTypeIds = append(delUserTypeIds, ore.UtpId)
		}
	}

	// changed
	for _, re := range userTypes {
		in := false
		for _, ore := range oldUserTypes {
			if ore.UtpId == re.UtpId {
				if ore.UtpCode != re.UtpCode ||
					ore.UtpDesc != re.UtpDesc ||
					ore.UtpType != re.UtpType ||
					ore.UtpActive != re.UtpActive ||
					ore.UtpContrast != re.UtpContrast {
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
					newUserTypes = append(newUserTypes, re)
					delUserTypeIds = append(delUserTypeIds, ore.UtpId)
				}
				in = true
			}
		}

		if !in {
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
			newUserTypes = append(newUserTypes, re)
		}
	}

	var delUserTypeIdsJson []byte
	var requestUserTypesJson []byte
	if len(delUserTypeIds) > 0 {
		GetSqlite().Where("dev_code in ?", delUserTypeIds).Delete(&UserType{})
		delUserTypeIdsJson, _ = json.Marshal(&delUserTypeIds)
	}

	if len(newUserTypes) > 0 {
		GetSqlite().Create(&newUserTypes)

	}
	if len(requestUserTypes) > 0 {
		requestUserTypesJson, _ = json.Marshal(&requestUserTypes)
	}

	if len(requestUserTypesJson) > 0 || len(delUserTypeIdsJson) > 0 {
		var res interface{}
		res, err = utils.SyncServices(path, fmt.Sprintf("delUserTypeIds=%s&requestUserTypes=%s", string(delUserTypeIdsJson), string(requestUserTypesJson)))
		if err != nil {
			logging.GetUserTypeLogger().Error(err)
		}

		logging.GetUserTypeLogger().Infof("职位数据同步提交返回信息:", res)
	}

	return
}
