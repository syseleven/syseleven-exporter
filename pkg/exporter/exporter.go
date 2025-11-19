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

package exporter

import (
	"time"

	"github.com/syseleven/syseleven-exporter/pkg/api"
	"github.com/syseleven/syseleven-exporter/pkg/auth"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	Username    string
	Password    string
	ProjectID   string
	UseAppCreds bool
}

func New(projectID string, useAppCreds bool, username string, password string) (*Exporter, error) {
	return &Exporter{
		ProjectID:   projectID,
		UseAppCreds: useAppCreds,
		Username:    username,
		Password:    password,
	}, nil
}

// set utils for API v3 for Quota and Usage Information

func SetUtilsV3(quota map[string]api.QuotaV3, usage map[string]api.CurrentUsageV3, exporter *Exporter) {
	computeCoresTotal.Reset()
	computeInstancesTotal.Reset()
	computeRamTotalMegabytes.Reset()
	dnsZonesTotal.Reset()
	networkFloatingIPsTotal.Reset()
	networkLoadbalancersTotal.Reset()
	s3SpaceTotalBytes.Reset()
	volumeSpaceTotalGigabytes.Reset()
	volumeVolumesTotalGigabytes.Reset()
	computeCoresUsed.Reset()
	computeInstancesUsed.Reset()
	computeRamUsedMegabytes.Reset()
	dnsZonesUsed.Reset()
	networkFloatingIPsUsed.Reset()
	networkLoadbalancersUsed.Reset()
	s3SpaceUsedBytes.Reset()
	volumeSpaceUsedGigabytes.Reset()
	volumeVolumesUsedGigabytes.Reset()
	computeFlavorsUsed.Reset()

	for k, v := range quota {
		computeCoresTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeCores)
		computeInstancesTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeInstances)
		computeRamTotalMegabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeRAMMb)
		dnsZonesTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.DNSZones)
		networkFloatingIPsTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.NetworkFloatingips)
		networkLoadbalancersTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.LoadbalancerLoadbalancers)
		volumeSpaceTotalGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeSpaceGb)
		volumeVolumesTotalGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeVolumes)

		for _, os := range v.Objectstorage {
			s3SpaceTotalBytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID, "type": os.Type}).Set(os.SpaceBytes)
		}

	}

	for k, v := range usage {
		computeCoresUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeCores)
		computeInstancesUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeInstances)
		computeRamUsedMegabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeRAMMb)
		dnsZonesUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.DNSZones)
		networkFloatingIPsUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.NetworkFloatingips)
		networkLoadbalancersUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.LoadbalancerLoadbalancers)
		volumeSpaceUsedGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeSpaceGb)
		volumeVolumesUsedGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeVolumes)

		for flavor := range v.ComputeFlavors {
			computeFlavorsUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID, "flavor": flavor}).Set(v.ComputeFlavors[flavor])
		}

		for _, os := range v.Objectstorage {
			s3SpaceUsedBytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID, "type": os.Type}).Set(os.SpaceBytes)
		}
	}

}

// set utils for API v1 for Quota and Usage Information

