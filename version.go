package main

// Build-time variables injected via ldflags.
// These are set during build with:
//
//	go build -ldflags "-X main.version=1.0.0 -X main.commit=abc1234 -X main.date=2025-01-01"
//
// If not set, defaults are used for development builds.
var (
	// version is the semantic version (e.g., "1.0.0", "1.0.0-beta.1")
	version = "dev"

	// commit is the git commit hash (short form)
	commit = "none"

	// date is the build date in ISO 8601 format
	date = "unknown"
)

// Version returns the full version string including build metadata.
func Version() string {
	if version == "dev" {
		return "dev (commit: " + commit + ", built: " + date + ")"
	}
	return version
}

// VersionShort returns just the semantic version.
func VersionShort() string {
	return version
}

// BuildInfo returns detailed build information.
func BuildInfo() (ver, com, dat string) {
	return version, commit, date
}
