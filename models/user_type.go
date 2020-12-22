package models

import (
	"encoding/json"
	"errors"
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

func UserTypeSync() error {

	appId := utils.GetAppInfoCache().Id
	appName := utils.GetAppInfoCache().Name

	if Sqlite == nil {
		logging.Err.Error("database is not init")
		return errors.New("database is not init")
	}

	query := "select ct_loc.loc_desc as loc_name, pa_patmas.pmi_name as name ,pa_adm.adm_in_pat_no,ct_hospital.hosp_desc as ct_hospital_name, pac_room.room_desc as pac_room_desc,pac_bed.bed_code as pac_bed_desc,dev_code ,dev_type,dev_active,dev_status  from cf_device "

	rows, err := Mysql.Raw(query).Rows()
	if err != nil {
		logging.Err.Error("mysql raw error :", err)
		return err
	}
	defer rows.Close()

	var userTypes []UserType
	for rows.Next() {
		var userType UserType
		Sqlite.ScanRows(rows, &userType)
		userTypes = append(userTypes, userType)
	}

	if len(userTypes) == 0 {
		return nil
	}

	var oldUserTypes []UserType
	Sqlite.Find(&oldUserTypes)

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
		Sqlite.Create(&newUserTypes)

		requestUserTypesJson, _ := json.Marshal(&requestUserTypes)
		var res interface{}
		res, err = utils.SyncServices(path, fmt.Sprintf("delUserTypeIds=%s&requestUserTypes=%s", "", requestUserTypesJson))
		if err != nil {
			logging.Err.Error("post common/v1/sync_remote get error ", err)
		}

		logging.Norm.Infof("数据提交返回信息:", res)

		return nil

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
		Sqlite.Where("dev_code in ?", delUserTypeIds).Delete(&UserType{})
	} else {
		delUserTypeIdsJson, _ = json.Marshal(&delUserTypeIds)
	}

	if len(newUserTypes) > 0 {
		Sqlite.Create(&newUserTypes)
	} else {
		requestUserTypesJson, _ = json.Marshal(&requestUserTypes)
	}

	var res interface{}
	res, err = utils.SyncServices(path, fmt.Sprintf("delUserTypeIds=%s&requestUserTypes=%s", string(delUserTypeIdsJson), string(requestUserTypesJson)))
	if err != nil {
		logging.Err.Error(err)
	}

	logging.Norm.Infof("数据提交返回信息:", res)

	return nil
}
