package helix

import (
	"strings"
)

var authPaths = map[string]string{
	"token":  "/token",
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
	token.RateLimit.Limit = resp.RateLimit.Limit
	token.RateLimit.Remaining = resp.RateLimit.Remaining
	token.RateLimit.Reset = resp.RateLimit.Reset
	token.Data.AccessToken = resp.Data.(*AccessCredentials).AccessToken
	token.Data.RefreshToken = resp.Data.(*AccessCredentials).RefreshToken
	token.Data.ExpiresIn = resp.Data.(*AccessCredentials).ExpiresIn
	token.Data.Scopes = resp.Data.(*AccessCredentials).Scopes

	return token, nil
}
