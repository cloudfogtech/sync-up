package data

import (
	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/cloudfogtech/sync-up/internal/utils"
)

func (db *DB) GetLog(serviceId, logId string) (Log, error) {
	var log Log
	result := db.orm.Where("service_id = ?", serviceId).Where("id = ?", logId).First(&log)
	return log, result.Error
}

func (db *DB) GetAllLogList(page int) ([]Log, int64, error) {
	return db.GetLogListByServiceId("", page)
}

func (db *DB) GetLogListByServiceId(serviceId string, page int) ([]Log, int64, error) {
	var logs []Log
	d := db.orm
	if serviceId != "" {
		d = d.Where("service_id = ?", serviceId)
	}
	result := d.Model(&Log{}).Offset((page - 1) * common.PageSize).
		Limit(common.PageSize).Order("start_at DESC").Find(&logs)
	return logs, db.GetLogCount(serviceId), result.Error
}

func (db *DB) GetLogCount(serviceId string) int64 {
	var count int64
	d := db.orm.Model(&Log{})
	if serviceId != "" {
		d = d.Where("service_id = ?", serviceId)
	}
	d.Count(&count)
	return count
}

func (db *DB) AddLog(request Log) (string, error) {
	request.ID = utils.UUID()
	result := db.orm.Create(&request)
	return request.ID, result.Error
}

func (db *DB) RemoveLog(id string) error {
	result := db.orm.Where("id = ?", id).Delete(&Log{})
	return result.Error
}
