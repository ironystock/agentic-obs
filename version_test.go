package main

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		commit     string
		date       string
		wantPrefix string
		wantSuffix string
	}{
		{
			name:       "dev version shows build info",
			version:    "dev",
			commit:     "none",
			date:       "unknown",
			wantPrefix: "dev",
			wantSuffix: "",
		},
		{
			name:       "release version shows only version",
			version:    "1.0.0",
			commit:     "abc1234",
			date:       "2025-01-15",
			wantPrefix: "1.0.0",
			wantSuffix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion := version
			origCommit := commit
			origDate := date
			defer func() {
				version = origVersion
				commit = origCommit
				date = origDate
			}()

			// Set test values
			version = tt.version
			commit = tt.commit
			date = tt.date

			got := Version()
			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf("Version() = %q, want prefix %q", got, tt.wantPrefix)
			}
		})
	}
}

func TestVersionShort(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "returns dev for development",
			version: "dev",
			want:    "dev",
		},
		{
			name:    "returns semantic version",
			version: "1.0.0",
			want:    "1.0.0",
		},
		{
			name:    "returns prerelease version",
			version: "1.0.0-beta.1",
			want:    "1.0.0-beta.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origVersion := version
			defer func() { version = origVersion }()

			version = tt.version
			if got := VersionShort(); got != tt.want {
				t.Errorf("VersionShort() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildInfo(t *testing.T) {
	// Save original values
	origVersion := version
	origCommit := commit
	origDate := date
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	// Set test values
	version = "1.2.3"
	commit = "abc1234"
	date = "2025-01-15T10:30:00Z"

	ver, com, dat := BuildInfo()

	if ver != "1.2.3" {
		t.Errorf("BuildInfo() version = %q, want %q", ver, "1.2.3")
	}
	if com != "abc1234" {
		t.Errorf("BuildInfo() commit = %q, want %q", com, "abc1234")
	}
	if dat != "2025-01-15T10:30:00Z" {
		t.Errorf("BuildInfo() date = %q, want %q", dat, "2025-01-15T10:30:00Z")
	}
}

func TestVersionDevFormat(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	version = "dev"
	commit = "abc1234"
	date = "2025-01-15"

	got := Version()

	// Dev version should include commit and date info
	if !strings.Contains(got, "dev") {
		t.Errorf("Version() = %q, should contain 'dev'", got)
	}
	if !strings.Contains(got, "abc1234") {
		t.Errorf("Version() = %q, should contain commit 'abc1234'", got)
	}
	if !strings.Contains(got, "2025-01-15") {
		t.Errorf("Version() = %q, should contain date '2025-01-15'", got)
	}
}
