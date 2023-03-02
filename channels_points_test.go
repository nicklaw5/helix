package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateCustomReward(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *ChannelCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ChannelCustomRewardsParams{
				BroadcasterID: "145328278",
				Title:         "game analysis 1v1",
				Cost:          50000,
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
		ctx:  context.Background(),
	}

	_, err := c.CreateCustomReward(&ChannelCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateCustomReward(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *UpdateChannelCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsParams{
				ID:            "123",
				BroadcasterID: "145328278",
				Title:         "game analysis 1v1",
				Cost:          50000,
			},
			`{"data": [{"broadcaster_name": "torpedo09","broadcaster_login": "torpedo09","broadcaster_id": "145328278","id": "afaa7e34-6b17-49f0-a19a-d1e76eaaf673","image": null,"background_color": "#00E5CB","is_enabled": true,"cost": 50000,"title": "game analysis 1v1","prompt": "","is_user_input_required": false,"max_per_stream_setting": {"is_enabled": false,"max_per_stream": 0},"max_per_user_per_stream_setting": {"is_enabled": false,"max_per_user_per_stream": 0},"global_cooldown_setting": {"is_enabled": false,"global_cooldown_seconds": 0},"is_paused": false,"is_in_stock": true,"default_image": {"url_1x": "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x": "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x": "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue": false,"redemptions_redeemed_current_stream": null,"cooldown_expires_at": null}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateCustomReward(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"id\""
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
		ctx:  context.Background(),
	}

	_, err := c.UpdateCustomReward(&UpdateChannelCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestDeleteCustomRewards(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *DeleteCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&DeleteCustomRewardsParams{
				BroadcasterID: "145328278",
				ID:            "84da6b13-efe1-4a82-91d0-25260aeb6a9b",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&DeleteCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.DeleteCustomRewards(testCase.params)
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
		ctx:  context.Background(),
	}

	_, err := c.DeleteCustomRewards(&DeleteCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetCustomRewards(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetCustomRewardsParams{
				BroadcasterID: "145328278",
			},
			`{"data":[{"broadcaster_name":"Scorfly","broadcaster_login":"scorfly","broadcaster_id":"145328278","id":"7e86e512-f423-45b2-9e4b-2ad99a814ffc","image":null,"background_color":"#421561","is_enabled":false,"cost":50000,"title":"game analysis 1v1","prompt":"","is_user_input_required":false,"max_per_stream_setting":{"is_enabled":false,"max_per_stream":0},"max_per_user_per_stream_setting":{"is_enabled":false,"max_per_user_per_stream":0},"global_cooldown_setting":{"is_enabled":false,"global_cooldown_seconds":0},"is_paused":false,"is_in_stock":true,"default_image":{"url_1x":"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x":"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x":"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue":false,"redemptions_redeemed_current_stream":null,"cooldown_expires_at":null},{"broadcaster_name":"Scorfly","broadcaster_login":"scorfly","broadcaster_id":"145328278","id":"c0d6687c-c543-4a3a-8494-502451e7fa45","image":null,"background_color":"#FFBF00","is_enabled":false,"cost":50000,"title":"game analysis 1v2","prompt":"","is_user_input_required":false,"max_per_stream_setting":{"is_enabled":false,"max_per_stream":0},"max_per_user_per_stream_setting":{"is_enabled":false,"max_per_user_per_stream":0},"global_cooldown_setting":{"is_enabled":false,"global_cooldown_seconds":0},"is_paused":false,"is_in_stock":true,"default_image":{"url_1x":"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x":"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x":"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue":false,"redemptions_redeemed_current_stream":null,"cooldown_expires_at":null}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetCustomRewards(testCase.params)
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
		ctx:  context.Background(),
	}

	_, err := c.GetCustomRewards(&GetCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateCustomRewardRedemptionStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *UpdateChannelCustomRewardsRedemptionStatusParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsRedemptionStatusParams{
				ID:            "17fa2df1-ad76-4804-bfa5-a40ef63efe63",
				BroadcasterID: "274637212",
				RewardID:      "92af127c-7326-4483-a52b-b0da0be61c01",
				Status:        "CANCELLED",
			},
			`{"data": [{"broadcaster_name": "torpedo09", "broadcaster_login": "torpedo09", "broadcaster_id": "274637212", "id": "17fa2df1-ad76-4804-bfa5-a40ef63efe63", "user_id": "274637212", "user_name": "torpedo09", "user_login": "torpedo09", "user_input": "", "status": "CANCELED", "redeemed_at": "2020-07-01T18:37:32Z", "reward": { "id": "92af127c-7326-4483-a52b-b0da0be61c01", "title": "game analysis", "prompt": "", "cost": 50000}}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsRedemptionStatusParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateChannelCustomRewardsRedemptionStatus(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if testCase.statusCode == http.StatusOK {
			numRedemptions := len(resp.Data.Redemptions)
			if numRedemptions != 1 {
				t.Errorf("expected 1 redemption, got %d", numRedemptions)
				continue
			}

			title := resp.Data.Redemptions[0].Reward.Title
			if title != "game analysis" {
				t.Errorf("expected reward title to be \"game analysis\", got \"%s\"", title)
			}
		}
	}
}
