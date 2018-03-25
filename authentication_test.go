package helix

import (
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
