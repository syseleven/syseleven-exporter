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

package api

type QuotaV3 struct {
	ComputeCores                   float64 `json:"compute.cores"`
	ComputeInstances               float64 `json:"compute.instances"`
	ComputeKeyPairs                float64 `json:"compute.key_pairs"`
	ComputeMetadataItems           float64 `json:"compute.metadata_items"`
	ComputeRAMMb                   float64 `json:"compute.ram_mb"`
	ComputeServerGroupMembers      float64 `json:"compute.server_group_members"`
	ComputeServerGroups            float64 `json:"compute.server_groups"`
	DNSZones                       float64 `json:"dns.zones"`
	LoadbalancerHealthmonitors     float64 `json:"loadbalancer.healthmonitors"`
	LoadbalancerListeners          float64 `json:"loadbalancer.listeners"`
	LoadbalancerLoadbalancers      float64 `json:"loadbalancer.loadbalancers"`
	LoadbalancerMembers            float64 `json:"loadbalancer.members"`
	LoadbalancerPools              float64 `json:"loadbalancer.pools"`
	NetworkFloatingips             float64 `json:"network.floatingips"`
	NetworkNetworks                float64 `json:"network.networks"`
	NetworkPorts                   float64 `json:"network.ports"`
	NetworkRbacPolicies            float64 `json:"network.rbac_policies"`
	NetworkRouters                 float64 `json:"network.routers"`
	NetworkSecurityGroupRules      float64 `json:"network.security_group_rules"`
	NetworkSecurityGroups          float64 `json:"network.security_groups"`
	NetworkSubnetPools             float64 `json:"network.subnet_pools"`
	NetworkSubnets                 float64 `json:"network.subnets"`
	NetworkTrunks                  float64 `json:"network.trunks"`
	NetworkVpnEndpointGroups       float64 `json:"network.vpn_endpoint_groups"`
	NetworkVpnIkepolicies          float64 `json:"network.vpn_ikepolicies"`
	NetworkVpnIpsecSiteConnections float64 `json:"network.vpn_ipsec_site_connections"`
	NetworkVpnIpsecpolicies        float64 `json:"network.vpn_ipsecpolicies"`
	NetworkVpnServices             float64 `json:"network.vpn_services"`
	Objectstorage                  []struct {
		SpaceBytes float64  `json:"space_bytes"`
		Type       string `json:"type"`
	} `json:"objectstorage"`
	VolumeBackupGb  							 float64 `json:"volume.backup_gb"`
	VolumeBackups   							 float64 `json:"volume.backups"`
	VolumeSnapshots 							 float64 `json:"volume.snapshots"`
	VolumeSpaceGb   							 float64 `json:"volume.space_gb"`
	VolumeVolumes   							 float64 `json:"volume.volumes"`
}

type CurrentUsageV3 struct {
	ComputeCores   								 float64 `json:"compute.cores"`
	ComputeFlavors                 map[string]float64 `json:"compute.flavors"`
	ComputeInstances    					 float64 `json:"compute.instances"`
	ComputeRAMMb        					 float64 `json:"compute.ram_mb"`
	ComputeServerGroups 					 float64 `json:"compute.server_groups"`
	DNSZones            					 float64 `json:"dns.zones"`
	ImageImages         					 float64 `json:"image.images"`
	ImageSpaceBytes     					 float64 `json:"image.space_bytes"`
	LoadbalancerFlavors 					 struct {
	} `json:"loadbalancer.flavors"`
	LoadbalancerHealthmonitors     float64 `json:"loadbalancer.healthmonitors"`
	LoadbalancerListeners          float64 `json:"loadbalancer.listeners"`
	LoadbalancerLoadbalancers      float64 `json:"loadbalancer.loadbalancers"`
	LoadbalancerMembers            float64 `json:"loadbalancer.members"`
	LoadbalancerPools              float64 `json:"loadbalancer.pools"`
	NetworkFloatingips             float64 `json:"network.floatingips"`
	NetworkNetworks                float64 `json:"network.networks"`
	NetworkPorts                   float64 `json:"network.ports"`
	NetworkRbacPolicies            float64 `json:"network.rbac_policies"`
	NetworkRouters                 float64 `json:"network.routers"`
	NetworkSecurityGroupRules      float64 `json:"network.security_group_rules"`
	NetworkSecurityGroups          float64 `json:"network.security_groups"`
	NetworkSubnetPools             float64 `json:"network.subnet_pools"`
	NetworkSubnets                 float64 `json:"network.subnets"`
	NetworkTrunks                  float64 `json:"network.trunks"`
	NetworkVpnEndpointGroups       float64 `json:"network.vpn_endpoint_groups"`
	NetworkVpnIkepolicies          float64 `json:"network.vpn_ikepolicies"`
	NetworkVpnIpsecSiteConnections float64 `json:"network.vpn_ipsec_site_connections"`
	NetworkVpnIpsecpolicies        float64 `json:"network.vpn_ipsecpolicies"`
	NetworkVpnServices             float64 `json:"network.vpn_services"`
	Objectstorage                  []struct {
		SpaceBytes float64    `json:"space_bytes"`
		Type       string `json:"type"`
	} `json:"objectstorage"`
	VolumeBackupGb  float64 `json:"volume.backup_gb"`
	VolumeBackups   float64 `json:"volume.backups"`
	VolumeSnapshots float64 `json:"volume.snapshots"`
	VolumeSpaceGb   float64 `json:"volume.space_gb"`
	VolumeVolumes   float64 `json:"volume.volumes"`
}

