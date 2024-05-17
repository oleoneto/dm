package migrator

import "context"

type MigrationGeneratorProtocol func(ctx context.Context, format, content, name, directory string) Migration
