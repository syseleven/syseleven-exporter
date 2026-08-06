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

package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

const (
	defaultIdentityEndpoint = "https://keystone.cloud.syseleven.net:5000/v3"
	defaultDomainName       = "Default"
)

var (
	// ErrNoAuthResult is returned when the provider client holds no auth result,
	// e.g. because the token was set manually with ProviderClient.SetToken().
	ErrNoAuthResult = errors.New("no AuthResult available")
	// ErrUnexpectedAuthResult is returned when the auth result has an unexpected type.
	ErrUnexpectedAuthResult = errors.New("unexpected AuthResult type")
)

// identityEndpoint returns the Keystone endpoint from OS_AUTH_URL, falling back
// to the SysEleven default.
func identityEndpoint() string {
	if v := os.Getenv("OS_AUTH_URL"); v != "" {
		return v
	}

	return defaultIdentityEndpoint
}

func GetToken(projectID, username, password string) (string, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpoint(),
		Username:         username,
		Password:         password,
		DomainName:       defaultDomainName,
		TenantID:         projectID,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return "", fmt.Errorf("authenticate: %w", err)
	}

	return provider.Token(), nil
}

func GetTokenAppCreds(_, applicationCredentialID, applicationCredentialSecret string) (string, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint:            identityEndpoint(),
		ApplicationCredentialID:     applicationCredentialID,
		ApplicationCredentialSecret: applicationCredentialSecret,
		DomainName:                  defaultDomainName,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return "", fmt.Errorf("authenticate: %w", err)
	}

	return provider.Token(), nil
}

func GetProject(applicationCredentialID, applicationCredentialSecret string) (*tokens.Project, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint:            identityEndpoint(),
		ApplicationCredentialID:     applicationCredentialID,
		ApplicationCredentialSecret: applicationCredentialSecret,
		DomainName:                  defaultDomainName,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	r := provider.GetAuthResult()
	if r == nil {
		return nil, ErrNoAuthResult
	}

	result, ok := r.(tokens.CreateResult)
	if !ok {
		return nil, fmt.Errorf("%w: %T", ErrUnexpectedAuthResult, r)
	}

	project, err := result.ExtractProject()
	if err != nil {
		return nil, fmt.Errorf("extract project: %w", err)
	}

	return project, nil
}
