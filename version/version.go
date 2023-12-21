package version

// version is set during build using -ldflags
var version string

// Version returns tagged release version if set, or unknown otherwise
func Version() string {
	if version == "" {
		return "unknown"
	}
	return version
}
