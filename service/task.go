package service

import (
	"time"

	"github.com/blacksheepaul/timelog/model"
)

// CreateTask 创建任务
func CreateTask(task *model.Task) error {
	dao := model.GetDao()
	return model.CreateTask(dao.Db(), task)
}

// GetTaskByID 根据ID获取任务
func GetTaskByID(id uint) (*model.Task, error) {
	dao := model.GetDao()
	return model.GetTaskByID(dao.Db(), id)
}

// GetAllTasks 获取所有任务
func GetAllTasks() ([]model.Task, error) {
	dao := model.GetDao()
	return model.GetAllTasks(dao.Db())
}

// GetTasksByDate 根据日期获取任务
func GetTasksByDate(date time.Time) ([]model.Task, error) {
	dao := model.GetDao()
	return model.GetTasksByDate(dao.Db(), date)
}

// GetTasksByDateRange 根据日期范围获取任务
func GetTasksByDateRange(startDate, endDate time.Time) ([]model.Task, error) {
	dao := model.GetDao()
	return model.GetTasksByDateRange(dao.Db(), startDate, endDate)
}

// UpdateTask 更新任务
func UpdateTask(task *model.Task) error {
	dao := model.GetDao()
	return model.UpdateTask(dao.Db(), task)
}

// DeleteTask 删除任务
func DeleteTask(id uint) error {
	dao := model.GetDao()
	return model.DeleteTask(dao.Db(), id)
}

// MarkTaskAsCompleted 标记任务为完成
func MarkTaskAsCompleted(taskID uint) error {
	dao := model.GetDao()
	return model.MarkTaskAsCompleted(dao.Db(), taskID)
}

// MarkTaskAsIncomplete 标记任务为未完成  
func MarkTaskAsIncomplete(taskID uint) error {
	dao := model.GetDao()
	return model.MarkTaskAsIncomplete(dao.Db(), taskID)
}

// GetCompletedTasksInDateRange 获取指定日期范围内的已完成任务
func GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]model.Task, error) {
	dao := model.GetDao()
	return model.GetCompletedTasksInDateRange(dao.Db(), startDate, endDate)
}

// GetTaskStats 获取任务统计信息
func GetTaskStats(date time.Time) (map[string]interface{}, error) {
	dao := model.GetDao()
	return model.GetTaskStats(dao.Db(), date)
}

// CompleteTaskWithTimelog 完成任务并创建时间记录
// 这是一个组合操作，将任务标记为完成，并可选地创建关联的时间记录
func CompleteTaskWithTimelog(taskID uint, createTimelog bool, timelogData *model.TimeLog) error {
	dao := model.GetDao()
	
	// 开始事务
	tx := dao.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	// 标记任务为完成
	if err := model.MarkTaskAsCompleted(tx.Db(), taskID); err != nil {
		tx.Rollback()
		return err
	}
	
	// 如果需要创建时间记录
	if createTimelog && timelogData != nil {
		timelogData.TaskID = &taskID
		if err := model.CreateTimeLog(tx.Db(), timelogData); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit()
}