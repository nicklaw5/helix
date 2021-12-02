package helix

import (
	"net/http"
	"strings"
)

var authPaths = map[string]string{
	"token":    "/token",
	"revoke":   "/revoke",
	"validate": "/validate",
}

type AuthorizationURLParams struct {
	ResponseType string   // (Required) Options: "code" or "token"
	Scopes       []string // (Required)
	State        string   // (Optional)
	ForceVerify  bool     // (Optional)
}

func (c *Client) GetAuthorizationURL(params *AuthorizationURLParams, opts ...Options) string {
	var options Options
	if len(opts) == 0 {
		options = *c.opts
	} else {
		options = opts[0]
	}

	url := AuthBaseURL + "/authorize"
	url += "?response_type=" + params.ResponseType
	url += "&client_id=" + options.ClientID
	url += "&redirect_uri=" + options.RedirectURI

	if params.State != "" {
		url += "&state=" + params.State
	}

	if params.ForceVerify {
		url += "&force_verify=true"
	}

	if len(params.Scopes) != 0 {
		url += "&scope=" + strings.Join(params.Scopes, "%20")
	}

	return url
}

type AccessCredentials struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scopes       []string `json:"scope"`
}

type AppAccessTokenResponse struct {
	ResponseCommon
	Data AccessCredentials
}

func (c *Client) RequestAppAccessToken(scopes []string, opts ...Options) (*AppAccessTokenResponse, error) {
	var options Options
	if len(opts) == 0 {
		options = *c.opts
	} else {
		options = opts[0]
	}

	data := &accessTokenRequestData{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		RedirectURI:  options.RedirectURI,
		GrantType:    "client_credentials",
		Scopes:       strings.Join(scopes, " "),
	}

	resp, err := c.post(authPaths["token"], &AccessCredentials{}, data)
	if err != nil {
		return nil, err
	}

	token := &AppAccessTokenResponse{}
	resp.HydrateResponseCommon(&token.ResponseCommon)
	token.Data.AccessToken = resp.Data.(*AccessCredentials).AccessToken
	token.Data.RefreshToken = resp.Data.(*AccessCredentials).RefreshToken
	token.Data.ExpiresIn = resp.Data.(*AccessCredentials).ExpiresIn
	token.Data.Scopes = resp.Data.(*AccessCredentials).Scopes

	return token, nil
}

type UserAccessTokenResponse struct {
	ResponseCommon
	Data AccessCredentials
}

type accessTokenRequestData struct {
	Code         string `query:"code"`
	ClientID     string `query:"client_id"`
	ClientSecret string `query:"client_secret"`
	RedirectURI  string `query:"redirect_uri"`
	GrantType    string `query:"grant_type"`
	Scopes       string `query:"scope"`
}

func (c *Client) RequestUserAccessToken(code string, opts ...Options) (*UserAccessTokenResponse, error) {
	var options Options
	if len(opts) == 0 {
		options = *c.opts
	} else {
		options = opts[0]
	}
	data := &accessTokenRequestData{
		Code:         code,
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		RedirectURI:  options.RedirectURI,
		GrantType:    "authorization_code",
	}

	resp, err := c.post(authPaths["token"], &AccessCredentials{}, data)
	if err != nil {
		return nil, err
	}

	token := &UserAccessTokenResponse{}
	resp.HydrateResponseCommon(&token.ResponseCommon)
	token.Data.AccessToken = resp.Data.(*AccessCredentials).AccessToken
	token.Data.RefreshToken = resp.Data.(*AccessCredentials).RefreshToken
	token.Data.ExpiresIn = resp.Data.(*AccessCredentials).ExpiresIn
	token.Data.Scopes = resp.Data.(*AccessCredentials).Scopes

	return token, nil
}

type RefreshTokenResponse struct {
	ResponseCommon
	Data AccessCredentials
}

type refreshTokenRequestData struct {
	ClientID     string `query:"client_id"`
	ClientSecret string `query:"client_secret"`
	GrantType    string `query:"grant_type"`
	RefreshToken string `query:"refresh_token"`
}

// RefreshUserAccessToken submits a request to have the longevity of an
// access token extended. Twitch OAuth2 access tokens have expirations.
// Token-expiration periods vary in length. You should build your applications
// in such a way that they are resilient to token authentication failures.
func (c *Client) RefreshUserAccessToken(refreshToken string, opts ...Options) (*RefreshTokenResponse, error) {
	var options Options
	if len(opts) == 0 {
		options = *c.opts
	} else {
		options = opts[0]
	}
	data := &refreshTokenRequestData{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	resp, err := c.post(authPaths["token"], &AccessCredentials{}, data)
	if err != nil {
		return nil, err
	}

	refresh := &RefreshTokenResponse{}
	resp.HydrateResponseCommon(&refresh.ResponseCommon)
	refresh.Data.AccessToken = resp.Data.(*AccessCredentials).AccessToken
	refresh.Data.RefreshToken = resp.Data.(*AccessCredentials).RefreshToken
	refresh.Data.ExpiresIn = resp.Data.(*AccessCredentials).ExpiresIn
	refresh.Data.Scopes = resp.Data.(*AccessCredentials).Scopes

	return refresh, nil
}

type RevokeAccessTokenResponse struct {
	ResponseCommon
}

type revokeAccessTokenRequestData struct {
	ClientID    string `query:"client_id"`
	AccessToken string `query:"token"`
}

// RevokeUserAccessToken submits a request to Twitch to have an access token revoked.
//
// Both successful requests and requests with bad tokens return 200 OK with
// no body. Requests with bad tokens return the same response, as there is no
// meaningful action a client can take after sending a bad token.
func (c *Client) RevokeUserAccessToken(accessToken string, opts ...Options) (*RevokeAccessTokenResponse, error) {
	var options Options
	if len(opts) == 0 {
		options = *c.opts
	} else {
		options = opts[0]
	}
	data := &revokeAccessTokenRequestData{
		ClientID:    options.ClientID,
		AccessToken: accessToken,
	}

	resp, err := c.post(authPaths["revoke"], nil, data)
	if err != nil {
		return nil, err
	}

	revoke := &RevokeAccessTokenResponse{}
	resp.HydrateResponseCommon(&revoke.ResponseCommon)

	return revoke, nil
}

type ValidateTokenResponse struct {
	ResponseCommon
	Data validateTokenDetails
}

type validateTokenDetails struct {
	ClientID string   `json:"client_id"`
	Login    string   `json:"login"`
	Scopes   []string `json:"scopes"`
	UserID   string   `json:"user_id"`
}

// ValidateToken - Validate access token
func (c *Client) ValidateToken(accessToken string) (bool, *ValidateTokenResponse, error) {
	opts := Options{
		UserAccessToken: accessToken,
	}

	var data validateTokenDetails
	resp, err := c.get(authPaths["validate"], &data, nil, opts)
	if err != nil {
		return false, nil, err
	}

	var isValid bool
	if resp.StatusCode == http.StatusOK {
		isValid = true
	}

	tokenResp := &ValidateTokenResponse{
		Data: data,
	}
	resp.HydrateResponseCommon(&tokenResp.ResponseCommon)

	return isValid, tokenResp, nil
}
