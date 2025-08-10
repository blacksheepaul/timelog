package router

import (
	"net/http"
	"strconv"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// RegisterTimeLogRoutes 注册 TimeLog 相关路由
func RegisterTimeLogRoutes(group *gin.RouterGroup) {
	group.POST("/timelogs", createTimeLogHandler)
	group.GET("/timelogs", listTimeLogsHandler)
	group.GET("/timelogs/:id", getTimeLogHandler)
	group.PUT("/timelogs/:id", updateTimeLogHandler)
	group.DELETE("/timelogs/:id", deleteTimeLogHandler)
}

// CreateTimeLogHandler godoc
// @Summary 创建时间日志
// @Description 新增一条时间日志
// @Tags timelog
// @Accept json
// @Produce json
// @Param data body model.TimeLog true "时间日志数据"
// @Success 200 {object} model.TimeLog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/timelogs [post]
func createTimeLogHandler(c *gin.Context) {
	var tl model.TimeLog
	if err := c.ShouldBindJSON(&tl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreateTimeLog(&tl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tl)
}

// ListTimeLogsHandler godoc
// @Summary 查询时间日志列表
// @Description 获取所有时间日志
// @Tags timelog
// @Produce json
// @Success 200 {array} model.TimeLog
// @Failure 500 {object} map[string]string
// @Router /api/timelogs [get]
func listTimeLogsHandler(c *gin.Context) {
	tls, err := service.ListTimeLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tls)
}

// GetTimeLogHandler godoc
// @Summary 查询单条时间日志
// @Description 根据ID获取时间日志
// @Tags timelog
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} model.TimeLog
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/timelogs/{id} [get]
func getTimeLogHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tl, err := service.GetTimeLogByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tl)
}

// UpdateTimeLogHandler godoc
// @Summary 更新时间日志
// @Description 根据ID更新时间日志
// @Tags timelog
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Param data body model.TimeLog true "时间日志数据"
// @Success 200 {object} model.TimeLog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/timelogs/{id} [put]
func updateTimeLogHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tl model.TimeLog
	if err := c.ShouldBindJSON(&tl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tl.ID = id
	if err := service.UpdateTimeLog(&tl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tl)
}

// DeleteTimeLogHandler godoc
// @Summary 删除时间日志
// @Description 根据ID删除时间日志
// @Tags timelog
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/timelogs/{id} [delete]
func deleteTimeLogHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteTimeLog(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

// parseUintParam 辅助函数
func parseUintParam(c *gin.Context, key string, out *uint) error {
	idStr := c.Param(key)
	var id64 uint64
	var err error
	if id64, err = parseUint(idStr); err != nil {
		return err
	}
	*out = uint(id64)
	return nil
}

func parseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
