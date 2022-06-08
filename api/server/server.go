package server

import (
	"github.com/cleopatrio/db-migrator-lib/api/controllers"
	"github.com/cleopatrio/db-migrator-lib/api/middleware"
	"github.com/gin-gonic/gin"
)

var (
	staticController     = controllers.StaticController{}
	migrationsController = controllers.MigrationsController{}
)

func API() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	app.GET("/", staticController.Ping)

	hg := app.Group("/health")
	hg.Use(middleware.ConfigurationMiddleware())
	{
		hg.GET("", staticController.Health)
	}

	mg := app.Group("/migrations")
	mg.Use(middleware.ConfigurationMiddleware())

	{
		mg.GET("", migrationsController.List)
		mg.GET("/applied", migrationsController.Applied)
		mg.GET("/pending", migrationsController.Pending)
		mg.POST("/migrate", migrationsController.Migrate)
		mg.POST("/rollback", migrationsController.Rollback)
	}

	return app
}
