package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/api/controllers"
	"github.com/cleopatrio/db-migrator-lib/api/middleware"
	"github.com/cleopatrio/db-migrator-lib/config"
	"github.com/gin-gonic/gin"
)

var (
	staticController     = controllers.StaticController{}
	migrationsController = controllers.MigrationsController{}
)

/*
Application routes:

GET 		/${API_VERSION}
GET 		/${API_VERSION}/docs
GET 		/${API_VERSION}/health
GET 		/${API_VERSION}/migrations
POST 		/${API_VERSION}/migrations
DELETE 	/${API_VERSION}/migrations
GET 		/${API_VERSION}/migrations/applied
GET 		/${API_VERSION}/migrations/pending
*/

func API(conf config.APIConfig) *gin.Engine {
	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	// CORS
	app.Use(middleware.CorsHeaders(conf.AllowedHost))

	// Static File Handling
	app.LoadHTMLGlob("./api/public/*.html")
	app.StaticFS("/static", http.Dir("./api/public/static"))
	app.StaticFile("/swagger.yaml", "./api/public/swagger.yaml")

	app.GET("/", staticController.Ping)

	versionedGroup := app.Group(fmt.Sprintf("/%v", sanitized(conf.Version)))
	versionedGroup.GET("/docs", staticController.Documentation)
	healthGroup := versionedGroup.Group("/health").Use(middleware.ConfigurationMiddleware(conf))
	namespacedGroup := versionedGroup.Group(fmt.Sprintf("/%v", sanitized(conf.Namespace))).Use(middleware.ConfigurationMiddleware(conf))

	{
		versionedGroup.GET("/", staticController.Ping)
		healthGroup.GET("", staticController.Health)
		namespacedGroup.GET("", migrationsController.List)
		namespacedGroup.POST("", migrationsController.Migrate)
		namespacedGroup.DELETE("", migrationsController.Rollback)
		namespacedGroup.GET("/applied", migrationsController.Applied)
		namespacedGroup.GET("/pending", migrationsController.Pending)
	}

	return app
}

func sanitized(value string) string {
	result := strings.TrimSpace(value)
	result = strings.TrimPrefix(result, "/")
	result = strings.TrimSuffix(result, "/")
	return result
}
