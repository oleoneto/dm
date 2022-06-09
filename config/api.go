package config

type APIConfig struct {
	/// CORS allowed origin
	AllowedHost string

	/// A connection string used to connect to the database (i.e. postgres://<user>:<password>@<host>:5432/database)
	ConnectionString string

	/// Enables server debug logs
	DebugMode bool

	/// The directory containing migration files
	Directory string

	/// The API resource namespace prefix (i.e. migrations)
	Namespace string

	/// The database table containing migration status information (i.e. schema_migrations)
	Table string

	/// The API version prefix (i.e. v1)
	Version string
}

func (t APIConfig) IsValid() ([]string, bool) {
	missing := []string{}

	if t.ConnectionString == "" {
		missing = append(missing, "DATABASE_URL")
	}

	if t.Directory == "" {
		missing = append(missing, "MIGRATIONS_DIRECTORY")
	}

	if t.Table == "" {
		missing = append(missing, "MIGRATIONS_TABLE")
	}

	return missing, len(missing) == 0
}
