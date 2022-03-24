package helix

import (
	"net/http"
	"testing"
)

func TestGetStreams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		First      int
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			`{"data":[{"id":"27833742640","user_id":"19571641","user_login":"ninja","user_name":"Ninja","game_id":"33214","game_name":"Tekken","tag_ids":[],"type":"live","title":"I have lost my voice D: | twitter.com/Ninja","viewer_count":72124,"started_at":"2018-03-06T15:07:45Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg","is_mature": false},{"id":"27834185424","user_id":"17337557","user_name":"DrDisrespect","game_id":"33214","game_name":"Tekken","tag_ids":[],"type":"live","title":"Turbo Treehouses || @DrDisRespect","viewer_count":29687,"started_at":"2018-03-06T16:05:00Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg","is_mature": true}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			101,
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetStreams(&StreamsParams{
			First: testCase.First,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "The parameter \"first\" was malformed: the value must be less than or equal to 100"
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Streams) != testCase.First {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.First, len(resp.Data.Streams))
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

	_, err := c.GetStreams(&StreamsParams{
		Language: []string{"en"},
	})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetFollowedStreams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		UserID     string
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"123456789",
			`{"data":[{"id":"27833742640","user_id":"19571641","user_login":"ninja","user_name":"Ninja","game_id":"33214","game_name":"Tekken","tag_ids":[],"type":"live","title":"I have lost my voice D: | twitter.com/Ninja","viewer_count":72124,"started_at":"2018-03-06T15:07:45Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg","is_mature": false},{"id":"27834185424","user_id":"17337557","user_name":"DrDisrespect","game_id":"33214","game_name":"Tekken","tag_ids":[],"type":"live","title":"Turbo Treehouses || @DrDisRespect","viewer_count":29687,"started_at":"2018-03-06T16:05:00Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg","is_mature": true}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"user_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetFollowedStream(&FollowedStreamsParams{
			UserID: testCase.UserID,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"user_id\""
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
	}

	_, err := c.GetFollowedStream(&FollowedStreamsParams{
		UserID: "123456",
	})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetStreamKeys(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		broadcasterID string
		Length        int
		respBody      string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"kivutar",
			1,
			`{"data":[{"stream_key":"live_695820277_TF1dAMbU4cQvGKyrk2Q88SvWNCw6Rs"}]}`,
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			"kivutar",
			0,
			`{"error":"Unauthorized","status":401,"message":"Invalid OAuth token"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetStreamKey(&StreamKeyParams{
			BroadcasterID: testCase.broadcasterID,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Invalid OAuth token"
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Data) != testCase.Length {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.Length, len(resp.Data.Data))
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

	_, err := c.GetStreamKey(&StreamKeyParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
