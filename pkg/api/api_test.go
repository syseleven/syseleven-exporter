package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRequest_ValidationError(t *testing.T) {
	const responseBody = `{"detail":[{"type":"string_pattern_mismatch","loc":["path","project_id"],"msg":"String should match pattern","input":"not-a-project"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(responseBody))
	}))
	defer srv.Close()

	_, err := MakeRequest(srv.URL, "token", "X-S11-CREDENTIAL")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if got := err.Error(); got != responseBody {
		t.Fatalf("error = %q, want raw response %q", got, responseBody)
	}
}

func TestGetS3InfoNCS_UsesAllRadosgwIdentifiers(t *testing.T) {
	const (
		orgID     = "org-id"
		projectID = "project-id"
		userID    = "user-id"
		secret    = "s11_orgsa_test"
	)
	quotaSizes := map[string]float64{"target-a": 2, "target-b": 3}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-S11-CREDENTIAL") != secret {
			t.Errorf("missing S3 credential header")
		}
		switch r.URL.Path {
		case "/v3/orgs/" + orgID + "/projects/" + projectID + "/s3-users/radosgw-identifiers":
			_, _ = w.Write([]byte(`["target-a","target-b"]`))
		case "/v3/orgs/" + orgID + "/projects/" + projectID + "/s3-users":
			_, _ = w.Write([]byte(`[{"name":"s3-user","id":"` + userID + `","description":"test"}]`))
		case "/v3/orgs/" + orgID + "/projects/" + projectID + "/s3-users/" + userID + "/quota":
			target := r.URL.Query().Get("target_object_storage")
			size, ok := quotaSizes[target]
			if !ok {
				t.Errorf("unexpected target_object_storage %q", target)
			}
			_, _ = w.Write([]byte(`{"size_kb":1,"size":` + fmt.Sprint(size) + `,"num_objects":3,"enabled":true,"check_on_raw":true,"max_size":4,"max_size_kb":5,"max_objects_user_quota":6,"max_objects_bucket_quota":7}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	t.Setenv("SYSELEVEN_IAM_API_ENDPOINT", srv.URL)
	t.Setenv("IAM_ORG_ID", orgID)
	t.Setenv("OS_APPLICATION_CREDENTIAL_SECRET", secret)

	usage, err := GetS3InfoNCS(projectID)
	if err != nil {
		t.Fatalf("GetS3InfoNCS returned error: %v", err)
	}
	if len(usage) != len(quotaSizes) {
		t.Fatalf("got %d S3 usage entries, want %d: %+v", len(usage), len(quotaSizes), usage)
	}
	for _, item := range usage {
		if got, want := item.Size, quotaSizes[item.Target]; got != want {
			t.Errorf("target %q size = %v, want %v", item.Target, got, want)
		}
	}
}

// TestGetQuotaV3_ErrorPropagation verifies that MakeRequest errors are
// returned to the caller with full context, not silently discarded.
//
// Regression: PR #90 introduced resp, _ := MakeRequest(...) at every call
// site, causing real API errors to be replaced by a generic json parse error
// "unexpected end of JSON input". This made failures impossible to diagnose.
func TestGetQuotaV3_ErrorPropagation(t *testing.T) {
	t.Run("non-2xx response returns real API error message", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"title":"Internal Server Error","detail":"Keystone authentication failed — token expired","type":"auth_error"}`))
		}))
		defer srv.Close()
		t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

		_, err := GetQuotaV3("proj-abc", "bad-token")

		if err == nil {
			t.Fatal("got nil error — API failure was completely undetected")
		}
		if !strings.Contains(err.Error(), "get quota v3:") {
			t.Fatalf("missing function context in error\n  want: contains %q\n  got:  %q", "get quota v3:", err.Error())
		}
		if !strings.Contains(err.Error(), "Keystone authentication failed") {
			t.Fatalf("missing real API message in error\n  want: contains %q\n  got:  %q", "Keystone authentication failed", err.Error())
		}
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			t.Fatalf("got misleading json parse error instead of real API error\n  got: %q", err.Error())
		}
	})

	t.Run("unreachable API returns network error not json error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// close connection immediately — simulates dead API
			hj, ok := w.(http.Hijacker)
			if !ok {
				w.WriteHeader(500)
				return
			}
			conn, _, _ := hj.Hijack()
			conn.Close()
		}))
		defer srv.Close()
		t.Setenv("SYSELEVEN_QUOTA_API_ENDPOINT", srv.URL)

		_, err := GetQuotaV3("proj-abc", "any-token")

		if err == nil {
			t.Fatal("got nil error — unreachable API was completely undetected")
		}
		if !strings.Contains(err.Error(), "get quota v3:") {
			t.Fatalf("missing function context in error\n  want: contains %q\n  got:  %q", "get quota v3:", err.Error())
		}
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			t.Fatalf("got misleading json parse error instead of network error\n  got: %q", err.Error())
		}
	})
}
