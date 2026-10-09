package exporter

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/syseleven/syseleven-exporter/pkg/api"
)

func TestSetS3Info_LabelsMetricsByTarget(t *testing.T) {
	exporter := &Exporter{ProjectID: "project-id"}
	usage := []api.S3UsageNCS{
		{
			S3UsersNCS: api.S3UsersNCS{Name: "s3-user", Description: "test"},
			S3InfoNCS:  api.S3InfoNCS{MaxSize: 10, Size: 2, NumObjects: 3, MaxObjectsUser: 4, MaxObjectsBucket: 5, Enabled: true, CheckOnRaw: true},
			Target:     "target-a",
		},
		{
			S3UsersNCS: api.S3UsersNCS{Name: "s3-user", Description: "test"},
			S3InfoNCS:  api.S3InfoNCS{MaxSize: 20, Size: 6, NumObjects: 7, MaxObjectsUser: 8, MaxObjectsBucket: 9},
			Target:     "target-b",
		},
	}

	SetS3Info(usage, exporter)

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	metricNames := []string{
		"syseleven_s3_space_max_bytes_ncs",
		"syseleven_s3_space_used_bytes_ncs",
		"syseleven_s3_enabled_ncs",
		"syseleven_s3_check_enabled_ncs",
		"syseleven_s3_num_objects_ncs",
		"syseleven_s3_max_objects_ncs",
		"syseleven_s3_max_objects_bucket_ncs",
	}
	for _, metricName := range metricNames {
		family := findMetricFamily(families, metricName)
		if family == nil {
			t.Fatalf("metric family %q not found", metricName)
		}
		targets := map[string]bool{}
		for _, metric := range family.Metric {
			for _, label := range metric.Label {
				if label.GetName() == "target" {
					targets[label.GetValue()] = true
				}
			}
		}
		if len(targets) != 2 || !targets["target-a"] || !targets["target-b"] {
			t.Fatalf("metric %q has targets %v, want target-a and target-b", metricName, targets)
		}
	}
}

func findMetricFamily(families []*dto.MetricFamily, name string) *dto.MetricFamily {
	for _, family := range families {
		if family.GetName() == name {
			return family
		}
	}
	return nil
}
