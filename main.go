package main

import (
	_ "github.com/joho/godotenv/autoload"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oleoneto/dm/cmd/cli"
	"github.com/oleoneto/dm/pkg/logger"
)

func main() {
	logger.NewLogger()

	cli.Execute()
}
