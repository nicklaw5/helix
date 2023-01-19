package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestSendShoutout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode         int
		options            *Options
		SendShoutoutParams *SendShoutoutParams
		respBody           string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendShoutoutParams{FromBroadcasterID: "100249558", ModeratorID: "100249558"},
			`{"error":"Bad Request","status":400,"message":"The parameter \"to_broadcaster_id\" was malformed: the value must be a valid"}`,
		},
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendShoutoutParams{FromBroadcasterID: "100249558", ModeratorID: "100249558", ToBroadcasterID: "80085"},
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendShoutout(testCase.SendShoutoutParams)
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

			expectedErrMsg := "The parameter \"to_broadcaster_id\" was malformed: the value must be a valid"
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
		ctx:  context.Background(),
	}

	_, err := c.SendShoutout(&SendShoutoutParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
