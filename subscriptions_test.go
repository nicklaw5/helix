package helix

import (
	"net/http"
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		UserID        []string
		respBody      string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"123",
			[]string{},
			`{"data": [{"broadcaster_id":"123","broadcaster_name":"test_user","is_gift":true,"tier":"3000","plan_name":"The Ninjas","user_id":"123","user_name":"test_user"}],"pagination":{"cursor":"xxxx"}}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"",
			[]string{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetSubscriptions(&SubscriptionsParams{
			BroadcasterID: testCase.BroadcasterID,
			UserID:        testCase.UserID,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if testCase.statusCode != http.StatusOK {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			errMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != errMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", errMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}
