/*
Copyright 2020, Staffbase GmbH and contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version_test

import (
	"testing"

	"github.com/syseleven/syseleven-exporter/pkg/version"
)

// Build-info values are package-level vars set via -ldflags. These tests write
// them directly, so they share global state and do not run in parallel.

func setBuildInfo(t *testing.T, ver, rev, branch, user, date, goVer string) {
	t.Helper()

	version.Version = ver
	version.Revision = rev
	version.Branch = branch
	version.BuildUser = user
	version.BuildDate = date
	version.GoVersion = goVer
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name    string
		version string
		branch  string
		rev     string
		want    string
	}{
		{
			name:    "all fields populated",
			version: "1.2.3",
			branch:  "main",
			rev:     "abc123",
			want:    "(version=1.2.3, branch=main, revision=abc123)",
		},
		{
			name:    "empty fields render as empty",
			version: "",
			branch:  "",
			rev:     "",
			want:    "(version=, branch=, revision=)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setBuildInfo(t, tc.version, tc.rev, tc.branch, "irrelevant", "irrelevant", "irrelevant")

			if got := version.Info(); got != tc.want {
				t.Fatalf("Info()\n  got:  %q\n  want: %q", got, tc.want)
			}
		})
	}
}

func TestBuildContext(t *testing.T) {
	tests := []struct {
		name    string
		goVer   string
		user    string
		date    string
		want    string
	}{
		{
			name:  "all fields populated",
			goVer: "go1.24.3",
			user:  "ci",
			date:  "2026-07-01",
			want:  "(go=go1.24.3, user=ci, date=2026-07-01)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setBuildInfo(t, "irrelevant", "irrelevant", "irrelevant", tc.user, tc.date, tc.goVer)

			if got := version.BuildContext(); got != tc.want {
				t.Fatalf("BuildContext()\n  got:  %q\n  want: %q", got, tc.want)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	setBuildInfo(t, "1.2.3", "abc123", "main", "ci", "2026-07-01", "go1.24.3")

	want := "SysEleven Exporter, version 1.2.3 (branch: main, revision: abc123)\n" +
		"  build user:       ci\n" +
		"  build date:       2026-07-01\n" +
		"  go version:       go1.24.3"

	got, err := version.Print("SysEleven Exporter")
	if err != nil {
		t.Fatalf("Print() returned unexpected error: %v", err)
	}

	if got != want {
		t.Fatalf("Print()\n  got:  %q\n  want: %q", got, want)
	}
}
