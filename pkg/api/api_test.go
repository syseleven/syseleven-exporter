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

package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/syseleven/syseleven-exporter/pkg/api"
)

const (
	testProjectID = "proj-123"
	testOrgID     = "org-456"
	testToken     = "test-token"
	testSecret    = "s11_orgsa_secret"
)

// pkg/api resolves the API and IAM endpoints into package globals that only
// GetQuota*/GetS3InfoNCS set from the environment; the dependent calls reuse
// them. Tests therefore follow the exporter's real call order.

func TestGetQuota_Success(t *testing.T) {
	tests := []struct {
		name      string
		wantPath  string
		response  string
		wantValue float64
		invoke    func() (float64, error)
	}{
		{
			name:      "GetQuotaV3",
			wantPath:  "/v3/projects/" + testProjectID + "/quota",
			response:  `{"region-a":{"compute.cores":4}}`,
			wantValue: 4,
			invoke: func() (float64, error) {
				q, err := api.GetQuotaV3(testProjectID, testToken)
				if err != nil {
					return 0, err
				}

				return q["region-a"].ComputeCores, nil
			},
		},
		{
			name:      "GetQuotaV1",
			wantPath:  "/v1/projects/" + testProjectID + "/quota",
			response:  `{"region-a":{"compute.cores":8}}`,
			wantValue: 8,
			invoke: func() (float64, error) {
				q, err := api.GetQuotaV1(testProjectID, testToken)
				if err != nil {
					return 0, err
				}

				return q["region-a"].ComputeCores, nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotPath, gotToken string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotToken = r.Header.Get("X-Auth-Token")

				if _, err := w.Write([]byte(tc.response)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer srv.Close()

			t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

			got, err := tc.invoke()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.wantValue {
				t.Errorf("parsed value: got %v, want %v", got, tc.wantValue)
			}

			if gotPath != tc.wantPath {
				t.Errorf("request path: got %q, want %q", gotPath, tc.wantPath)
			}

			if gotToken != testToken {
				t.Errorf("auth header X-Auth-Token: got %q, want %q", gotToken, testToken)
			}
		})
	}
}

// TestQuotaThenUsage exercises the real exporter flow: a quota call resolves the
// endpoint, then a current-usage call reuses it.
func TestQuotaThenUsage(t *testing.T) {
	tests := []struct {
		name          string
		quotaResp     string
		usageResp     string
		wantUsagePath string
		invoke        func() (float64, error)
	}{
		{
			name:          "v3",
			quotaResp:     `{"region-a":{"compute.cores":4}}`,
			usageResp:     `{"region-a":{"compute.cores":2}}`,
			wantUsagePath: "/v3/projects/" + testProjectID + "/current_usage",
			invoke: func() (float64, error) {
				if _, err := api.GetQuotaV3(testProjectID, testToken); err != nil {
					return 0, err
				}

				u, err := api.GetCurrentUsageV3(testProjectID, testToken)
				if err != nil {
					return 0, err
				}

				return u["region-a"].ComputeCores, nil
			},
		},
		{
			name:          "v1",
			quotaResp:     `{"region-a":{"compute.cores":8}}`,
			usageResp:     `{"region-a":{"compute.cores":6}}`,
			wantUsagePath: "/v1/projects/" + testProjectID + "/current_usage",
			invoke: func() (float64, error) {
				if _, err := api.GetQuotaV1(testProjectID, testToken); err != nil {
					return 0, err
				}

				u, err := api.GetCurrentUsageV1(testProjectID, testToken)
				if err != nil {
					return 0, err
				}

				return u["region-a"].ComputeCores, nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotUsagePath string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body := tc.quotaResp
				if strings.HasSuffix(r.URL.Path, "/current_usage") {
					gotUsagePath = r.URL.Path
					body = tc.usageResp
				}

				if _, err := w.Write([]byte(body)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer srv.Close()

			t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

			got, err := tc.invoke()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != 2 && got != 6 {
				t.Errorf("parsed usage value: got %v, want 2 (v3) or 6 (v1)", got)
			}

			if gotUsagePath != tc.wantUsagePath {
				t.Errorf("usage request path: got %q, want %q", gotUsagePath, tc.wantUsagePath)
			}
		})
	}
}

func TestGetQuotaV3_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte(`{not valid json`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer srv.Close()

	t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

	_, err := api.GetQuotaV3(testProjectID, testToken)
	if err == nil {
		t.Fatal("got nil error for malformed JSON body")
	}

	if !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("expected a JSON parse error, got: %q", err.Error())
	}
}

// TestGetS3InfoNCS_Success drives the full NCS S3 flow (list users, then fetch
// each user's quota) and covers the IAM endpoint override.
func TestGetS3InfoNCS_Success(t *testing.T) {
	var gotUsersPath, gotQuotaCred string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string

		switch {
		case strings.HasSuffix(r.URL.Path, "/s3-users"):
			gotUsersPath = r.URL.Path
			body = `[{"name":"user-a","id":"id-a","description":"desc-a"}]`
		case strings.HasSuffix(r.URL.Path, "/quota"):
			gotQuotaCred = r.Header.Get("X-S11-CREDENTIAL")
			body = `{"size":100,"max_size":1000,"num_objects":5,"max_objects":50,"enabled":true,"check_on_raw":true}`
		default:
			t.Errorf("unexpected request path: %q", r.URL.Path)
		}

		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer srv.Close()

	t.Setenv("SYSELEVEN_IAM_API_ENDPOINT", srv.URL)
	t.Setenv("IAM_ORG_ID", testOrgID)
	t.Setenv("OS_APPLICATION_CREDENTIAL_SECRET", testSecret)

	usage, err := api.GetS3InfoNCS(testProjectID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(usage) != 1 {
		t.Fatalf("expected 1 s3 usage entry, got %d", len(usage))
	}

	if usage[0].Name != "user-a" {
		t.Errorf("user name: got %q, want %q", usage[0].Name, "user-a")
	}

	if usage[0].Size != 100 || !usage[0].Enabled {
		t.Errorf("s3 info: got size=%v enabled=%v, want size=100 enabled=true", usage[0].Size, usage[0].Enabled)
	}

	wantUsersPath := "/v3/orgs/" + testOrgID + "/projects/" + testProjectID + "/s3-users"
	if gotUsersPath != wantUsersPath {
		t.Errorf("s3-users request path: got %q, want %q", gotUsersPath, wantUsersPath)
	}

	if gotQuotaCred != testSecret {
		t.Errorf("quota auth header X-S11-CREDENTIAL: got %q, want %q", gotQuotaCred, testSecret)
	}
}

// Regression (PR #90): call sites used `resp, _ := MakeRequest(...)`, so real API
// errors were masked by "unexpected end of JSON input". Errors must carry the
// function context and the upstream message.
func TestGetQuotaV3_ErrorPropagation(t *testing.T) {
	t.Run("non-2xx response returns real API error message", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)

			if err := json.NewEncoder(w).Encode(api.Error{
				Title:  "Internal Server Error",
				Detail: "Keystone authentication failed — token expired",
				Type:   "auth_error",
			}); err != nil {
				t.Errorf("encode response: %v", err)
			}
		}))
		defer srv.Close()

		t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

		_, err := api.GetQuotaV3("proj-abc", "bad-token")
		if err == nil {
			t.Fatal("got nil error — API failure was completely undetected")
		}

		if !strings.Contains(err.Error(), "get quota v3:") {
			t.Fatalf("missing function context: %q", err.Error())
		}

		if !strings.Contains(err.Error(), "Keystone authentication failed") {
			t.Fatalf("missing upstream API message: %q", err.Error())
		}

		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			t.Fatalf("got misleading json parse error instead of real API error: %q", err.Error())
		}
	})

	t.Run("unreachable API returns network error not json error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			hj, ok := w.(http.Hijacker)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			conn, _, _ := hj.Hijack()
			conn.Close()
		}))
		defer srv.Close()

		t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

		_, err := api.GetQuotaV3("proj-abc", "any-token")
		if err == nil {
			t.Fatal("got nil error — unreachable API was completely undetected")
		}

		if !strings.Contains(err.Error(), "get quota v3:") {
			t.Fatalf("missing function context: %q", err.Error())
		}

		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			t.Fatalf("got misleading json parse error instead of network error: %q", err.Error())
		}
	})
}
