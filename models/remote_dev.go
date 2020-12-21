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
	ID           uint   `gorm:"primarykey"`
	Name         string `json:"name"`           // 患者名称
	AdmInPatNo   string `json:"adm_in_pat_no"`  // 患者住院号
	DevCode      string `json:"dev_code"`       // 设备代码
	DevType      int64  `json:"dev_type"`       // 设备类型 2 床旁
	PacRoomDesc  string `json:"pac_room_desc"`  // 房号
	PacBedDesc   string `json:"pac_bed_desc"`   // 床号
	CtHospitalId string `json:"ct_hospital_id"` // 医院id
	DevStatus    string `json:"dev_status"`     // 设备状态
	DevAction    string `json:"dev_active"`     // 设备状态
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RequestRemoteDev struct {
	Name        string `json:"name"`          // 患者名称
	AdmInPatNo  string `json:"adm_in_pat_no"` // 患者住院号
	DevCode     string `json:"dev_code"`      // 设备代码
	PacRoomDesc string `json:"pac_room_desc"` // 房号
	PacBedDesc  string `json:"pac_bed_desc"`  // 床号
	DevStatus   string `json:"dev_status"`    // 设备状态
	DevAction   string `json:"dev_active"`    // 设备状态
}

func Sync() error {
	if err := utils.GetToken(); err != nil {
		return err
	}
	if err := utils.GetAppInfo(); err != nil {
		return err
	}

	if Sqlite == nil {
		return errors.New("database is not init")
	}

	query := "select pa_patmas.pmi_name as name ,pa_adm.adm_in_pat_no,pa_adm.ct_hospital_id, pac_room.room_desc as pac_room_desc,pac_bed.bed_code as pac_bed_desc,dev_code ,dev_type,dev_active,dev_status  from cf_device "
	query += " left join pa_adm on pa_adm.pac_bed_id = cf_device.pac_bed_id"
	query += " left join pa_patmas on pa_patmas.pmi_id = pa_adm.pa_patmas_id"
	query += " left join pac_room on pac_room.room_id = cf_device.pac_room_id"
	query += " left join pac_bed on pac_bed.bed_id = cf_device.pac_bed_id"
	query += fmt.Sprintf(" where cf_device.dev_type = 2 and pa_adm.ct_hospital_id = %d", utils.GetAppInfoCache().CtHospitalId)

	rows, err := Mysql.Raw(query).Rows()
	if err != nil {
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
	if len(oldRemoteDevs) == 0 {
		newRemoteDevs = remoteDevs
		for _, re := range remoteDevs {
			requestRemoteDev := &RequestRemoteDev{
				Name:        re.Name,
				AdmInPatNo:  re.AdmInPatNo,
				DevCode:     re.DevCode,
				PacRoomDesc: re.PacRoomDesc,
				PacBedDesc:  re.PacBedDesc,
				DevStatus:   re.DevStatus,
				DevAction:   re.DevAction,
			}
			requestRemoteDevs = append(requestRemoteDevs, requestRemoteDev)
		}
		Sqlite.Create(&newRemoteDevs)

		requestRemoteDevsJson, _ := json.Marshal(&requestRemoteDevs)
		var res interface{}
		res, err = utils.SyncServices("api/v1/sync_remote", fmt.Sprintf("delDevCodes=%s&requestRemoteDevs=%s", "", requestRemoteDevsJson))
		if err != nil {
			logging.Err.Error(err)
		}

		logging.Norm.Infof("数据提交返回信息:%+v", res)

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
					ore.PacRoomDesc != re.PacRoomDesc ||
					ore.PacBedDesc != re.PacBedDesc ||
					ore.CtHospitalId != re.CtHospitalId ||
					ore.DevStatus != re.DevStatus ||
					ore.DevAction != re.DevAction {
					requestRemoteDev := &RequestRemoteDev{
						Name:        re.Name,
						AdmInPatNo:  re.AdmInPatNo,
						DevCode:     re.DevCode,
						PacRoomDesc: re.PacRoomDesc,
						PacBedDesc:  re.PacBedDesc,
						DevStatus:   re.DevStatus,
						DevAction:   re.DevAction,
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
				Name:        re.Name,
				AdmInPatNo:  re.AdmInPatNo,
				DevCode:     re.DevCode,
				PacRoomDesc: re.PacRoomDesc,
				PacBedDesc:  re.PacBedDesc,
				DevStatus:   re.DevStatus,
				DevAction:   re.DevAction,
			}
			requestRemoteDevs = append(requestRemoteDevs, requestRemoteDev)
			newRemoteDevs = append(newRemoteDevs, re)
		}
	}

	var delDevCodesJson []byte
	var requestRemoteDevsJson []byte
	if len(delDevCodes) > 0 {
		Sqlite.Where("dev_code in ?", delDevCodes).Delete(&RemoteDev{})
	} else {
		delDevCodesJson, _ = json.Marshal(&delDevCodes)
	}

	if len(newRemoteDevs) > 0 {
		Sqlite.Create(&newRemoteDevs)
	} else {
		requestRemoteDevsJson, _ = json.Marshal(&requestRemoteDevs)
	}

	var res interface{}
	res, err = utils.SyncServices("platform/report/syncdevice", fmt.Sprintf("delDevCodes=%s&requestRemoteDevs=%s", string(delDevCodesJson), string(requestRemoteDevsJson)))
	if err != nil {
		logging.Err.Error(err)
	}

	logging.Norm.Infof("数据提交返回信息:%+v", res)

	return nil
}