func SetUtilsV1(quota map[string]api.QuotaV1, usage map[string]api.CurrentUsageV1, exporter *Exporter) {
	computeCoresTotal.Reset()
	computeInstancesTotal.Reset()
	computeRamTotalMegabytes.Reset()
	dnsZonesTotal.Reset()
	networkFloatingIPsTotal.Reset()
	networkLoadbalancersTotal.Reset()
	s3SpaceTotalBytes.Reset()
	volumeSpaceTotalGigabytes.Reset()
	volumeVolumesTotalGigabytes.Reset()
	computeCoresUsed.Reset()
	computeInstancesUsed.Reset()
	computeRamUsedMegabytes.Reset()
	dnsZonesUsed.Reset()
	networkFloatingIPsUsed.Reset()
	networkLoadbalancersUsed.Reset()
	s3SpaceUsedBytes.Reset()
	volumeSpaceUsedGigabytes.Reset()
	volumeVolumesUsedGigabytes.Reset()
	computeFlavorsUsed.Reset()

	for k, v := range quota {
		computeCoresTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeCores)
		computeInstancesTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeInstances)
		computeRamTotalMegabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeRAMMb)
		dnsZonesTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.DNSZones)
		networkFloatingIPsTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.NetworkFloatingips)
		networkLoadbalancersTotal.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.NetworkLoadbalancers)
		s3SpaceTotalBytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID, "type": "quobyte"}).Set(v.S3SpaceBytes)
		volumeSpaceTotalGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeSpaceGb)
		volumeVolumesTotalGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeVolumes)
	}

	for k, v := range usage {
		computeCoresUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeCores)
		computeInstancesUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeInstances)
		computeRamUsedMegabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.ComputeRAMMb)
		dnsZonesUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.DNSZones)
		networkFloatingIPsUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.NetworkFloatingips)
		networkLoadbalancersUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.NetworkLoadbalancers)
		s3SpaceUsedBytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID, "type": "quobyte"}).Set(v.S3SpaceBytes)
		volumeSpaceUsedGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeSpaceGb)
		volumeVolumesUsedGigabytes.With(prometheus.Labels{"region": k, "project": exporter.ProjectID}).Set(v.VolumeVolumes)

		for flavor := range v.ComputeFlavors {
			computeFlavorsUsed.With(prometheus.Labels{"region": k, "project": exporter.ProjectID, "flavor": flavor}).Set(v.ComputeFlavors[flavor])
		}
	}
}

func SetS3Info(s3infoncs []api.S3UsageNCS, exporter *Exporter) {
	s3SpaceMaxBytesNcs.Reset()
	s3SpaceUsedBytesNcs.Reset()
	s3EnabledNcs.Reset()
	s3CheckEnabledNcs.Reset()
	s3NumObjectsNcs.Reset()
	s3MaxObjectsNcs.Reset()
	for _, v := range s3infoncs {
		s3SpaceMaxBytesNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(v.MaxSize)
		s3SpaceUsedBytesNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(v.Size)
		s3NumObjectsNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(v.NumObjects)
		s3MaxObjectsNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(v.MaxObjects)
		if v.Enabled {
			s3EnabledNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(1)
		} else {
			s3EnabledNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(0)
		}
		if v.CheckOnRaw {
			s3CheckEnabledNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(1)
		} else {
			s3CheckEnabledNcs.With(prometheus.Labels{"project": exporter.ProjectID, "s3username": v.Name, "description": v.Description}).Set(0)
		}
	}
}

func setQuotaAndUsage(apiVersion string, exporter *Exporter) {
	var token string
	var err error

	log.Infof("Scrape Quota and Usage Metrics")
	if !(exporter.UseAppCreds) {
		token, err = auth.GetToken(exporter.ProjectID, exporter.Username, exporter.Password)
	} else {
		token, err = auth.GetTokenAppCreds(exporter.ProjectID, exporter.Username, exporter.Password)
	}
	if err != nil {
		log.WithError(err).Error("Could not get API Token")
	}

	switch apiVersion {
	case "v3":
		quota, err := api.GetQuotaV3(exporter.ProjectID, token)
		if err != nil {
			log.WithError(err).Error("Could not get quota")
		}

		usage, err := api.GetCurrentUsageV3(exporter.ProjectID, token)
		if err != nil {
			log.WithError(err).Error("Could not get current usage")
		}

		SetUtilsV3(quota, usage, exporter)

	default:
		quota, err := api.GetQuotaV1(exporter.ProjectID, token)
		if err != nil {
			log.WithError(err).Error("Could not get quota")
		}

		usage, err := api.GetCurrentUsageV1(exporter.ProjectID, token)
		if err != nil {
			log.WithError(err).Error("Could not get current usage")
		}

		SetUtilsV1(quota, usage, exporter)
	}
}

func SetS3StatsNCS(exporter *Exporter) {
	log.Infof("Fetching S3 info from NCS")
	s3infoncs, err := api.GetS3InfoNCS(exporter.ProjectID)
	if err != nil {
		log.WithError(err).Error("Could not get current usage")
	}
	SetS3Info(s3infoncs, exporter)
}
func Run(interval int64, apiVersion string, s3StatsNCS bool, exporter *Exporter) {
	for {
		go setQuotaAndUsage(apiVersion, exporter)
		if s3StatsNCS {
			go SetS3StatsNCS(exporter)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
