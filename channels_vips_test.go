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
		params     *ChannelsVipsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ChannelsVipsParams{
				UserID: "123",
			},
			`{ "data": [{ "user_id": "11111", "user_name": "UserDisplayName", "user_login": "userloginname" }], "pagination": { "cursor": "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6NX19" } }`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ChannelsVipsParams{
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

	_, err := c.GetChannelVips(&ChannelsVipsParams{})
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
		params     *AddChannelsVipsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&AddChannelsVipsParams{
				BroadcasterID: "123",
				UserID:        "456",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&AddChannelsVipsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.AddChannelVips(testCase.params)
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

	_, err := c.AddChannelVips(&AddChannelsVipsParams{})
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
		params     *RemoveChannelsVipsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&RemoveChannelsVipsParams{
				BroadcasterID: "123",
				UserID:        "456",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&RemoveChannelsVipsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RemoveChannelVips(testCase.params)
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

	_, err := c.RemoveChannelVips(&RemoveChannelsVipsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
