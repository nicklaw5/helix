package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestSendUserWhisper(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode            int
		options               *Options
		SendUserWhisperParams *SendUserWhisperParams
		respBody              string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendUserWhisperParams{ToUserID: "100249558", FromUserID: "100249559", Message: ""},
			`{"error":"Bad Request","status":400,"message":"The parameter \"Color\" was malformed: the value must be a valid color"}`,
		},
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendUserWhisperParams{ToUserID: "100249558", FromUserID: "100249559", Message: "hello twitch chat"},
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendUserWhisper(testCase.SendUserWhisperParams)
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
		ctx:  context.Background(),
	}

	_, err := c.SendUserWhisper(&SendUserWhisperParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
