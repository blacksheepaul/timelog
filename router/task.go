package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// 添加任务相关路由
func setupTaskRoutes(group *gin.RouterGroup) {
	group.GET("/tasks", listTasksHandler)
	group.POST("/tasks", createTaskHandler)
	group.GET("/tasks/:id", getTaskHandler)
	group.PUT("/tasks/:id", updateTaskHandler)
	group.DELETE("/tasks/:id", deleteTaskHandler)
	group.POST("/tasks/:id/complete", completeTaskHandler)
	group.POST("/tasks/:id/incomplete", incompleteTaskHandler)
	group.GET("/tasks/stats/:date", getTaskStatsHandler)
}

// CreateTaskHandler godoc
// @Summary 创建任务
// @Description 新增一项任务
// @Tags task
// @Accept json
// @Produce json
// @Param data body model.Task true "任务数据"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks [post]
func createTaskHandler(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	
	if err := service.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	// 重新查询以获取完整的Tag信息
	if createdTask, err := service.GetTaskByID(task.ID); err == nil {
		c.JSON(http.StatusOK, SuccessResponse(createdTask, "Task created successfully"))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(task, "Task created successfully"))
	}
}

// ListTasksHandler godoc
// @Summary 获取任务列表
// @Description 获取所有任务，支持按日期过滤
// @Tags task
// @Produce json
// @Param date query string false "日期过滤 (YYYY-MM-DD格式)"
// @Success 200 {array} model.Task
// @Failure 500 {object} map[string]string
// @Router /api/tasks [get]
func listTasksHandler(c *gin.Context) {
	dateStr := c.Query("date")
	
	var tasks []model.Task
	var err error
	
	if dateStr != "" {
		// 解析日期
		if date, parseErr := time.Parse("2006-01-02", dateStr); parseErr != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD"))
			return
		} else {
			tasks, err = service.GetTasksByDate(date)
		}
	} else {
		tasks, err = service.GetAllTasks()
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse(tasks, "Tasks retrieved successfully"))
}

// GetTaskHandler godoc
// @Summary 获取单个任务
// @Description 根据ID获取任务详情
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [get]
func getTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	
	task, err := service.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse(task, "Task retrieved successfully"))
}

// UpdateTaskHandler godoc
// @Summary 更新任务
// @Description 更新任务信息
// @Tags task
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param data body model.Task true "任务数据"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [put]
func updateTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	
	// 先检查任务是否存在
	existingTask, err := service.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}
	
	var updateData model.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	
	// 保持ID不变
	updateData.ID = existingTask.ID
	updateData.CreatedAt = existingTask.CreatedAt
	
	if err := service.UpdateTask(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	// 重新查询以获取完整信息
	if updatedTask, err := service.GetTaskByID(uint(id)); err == nil {
		c.JSON(http.StatusOK, SuccessResponse(updatedTask, "Task updated successfully"))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(updateData, "Task updated successfully"))
	}
}

// DeleteTaskHandler godoc
// @Summary 删除任务
// @Description 删除指定任务
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [delete]
func deleteTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	
	// 先检查任务是否存在
	if _, err := service.GetTaskByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}
	
	if err := service.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse(nil, "Task deleted successfully"))
}

// CompleteTaskHandler godoc
// @Summary 标记任务为完成
// @Description 将任务标记为完成状态
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/complete [post]
func completeTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	
	if err := service.MarkTaskAsCompleted(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse(nil, "Task marked as completed"))
}

// IncompleteTaskHandler godoc
// @Summary 标记任务为未完成
// @Description 将任务标记为未完成状态
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/incomplete [post]
func incompleteTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	
	if err := service.MarkTaskAsIncomplete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse(nil, "Task marked as incomplete"))
}

// GetTaskStatsHandler godoc
// @Summary 获取任务统计
// @Description 获取指定日期的任务完成统计
// @Tags task
// @Produce json
// @Param date path string true "日期 (YYYY-MM-DD格式)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/stats/{date} [get]
func getTaskStatsHandler(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD"))
		return
	}
	
	stats, err := service.GetTaskStats(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse(stats, "Task stats retrieved successfully"))
}