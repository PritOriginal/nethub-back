package handler

import (
	"log/slog"

	"github.com/PritOriginal/nethub-back/docs"
	_ "github.com/PritOriginal/nethub-back/docs"
	"github.com/PritOriginal/nethub-back/internal/config"
	"github.com/PritOriginal/nethub-back/pkg/logger"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter(log *slog.Logger, env logger.Environment) *gin.Engine {
	r := gin.New()

	if env == logger.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	if env == logger.Local {
		r.Use(gin.Logger())
	} else {
		r.Use(sloggin.New(log))
	}

	r.Use(gin.Recovery())

	return r
}

func SetSwagger(r *gin.Engine, cfg *config.Config) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/"
}
