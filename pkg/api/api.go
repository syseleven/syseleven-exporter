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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultQuotaEndpoint = "https://api.cloud.syseleven.net:5001"
	defaultIamEndpoint   = "https://iam.apis.syseleven.de"
)

// ErrUnsupportedURLScheme is returned when a request URL does not use http or https.
var ErrUnsupportedURLScheme = errors.New("unsupported URL scheme")

type Error struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Type   string `json:"type"`
}

// Error implements the error interface, preserving the upstream API error format.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s (%s)", e.Title, e.Detail, e.Type)
}

// quotaEndpoint returns the quota API endpoint from SYSELEVEN_QUOTA_API_ENDPOINT,
// falling back to the SysEleven default.
func quotaEndpoint() string {
	if v := os.Getenv("SYSELEVEN_QUOTA_API_ENDPOINT"); v != "" {
		return v
	}

	return defaultQuotaEndpoint
}

// iamEndpoint returns the IAM API endpoint from SYSELEVEN_IAM_API_ENDPOINT,
// falling back to the SysEleven default.
func iamEndpoint() string {
	if v := os.Getenv("SYSELEVEN_IAM_API_ENDPOINT"); v != "" {
		return v
	}

	return defaultIamEndpoint
}

func MakeRequest(rawURL string, token string, header string) ([]byte, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse request url: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedURLScheme, parsed.Scheme)
	}

	// The URL host intentionally comes from operator-set environment
	// configuration (SYSELEVEN_QUOTA_API_ENDPOINT / SYSELEVEN_IAM_API_ENDPOINT),
	// never from request input, and the scheme is validated above — not an SSRF
	// vector.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, parsed.String(), nil) //nolint:gosec // G704: see above
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set(header, token)

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // G704: see above
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError Error

		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, fmt.Errorf("decode api error: %w", err)
		}

		return nil, apiError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return body, nil
}

// use API v3 for Quota and Usage Information

func GetQuotaV3(projectID, token string) (map[string]QuotaV3, error) {
	requestURL := fmt.Sprintf("%s/v3/projects/%s/quota", quotaEndpoint(), projectID)

	resp, err := MakeRequest(requestURL, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get quota v3: %w", err)
	}

	quotas := make(map[string]QuotaV3)

	if err := json.Unmarshal(resp, &quotas); err != nil {
		return nil, fmt.Errorf("unmarshal quota v3: %w", err)
	}

	return quotas, nil
}

func GetCurrentUsageV3(projectID, token string) (map[string]CurrentUsageV3, error) {
	requestURL := fmt.Sprintf("%s/v3/projects/%s/current_usage", quotaEndpoint(), projectID)

	resp, err := MakeRequest(requestURL, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get current usage v3: %w", err)
	}

	currentUsages := make(map[string]CurrentUsageV3)

	if err := json.Unmarshal(resp, &currentUsages); err != nil {
		return nil, fmt.Errorf("unmarshal current usage v3: %w", err)
	}

	return currentUsages, nil
}

// use API v3 for Quota and Usage Information

func GetQuotaV1(projectID, token string) (map[string]QuotaV1, error) {
	requestURL := fmt.Sprintf("%s/v1/projects/%s/quota", quotaEndpoint(), projectID)

	resp, err := MakeRequest(requestURL, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get quota v1: %w", err)
	}

	quotas := make(map[string]QuotaV1)

	err = json.Unmarshal(resp, &quotas)
	if err != nil {
		return nil, fmt.Errorf("unmarshal quota v1: %w", err)
	}

	return quotas, nil
}

func GetCurrentUsageV1(projectID, token string) (map[string]CurrentUsageV1, error) {
	requestURL := fmt.Sprintf("%s/v1/projects/%s/current_usage", quotaEndpoint(), projectID)

	resp, err := MakeRequest(requestURL, token, "X-Auth-Token")
	if err != nil {
		return nil, fmt.Errorf("get current usage v1: %w", err)
	}

	currentUsages := make(map[string]CurrentUsageV1)

	if err := json.Unmarshal(resp, &currentUsages); err != nil {
		return nil, fmt.Errorf("unmarshal current usage v1: %w", err)
	}

	return currentUsages, nil
}

func GetS3InfoNCS(projectID string) ([]S3UsageNCS, error) {
	orgID := os.Getenv("IAM_ORG_ID")
	secret := os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET")

	s3Users, err := GetS3Users(orgID, projectID, secret)
	if err != nil {
		return nil, err
	}

	s3Usage := []S3UsageNCS{}

	for _, t := range s3Users {
		requestURL := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users/%s/quota", iamEndpoint(), orgID, projectID, t.ID)

		resp, err := MakeRequest(requestURL, secret, "X-S11-CREDENTIAL")
		if err != nil {
			return nil, fmt.Errorf("get s3 info ncs user %s: %w", t.ID, err)
		}

		var currentUsage S3InfoNCS

		if err := json.Unmarshal(resp, &currentUsage); err != nil {
			return nil, fmt.Errorf("unmarshal s3 info ncs user %s: %w", t.ID, err)
		}

		s3Usage = append(s3Usage, S3UsageNCS{S3UsersNCS: t, S3InfoNCS: currentUsage})
	}

	return s3Usage, nil
}

func GetS3Users(orgID, projectID, secret string) ([]S3UsersNCS, error) {
	requestURL := fmt.Sprintf("%s/v3/orgs/%s/projects/%s/s3-users", iamEndpoint(), orgID, projectID)

	resp, err := MakeRequest(requestURL, secret, "X-S11-CREDENTIAL")
	if err != nil {
		return nil, fmt.Errorf("get s3 users: %w", err)
	}

	var s3users []S3UsersNCS

	if err := json.Unmarshal(resp, &s3users); err != nil {
		return nil, fmt.Errorf("unmarshal s3 users: %w", err)
	}

	return s3users, nil
}
