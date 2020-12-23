package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/snowlyg/RemoteSync/logging"
	"github.com/snowlyg/RemoteSync/utils"
	"time"
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

func RemoteSync() error {
	appId := utils.GetAppInfoCache().Id
	appName := utils.GetAppInfoCache().Name

	if Sqlite == nil {
		logging.Err.Error("database is not init")
		return errors.New("database is not init")
	}

	query := "select ct_loc.loc_desc as loc_name, pa_patmas.pmi_name as name ,pa_adm.adm_in_pat_no,ct_hospital.hosp_desc as ct_hospital_name, pac_room.room_desc as pac_room_desc,pac_bed.bed_code as pac_bed_desc,dev_code,dev_desc,dev_type,dev_active,dev_status,dev_video_status  from cf_device "
	query += " left join pa_adm on pa_adm.pac_bed_id = cf_device.pac_bed_id"
	query += " left join ct_loc on ct_loc.loc_id = cf_device.ct_loc_id"
	query += " left join pa_patmas on pa_patmas.pmi_id = pa_adm.pa_patmas_id"
	query += " left join pac_room on pac_room.room_id = pa_adm.pac_room_id"
	query += " left join pac_bed on pac_bed.bed_id = cf_device.pac_bed_id"
	query += " left join ct_hospital on pa_adm.ct_hospital_id = ct_hospital.hosp_id"
	query += fmt.Sprintf(" where cf_device.dev_type = %s ", utils.Config.DevType)

	rows, err := Mysql.Raw(query).Rows()
	if err != nil {
		logging.Err.Error("mysql raw error :", err)
		return err
	}
	defer rows.Close()

	var remoteDevs []RemoteDev
	for rows.Next() {
		var remoteDev RemoteDev
		Sqlite.ScanRows(rows, &remoteDev)
		remoteDevs = append(remoteDevs, remoteDev)
	}

	if len(remoteDevs) == 0 {
		return nil
	}

	var oldRemoteDevs []RemoteDev
	Sqlite.Find(&oldRemoteDevs)

	var delDevCodes []string
	var newRemoteDevs []RemoteDev
	var requestRemoteDevs []*RequestRemoteDev

	// 没有旧数据
	path := "common/v1/data_sync/remote"
	if len(oldRemoteDevs) == 0 {
		newRemoteDevs = remoteDevs
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
				ApplicationId:   appId,
				ApplicationName: appName,
			}
			requestRemoteDevs = append(requestRemoteDevs, requestRemoteDev)
		}
		Sqlite.Create(&newRemoteDevs)

		requestRemoteDevsJson, _ := json.Marshal(&requestRemoteDevs)

		var res interface{}
		postdata := fmt.Sprintf("delDevCodes=%s&requestRemoteDevs=%s", "", string(requestRemoteDevsJson))
		logging.Dbug.Infof("data len %d", len(postdata))
		res, err = utils.SyncServices(path, postdata)
		if err != nil {
			logging.Err.Error("post common/v1/sync_remote get error ", err)
		}

		logging.Norm.Infof("数据提交返回信息:", res)

		return nil

	}

	// not in new
	for _, ore := range oldRemoteDevs {
		in := false
		for _, re := range remoteDevs {
			if ore.DevCode == re.DevCode {
				in = true
			}
		}
		if !in {
			delDevCodes = append(delDevCodes, ore.DevCode)
		}
	}

	// changed
	for _, re := range remoteDevs {
		in := false
		for _, ore := range oldRemoteDevs {
			if ore.DevCode == re.DevCode {
				if ore.Name != re.Name ||
					ore.AdmInPatNo != re.AdmInPatNo ||
					ore.DevDesc != re.DevDesc ||
					ore.PacRoomDesc != re.PacRoomDesc ||
					ore.PacBedDesc != re.PacBedDesc ||
					ore.CtHospitalName != re.CtHospitalName ||
					ore.DevVideoStatus != re.DevVideoStatus ||
					ore.DevStatus != re.DevStatus ||
					ore.DevActive != re.DevActive {
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
						ApplicationId:   appId,
						ApplicationName: appName,
					}
					requestRemoteDevs = append(requestRemoteDevs, requestRemoteDev)
					newRemoteDevs = append(newRemoteDevs, re)
					delDevCodes = append(delDevCodes, ore.DevCode)
				}
				in = true
			}
		}

		if !in {
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
				ApplicationId:   appId,
				ApplicationName: appName,
			}
			requestRemoteDevs = append(requestRemoteDevs, requestRemoteDev)
			newRemoteDevs = append(newRemoteDevs, re)
		}
	}

	var delDevCodesJson []byte
	var requestRemoteDevsJson []byte
	if len(delDevCodes) > 0 {
		Sqlite.Where("dev_code in ?", delDevCodes).Delete(&RemoteDev{})
		delDevCodesJson, _ = json.Marshal(&delDevCodes)
	}

	if len(newRemoteDevs) > 0 {
		Sqlite.Create(&newRemoteDevs)
	}
	if len(requestRemoteDevs) > 0 {
		requestRemoteDevsJson, _ = json.Marshal(&requestRemoteDevs)
	}

	postdata := fmt.Sprintf("delDevCodes=%s&requestRemoteDevs=%s", string(delDevCodesJson), string(requestRemoteDevsJson))
	logging.Dbug.Infof("data len %d", len(postdata))
	var res interface{}
	res, err = utils.SyncServices(path, postdata)
	if err != nil {
		logging.Err.Error(err)
	}

	logging.Norm.Infof("数据提交返回信息:", res)

	return nil
}
