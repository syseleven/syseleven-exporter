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

package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/syseleven/syseleven-exporter/pkg/auth"
)

// These cover what is stable in the gophercloud-backed auth flow: the OS_AUTH_URL
// override is honored and a Keystone failure surfaces as an error. The success
// path (including the unchecked tokens.CreateResult assertion in GetProject) is a
// gap until #114 makes the client injectable.

func keystone401(t *testing.T) (url string, hit *bool) {
	t.Helper()

	var contacted bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		contacted = true
		w.WriteHeader(http.StatusUnauthorized)

		if _, err := w.Write([]byte(`{"error":{"code":401,"message":"denied"}}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)

	return srv.URL + "/v3", &contacted
}

func TestAuthFunctions_UseOverrideAndSurfaceError(t *testing.T) {
	tests := []struct {
		name   string
		invoke func() error
	}{
		{
			name:   "GetToken",
			invoke: func() error { _, err := auth.GetToken("project-id", "user", "pass"); return err },
		},
		{
			name:   "GetTokenAppCreds",
			invoke: func() error { _, err := auth.GetTokenAppCreds("project-id", "cred-id", "cred-secret"); return err },
		},
		{
			name:   "GetProject",
			invoke: func() error { _, err := auth.GetProject("cred-id", "cred-secret"); return err },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url, contacted := keystone401(t)
			t.Setenv("OS_AUTH_URL", url)

			if err := tc.invoke(); err == nil {
				t.Fatal("expected an error from a 401 Keystone, got nil")
			}

			if !*contacted {
				t.Error("OS_AUTH_URL override was not honored: stub Keystone was never contacted")
			}
		})
	}
}
