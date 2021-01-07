package models

import (
	"encoding/json"
	"fmt"
	"github.com/snowlyg/RemoteSync/logging"
	"github.com/snowlyg/RemoteSync/utils"
	"time"
)

type Loc struct {
	ID           uint   `gorm:"primarykey"`
	LocId        int64  `json:"loc_id"`
	LocDesc      string `json:"loc_desc"`
	LocWardFlag  int64  `json:"loc_ward_flag"`
	CtHospitalId int64  `json:"ct_hospital_id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RequestLoc struct {
	LocId           int64  `json:"loc_id"`
	LocDesc         string `json:"loc_desc"`
	LocWardFlag     int64  `json:"loc_ward_flag"`
	CtHospitalId    int64  `json:"ct_hospital_id"`
	ApplicationName string `json:"application_name"`
	ApplicationId   int64  `json:"application_id"`
}

func LocSync() {
	appId := utils.GetAppID()
	appName := utils.GetAppName()

	query := "select loc_id,loc_desc,loc_ward_flag,loc_active_flag,ct_hospital_id from ct_loc where loc_active_flag = 1"

	rows, err := GetMysql().Raw(query).Rows()
	if err != nil {
		logging.GetLocLogger().Error("mysql raw error :", err)
	}
	defer rows.Close()

	var locs []Loc
	for rows.Next() {
		var loc Loc
		GetMysql().ScanRows(rows, &loc)
		locs = append(locs, loc)
	}

	if len(locs) == 0 {
		return
	}

	//var oldLocs []Loc
	//GetSqlite().Find(&oldLocs)
	//
	//var delLocIds []int64
	//var newLocs []Loc
	var requestLocs []*RequestLoc

	// 没有旧数据
	path := "common/v1/data_sync/loc"
	//if len(oldLocs) == 0 {
	//	newLocs = locs
	for _, re := range locs {
		requestLoc := &RequestLoc{
			LocId:           re.LocId,
			LocDesc:         re.LocDesc,
			LocWardFlag:     re.LocWardFlag,
			CtHospitalId:    re.CtHospitalId,
			ApplicationId:   appId,
			ApplicationName: appName,
		}
		requestLocs = append(requestLocs, requestLoc)
	}
	//GetSqlite().Create(&newLocs)

	var requestLocsJson []byte
	requestLocsJson, err = json.Marshal(&requestLocs)

	if len(requestLocsJson) > 0 {
		var res interface{}
		res, err = utils.SyncServices(path, fmt.Sprintf("delLocIds=%s&requestLocs=%s", "", string(requestLocsJson)))
		if err != nil {
			logging.GetLocLogger().Error(err)
		}

		logging.GetLocLogger().Infof("科室数据同步提交返回信息:", res)
	}

	return

	//}
	//
	//// not in new
	//for _, ore := range oldLocs {
	//	in := false
	//	for _, re := range locs {
	//		if ore.LocId == re.LocId {
	//			in = true
	//		}
	//	}
	//	if !in {
	//		delLocIds = append(delLocIds, ore.LocId)
	//	}
	//}
	//
	//// changed
	//for _, re := range locs {
	//	in := false
	//	for _, ore := range oldLocs {
	//		if ore.LocId == re.LocId {
	//			if ore.LocWardFlag != re.LocWardFlag ||
	//				ore.LocDesc != re.LocDesc ||
	//				ore.CtHospitalId != re.CtHospitalId {
	//				requestLoc := &RequestLoc{
	//					LocId:           re.LocId,
	//					LocDesc:         re.LocDesc,
	//					LocWardFlag:     re.LocWardFlag,
	//					CtHospitalId:    re.CtHospitalId,
	//					ApplicationId:   appId,
	//					ApplicationName: appName,
	//				}
	//				requestLocs = append(requestLocs, requestLoc)
	//				newLocs = append(newLocs, re)
	//				delLocIds = append(delLocIds, ore.LocId)
	//			}
	//			in = true
	//		}
	//	}
	//
	//	if !in {
	//		requestLoc := &RequestLoc{
	//			LocId:           re.LocId,
	//			LocDesc:         re.LocDesc,
	//			LocWardFlag:     re.LocWardFlag,
	//			CtHospitalId:    re.CtHospitalId,
	//			ApplicationId:   appId,
	//			ApplicationName: appName,
	//		}
	//		requestLocs = append(requestLocs, requestLoc)
	//		newLocs = append(newLocs, re)
	//	}
	//}
	//
	//var delLocIdsJson []byte
	//var requestLocsJson []byte
	//if len(delLocIds) > 0 {
	//	GetSqlite().Where("loc_id in ?", delLocIds).Delete(&Loc{})
	//	delLocIdsJson, _ = json.Marshal(&delLocIds)
	//}
	//
	//if len(newLocs) > 0 {
	//	GetSqlite().Create(&newLocs)
	//}
	//
	//if len(requestLocs) > 0 {
	//	requestLocsJson, _ = json.Marshal(&requestLocs)
	//}
	//
	//if len(delLocIdsJson) > 0 || len(requestLocsJson) > 0 {
	//	var res interface{}
	//	res, err = utils.SyncServices(path, fmt.Sprintf("delLocIds=%s&requestLocs=%s", string(delLocIdsJson), string(requestLocsJson)))
	//	if err != nil {
	//		logging.GetLocLogger().Error(err)
	//	}
	//
	//	logging.GetLocLogger().Infof("科室数据同步提交返回信息:", res)
	//}
	//
	//return
}
