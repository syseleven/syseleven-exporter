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
	"net/http"
	"os"
)

type Error struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Type   string `json:"type"`
}

var endpoint string

// use API v3 for Quota and Usage Information

func GetQuotaV3(projectID, token string) (map[string]QuotaV3, error) {
	if os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT") == "" {
		endpoint = "https://api.cloud.syseleven.net:5001"
	} else {
		endpoint = os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/projects/%s/quota", endpoint, projectID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError Error

		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s (%s)", apiError.Title, apiError.Detail, apiError.Type)
	}

	var quotas map[string]QuotaV3
	quotas = make(map[string]QuotaV3)

	err = json.NewDecoder(resp.Body).Decode(&quotas)

	if err != nil {
		return nil, err
	}

	return quotas, nil
}

func GetCurrentUsageV3(projectID, token string) (map[string]CurrentUsageV3, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/projects/%s/current_usage", endpoint, projectID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError Error

		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s (%s)", apiError.Title, apiError.Detail, apiError.Type)
	}

	var currentUsages map[string]CurrentUsageV3
	currentUsages = make(map[string]CurrentUsageV3)

	err = json.NewDecoder(resp.Body).Decode(&currentUsages)

	if err != nil {
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

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/projects/%s/quota", endpoint, projectID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError Error

		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s (%s)", apiError.Title, apiError.Detail, apiError.Type)
	}

	var quotas map[string]QuotaV1
	quotas = make(map[string]QuotaV1)

	err = json.NewDecoder(resp.Body).Decode(&quotas)

	if err != nil {
		return nil, err
	}

	return quotas, nil
}

func GetCurrentUsageV1(projectID, token string) (map[string]CurrentUsageV1, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/projects/%s/current_usage", endpoint, projectID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError Error

		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s (%s)", apiError.Title, apiError.Detail, apiError.Type)
	}

	var currentUsages map[string]CurrentUsageV1
	currentUsages = make(map[string]CurrentUsageV1)

	err = json.NewDecoder(resp.Body).Decode(&currentUsages)

	if err != nil {
		return nil, err
	}

	return currentUsages, nil
}
