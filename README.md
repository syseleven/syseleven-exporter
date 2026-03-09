# SysEleven Exporter

> This project is a fork of the original [syseleven-exporter](https://github.com/Staffbase/syseleven-exporter) created by [Staffbase](https://github.com/Staffbase/).

Export your quota and current usage from the SysEleven API as Prometheus metrics. The exporter uses the `https://api.cloud.syseleven.net:5001` API endpoint to get the quota and usage statistics for all SysEleven resources.

The SysEleven Exporter can be deployed using the [Helm chart](https://github.com/syseleven/syseleven-exporter/tree/main/charts/syseleven-exporter-chart).

## Usage

Clone the repository and build the binary:

```sh
git clone git@github.com:syseleven/syseleven-exporter.git
make build
```

Set the environment variables for authentication against the SysEleven API.
There are two authentication options available:

> [!NOTE]
> Fetching NCS S3 metrics only works with application credentials

- Username and password

  ```sh
  export OS_USERNAME=
  export OS_PASSWORD=
  # OS_PROJECT_ID could be a comma separated list of project IDs
  export OS_PROJECT_ID=
  ```

- Application credentials (do **not** set `OS_PROJECT_ID`)

  ```sh
  export OS_APPLICATION_CREDENTIAL_ID=
  export OS_APPLICATION_CREDENTIAL_SECRET=
  ```

  *Note that when using application credentials you cannot specify `OS_PROJECT_ID`.*
  *Else the authentication won't work. Also this means that you can only scrape metrics for one project.*
  *This will be the project the application credentials are created in and scoped to.*

  ---

  For NCS (HAM1/DUS2), the application credentials look like following:

  ```sh
  export OS_APPLICATION_CREDENTIAL_ID=s11auth:<PROJECT_ID>
  export OS_APPLICATION_CREDENTIAL_SECRET=s11_orgsa_<...>
  ```

- In order for the exporter to fetch S3 metrics from NCS, you need to set the `IAM_ORG_ID` variable:

  ```sh
  export IAM_ORG_ID=
  ```
  You can find this ID from the URL in SysEleven Dashboard, e.g. `https://dashboard.syseleven.de/dashboard/iam/organizations/00000000-0000-0000-0000-000000000000`

  ---

- Optional: Configure API endpoints for non standard (not `https://keystone.cloud.syseleven.net:5000/v3`,
`https://api.cloud.syseleven.net:5001` or `https://iam.apis.syseleven.de`)

  ```sh
  export OS_AUTH_URL="https://api.example.syseleven.de:5000"
  export SYSELEVEN_QUOTA_API_ENDPOINT="https://api.example.syseleven.de:5001"

  # When using S3 in NCS
  export SYSELEVEN_IAM_API_ENDPOINT="https://iam.example.syseleven.de"
  ```

Then run the exporter:

```sh
./bin/syselevenexporter
```

The exporter uses the API version v1 by default. If you want to change to the current API version v3, you can run the exporter with:

```sh
./bin/syselevenexporter  --api-version v3
```

See here for more information about the [API for Quota and Usage Information](https://docs.syseleven.de/syseleven-stack/en/reference/get-quota-info).

A Docker image is available at `syseleven/syseleven-exporter:<TAG>` and can be retrieved via:

```sh
docker pull syseleven/syseleven-exporter:<TAG>
```

## Metrics

| Metric | Description |
| ------ | ----------- |
| syseleven_compute_cores_total | Quota for number of compute cores per `region` and `project` |
| syseleven_compute_cores_used | Number of used compute cores per `region` and `project` |
| syseleven_compute_instances_total | Quota for number of compute instances per `region` and `project` |
| syseleven_compute_instances_used | Number of used compute instances per `region` and `project` |
| syseleven_compute_flavors_used | Number of used compute flavors per `type` and `region` and `project` |
| syseleven_compute_ram_total_megabytes | Quota for ram per `region` and `project` in megabytes |
| syseleven_compute_ram_used_megabytes | Used ram per `region` and `project` in megabytes |
| syseleven_dns_zones_total | Quota for number of DNS zones per `region` and `project` |
| syseleven_dns_zones_used | Number of used DNS zones per `region` and `project` |
| syseleven_network_floating_ips_total | Quota for number of floating IPs per `region` and `project` |
| syseleven_network_floating_ips_used | Number of used floating IPs per `region` and `project` |
| syseleven_network_loadbalancers_total | Quota for number of load balancers per `region` and `project` |
| syseleven_network_loadbalancers_used | Number of used load balancers per `region` and `project` |
| syseleven_s3_space_total_bytes | Quota for S3 space per `region`, `project` and `type` in bytes |
| syseleven_s3_space_used_bytes | Used S3 space per `region` and `project` in bytes |
| syseleven_s3_space_max_bytes_ncs | Quota for S3 space in ncs regions in bytes |
| syseleven_s3_space_used_bytes_ncs | Used S3 space in ncs regions in bytes |
| syseleven_s3_num_objects_ncs | Number of objects stored in S3 |
| syseleven_s3_max_objects_ncs | Maximal number of objects stored in S3 |
| syseleven_s3_enabled_ncs | Checks if s3 space is enabled for user or not |
| syseleven_s3_check_enabled_ncs | Checks if check on raw is enabled for user or not |
| syseleven_volume_space_total_gigabytes | Quota for volume space per `region` and `project` in gigabytes |
| syseleven_volume_space_used_gigabytes | Number of used volume space per `region` and `project` in gigabytes |
| syseleven_volume_volumes_total | Quota for number of volumes per `region` and `project` |
| syseleven_volume_volumes_used | Number of used volumes per `region` and `project` |

## Credits

- [Staffbase](https://github.com/Staffbase/) - Original Author
- [vfm](https://github.com/vfm) - Original Author of the Helm chart
