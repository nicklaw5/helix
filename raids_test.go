package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestStartRaid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode        int
		options           *Options
		FromBroadcasterID string
		ToBroadcasterID   string
		respBody          string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"237509153",
			"237509153",
			`{"data":[{"created_at": "2022-02-18T07:20:50.52Z","is_mature": false}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"237509153",
			"237509153",
			`{"error":"Bad Request","status":400,"message":"The IDs in \"from_broadcaster_id\" and \"to_broadcaster_id\" cannot be the same ID."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.StartRaid(&StartRaidParams{
			FromBroadcasterID: testCase.FromBroadcasterID,
			ToBroadcasterID:   testCase.ToBroadcasterID,
		})
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

			expectedErrMsg := "The IDs in \"from_broadcaster_id\" and \"to_broadcaster_id\" cannot be the same ID."
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

	_, err := c.StartRaid(&StartRaidParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestCancelRaid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		respBody      string
	}{
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id"},
			"237509153",
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"237509153",
			`{"error":"Bad Request","status":400,"message":"The ID in the \"broadcaster_id query\" parameter is not valid."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CancelRaid(&CancelRaidParams{
			BroadcasterID: testCase.BroadcasterID,
		})
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

			expectedErrMsg := "The ID in the \"broadcaster_id query\" parameter is not valid."
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

	_, err := c.CancelRaid(&CancelRaidParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
