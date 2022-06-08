package server

import (
	"fmt"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/api/controllers"
	"github.com/cleopatrio/db-migrator-lib/api/middleware"
	"github.com/gin-gonic/gin"
)

var (
	staticController     = controllers.StaticController{}
	migrationsController = controllers.MigrationsController{}
)

/*
Application routes:

GET /${API_VERSION}
GET /${API_VERSION}/health
GET /${API_VERSION}/${API_NAMESPACE}
GET /${API_VERSION}/${API_NAMESPACE}/applied
GET /${API_VERSION}/${API_NAMESPACE}/pending

POST /${API_VERSION}/${API_NAMESPACE}/migrate
POST /${API_VERSION}/${API_NAMESPACE}/rollback
*/

func API(version, namespace string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	app.GET("/", staticController.Ping)

	versionedGroup := app.Group(fmt.Sprintf("/%v", sanitized(version)))
	healthGroup := versionedGroup.Group("/health").Use(middleware.ConfigurationMiddleware())
	namespacedGroup := versionedGroup.Group(fmt.Sprintf("/%v", sanitized(namespace))).Use(middleware.ConfigurationMiddleware())

	{
		versionedGroup.GET("/", staticController.Ping)
		healthGroup.GET("", staticController.Health)
		namespacedGroup.GET("", migrationsController.List)
		namespacedGroup.GET("/applied", migrationsController.Applied)
		namespacedGroup.GET("/pending", migrationsController.Pending)
		namespacedGroup.POST("/migrate", migrationsController.Migrate)
		namespacedGroup.POST("/rollback", migrationsController.Rollback)
	}

	return app
}

func sanitized(value string) string {
	result := strings.TrimSpace(value)
	result = strings.TrimPrefix(result, "/")
	result = strings.TrimSuffix(result, "/")
	return result
}
