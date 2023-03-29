package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vesicash/payment-ms/internal/config"
)

func SetupMetrics(appConfiguration *config.App) *gin.Engine {
	if appConfiguration.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
