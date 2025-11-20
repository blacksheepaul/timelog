package router

import (
	"net/http"
	"strconv"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// ApiResponse 统一API响应结构
type ApiResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"status"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}, message ...string) ApiResponse {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}
	return ApiResponse{
		Data:    data,
		Message: msg,
		Status:  http.StatusOK,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(status int, message string) ApiResponse {
	return ApiResponse{
		Data:    nil,
		Message: message,
		Status:  status,
	}
}

// RegisterTimeLogRoutes 注册 TimeLog 相关路由
func RegisterTimeLogRoutes(group *gin.RouterGroup) {
	group.POST("/timelogs", createTimeLogHandler)
	group.GET("/timelogs", listTimeLogsHandler)
	group.GET("/timelogs/:id", getTimeLogHandler)
	group.PUT("/timelogs/:id", updateTimeLogHandler)
	group.DELETE("/timelogs/:id", deleteTimeLogHandler)

	// Tag 相关路由
	group.GET("/tags", listTagsHandler)
	group.POST("/tags", createTagHandler)
	group.GET("/tags/:id", getTagHandler)
	group.PUT("/tags/:id", updateTagHandler)
	group.DELETE("/tags/:id", deleteTagHandler)
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
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	if err := service.CreateTimeLog(&tl); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整的Tag信息
	createdLog, err := service.GetTimeLogByID(tl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(createdLog, "Time log created successfully"))
}

// ListTimeLogsHandler godoc
// @Summary 查询时间日志列表
// @Description 获取所有时间日志
// @Tags timelog
// @Produce json
// @Param limit query int false "Limit number of results"
// @Param order query string false "Order by field (default: created_at DESC)"
// @Success 200 {array} model.TimeLog
// @Failure 500 {object} map[string]string
// @Router /api/timelogs [get]
func listTimeLogsHandler(c *gin.Context) {
	limitStr := c.Query("limit")
	orderBy := c.Query("order")

	var tls []model.TimeLog
	var err error

	if limitStr != "" || orderBy != "" {
		limit := 0
		if limitStr != "" {
			if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 {
				limit = l
			}
		}

		if orderBy == "" {
			orderBy = "created_at DESC"
		}

		tls, err = service.ListTimeLogsWithOptions(limit, orderBy)
	} else {
		tls, err = service.ListTimeLogs()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(tls, "Time logs retrieved successfully"))
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
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	tl, err := service.GetTimeLogByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(tl, "Time log retrieved successfully"))
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
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	var tl model.TimeLog
	if err := c.ShouldBindJSON(&tl); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	tl.ID = id
	if err := service.UpdateTimeLog(&tl); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整的Tag信息
	updatedLog, err := service.GetTimeLogByID(tl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(updatedLog, "Time log updated successfully"))
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
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	if err := service.DeleteTimeLog(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil, "Time log deleted successfully"))
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

// --- Tag Handlers ---

// listTagsHandler godoc
// @Summary 查询标签列表
// @Description 获取所有标签
// @Tags tag
// @Produce json
// @Success 200 {array} model.Tag
// @Failure 500 {object} map[string]string
// @Router /api/tags [get]
func listTagsHandler(c *gin.Context) {
	tags, err := service.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(tags, "Tags retrieved successfully"))
}

// createTagHandler godoc
// @Summary 创建标签
// @Description 新增一个标签
// @Tags tag
// @Accept json
// @Produce json
// @Param data body model.Tag true "标签数据"
// @Success 200 {object} model.Tag
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags [post]
func createTagHandler(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	if err := service.CreateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(tag, "Tag created successfully"))
}

// getTagHandler godoc
// @Summary 查询单个标签
// @Description 根据ID获取标签
// @Tags tag
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} model.Tag
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tags/{id} [get]
func getTagHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	tag, err := service.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(tag, "Tag retrieved successfully"))
}

// updateTagHandler godoc
// @Summary 更新标签
// @Description 根据ID更新标签
// @Tags tag
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param data body model.Tag true "标签数据"
// @Success 200 {object} model.Tag
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags/{id} [put]
func updateTagHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	tag.ID = id
	if err := service.UpdateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(tag, "Tag updated successfully"))
}

// deleteTagHandler godoc
// @Summary 删除标签
// @Description 根据ID删除标签
// @Tags tag
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags/{id} [delete]
func deleteTagHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	if err := service.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil, "Tag deleted successfully"))
}
