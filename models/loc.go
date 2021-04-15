package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chindeo/RemoteSync/utils"
	"github.com/chindeo/pkg/logging"
	"github.com/chindeo/pkg/net"
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
	druation := utils.GetDuration(utils.Config.Timeduration, utils.Config.Timetype)
	logger := logging.GetMyLogger("loc")
	var locs []Loc
	var requestLocsJSON []byte
	go func() {
		for {
			loc(locs, requestLocsJSON, logger)
			time.Sleep(druation)
		}
	}()
}

func loc(locs []Loc, requestLocsJson []byte, logger *logging.Logger) {
	mysql, err := GetMysql()
	if err != nil {
		fmt.Println(fmt.Sprintf("get mysql error %+v", err))
		return
	}
	defer CloseMysql(mysql)
	// appId := utils.GetAppID()
	// appName := utils.GetAppName()

	query := "select loc_id,loc_desc,loc_ward_flag,loc_active_flag,ct_hospital_id from ct_loc where loc_active_flag = 1"

	rows, err := mysql.Raw(query).Rows()
	if err != nil {
		logger.Error("mysql raw error :", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var loc Loc
		mysql.ScanRows(rows, &loc)
		locs = append(locs, loc)
	}

	if len(locs) == 0 {
		return
	}
	appinfo, err := utils.GetAppInfo()
	if err != nil {
		fmt.Printf("GetAppInfo : %v", err)
		return
	}
	var requestLocs []*RequestLoc
	for _, re := range locs {
		requestLoc := &RequestLoc{
			LocId:           re.LocId,
			LocDesc:         re.LocDesc,
			LocWardFlag:     re.LocWardFlag,
			CtHospitalId:    re.CtHospitalId,
			ApplicationId:   appinfo.Id,
			ApplicationName: appinfo.Name,
		}
		requestLocs = append(requestLocs, requestLoc)
	}

	requestLocsJson, err = json.Marshal(&requestLocs)
	if len(requestLocsJson) > 0 {
		// var res interface{}
		// res, err = utils.SyncServices(path, fmt.Sprintf("delLocIds=%s&requestLocs=%s", "", string(requestLocsJson)))
		// if err != nil {
		// 	logger.Error(err)
		// }
		// 没有旧数据
		path := "/common/v1/data_sync/loc"
		serviceResponse := &net.ServerResponse{
			FullPath:     utils.Config.Host + path,
			Auth:         true,
			ResponseInfo: &net.ResponseInfo{},
		}
		_, err = net.NetClient.POSTNet(serviceResponse, fmt.Sprintf("delLocIds=%s&requestLocs=%s", "", string(requestLocsJson)))
		if err != nil {
			logger.Error(err)
		}
		logger.Infof("科室数据同步提交返回信息:", serviceResponse.ResponseInfo)
	}
	locs = nil
	requestLocsJson = nil
}
