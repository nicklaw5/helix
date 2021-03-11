package helix

import (
	"net/http"
	"testing"
)

func TestGetEventSubSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *EventSubSubscriptionsParams
		count		int
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&EventSubSubscriptionsParams{},
			0,
			`{"error":"Unauthorized","status":401,"message":"OAuth token is missing"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&EventSubSubscriptionsParams{},
			2,
			`{"total":2,"data":[{"id":"832389eb-0d0b-41f8-b564-da039f6c4c75","status":"enabled","type":"channel.follow","version":"1","condition":{"broadcaster_user_id":"12345678"},"created_at":"2021-03-09T10:37:32.308415339Z","transport":{"method":"webhook","callback":"https://example.com/eventsub/follow"},"cost":1},{"id":"832389eb-0d0b-41f8-b564-da039f6c4c73","status":"enabled","type":"channel.follow","version":"1","condition":{"broadcaster_user_id":"12345679"},"created_at":"2021-03-09T10:37:32.308415339Z","transport":{"method":"webhook","callback":"https://example.com/eventsub/follow"},"cost":1}],"limit":100000000,"max_total_cost":10000,"total_cost":2,"pagination":{}}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&EventSubSubscriptionsParams{Status: "webhook_callback_verification_failed"},
			0,
			`{"total":1,"data":[],"limit":100000000,"max_total_cost":10000,"total_cost":1,"pagination":{}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetEventSubSubscriptions(testCase.params)
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

			expectedErrMsg := "OAuth token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.EventSubSubscriptions) != testCase.count {
			t.Errorf("expected result length to be \"%d\", got \"%d\"", testCase.count, len(resp.Data.EventSubSubscriptions))
		}
	}
}

func TestRemoveEventSubSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     string
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			"",
			`{"error":"Unauthorized","status":401,"message":"OAuth token is missing"}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"id\""}`,
		},
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id"},
			"832389eb-0d0b-41f8-b564-da039f6c4c75",
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RemoveEventSubSubscription(testCase.params)
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

			expectedErrMsg := "OAuth token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			expectedErrMsg := `Missing required parameter "id"`
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}

func TestCreateEventSubSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *EventSubSubscription
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&EventSubSubscription{},
			`{"error":"Unauthorized","status":401,"message":"OAuth token is missing"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&EventSubSubscription{
				Type:      "channel.follow",
				Version:   "1",
				Condition: EventSubCondition{
					BroadcasterUserID:     "12345678",
				},
				Transport: EventSubTransport{
					Method:   "webhook",
					Callback: "https://example.com/eventsub/follow",
					Secret:   "s3cr37w0rd",
				},
			},
			`{"data":[{"id":"4d06fabc-4cf4-4e99-a60f-b457d5c69305","status":"webhook_callback_verification_pending","type":"channel.follow","version":"1","condition":{"broadcaster_user_id":"12345678"},"created_at":"2021-03-10T23:38:50.311154721Z","transport":{"method":"webhook","callback":"https://example.com/eventsub/follow"},"cost":1}],"limit":10000,"total":1,"max_total_cost":10000,"total_cost":1}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreateEventSubSubscription(testCase.params)
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

			expectedErrMsg := "OAuth token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
		if len(resp.Data.EventSubSubscriptions) != 1 {
			t.Errorf("expected result length to be \"%d\", got \"%d\"", 1, len(resp.Data.EventSubSubscriptions))

			continue
		}
		if resp.Data.EventSubSubscriptions[0].Transport.Method != "webhook" {
			t.Errorf("expected result transport method to be \"%s\", got \"%s\"", "webhook", resp.Data.EventSubSubscriptions[0].Transport.Method)
		}
	}
}

func TestVerifyEventSubNotification(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		messageID	string
		messageSignature	string
		messageTimestamp	string
		respBody   string
		secret	string
	}{
		{
			"e76c6bd4-55c9-4987-8304-da1588d8988b",
			"sha256=7e5a96480c29cdf834b371e7a5b049638cba6e425ea51b9b2a9fabf69bc5d227",
			"2019-11-16T10:11:12.123Z",
			`{"challenge":"pogchamp-kappa-360noscope-vohiyo","subscription":{"id":"f1c2a387-161a-49f9-a165-0f21d7a4e1c4","status":"webhook_callback_verification_pending","type":"channel.follow","version":"1","condition":{"broadcaster_user_id":"12826"},"transport":{"method":"webhook","callback":"https://example.com/webhooks/callback"},"created_at":"2019-11-16T10:11:12.123Z"}}`,
			"s3cRe7",
		},
	}

	for _, testCase := range testCases {
		header := http.Header{}
		header.Add("Twitch-Eventsub-Message-Id", testCase.messageID)
		header.Add("Twitch-Eventsub-Message-Signature", testCase.messageSignature)
		header.Add("Twitch-Eventsub-Message-Timestamp", testCase.messageTimestamp)
		signatureOk := VerifyEventSubNotification(testCase.secret, header, testCase.respBody)
		if !signatureOk {
			t.Errorf("expected signature to match \"%s\", but it didn't", testCase.messageSignature)
		}
	}
}