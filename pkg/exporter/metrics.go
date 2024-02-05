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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "syseleven"

var (
	computeCoresTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_cores_total",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	computeCoresUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_cores_used",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	computeInstancesTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_instances_total",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	computeInstancesUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_instances_used",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	computeFlavorsUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_flavors_used",
	}, []string{REGION_LABEL, PROJECT_LABEL, FLAVOR_LABEL})

	computeRamTotalMegabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_ram_total_megabytes",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	computeRamUsedMegabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_ram_used_megabytes",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	dnsZonesTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "dns_zones_total",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	dnsZonesUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "dns_zones_used",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	networkFloatingIPsTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_floating_ips_total",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	networkFloatingIPsUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_floating_ips_used",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	networkLoadbalancersTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_loadbalancers_total",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	networkLoadbalancersUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_loadbalancers_used",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	s3SpaceTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "s3_space_total_bytes",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	s3SpaceUsedBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "s3_space_used_bytes",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	volumeSpaceTotalGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_space_total_gigabytes",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	volumeSpaceUsedGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_space_used_gigabytes",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	volumeVolumesTotalGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_volumes_total",
	}, []string{REGION_LABEL, PROJECT_LABEL})

	volumeVolumesUsedGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_volumes_used",
	}, []string{REGION_LABEL, PROJECT_LABEL})
)
