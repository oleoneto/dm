package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

type BadRequestError struct {
	Message string
}

func ConfigurationMiddleware() gin.HandlerFunc {
	/*
		- A connection string used to connect to the database (i.e. postgres://<user>:<password>@<host>:5432/database)
		DATABASE_URL string

		- The directory containing migration files
		MIGRATIONS_DIRECTORY string

		- The database table containing migration status information (i.e. schema_migrations)
		MIGRATIONS_TABLE string
	*/

	return func(ctx *gin.Context) {
		ConfigurationValidatorMiddleware(ctx)
		ArgumentsSetterMiddleware(ctx)

		ctx.Next()
	}
}

func ConfigurationValidatorMiddleware(ctx *gin.Context) {
	missingRequiredFlags := os.Getenv("DATABASE_URL") == "" || os.Getenv("MIGRATIONS_DIRECTORY") == "" || os.Getenv("MIGRATIONS_TABLE") == ""

	if missingRequiredFlags {
		ctx.AbortWithStatusJSON(500, BadRequestError{Message: "Missing configuration variables"})
	}
}

func ArgumentsSetterMiddleware(ctx *gin.Context) {
	flags := []string{
		"--output-format", "json",
		"--directory", os.Getenv("MIGRATIONS_DIRECTORY"),
		"--table", os.Getenv("MIGRATIONS_TABLE"),
	}

	ctx.Set("command_flags", flags)
}
