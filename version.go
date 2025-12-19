package main

// Build-time version injection using ldflags.
// Example: go build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

var (
	// version is the semantic version (set via ldflags)
	version = "dev"

	// commit is the git commit hash (set via ldflags)
	commit = "none"

	// date is the build date in RFC3339 format (set via ldflags)
	date = "unknown"
)

// Version returns the full version string including commit and date.
func Version() string {
	if commit != "none" && len(commit) > 7 {
		return version + " (" + commit[:7] + ", " + date + ")"
	}
	return version
}

// VersionShort returns just the semantic version.
func VersionShort() string {
	return version
}

// BuildInfo returns the individual build components.
func BuildInfo() (ver, com, dat string) {
	return version, commit, date
}
