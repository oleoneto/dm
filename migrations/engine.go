package migrations

type Engine interface {
	MigrationRunner
	Tracker
	Validator
}

type MigrationRunner interface {
	/**
	* @brief Runs migrations
	 */
	Up(changes MigrationList) error

	/**
	 * @brief Reverts migrations
	 */
	Down(changes MigrationList) error
}

type Tracker interface {
	/**
	 * @brief Prepares database for migration tracking
	 */
	StartTracking() error

	/**
	 * @brief Stops tracking database migrations
	 */
	StopTracking() error

	/**
	 * @brief Return the version of the last applied migration. The returned boolean should indicate if the database is being tracked
	 */
	Version() (string, bool)

	/**
	 * @brief Indicator of whether migrations are current or up-to-date
	 */
	IsUpToDate(changes MigrationList) bool

	// IsTracked - Indicator of whether the database is being managed by this tool
	IsTracked() bool


	// IsEmpty - Indicator of whether the database has any migrations
	IsEmpty() bool

	// AppliedMigrations - Returns all applied/saved migrations
	AppliedMigrations() map[string]Migration
}

type Validator interface {
	// Validate - Given a set of migrations, this method should return whether or not the migrations are valid.
	Validate(changes MigrationList) bool
}
