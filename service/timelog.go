package service

import (
	"github.com/blacksheepaul/timelog/model"
)

// --- TimeLog Service ---

// CreateTimeLog 新增一条时间日志
func CreateTimeLog(tl *model.TimeLog) error {
	db := model.GetDao().Db()
	return model.CreateTimeLog(db, tl)
}

// GetTimeLogByID 根据ID获取时间日志
func GetTimeLogByID(id uint) (*model.TimeLog, error) {
	db := model.GetDao().Db()
	return model.GetTimeLogByID(db, id)
}

// ListTimeLogs 查询时间日志（可扩展条件）
func ListTimeLogs(conds ...interface{}) ([]model.TimeLog, error) {
	db := model.GetDao().Db()
	return model.ListTimeLogs(db, conds...)
}

// UpdateTimeLog 更新一条时间日志
func UpdateTimeLog(tl *model.TimeLog) error {
	db := model.GetDao().Db()
	return model.UpdateTimeLog(db, tl)
}

// DeleteTimeLog 删除一条时间日志
func DeleteTimeLog(id uint) error {
	db := model.GetDao().Db()
	return model.DeleteTimeLog(db, id)
}
