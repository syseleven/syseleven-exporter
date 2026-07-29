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

package exporter_test

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/syseleven/syseleven-exporter/pkg/api"
	"github.com/syseleven/syseleven-exporter/pkg/exporter"
)

// The exporter's GaugeVec metrics are package-level globals on the default
// registry, and every Set* call Reset()s and rewrites them. These tests share
// that state and do not run in parallel; the #114 refactor (injected registry)
// is what would let them use a per-test registry and parallelize.

func gatherEqual(t *testing.T, expected string, names ...string) {
	t.Helper()

	if err := testutil.GatherAndCompare(prometheus.DefaultGatherer, strings.NewReader(expected), names...); err != nil {
		t.Fatalf("metric exposition mismatch:\n%v", err)
	}
}

func TestSetUtilsV3_SetsGauges(t *testing.T) {
	exp := &exporter.Exporter{ProjectID: "project-a"}

	exporter.SetUtilsV3(
		map[string]api.QuotaV3{
			"region-a": {
				ComputeCores: 4,
				ComputeRAMMb: 2048,
				Objectstorage: []struct {
					SpaceBytes float64 `json:"space_bytes"`
					Type       string  `json:"type"`
				}{
					{SpaceBytes: 5000, Type: "s3"},
				},
			},
		},
		map[string]api.CurrentUsageV3{
			"region-a": {ComputeCores: 2},
		},
		exp,
	)

	expected := `
# HELP syseleven_compute_cores_total Quota for number of compute cores per region and project
# TYPE syseleven_compute_cores_total gauge
syseleven_compute_cores_total{project="project-a",region="region-a"} 4
# HELP syseleven_compute_cores_used Number of used compute cores per region and project
# TYPE syseleven_compute_cores_used gauge
syseleven_compute_cores_used{project="project-a",region="region-a"} 2
# HELP syseleven_compute_ram_total_megabytes Quota for ram per region and project in megabytes
# TYPE syseleven_compute_ram_total_megabytes gauge
syseleven_compute_ram_total_megabytes{project="project-a",region="region-a"} 2048
# HELP syseleven_s3_space_total_bytes Quota for S3 space per region and project in bytes
# TYPE syseleven_s3_space_total_bytes gauge
syseleven_s3_space_total_bytes{project="project-a",region="region-a",type="s3"} 5000
`

	gatherEqual(t, expected,
		"syseleven_compute_cores_total",
		"syseleven_compute_cores_used",
		"syseleven_compute_ram_total_megabytes",
		"syseleven_s3_space_total_bytes",
	)
}

// V1 reports S3 space under a hardcoded type="quobyte" label and reads load
// balancers from the NetworkLoadbalancers field. Both are easy to break silently.
func TestSetUtilsV1_QuobyteAndLoadbalancers(t *testing.T) {
	exp := &exporter.Exporter{ProjectID: "project-b"}

	exporter.SetUtilsV1(
		map[string]api.QuotaV1{
			"region-b": {
				NetworkLoadbalancers: 3,
				S3SpaceBytes:         9000,
			},
		},
		nil,
		exp,
	)

	expected := `
# HELP syseleven_network_loadbalancers_total Quota for number of load balancers per region and project
# TYPE syseleven_network_loadbalancers_total gauge
syseleven_network_loadbalancers_total{project="project-b",region="region-b"} 3
# HELP syseleven_s3_space_total_bytes Quota for S3 space per region and project in bytes
# TYPE syseleven_s3_space_total_bytes gauge
syseleven_s3_space_total_bytes{project="project-b",region="region-b",type="quobyte"} 9000
`

	gatherEqual(t, expected,
		"syseleven_network_loadbalancers_total",
		"syseleven_s3_space_total_bytes",
	)
}

// Current behavior on main: a nil quota+usage scrape still hits the unconditional
// Reset(), wiping all values. This is a smell (a transient API failure blanks the
// dashboards); update this assertion once #114 adds a guard to retain values.
func TestSetUtilsV3_NilScrapeClearsGauges(t *testing.T) {
	exp := &exporter.Exporter{ProjectID: "project-a"}

	exporter.SetUtilsV3(
		map[string]api.QuotaV3{"region-a": {ComputeCores: 4}},
		map[string]api.CurrentUsageV3{"region-a": {ComputeCores: 2}},
		exp,
	)

	exporter.SetUtilsV3(nil, nil, exp)

	// Empty expectation: after the nil scrape both gauges have no series.
	gatherEqual(t, "",
		"syseleven_compute_cores_total",
		"syseleven_compute_cores_used",
	)
}

// Enabled and CheckOnRaw map to gauge values 1 or 0.
func TestSetS3Info_EnabledFlags(t *testing.T) {
	exp := &exporter.Exporter{ProjectID: "project-c"}

	info := []api.S3UsageNCS{
		{
			S3UsersNCS: api.S3UsersNCS{Name: "u", Description: "d"},
			S3InfoNCS: api.S3InfoNCS{
				Size:       50,
				MaxSize:    200,
				NumObjects: 3,
				MaxObjects: 30,
				Enabled:    true,
				CheckOnRaw: false,
			},
		},
	}

	exporter.SetS3Info(info, exp)

	expected := `
# HELP syseleven_s3_enabled_ncs Checks if s3 space is enabled for user or not
# TYPE syseleven_s3_enabled_ncs gauge
syseleven_s3_enabled_ncs{description="d",project="project-c",s3username="u"} 1
# HELP syseleven_s3_check_enabled_ncs Checks if check on raw is enabled for user or not
# TYPE syseleven_s3_check_enabled_ncs gauge
syseleven_s3_check_enabled_ncs{description="d",project="project-c",s3username="u"} 0
`

	gatherEqual(t, expected,
		"syseleven_s3_enabled_ncs",
		"syseleven_s3_check_enabled_ncs",
	)
}
