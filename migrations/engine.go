package migrations

type Engine interface {
	/**
	 * @brief Runs migrations
	 */
	Up(changes []Migration) error

	/**
	 * @brief Reverts migrations
	 */
	Down(changes []Migration) error

	/**
	 * @brief Indicator of whether migrations are current or up-to-date
	 */
	IsUpToDate() bool
}
