# SysEleven Exporter Helm-Chart
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Release Charts](https://github.com/vfm/syseleven-exporter-chart/workflows/Release%20Charts/badge.svg?branch=main)
![Badge](https://img.shields.io/badge/made%20with%20%E2%9D%A4%20by-vfm-002C5B)


A Helm Chart to deploy the [SysEleven Exporter](https://github.com/syseleven/syseleven-exporter) in Kubernetes.

## Usage

Add this Github-Repository as Helm Repository.

```bash
helm repo add syseleven-exporter https://syseleven.github.io/syseleven-exporter
helm repo update
```

Now you can install the Chart.

```bash
helm install syseleven-exporter/syseleven-exporter --generate-name
```

Add your Openstack credentials to allow *read-only* access to the API-Endpoint `https://api.cloud.syseleven.net:5001`

```bash
openstack.username
openstack.password
openstack.projectId
```

The exporter uses the API version v1 by default. If you want to change to the current API version v3, you can adjust the API version with `openstack.apiversion`. See here for more information about the [API for Quota and Usage Information](https://docs.syseleven.de/syseleven-stack/en/reference/get-quota-info).

### Prometheus

If you want a [Prometheus ServiceMonitor](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#servicemonitor) you can use the `prometheus.serviceMonitor.enabled` flag, which is false by default. The default `serviceMonitor.interval` is 10 minutes, since the API is only polled once per hour by the exporter pod.

```bash
prometheus.serviceMonitor.enabled
prometheus.serviceMonitor.interval
```

The default Namespace for the ServiceMonitor is the Helm Release Namespace.
If you need the ServiceMonitor in another namespace, you have to specify the `prometheus.serviceMonitor.namespace`.

```bash
prometheus.serviceMonitor.namespace
```

#### Rules

To use a default set of Prometheus Alertrules, you have to enable the `prometheus.rules.enabled` flag. Similar to the ServiceMonitor, you can adjust the default Namespace with `prometheus.rules.namespace`.

To see the full list of all options have a look at [Values.yaml](charts/syseleven-exporter-chart/values.yaml).


#### Dashboard

To use a default Grafana Dashboard, you have to enable the `prometheus.dashboards.enabled` flag. Similar to the Rules, you can adjust the Namespace with `prometheus.dashboards.namespace`.

To see the full list of all options have a look at [Values.yaml](charts/syseleven-exporter-chart/values.yaml).
