package helix

import (
	"net/http"
	"testing"
)

func TestGetCharityCampaigns(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		PollsParams *CharityCampaignsParams
		respBody    string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&CharityCampaignsParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&CharityCampaignsParams{BroadcasterID: "123456"},
			`{ "data": [{ "id": "123-abc-456-def", "broadcaster_id": "123456", "broadcaster_name": "SunnySideUp", "broadcaster_login": "sunnysideup", "charity_name": "Example name", "charity_description": "Example description", "charity_logo": "https://example.url/logo.png", "charity_website": "https://www.example.com", "current_amount": { "value": 86000, "decimal_places": 2, "currency": "USD" }, "target_amount": { "value": 1500000, "decimal_places": 2, "currency": "USD" } }] }`},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetCharityCampaigns(testCase.PollsParams)
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

		if len(resp.Data.Campaigns) != 1 {
			t.Errorf("expected charity campaigns len to be 1, got %d", len(resp.Data.Campaigns))
		}
		if resp.Data.Campaigns[0].Name != "Example name" {
			t.Errorf("invalid charity name %q, expected Example name", resp.Data.Campaigns[0].Name)
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

	_, err := c.GetCharityCampaigns(&CharityCampaignsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
