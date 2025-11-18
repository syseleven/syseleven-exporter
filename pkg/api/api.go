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

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Error struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Type   string `json:"type"`
}

var endpoint string
var endpointIam string

func MakeRequest(url string, token string, header string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(header, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError Error

		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s (%s)", apiError.Title, apiError.Detail, apiError.Type)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// use API v3 for Quota and Usage Information

func GetQuotaV3(projectID, token string) (map[string]QuotaV3, error) {
	if os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT") == "" {
		endpoint = "https://api.cloud.syseleven.net:5001"
	} else {
		endpoint = os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT")
	}
	url := fmt.Sprintf("%s/v1/projects/%s/quota", endpoint, projectID)
	resp, _ := MakeRequest(url, token, "X-Auth-Token")

	var quotas = make(map[string]QuotaV3)

	if err := json.Unmarshal(resp, &quotas); err != nil {
		return nil, err
	}

	return quotas, nil
}

func GetCurrentUsageV3(projectID, token string) (map[string]CurrentUsageV3, error) {
	url := fmt.Sprintf("%s/v3/projects/%s/current_usage", endpoint, projectID)
	resp, _ := MakeRequest(url, token, "X-Auth-Token")

	var currentUsages = make(map[string]CurrentUsageV3)

	if err := json.Unmarshal(resp, &currentUsages); err != nil {
		return nil, err
	}

	return currentUsages, nil
}

// use API v3 for Quota and Usage Information

func GetQuotaV1(projectID, token string) (map[string]QuotaV1, error) {
	if os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT") == "" {
		endpoint = "https://api.cloud.syseleven.net:5001"
	} else {
		endpoint = os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT")
	}

	url := fmt.Sprintf("%s/v1/projects/%s/quota", endpoint, projectID)
	resp, _ := MakeRequest(url, token, "X-Auth-Token")

	var quotas = make(map[string]QuotaV1)

	err := json.Unmarshal(resp, &quotas)

	if err != nil {
		return nil, err
	}

	return quotas, nil
}

func GetCurrentUsageV1(projectID, token string) (map[string]CurrentUsageV1, error) {
	url := fmt.Sprintf("%s/v1/projects/%s/current_usage", endpoint, projectID)
	resp, _ := MakeRequest(url, token, "X-Auth-Token")

	var currentUsages = make(map[string]CurrentUsageV1)

	if err := json.Unmarshal(resp, &currentUsages); err != nil {
		return nil, err
	}

	return currentUsages, nil
}

func GetS3InfoNCS(projectID string) ([]S3UsageNCS, error) {
	if os.Getenv("SYSELEVEN_IAM_API_ENDPOINT") == "" {
		endpointIam = "https://iam.apis.syseleven.de"
	} else {
		endpointIam = os.Getenv("SYSELEVEN_IAM_API_ENDPOINT")
	}

	orgID := os.Getenv("IAM_ORG_ID")
	secret := os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET")
	s3Users, err := GetS3Users(orgID, projectID, secret)
	if err != nil {
		return nil, err
	}
	s3Usage := []S3UsageNCS{}
	for _, t := range s3Users {
		url := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users/%s/quota", endpointIam, orgID, projectID, t.Id)
		resp, _ := MakeRequest(url, secret, "X-S11-CREDENTIAL")

		var currentUsage S3InfoNCS

		err := json.Unmarshal(resp, &currentUsage)

		if err != nil {
			return nil, err
		}
		s3Usage = append(s3Usage, S3UsageNCS{S3UsersNCS: t, S3InfoNCS: currentUsage})
	}
	return s3Usage, nil
}

func GetS3Users(orgID, projectID, secret string) ([]S3UsersNCS, error) {
	url := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users", endpointIam, orgID, projectID)
	resp, _ := MakeRequest(url, secret, "X-S11-CREDENTIAL")

	var s3users []S3UsersNCS
	err := json.Unmarshal(resp, &s3users)

	return s3users, err
}
