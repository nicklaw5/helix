package helix

import (
	"net/http"
	"testing"
)

func TestCreateCustomReward(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		params      *ChannelCustomRewardsParams
		respBody    string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ChannelCustomRewardsParams{
				BroadcasterID : "145328278",
				Title         : "game analysis 1v1",
				Cost          : 50000,
			},
			`{"data": [{"broadcaster_name": "torpedo09","broadcaster_login": "torpedo09","broadcaster_id": "145328278","id": "afaa7e34-6b17-49f0-a19a-d1e76eaaf673","image": null,"background_color": "#00E5CB","is_enabled": true,"cost": 50000,"title": "game analysis 1v1","prompt": "","is_user_input_required": false,"max_per_stream_setting": {"is_enabled": false,"max_per_stream": 0},"max_per_user_per_stream_setting": {"is_enabled": false,"max_per_user_per_stream": 0},"global_cooldown_setting": {"is_enabled": false,"global_cooldown_seconds": 0},"is_paused": false,"is_in_stock": true,"default_image": {"url_1x": "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x": "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x": "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue": false,"redemptions_redeemed_current_stream": null,"cooldown_expires_at": null}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ChannelCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreateCustomReward(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}

    // Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
	}

	_, err := c.CreateCustomReward(&ChannelCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
