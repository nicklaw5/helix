package helix

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var authPaths = map[string]string{
	"token":    "/token",
	"revoke":   "/revoke",
	"validate": "/validate",
	"userinfo": "/userinfo",
}

type OIDCAuth struct {
	oauth2.Token
	RawIDToken string        `json:"raw_id_token"`
	IDToken    *oidc.IDToken `json:"id_token"`
	Claims     *OIDCClaims   `json:"claims"`
}

type IDToken struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type OIDCClaims struct {
	IDToken  IDToken       `json:"id_token"`
	UserInfo UserInfoClaim `json:"user_info"`
}

type UserInfoClaim struct {
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
	Picture  string `json:"picture"`
	Updated  string `json:"updated_at"`
}

type AuthorizationURLParams struct {
	ResponseType string   // (Required) Options: "code" or "token"
	Scopes       []string // (Required)
	State        string   // (Optional)
	ForceVerify  bool     // (Optional)
}

func (c *Client) GetAuthorizationURL(params *AuthorizationURLParams) string {
	url := c.opts.AuthAPIBaseURL + "/authorize"
	url += "?response_type=" + params.ResponseType
	url += "&client_id=" + c.opts.ClientID
	url += "&redirect_uri=" + c.opts.RedirectURI

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

func (c *Client) RequestAppAccessToken(scopes []string) (*AppAccessTokenResponse, error) {
	opts := c.opts
	data := &accessTokenRequestData{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURI:  opts.RedirectURI,
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

func (c *Client) RequestUserAccessToken(code string) (*UserAccessTokenResponse, error) {
	opts := c.opts
	data := &accessTokenRequestData{
		Code:         code,
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURI:  opts.RedirectURI,
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
func (c *Client) RefreshUserAccessToken(refreshToken string) (*RefreshTokenResponse, error) {
	opts := c.opts
	data := &refreshTokenRequestData{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
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
func (c *Client) RevokeUserAccessToken(accessToken string) (*RevokeAccessTokenResponse, error) {
	data := &revokeAccessTokenRequestData{
		ClientID:    c.opts.ClientID,
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
	ClientID  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserID    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}

// ValidateToken - Validate access token
func (c *Client) ValidateToken(accessToken string) (bool, *ValidateTokenResponse, error) {
	// Reset to original token after request
	currentToken := c.opts.UserAccessToken
	c.SetUserAccessToken(accessToken)
	defer c.SetUserAccessToken(currentToken)

	var data validateTokenDetails
	resp, err := c.get(authPaths["validate"], &data, nil)
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

func (c *Client) RequestUserOIDCAccessToken(
	code string,
	scopes []string,
) (
	resp *OIDCAuth,
	err error,
) {
	resp = &OIDCAuth{}

	oauth2Config := oauth2.Config{
		ClientID:     c.opts.ClientID,
		ClientSecret: c.opts.ClientSecret,
		RedirectURL:  c.opts.RedirectURI,
		Endpoint:     c.opts.OidcProvider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: append([]string{oidc.ScopeOpenID}, scopes...),
	}

	verifier := c.opts.OidcProvider.Verifier(&oidc.Config{ClientID: c.opts.ClientID})

	oauth2Token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return
	}
	resp.Token = *oauth2Token

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		err = fmt.Errorf("id_token is MISSING validate twitch auth req")
		return
	}
	resp.RawIDToken = rawIDToken

	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		err = fmt.Errorf("failed to validate oidc token:%s err:%w", rawIDToken, err)
		return
	}
	resp.IDToken = idToken

	claims := &OIDCClaims{}
	err = idToken.Claims(&claims.IDToken)
	if err != nil {
		return
	}
	resp.Claims = claims

	return
}

func (c *Client) UserInfoFromOIDCAccessToken(
	token string,
) (
	claim UserInfoClaim,
	err error,
) {
	userInfo, err := c.opts.OidcProvider.UserInfo(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken:  token,
			TokenType:    "",
			RefreshToken: "",
			Expiry:       time.Time{},
		}),
	)
	if err != nil {
		return
	}

	// unmarshal user info claims
	err = userInfo.Claims(&claim)

	return
}
