package helix

import (
	"strings"
)

var authPaths = map[string]string{
	"token":  "/token",
	"revoke": "/revoke",
}

// GetAuthorizationURL ...
func (c *Client) GetAuthorizationURL(state string, forceVerify bool) string {
	url := AuthBaseURL + "/authorize?response_type=code"
	url += "&client_id=" + c.clientID
	url += "&redirect_uri=" + c.redirectURI

	if state != "" {
		url += "&state=" + state
	}

	if forceVerify {
		url += "&force_verify=true"
	}

	if len(c.scopes) > 0 {
		url += "&scope=" + strings.Join(c.scopes, "%20")
	}

	return url
}

// AccessCredentials ...
type AccessCredentials struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scopes       []string `json:"scope"`
}

// AccessTokenResponse ...
type AccessTokenResponse struct {
	ResponseCommon
	Data AccessCredentials
}

type accessTokenRequestData struct {
	Code         string `query:"code"`
	ClientID     string `query:"client_id"`
	ClientSecret string `query:"client_secret"`
	RedirectURI  string `query:"redirect_uri"`
	GrantType    string `query:"grant_type"`
}

// GetAccessToken ...
func (c *Client) GetAccessToken(code string) (*AccessTokenResponse, error) {
	data := &accessTokenRequestData{
		Code:         code,
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		RedirectURI:  c.redirectURI,
		GrantType:    "authorization_code",
	}

	resp, err := c.post(authPaths["token"], &AccessCredentials{}, data)
	if err != nil {
		return nil, err
	}

	token := &AccessTokenResponse{}
	token.StatusCode = resp.StatusCode
	token.Error = resp.Error
	token.ErrorStatus = resp.ErrorStatus
	token.ErrorMessage = resp.ErrorMessage
	token.Data.AccessToken = resp.Data.(*AccessCredentials).AccessToken
	token.Data.RefreshToken = resp.Data.(*AccessCredentials).RefreshToken
	token.Data.ExpiresIn = resp.Data.(*AccessCredentials).ExpiresIn
	token.Data.Scopes = resp.Data.(*AccessCredentials).Scopes

	return token, nil
}

// RevokeAccessTokenResponse ...
type RevokeAccessTokenResponse struct {
	ResponseCommon
}

type revokeAccessTokenRequestData struct {
	ClientID    string `query:"client_id"`
	AccessToken string `query:"token"`
}

// RevokeAccessToken submits a request to Twitch to have an access token revoked.
//
// Both successful requests and requests with bad tokens return 200 OK with
// no body. Requests with bad tokens return the same response, as there is no
// meaningful action a client can take after sending a bad token.
func (c *Client) RevokeAccessToken(accessToken string) (*RevokeAccessTokenResponse, error) {
	data := &revokeAccessTokenRequestData{
		ClientID:    c.clientID,
		AccessToken: accessToken,
	}

	resp, err := c.post(authPaths["revoke"], nil, data)
	if err != nil {
		return nil, err
	}

	revoke := &RevokeAccessTokenResponse{}
	revoke.StatusCode = resp.StatusCode
	revoke.Error = resp.Error
	revoke.ErrorStatus = resp.ErrorStatus
	revoke.ErrorMessage = resp.ErrorMessage

	return revoke, nil
}
