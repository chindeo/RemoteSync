package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chindeo/RemoteSync/utils"
	"github.com/chindeo/pkg/logging"
	"github.com/chindeo/pkg/net"
)

type RemoteDev struct {
	ID             uint   `gorm:"primarykey"`
	Name           string `json:"name"`             // 患者名称
	AdmInPatNo     string `json:"adm_in_pat_no"`    // 患者住院号
	DevCode        string `json:"dev_code"`         // 设备代码
	DevDesc        string `json:"dev_desc"`         // 设备名称
	DevType        int64  `json:"dev_type"`         // 设备类型 2 床旁
	PacRoomDesc    string `json:"pac_room_desc"`    // 房号
	PacBedDesc     string `json:"pac_bed_desc"`     // 床号
	CtHospitalName string `json:"ct_hospital_name"` // 院区
	LocName        string `json:"loc_name"`         // 科室
	DevStatus      string `json:"dev_status"`       // 设备状态
	DevActive      string `json:"dev_active"`       // 设备状态
	DevVideoStatus string `json:"dev_video_status"` // 探视状态
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RequestRemoteDev struct {
	Name            string `json:"name"`          // 患者名称
	AdmInPatNo      string `json:"adm_in_pat_no"` // 患者住院号
	DevCode         string `json:"dev_code"`      // 设备代码
	DevDesc         string `json:"dev_desc"`
	PacRoomDesc     string `json:"pac_room_desc"`    // 房号
	PacBedDesc      string `json:"pac_bed_desc"`     // 床号
	DevStatus       string `json:"dev_status"`       // 设备状态
	DevActive       string `json:"dev_active"`       // 设备状态
	DevVideoStatus  string `json:"dev_video_status"` // 探视状态
	CtHospitalName  string `json:"ct_hospital_name"` // 院区
	LocName         string `json:"loc_name"`         // 科室
	ApplicationName string `json:"application_name"`
	ApplicationId   int64  `json:"application_id"`
}

func RemoteSync() {
	druation := utils.GetDuration(utils.Config.Timeduration, utils.Config.Timetype)
	logger := logging.GetMyLogger("remote")
	var remoteDevs []RemoteDev
	var requestRemoteDevsJSON []byte
	go func() {
		for {
			remote(remoteDevs, requestRemoteDevsJSON, logger)
			time.Sleep(druation)
		}
	}()
}

func remote(remoteDevs []RemoteDev, requestRemoteDevsJson []byte, logger *logging.Logger) {
	mysql, err := GetMysql()
	if err != nil {
		fmt.Println(fmt.Sprintf("get mysql error %+v", err))
		return
	}
	defer CloseMysql(mysql)

	query := "select ct_loc.loc_desc as loc_name, pa_patmas.pmi_name as name ,pa_adm.adm_in_pat_no,ct_hospital.hosp_desc as ct_hospital_name, pac_room.room_desc as pac_room_desc,pac_bed.bed_code as pac_bed_desc,dev_code,dev_desc,dev_type,dev_active,dev_status,dev_video_status  from cf_device "
	query += " left join pa_adm on pa_adm.pac_bed_id = cf_device.pac_bed_id"
	query += " left join ct_loc on ct_loc.loc_id = cf_device.ct_loc_id"
	query += " left join pa_patmas on pa_patmas.pmi_id = pa_adm.pa_patmas_id"
	query += " left join pac_room on pac_room.room_id = pa_adm.pac_room_id"
	query += " left join pac_bed on pac_bed.bed_id = cf_device.pac_bed_id"
	query += " left join ct_hospital on pa_adm.ct_hospital_id = ct_hospital.hosp_id"
	query += fmt.Sprintf(" where cf_device.dev_type = %s ", utils.Config.DevType)

	rows, err := mysql.Raw(query).Rows()
	if err != nil {
		logger.Error("mysql raw error :", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var remoteDev RemoteDev
		mysql.ScanRows(rows, &remoteDev)
		remoteDevs = append(remoteDevs, remoteDev)
	}

	if len(remoteDevs) == 0 {
		return
	}
	appinfo, err := utils.GetAppInfo()
	if err != nil {
		fmt.Printf("GetAppInfo : %v", err)
		return
	}
	var requestRemoteDevs []*RequestRemoteDev
	for _, re := range remoteDevs {
		requestRemoteDev := &RequestRemoteDev{
			Name:            re.Name,
			AdmInPatNo:      re.AdmInPatNo,
			DevCode:         re.DevCode,
			DevDesc:         re.DevDesc,
			PacRoomDesc:     re.PacRoomDesc,
			PacBedDesc:      re.PacBedDesc,
			DevStatus:       re.DevStatus,
			DevActive:       re.DevActive,
			DevVideoStatus:  re.DevVideoStatus,
			CtHospitalName:  re.CtHospitalName,
			LocName:         re.LocName,
			ApplicationId:   appinfo.Id,
			ApplicationName: appinfo.Name,
		}
		requestRemoteDevs = append(requestRemoteDevs, requestRemoteDev)
	}

	requestRemoteDevsJson, err = json.Marshal(&requestRemoteDevs)
	if len(requestRemoteDevsJson) > 0 {
		// 没有旧数据
		path := "/common/v1/data_sync/remote"
		serviceResponse := &net.ServerResponse{
			FullPath:     utils.Config.Host + path,
			Auth:         true,
			ResponseInfo: &net.ResponseInfo{},
		}
		_, err = net.NetClient.POSTNet(serviceResponse, fmt.Sprintf("&requestRemoteDevs=%s", string(requestRemoteDevsJson)))
		if err != nil {
			logger.Error(err)
		}
		logger.Infof("探视数据同步提交返回信息:", serviceResponse.ResponseInfo)
	}

	remoteDevs = nil
	requestRemoteDevsJson = nil

	return
}
