package helix

import (
	"net/http"
	"testing"
)

func TestGetCreatorGoals(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		respBody      string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"123",
			`{ "data": [ { "id": "1woowvbkiNv8BRxEWSqmQz6Zk92", "broadcaster_id": "141981764", "broadcaster_name": "TwitchDev", "broadcaster_login": "twitchdev", "type": "follower", "description": "Follow goal for Helix testing", "current_amount": 27062, "target_amount": 30000, "created_at": "2021-08-16T17:22:23Z" } ] }`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetCreatorGoals(&GetCreatorGoalsParams{
			BroadcasterID: testCase.BroadcasterID,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Goals) != 1 {
			t.Errorf("expected %d goals got %d", 1, len(resp.Data.Goals))
		}
	}
}