type QuotaV1 struct {
	ComputeCores                   float64 `json:"compute.cores"`
	ComputeInstances               float64 `json:"compute.instances"`
	ComputeRAMMb                   float64 `json:"compute.ram_mb"`
	DNSZones                       float64 `json:"dns.zones"`
	NetworkFloatingips             float64 `json:"network.floatingips"`
	NetworkLbHealthmonitors        float64 `json:"network.lb_healthmonitors"`
	NetworkLbListeners             float64 `json:"network.lb_listeners"`
	NetworkLbMembers               float64 `json:"network.lb_members"`
	NetworkLbPools                 float64 `json:"network.lb_pools"`
	NetworkLoadbalancers           float64 `json:"network.loadbalancers"`
	NetworkVpnEndpointGroups       float64 `json:"network.vpn_endpoint_groups"`
	NetworkVpnIkepolicies          float64 `json:"network.vpn_ikepolicies"`
	NetworkVpnIpsecSiteConnections float64 `json:"network.vpn_ipsec_site_connections"`
	NetworkVpnIpsecpolicies        float64 `json:"network.vpn_ipsecpolicies"`
	NetworkVpnServices             float64 `json:"network.vpn_services"`
	S3SpaceBytes                   float64 `json:"s3.space_bytes"`
	VolumeSpaceGb                  float64 `json:"volume.space_gb"`
	VolumeVolumes                  float64 `json:"volume.volumes"`
}

type CurrentUsageV1 struct {
	ComputeCores                   float64            `json:"compute.cores"`
	ComputeFlavors                 map[string]float64 `json:"compute.flavors"`
	ComputeInstances               float64            `json:"compute.instances"`
	ComputeRAMMb                   float64            `json:"compute.ram_mb"`
	DNSZones                       float64            `json:"dns.zones"`
	NetworkFloatingips             float64            `json:"network.floatingips"`
	NetworkLbHealthmonitors        float64            `json:"network.lb_healthmonitors"`
	NetworkLbListeners             float64            `json:"network.lb_listeners"`
	NetworkLbMembers               float64            `json:"network.lb_members"`
	NetworkLbPools                 float64            `json:"network.lb_pools"`
	NetworkLoadbalancers           float64            `json:"network.loadbalancers"`
	NetworkVpnEndpointGroups       float64            `json:"network.vpn_endpoint_groups"`
	NetworkVpnIkepolicies          float64            `json:"network.vpn_ikepolicies"`
	NetworkVpnIpsecSiteConnections float64            `json:"network.vpn_ipsec_site_connections"`
	NetworkVpnIpsecpolicies        float64            `json:"network.vpn_ipsecpolicies"`
	NetworkVpnServices             float64            `json:"network.vpn_services"`
	S3SpaceBytes                   float64            `json:"s3.space_bytes"`
	VolumeSpaceGb                  float64            `json:"volume.space_gb"`
	VolumeVolumes                  float64            `json:"volume.volumes"`
}
