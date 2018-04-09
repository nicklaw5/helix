package helix

import (
	"net/http"
	"testing"
)

func TestGetAuthorizationURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		state       string
		forceVerify bool
		options     *Options
		expectedURL string
	}{
		{
			"",
			false,
			&Options{
				ClientID:    "my-client-id",
				RedirectURI: "https://example.com/auth/callback",
			},
			"https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=my-client-id&redirect_uri=https://example.com/auth/callback",
		},
		{
			"some-state",
			true,
			&Options{
				ClientID:    "my-client-id",
				RedirectURI: "https://example.com/auth/callback",
				Scopes:      []string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"},
			},
			"https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=my-client-id&redirect_uri=https://example.com/auth/callback&state=some-state&force_verify=true&scope=analytics:read:games%20bits:read%20clips:edit%20user:edit%20user:read:email",
		},
	}

	for _, testCase := range testCases {

		client, err := NewClient(testCase.options)
		if err != nil {
			t.Errorf("Did not expect an error, got \"%s\"", err.Error())
		}

		url := client.GetAuthorizationURL(testCase.state, testCase.forceVerify)

		if url != testCase.expectedURL {
			t.Errorf("expected url to be \"%s\", got \"%s\"", testCase.expectedURL, url)
		}
	}
}

func TestGetAppAccessToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusBadRequest,
			&Options{
				ClientID:     "invalid-client-id", // invalid client id
				ClientSecret: "valid-client-secret",
			},
			`{"status":400,"message":"invalid client"}`,
			"invalid client",
		},
		{
			http.StatusForbidden,
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "invalid-client-secret", // invalid client secret
			},
			`{"status":403,"message":"invalid client secret"}`,
			"invalid client secret",
		},
		{
			http.StatusOK,
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
			},
			`{"access_token":"ajsdfloehfoihsdfhoasjfdpoiqh","expires_in":4999199}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetAppAccessToken()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if resp.StatusCode != http.StatusOK {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		// Test success cases
		if resp.Data.AccessToken == "" {
			t.Errorf("expected an access token but got an empty string")
		}

		if resp.Data.ExpiresIn == 0 {
			t.Errorf("expected ExpiresIn to not be \"0\"")
		}
	}
}

func TestGetUserAccessToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		code           string
		scopes         []string
		options        *Options
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusBadRequest,
			"invalid-auth-code", // invalid auth code
			[]string{"user:read:email"},
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
				RedirectURI:  "https://example.com/auth/callback",
			},
			`{"status":400,"message":"Invalid authorization code"}`,
			"Invalid authorization code",
		},
		{
			http.StatusBadRequest,
			"valid-auth-code",
			[]string{"user:read:email"},
			&Options{
				ClientID:     "invalid-client-id", // invalid client id
				ClientSecret: "valid-client-secret",
				RedirectURI:  "https://example.com/auth/callback",
			},
			`{"status":400,"message":"invalid client"}`,
			"invalid client",
		},
		{
			http.StatusForbidden,
			"valid-auth-code",
			[]string{"user:read:email"},
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "invalid-client-secret", // invalid client secret
				RedirectURI:  "https://example.com/auth/callback",
			},
			`{"status":403,"message":"invalid client secret"}`,
			"invalid client secret",
		},
		{
			http.StatusBadRequest,
			"valid-auth-code",
			[]string{"user:read:email"},
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
				RedirectURI:  "https://example.com/invalid/callback", // invalid redirect uri
			},
			`{"status":400,"message":"Parameter redirect_uri does not match registeredURI"}`,
			"Parameter redirect_uri does not match registeredURI",
		},
		{
			http.StatusOK,
			"valid-auth-code",
			[]string{}, // no scopes
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
				RedirectURI:  "https://example.com/auth/callback",
			},
			`{"access_token":"kagsfkgiuowegfkjsbdcuiwebf","expires_in":14146,"refresh_token":"fiuhgaofohofhohdflhoiwephvlhowiehfoi"}`,
			"",
		},
		{
			http.StatusOK,
			"valid-auth-code",
			[]string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"},
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
				RedirectURI:  "https://example.com/auth/callback",
			},
			`{"access_token":"kagsfkgiuowegfkjsbdcuiwebf","expires_in":14154,"refresh_token":"fiuhgaofohofhohdflhoiwephvlhowiehfoi","scope":["analytics:read:games","bits:read","clips:edit","user:edit","user:read:email"]}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUserAccessToken(testCase.code)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if resp.StatusCode != http.StatusOK {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		// Test success cases
		if resp.Data.AccessToken == "" {
			t.Errorf("expected an access token but got an empty string")
		}

		if resp.Data.RefreshToken == "" {
			t.Errorf("expected a refresh token but got an empty string")
		}

		if resp.Data.ExpiresIn == 0 {
			t.Errorf("expected ExpiresIn to not be \"0\"")
		}

		if len(resp.Data.Scopes) != len(testCase.scopes) {
			t.Errorf("expected number of scope to be \"%d\", got \"%d\"", len(testCase.scopes), len(resp.Data.Scopes))
		}
	}
}

func TestRefreshUserAccessToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		refreshToken   string
		options        *Options
		respBody       string
		expectedErrMsg string
		expectedScopes []string
	}{
		{
			http.StatusBadRequest,
			"", // no refresh token
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
			},
			`{"status":400,"message":"missing refresh token"}`,
			"missing refresh token",
			[]string{},
		},
		{
			http.StatusBadRequest,
			"invalid-refresh-token", // invalid refresh token
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
			},
			`{"status":400,"message":"Invalid refresh token"}`,
			"Invalid refresh token",
			[]string{},
		},
		{
			http.StatusBadRequest,
			"valid-refresh-token",
			&Options{
				ClientID:     "invalid-client-id", // invalid client id
				ClientSecret: "valid-client-secret",
			},
			`{"status":400,"message":"invalid client"}`,
			"invalid client",
			[]string{},
		},
		{
			http.StatusForbidden,
			"valid-refresh-token",
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "invalid-client-secret", // invalid client secret
			},
			`{"status":403,"message":"invalid client secret"}`,
			"invalid client secret",
			[]string{},
		},
		{
			http.StatusBadRequest,
			"valid-refresh-token",
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
			},
			`{"status":400,"message":"invalid scope requested: 'invalid:scope'"}`,
			"invalid scope requested: 'invalid:scope'",
			[]string{},
		},
		{
			http.StatusOK,
			"valid-refresh-token",
			&Options{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
			},
			`{"access_token":"oihhkfhsajkhfjksahfkjahsf","expires_in":13669,"refresh_token":"oihhkfhsajkhfjksahfkjahsfahsldhasld","scope":["analytics:read:games","bits:read","clips:edit","user:edit","user:read:email"]}`,
			"",
			[]string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RefreshUserAccessToken(testCase.refreshToken)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if resp.StatusCode != http.StatusOK {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		// // Test success cases
		if resp.Data.AccessToken == "" {
			t.Errorf("expected an access token but got an empty string")
		}

		if resp.Data.RefreshToken == "" {
			t.Errorf("expected a refresh token but got an empty string")
		}

		if resp.Data.ExpiresIn == 0 {
			t.Errorf("expected ExpiresIn to not be \"0\"")
		}

		if len(resp.Data.Scopes) != len(testCase.expectedScopes) {
			t.Errorf("expected number of scope to be \"%d\", got \"%d\"", len(testCase.expectedScopes), len(resp.Data.Scopes))
		}
	}
}

func TestRevokeUserAccessToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		accessToken    string
		options        *Options
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusBadRequest,
			"valid-access-token",
			&Options{ClientID: "invalid-client-id"}, // invalid client id
			`{"status":400,"message":"Invalid client_id: invalid-client-id"}`,
			"Invalid client_id: invalid-client-id",
		},
		{
			http.StatusBadRequest,
			"", // no access token
			&Options{ClientID: "valid-client-id"},
			`{"status":400,"message":"missing oauth token"}`,
			"missing oauth token",
		},
		{
			http.StatusOK,
			"invalid-access-token", // invalid token still returns 200 OK response
			&Options{ClientID: "valid-client-id"},
			"",
			"",
		},
		{
			http.StatusOK,
			"valid-access-token",
			&Options{ClientID: "valid-client-id"},
			"",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RevokeUserAccessToken(testCase.accessToken)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if resp.StatusCode != http.StatusOK {
			if testCase.expectedErrMsg != "" && resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			continue
		}
	}
}
