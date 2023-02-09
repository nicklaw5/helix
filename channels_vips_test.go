package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetChannelVips(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetChannelVipsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetChannelVipsParams{
				UserID: "123",
			},
			`{ "data": [{ "user_id": "11111", "user_name": "UserDisplayName", "user_login": "userloginname" }], "pagination": { "cursor": "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6NX19" } }`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChannelVipsParams{
				UserID: "",
			},
			`{"error":"Bad Request","status":400,"message":"the user id was not provided"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChannelVips(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			broadcasterIDErrStr := "the user id was not provided"

			if resp.ErrorMessage != broadcasterIDErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", broadcasterIDErrStr, resp.ErrorMessage)
				continue
			}
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

	_, err := c.GetChannelVips(&GetChannelVipsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestAddChannelVips(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *AddChannelVipParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&AddChannelVipParams{
				BroadcasterID: "123",
				UserID:        "456",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&AddChannelVipParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.AddChannelVip(testCase.params)
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

	_, err := c.AddChannelVip(&AddChannelVipParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestRemoveChannelVips(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *RemoveChannelVipParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&RemoveChannelVipParams{
				BroadcasterID: "123",
				UserID:        "456",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&RemoveChannelVipParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RemoveChannelVip(testCase.params)
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

	_, err := c.RemoveChannelVip(&RemoveChannelVipParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
