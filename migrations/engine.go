package migrations

type Engine interface {
	/**
	 * @brief Prepares database for migration tracking
	 */
	StartTracking() error

	/**
	 * @brief Stops tracking database migrations
	 */
	StopTracking() error

	/**
	 * @brief Runs migrations
	 */
	Up(changes Migrations) error

	/**
	 * @brief Reverts migrations
	 */
	Down(changes Migrations) error

	/**
	 * @brief Return the version of the last applied migration. The returned boolean should indicate if the database is being tracked
	 */
	Version() (string, bool)

	/**
	 * @brief Indicator of whether migrations are current or up-to-date
	 */
	IsUpToDate(changes Migrations) bool

	/**
	 * @brief Indicator of whether the database is being managed by this tool
	 */
	IsTracked() bool

	/**
	 * @brief Indicator of whether the database has any migrations
	 */
	IsEmpty() bool
}
