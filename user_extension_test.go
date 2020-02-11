package helix

import (
	"net/http"
	"testing"
)

func TestGetUserExtensions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			`{"data":[{"id":"wi08ebtatdc7oj83wtl9uxwz807l8b","version":"1.1.8","name":"Streamlabs Leaderboard","can_activate":true,"type":["panel"]},{"id":"d4uvtfdr04uq6raoenvj7m86gdk16v","version":"2.0.2","name":"Prime Subscription and Loot Reminder","can_activate":true,"type":["overlay"]}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUserExtensions()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode == http.StatusOK && len(resp.Data.UserExtensions) == 0 {
			t.Error("failed to parse successful UserExtension response data")
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}
}
