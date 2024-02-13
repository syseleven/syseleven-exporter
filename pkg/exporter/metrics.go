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
		Help:      "Quota for number of compute cores per region and project",
	}, []string{"region", "project"})

	computeCoresUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_cores_used",
		Help:      "Number of used compute cores per region and project",
	}, []string{"region", "project"})

	computeInstancesTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_instances_total",
		Help:      "Quota for number of compute instances per region and project",
	}, []string{"region", "project"})

	computeInstancesUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_instances_used",
		Help:      "Number of used compute instances per region and project",
	}, []string{"region", "project"})

	computeFlavorsUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_flavors_used",
		Help:      "Number of used compute flavors per type and region and project",
	}, []string{"region", "project", "flavor"})

	computeRamTotalMegabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_ram_total_megabytes",
		Help:      "Quota for ram per region and project in megabytes",
	}, []string{"region", "project"})

	computeRamUsedMegabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "compute_ram_used_megabytes",
		Help:      "Used ram per region and project in megabytes",
	}, []string{"region", "project"})

	dnsZonesTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "dns_zones_total",
		Help:      "Quota for number of DNS zones per region and project",
	}, []string{"region", "project"})

	dnsZonesUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "dns_zones_used",
		Help:      "Number of used DNS zones per region and project",
	}, []string{"region", "project"})

	networkFloatingIPsTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_floating_ips_total",
		Help:      "Quota for number of floating IPs per region and project",
	}, []string{"region", "project"})

	networkFloatingIPsUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_floating_ips_used",
		Help:      "Number of used floating IPs per region and project",
	}, []string{"region", "project"})

	networkLoadbalancersTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_loadbalancers_total",
		Help:      "Quota for number of load balancers per region and project",
	}, []string{"region", "project"})

	networkLoadbalancersUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "network_loadbalancers_used",
		Help:      "Number of used load balancers per region and project",
	}, []string{"region", "project"})

	s3SpaceTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "s3_space_total_bytes",
		Help:      "Quota for S3 space per region and project in bytes",
	}, []string{"region", "project"})

	s3SpaceUsedBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "s3_space_used_bytes",
		Help:      "Used S3 space per region and project in bytes",
	}, []string{"region", "project"})

	volumeSpaceTotalGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_space_total_gigabytes",
		Help:      "Quota for volume space per region and project in gigabytes",
	}, []string{"region", "project"})

	volumeSpaceUsedGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_space_used_gigabytes",
		Help:      "Number of used volume space per region and project in gigabytes",
	}, []string{"region", "project"})

	volumeVolumesTotalGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_volumes_total",
		Help:      "Quota for number of volumes per region and project",
	}, []string{"region", "project"})

	volumeVolumesUsedGigabytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "volume_volumes_used",
		Help:      "Number of used volumes per region and project",
	}, []string{"region", "project"})
)
