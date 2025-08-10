package router

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/router/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/blacksheepaul/timelog/docs"
)

var GinLogger = gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
	return fmt.Sprintf(`{"time":"%s","client":"%s","method":"%s","path":"%s","latency":"%s","status":%d,"emsg":"%s"}`+"\n",
		p.TimeStamp.Format(time.DateTime),
		p.ClientIP,
		p.Method,
		p.Path,
		p.Latency,
		p.StatusCode,
		p.ErrorMessage,
	)
})
var log logger.Logger

func Register(r *gin.Engine, cfg *config.Config, l logger.Logger) *gin.Engine {
	log = l

	r.Use(GinLogger)
	r.Use(middleware.Cors(cfg))

	api := r.Group("/api")   // api
	auth := api.Group("/v1") // api/v1
	auth.Use(middleware.Auth())

	// 注册 TimeLog 路由
	RegisterTimeLogRoutes(api)

	// 注册 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func LaunchServer(ctx context.Context, wg *sync.WaitGroup, r *gin.Engine, cfg *config.Config) {
	defer wg.Done()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed", err)
		}
	}()

	<-ctx.Done()
	log.Info("Server received stop signal, shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server shutdown failed", err)
	}
	log.Info("Server exited gracefully.")
}
