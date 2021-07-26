package helix

import (
	"net/http"
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		UserID        []string
		respBody      string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"123",
			[]string{},
			`{"data": [{"broadcaster_id":"123","broadcaster_login":"test_user","broadcaster_name":"test_user","is_gift":true,"gifter_id":"456","gifter_login":"another_user","gifter_name":"Another_User","tier":"3000","plan_name":"The Ninjas","user_id":"123","user_id":"123","user_login":"test_user","user_name":"test_user"}],"pagination":{"cursor":"xxxx"},"total":1}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"",
			[]string{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetSubscriptions(&SubscriptionsParams{
			BroadcasterID: testCase.BroadcasterID,
			UserID:        testCase.UserID,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if testCase.statusCode != http.StatusOK {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			errMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != errMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", errMsg, resp.ErrorMessage)
			}

			continue
		}

		if resp.Data.Total != 1 {
			t.Errorf("expected total field to be 1 got %d", resp.Data.Total)
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

	_, err := c.GetSubscriptions(&SubscriptionsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestChechUserSubscription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		UserID        string
		respBody      string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"123",
			"123",
			`{"data":[{"broadcaster_id":"122330828","broadcaster_name":"test_user","broadcaster_login":"test_user","is_gift":false,"tier":"1000"}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			"123",
			"",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"user_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CheckUserSubsription(&UserSubscriptionsParams{
			BroadcasterID: testCase.BroadcasterID,
			UserID:        testCase.UserID,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if testCase.statusCode != http.StatusOK {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			errMsg := "Missing required parameter \"user_id\""
			if resp.ErrorMessage != errMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", errMsg, resp.ErrorMessage)
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

	_, err := c.CheckUserSubsription(&UserSubscriptionsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
