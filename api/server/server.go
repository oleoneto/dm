package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
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

	//go:embed public/*
	assets embed.FS

	//go:embed templates/*
	templates embed.FS
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
	loadHTMLFromFS(app, templates, "templates/*")
	app.StaticFS("static", assetsFS(assets))

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

// Source: https://github.com/gin-gonic/gin/issues/2795
func loadHTMLFromFS(engine *gin.Engine, embedFS embed.FS, pattern string) {
	templ := template.Must(template.ParseFS(embedFS, pattern))
	engine.SetHTMLTemplate(templ)
}

// Source: https://dev.to/fareez/embedding-a-react-application-in-go-binary-188n
func assetsFS(embed embed.FS) http.FileSystem {
	fsystem, err := fs.Sub(embed, "public")

	if err != nil {
		panic(err)
	}

	return http.FS(fsystem)
}
