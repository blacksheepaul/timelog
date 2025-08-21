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

// ListTimeLogsWithOptions 查询时间日志（支持排序和限制）
func ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]model.TimeLog, error) {
	db := model.GetDao().Db()
	return model.ListTimeLogsWithOptions(db, limit, orderBy, conds...)
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

// --- Tag Service ---

// CreateTag 新增一个标签
func CreateTag(tag *model.Tag) error {
	db := model.GetDao().Db()
	return model.CreateTag(db, tag)
}

// GetTagByID 根据ID获取标签
func GetTagByID(id uint) (*model.Tag, error) {
	db := model.GetDao().Db()
	return model.GetTagByID(db, id)
}

// GetTagByName 根据名称获取标签
func GetTagByName(name string) (*model.Tag, error) {
	db := model.GetDao().Db()
	return model.GetTagByName(db, name)
}

// ListTags 查询标签列表
func ListTags(conds ...interface{}) ([]model.Tag, error) {
	db := model.GetDao().Db()
	return model.ListTags(db, conds...)
}

// UpdateTag 更新一个标签
func UpdateTag(tag *model.Tag) error {
	db := model.GetDao().Db()
	return model.UpdateTag(db, tag)
}

// DeleteTag 删除一个标签
func DeleteTag(id uint) error {
	db := model.GetDao().Db()
	return model.DeleteTag(db, id)
}
