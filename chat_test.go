package helix

import (
	"net/http"
	"testing"
)

func TestGetChannelChatBadges(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode         int
		options            *Options
		GetChatBadgeParams *GetChatBadgeParams
		respBody           string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChatBadgeParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetChatBadgeParams{BroadcasterID: "121445595"},
			`{"data": [{"set_id": "bits","versions": [{"id": "1","image_url_1x": "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/1","image_url_2x": "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/2","image_url_4x": "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/3"}]},{"set_id": "subscriber","versions": [{"id": "0","image_url_1x": "https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/1","image_url_2x": "https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/2","image_url_4x": "https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/3"}]}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChannelChatBadges(testCase.GetChatBadgeParams)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be %s, got %s", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
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

	_, err := c.GetChannelChatBadges(&GetChatBadgeParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetGlobalChatBadges(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			`{"error":"Unauthorized","status":401,"message":"OAuth token is missing"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-user-access-token"},
			`{"data": [{"set_id": "vip","versions": [{"id": "1","image_url_1x": "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/1","image_url_2x": "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/2","image_url_4x": "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/3"}]}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetGlobalChatBadges()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be %s, got %s", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "OAuth token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
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

	_, err := c.GetGlobalChatBadges()
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
