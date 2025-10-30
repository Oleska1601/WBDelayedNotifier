package controller

import (
	_ "github.com/Oleska1601/WBDelayedNotifier/docs"

	"github.com/wb-go/wbf/ginext"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRouter() {
	ginMode := ""
	engine := ginext.New(ginMode)

	notifyGroup := engine.Group("/notify")
	{
		notifyGroup.GET("/:notification_id", s.GetNotificationStatusHandler)
		notifyGroup.POST("", s.CreateNotificationHandler)
		notifyGroup.DELETE("/:notification_id", s.DeleteNotificationHandler)
	}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Srv.Handler = engine
}
