# Database Migrator, DM
A database migration tool.

[![Build and Test](https://github.com/cleopatrio/db-migrator-lib/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/cleopatrio/db-migrator-lib/actions/workflows/go.yml)

**Table of Contents**
- [Commands](#commands)
  - [dm](#dm)
  - [Migrate](#migrate)
  - [Rollback](#rollback)
  - [Validate](#validate)
  - [Show](#show)
- [To Do](#to-do)

## Commands
Assume the installed binary is called `dm`.

### dm
```
Usage:
  dm [flags]
  dm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  migrate     Run migration(s)
  rollback    Rollback migration(s)
  show        Shows the state of applied and pending migrations
  validate    Validate configuration of migration files

Flags:
      --config string         config file
      --database-url string   database url (default "postgres://***:****@***:5432/****")
      --directory string      migrations directory (default "./migrations")
      --engine string         database engine (default "postgresql")
  -h, --help                  help for dm
      --table string          table wherein migrations are tracked (default "_migrations")

Use "dm [command] --help" for more information about a command.
```

---

### Migrate
```
Run migration(s)

Usage:
  dm migrate [flags]

Flags:
  -h, --help             help for migrate
      --version string   run migrations up do this version
```
Since both the version and the name of a migration are validated for uniqueness, the `version` flag can take either value. So, both `20220422101345` and `create_user` are valid arguments for the flag.

Note that this and other commands can load migrations from anywhere in your file system. Just point the `--directory` flag to where your files are.

---

### Rollback
```
Rollback migration(s)

Usage:
  dm rollback [flags]

Flags:
  -h, --help             help for rollback
      --version string   rollback this version (and everything applied after it)
```

---

### Validate
```
Validate the configuration of migration files

Usage:
  dm validate [flags]

Flags:
  -h, --help   help for validate
```

The validator checks for duplicate timestamps in the file name, duplicate/mismatched file and/or schema names and a few other things. Take a look at [validations.go](migrations/validations.go) for a better understanding of what the validator takesi into account at this point.

---

### Show
```
Shows the state of applied and pending migrations

Usage:
  dm show [flags]
  dm show [command]

Available Commands:
  all         List all migrations for a given application
  applied     List only applied migrations
  pending     List only pending migrations
  version     Shows the most recently applied migration
```

## To Do
[Check out open issues](https://github.com/cleopatrio/issues).
