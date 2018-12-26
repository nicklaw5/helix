package helix

import (
	"net/http"
	"testing"
)

func TestGetWebhookSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		First      int
		respBody   string
	}{
		{
			http.StatusUnauthorized,
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

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusUnauthorized, resp.ErrorStatus)
			}

			expectedErrMsg := "Must provide valid app token."
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

func TestPostWebhookSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusAccepted,
			&Options{ClientID: "my-client-id", AppAccessToken: "valid-app-access-token"},
			"",
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id", AppAccessToken: "invalid-app-access-token"},
			`{"error":"Unauthorized","status":401,"message":"Must provide valid app token."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.PostWebhookSubscription(&WebhookSubscriptionPayload{
			Callback:     "https://my.callback/url",
			Mode:         "subscribe",
			Topic:        "https://api.twitch.tv/helix/users/follows?first=1&from_id=1111&to_id=2222",
			LeaseSeconds: 0,
			Secret:       "53cr3t",
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusUnauthorized, resp.ErrorStatus)
			}

			expectedErrMsg := "Must provide valid app token."
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}
