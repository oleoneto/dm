# Database Migrator, DM
A database migration tool.

[![Build and Test](https://github.com/cleopatrio/db-migrator-lib/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/cleopatrio/db-migrator-lib/actions/workflows/go.yml)

**Table of Contents**
- [Database Migrator, DM](#database-migrator-dm)
  - [Commands](#commands)
    - [dm](#dm)
    - [Migrate](#migrate)
    - [Rollback](#rollback)
    - [Generate](#generate)
    - [Validate](#validate)
    - [Show](#show)
    - [API](#api)
  - [To Do](#to-do)

## Commands
Assume the installed binary is called `dm`.

### dm
```
DM, short for Database Migrator is a migration management tool.

Usage:
  dm [flags]
  dm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  migrate     Run migration(s)
  rollback    Rollback migration(s)
  show        Shows the state of applied and pending migrations
  validate    Validate the configuration of migration files

Flags:
  -a, --adapter string     database adapter (default "postgresql")
      --config string      config file
  -d, --directory string   migrations directory (default "./migrations")
  -h, --help               help for dm
  -t, --table string       table wherein migrations are tracked (default "_migrations")

Use "dm [command] --help" for more information about a command.
```

---

### Migrate
```
Run migration(s)

Usage:
  dm migrate NAME|VERSION [flags]

Aliases:
  migrate, m

Flags:
  -u, --database-url string   database url (default "postgres://****:**@**:5432/******")
  -h, --help                  help for migrate
```
When an argument for NAME|VERSION is provided, the command will execute this migration and everything that comes before it. This is done to ensure schema consistency.

Since both the version and the name of a migration are validated for uniqueness, the command can take a single argument for either value. In other words, both `20220422101345` and `create_user` are valid arguments.

Note that this and other commands can load migrations from anywhere in your file system. Just point the `--directory` flag to where your files are.

---

### Rollback
```
Rollback migration(s)

Usage:
  dm rollback NAME|VERSION [flags]

Aliases:
  rollback, r

Flags:
  -u, --database-url string   database url (default "postgres://****:**@**:5432/******")
  -h, --help             help for rollback
```

When an argument for NAME|VERSION is provided, the command will remove this migration and everything applied after it
This is done to ensure schema consistency.

Since both the version and the name of a migration are validated for uniqueness, the command can take a single argument for either value. In other words, both `20220422101345` and `create_user` are valid arguments.

---

### Generate
```
Generate a database migration file

Usage:
  dm generate NAME [flags]

Aliases:
  generate, g

Flags:
  -h, --help   help for generate

Global Flags:
  -d, --directory string   migrations directory (default "./migrations")
```

If the provided migration name passes validation, this command will create a migration file and save it in the migrations directory.
The file will be created using the schema in use by the running version of the CLI. Check the [examples directory](examples) for examples schemas.

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

### API
Beginning in version v2.0.0, the CLI now features a server that exposes some of its functionality as RESTful endpoints. 

```
Run a RESTful API

Usage:
  dm api [flags]

Flags:
  -u, --database-url string   database url (default "postgres://****:**@**:5432/******")
      --debug                 shows server debug output
  -h, --help                  help for api
  -n, --namespace string      api resource namespace param (default "migrations")
  -p, --port int              server port (default 3809)
  -v, --version string        api version param (default "v1")
```

The API server requires you to set three variables: `DATABASE_URL`, `MIGRATIONS_DIRECTORY`, and `MIGRATIONS_TABLE`. These can be set via environment variables or by setting their respective flags in the server executable. The default server port is `3809`. You can also specify both an `API_VERSION` and an `API_NAMESPACE` to configure the API endpoints.

**Endpoints**
```
GET /${API_VERSION}
GET /${API_VERSION}/health
GET /${API_VERSION}/${API_NAMESPACE}
GET /${API_VERSION}/${API_NAMESPACE}/applied
GET /${API_VERSION}/${API_NAMESPACE}/pending

POST /${API_VERSION}/${API_NAMESPACE}/migrate
POST /${API_VERSION}/${API_NAMESPACE}/rollback
```


## To Do
[Check out open issues](https://github.com/cleopatrio/db-migrator-lib/issues).
