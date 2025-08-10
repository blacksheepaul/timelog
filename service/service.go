package service

import (
	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
)

var log logger.Logger
var cfg *config.Config

func InitService(loggerInstance logger.Logger, config *config.Config) {
	log = loggerInstance
	cfg = config
}

type Response struct {
	Items []any `json:"items"`
	Pages
}

type Pages struct {
	Page  int `form:"page" json:"page"`
	Size  int `form:"size" json:"size"`
	Total int `form:"total" json:"total"`
}
