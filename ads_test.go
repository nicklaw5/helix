package helix

import (
	"net/http"
	"testing"
)

func TestClient_StartCommercial(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		count          int
		broadcasterID  string
		adLength       AdLengthEnum
		respBody       string
		expectedErrMsg string
	}{
		// TODO: expand with other test cases

		// Failure - broadcaster not live
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			1,
			"89754791",
			AdLen30,
			`{"error":"Bad Request","status":400,"message":"the channel 'codingsloth' is not currently live and needs to be in order to start commercials."}`,
			"the channel 'codingsloth' is not currently live and needs to be in order to start commercials.",
		},
		// success
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			1,
			"41245072",
			AdLen60,
			`{"data":[{"length":60,"message":"","retry_after":480}]}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		params := &StartCommercialParams{
			BroadcasterID: testCase.broadcasterID,
			Length:        testCase.adLength,
		}

		resp, err := c.StartCommercial(params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if resp.StatusCode != http.StatusOK {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		// Test success cases
		if len(resp.Data.AdDetails) != testCase.count {
			t.Errorf("expected number of results to be \"%d\", got \"%d\"", testCase.count, len(resp.Data.AdDetails))
		}

		commercialDetails := resp.Data.AdDetails[0]
		if commercialDetails.Length != testCase.adLength {
			t.Error("expected ad length \"#{testCase.adLength}\", got \"#{commercialDetails.Length}\"")
		}

		if commercialDetails.Message != "" {
			t.Error("expected an empty error message, got \"#{commercialDetails.Message}\"")
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

	_, err := c.StartCommercial(&StartCommercialParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
