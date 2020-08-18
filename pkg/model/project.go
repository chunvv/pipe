// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"crypto/subtle"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var (
	githubScopes = []string{"read:org"}
)

// RedactSensitiveData redacts sensitive data.
func (p *Project) RedactSensitiveData() {
	if p.StaticAdmin != nil {
		p.StaticAdmin.RedactSensitiveData()
	}
	if p.Sso != nil {
		p.Sso.RedactSensitiveData()
	}
}

// RedactSensitiveData redacts sensitive data.
func (p *ProjectStaticUser) RedactSensitiveData() {
	p.PasswordHash = redactedMessage
}

// Update updates ProjectStaticUser with given data.
func (p *ProjectStaticUser) Update(username, password string) error {
	if username != "" {
		p.Username = username
	}
	if password != "" {
		encoded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		p.PasswordHash = string(encoded)
	}
	return nil
}

// Auth confirms username and password.
func (p *ProjectStaticUser) Auth(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if subtle.ConstantTimeCompare([]byte(p.Username), []byte(username)) != 1 {
		return fmt.Errorf("wrong username %q", username)
	}
	if password == "" {
		return fmt.Errorf("password is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("wrong password for username %q: %v", username, err)
	}
	return nil
}

// RedactSensitiveData redacts sensitive data.
func (p *ProjectSingleSignOn) RedactSensitiveData() {
	if p.Github != nil {
		p.Github.RedactSensitiveData()
	}
	if p.Google != nil {
	}
}

// Update updates ProjectSingleSignOn with given data.
func (p *ProjectSingleSignOn) Update(sso *ProjectSingleSignOn) {
	p.Provider = sso.Provider
	if sso.Github != nil {
		if p.Github == nil {
			p.Github = &ProjectSingleSignOn_GitHub{}
		}
		p.Github.Update(sso.Github)
	}
	if sso.Google != nil {
	}
}

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSingleSignOn) GenerateAuthCodeURL(project, apiURL, callbackPath, state string) (string, error) {
	switch p.Provider {
	case ProjectSingleSignOnProvider_GITHUB:
		if p.Github == nil {
			return "", fmt.Errorf("missing GitHub oauth in the SSO configuration")
		}
		return p.Github.GenerateAuthCodeURL(project, apiURL, callbackPath, state)
	default:
		return "", fmt.Errorf("not implemented")
	}
}

// RedactSensitiveData redacts sensitive data.
func (p *ProjectSingleSignOn_GitHub) RedactSensitiveData() {
	p.ClientId = redactedMessage
	p.ClientSecret = redactedMessage
}

// Update updates ProjectSingleSignOn with given data.
func (p *ProjectSingleSignOn_GitHub) Update(input *ProjectSingleSignOn_GitHub) {
	if input.ClientId != "" {
		p.ClientId = input.ClientId
	}
	if input.ClientSecret != "" {
		p.ClientSecret = input.ClientSecret
	}
	if input.BaseUrl != "" {
		p.BaseUrl = input.ClientSecret
	}
	if input.UploadUrl != "" {
		p.UploadUrl = input.UploadUrl
	}
	if input.Org != "" {
		p.Org = input.Org
	}
	if input.AdminTeam != "" {
		p.AdminTeam = input.AdminTeam
	}
	if input.EditorTeam != "" {
		p.EditorTeam = input.EditorTeam
	}
	if input.ViewerTeam != "" {
		p.ViewerTeam = input.ViewerTeam
	}
}

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSingleSignOn_GitHub) GenerateAuthCodeURL(project, apiURL, callbackPath, state string) (string, error) {
	u, err := url.Parse(p.BaseUrl)
	if err != nil {
		return "", err
	}

	cfg := oauth2.Config{
		ClientID: p.ClientId,
		Endpoint: oauth2.Endpoint{AuthURL: fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, "/login/oauth/authorize")},
	}

	cfg.Scopes = githubScopes
	apiURL = strings.TrimSuffix(apiURL, "/")
	cfg.RedirectURL = fmt.Sprintf("%s%s?project=%s", apiURL, callbackPath, project)
	authURL := cfg.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOnline)

	return authURL, nil
}