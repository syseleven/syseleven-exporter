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
	"net/url"
	"os"
)

var endpoint string
var endpointIam string

func GetRadosgwIdentifiers(orgID, projectID, secret string) ([]string, error) {
	requestURL := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users/radosgw-identifiers", endpointIam, orgID, projectID)
	resp, err := MakeRequest(requestURL, secret, "X-S11-CREDENTIAL")
	if err != nil {
		return nil, fmt.Errorf("get radosgw identifiers: %w", err)
	}

	var identifiers []string
	if err := json.Unmarshal(resp, &identifiers); err != nil {
		return nil, fmt.Errorf("decode radosgw identifiers: %w", err)
	}
	return identifiers, nil
}

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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body from %s: %w", url, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s returned HTTP %d: %s", url, resp.StatusCode, body)
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
	url := fmt.Sprintf("%s/v3/projects/%s/quota", endpoint, projectID)
	resp, err := MakeRequest(url, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get quota v3: %w", err)
	}

	var quotas = make(map[string]QuotaV3)

	if err := json.Unmarshal(resp, &quotas); err != nil {
		return nil, err
	}

	return quotas, nil
}

func GetCurrentUsageV3(projectID, token string) (map[string]CurrentUsageV3, error) {
	url := fmt.Sprintf("%s/v3/projects/%s/current_usage", endpoint, projectID)
	resp, err := MakeRequest(url, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get current usage v3: %w", err)
	}

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
	resp, err := MakeRequest(url, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get quota v1: %w", err)
	}

	var quotas = make(map[string]QuotaV1)

	err = json.Unmarshal(resp, &quotas)

	if err != nil {
		return nil, err
	}

	return quotas, nil
}

func GetCurrentUsageV1(projectID, token string) (map[string]CurrentUsageV1, error) {
	url := fmt.Sprintf("%s/v1/projects/%s/current_usage", endpoint, projectID)
	resp, err := MakeRequest(url, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get current usage v1: %w", err)
	}

	var currentUsages = make(map[string]CurrentUsageV1)

	if err := json.Unmarshal(resp, &currentUsages); err != nil {
		return nil, err
	}

	return currentUsages, nil
}

func GetS3InfoNCS(projectID string) ([]S3UsageNCS, error) {
	endpointIam = "https://iam.apis.syseleven.de"
	if len(os.Getenv("SYSELEVEN_IAM_API_ENDPOINT")) > 0 {
		endpointIam = os.Getenv("SYSELEVEN_IAM_API_ENDPOINT")
	}

	orgID := os.Getenv("IAM_ORG_ID")
	secret := os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET")
	identifiers, err := GetRadosgwIdentifiers(orgID, projectID, secret)
	if err != nil {
		return nil, err
	}
	s3Users, err := GetS3Users(orgID, projectID, secret)
	if err != nil {
		return nil, err
	}
	s3Usage := make([]S3UsageNCS, 0, len(s3Users)*len(identifiers))
	for _, user := range s3Users {
		for _, target := range identifiers {
			requestURL := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users/%s/quota?target_object_storage=%s", endpointIam, orgID, projectID, user.Id, url.QueryEscape(target))
			resp, err := MakeRequest(requestURL, secret, "X-S11-CREDENTIAL")
			if err != nil {
				return nil, fmt.Errorf("get s3 info ncs user %s target %s: %w", user.Id, target, err)
			}

			var currentUsage S3InfoNCS
			if err := json.Unmarshal(resp, &currentUsage); err != nil {
				return nil, fmt.Errorf("decode s3 info ncs user %s target %s: %w", user.Id, target, err)
			}
			s3Usage = append(s3Usage, S3UsageNCS{S3UsersNCS: user, S3InfoNCS: currentUsage, Target: target})
		}
	}
	return s3Usage, nil
}

func GetS3Users(orgID, projectID, secret string) ([]S3UsersNCS, error) {
	url := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users", endpointIam, orgID, projectID)
	resp, err := MakeRequest(url, secret, "X-S11-CREDENTIAL")
	if err != nil {
		return nil, fmt.Errorf("get s3 users: %w", err)
	}

	var s3users []S3UsersNCS
	err = json.Unmarshal(resp, &s3users)

	return s3users, err
}
