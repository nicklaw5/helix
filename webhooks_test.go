package helix

import (
	"net/http"
	"reflect"
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

func newMockWebhookRequest(header string) *http.Request {
	return &http.Request{Header: map[string][]string{"Link": []string{header}}}
}

func TestGetWebhookTopicFromRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		req   *http.Request
		topic WebhookTopic
	}{
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=1234&to_id=5678>; rel=\"self\""),
			UserFollowsTopic,
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=1234>; rel=\"self\""),
			UserFollowsTopic,
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&to_id=1234>; rel=\"self\""),
			UserFollowsTopic,
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/streams?user_id=1234>; rel=\"self\""),
			StreamChangedTopic,
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users?id=1234>; rel=\"self\""),
			UserChangedTopic,
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?game_id=1a2b3c4d>; rel=\"self\""),
			GameAnalyticsTopic,
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?extension_id=1a2b3c4d>; rel=\"self\""),
			ExtensionAnalyticsTopic,
		},
		{
			newMockWebhookRequest("bad header value"),
			-1,
		},
		{
			newMockWebhookRequest(""),
			-1,
		},
		{
			&http.Request{},
			-1,
		},
	}

	for _, testCase := range testCases {
		topic := GetWebhookTopicFromRequest(testCase.req)
		if topic != testCase.topic {
			t.Errorf("expected webhook topic to be \"%d\", got \"%d\"", testCase.topic, topic)
		}
	}
}

func TestGetWebhookTopicValuesFromRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		req    *http.Request
		topic  WebhookTopic
		values map[string]string
	}{
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=1234&to_id=5678>; rel=\"self\""),
			UserFollowsTopic,
			map[string]string{"from_id": "1234", "to_id": "5678"},
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=1234>; rel=\"self\""),
			UserFollowsTopic,
			map[string]string{"from_id": "1234", "to_id": ""},
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&to_id=5678>; rel=\"self\""),
			UserFollowsTopic,
			map[string]string{"from_id": "", "to_id": "5678"},
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/streams?user_id=1234>; rel=\"self\""),
			StreamChangedTopic,
			map[string]string{"user_id": "1234"},
		},

		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users?id=1234>; rel=\"self\""),
			UserChangedTopic,
			map[string]string{"id": "1234"},
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?game_id=1a2b3c4d>; rel=\"self\""),
			GameAnalyticsTopic,
			map[string]string{"game_id": "1a2b3c4d"},
		},
		{
			newMockWebhookRequest("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?extension_id=1a2b3c4d>; rel=\"self\""),
			ExtensionAnalyticsTopic,
			map[string]string{"extension_id": "1a2b3c4d"},
		},
		{
			newMockWebhookRequest("bad header value"),
			-1,
			make(map[string]string),
		},
		{
			newMockWebhookRequest(""),
			-1,
			make(map[string]string),
		},
		{
			&http.Request{},
			-1,
			make(map[string]string),
		},
	}

	for _, testCase := range testCases {
		values := GetWebhookTopicValuesFromRequest(testCase.req, testCase.topic)
		if !reflect.DeepEqual(values, testCase.values) {
			t.Errorf("expected webhook values to be \"%v\", got \"%v\"", testCase.values, values)
		}
	}
}
