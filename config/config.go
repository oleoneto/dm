package config

type DMConfig struct {
	/// A connection string used to connect to the database (i.e. postgres://<user>:<password>@<host>:5432/database)
	ConnectionString string

	/// The directory containing migration files
	Directory string

	/// The database table containing migration status information (i.e. schema_migrations)
	Table string
}

func (t DMConfig) IsValid() ([]string, bool) {
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
