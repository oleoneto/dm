package cmd

type Version struct {
	Version string
}

// MARK: - Formattable
func (v Version) Description() string {
	return v.Version
}

var version = Version{
	Version: "1.1.0",
}
