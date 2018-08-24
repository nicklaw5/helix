package helix

import (
	"net/http"
	"testing"
)

func GetWebhookSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		First      int
		respBody   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			2,
			`{"error":"Unauthorized","status":401,"message":"Must provide valid app token."}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			`{"total":2,"data":[{"topic":"https://api.twitch.tv/helix/streams?user_id=123","callback":"http://example.com/your_callback","expires_at":"2018-07-30T20:00:00Z"},{"topic":"https://api.twitch.tv/helix/streams?user_id=345","callback":"http://example.com/your_callback","expires_at":"2018-07-30T20:03:00Z"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IkFYc2laU0k2TVN3aWFTSTZNWDAifX0"}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetWebhookSubscriptions(&WebhookSubscriptionsParams{
			First: testCase.First,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Must provide either from_id or to_id"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.WebhookSubscriptions) != testCase.First {
			t.Errorf("expected result length to be \"%d\", got \"%d\"", testCase.First, len(resp.Data.WebhookSubscriptions))
		}
	}
}
