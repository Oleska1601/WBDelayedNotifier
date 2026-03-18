package api

import (
	"github.com/Oleska1601/WBDelayedNotifier/config"
	_ "github.com/Oleska1601/WBDelayedNotifier/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wb-go/wbf/ginext"
)

const (
	APIGroupURI = "/api"
)

type HTTPController interface {
	RegisterHandlers(*ginext.RouterGroup)
}

func Register(gin *config.GinConfig, controller HTTPController) *ginext.Engine {
	engine := ginext.New(gin.Mode)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	group := engine.Group(APIGroupURI)
	controller.RegisterHandlers(group)
	engine.Static("/static", "./web")
	return engine
}
