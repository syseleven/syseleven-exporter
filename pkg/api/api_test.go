package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			json.NewEncoder(w).Encode(Error{
				Title:  "Internal Server Error",
				Detail: "Keystone authentication failed — token expired",
				Type:   "auth_error",
			})
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
