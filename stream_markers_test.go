package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetStreamMarkers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		userID      string
		markerCount int
		first       int
		respBody    string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"123",
			1,
			5,
			`{"data":[{"user_id":"123","user_name":"displayname","user_name":"DisplayName","videos":[{"video_id":"456","markers":[{"id":"106b8d6243a4f883d25ad75e6cdffdc4","created_at":"2018-08-20T20:10:03Z","description":"hello,thisisamarker!","position_seconds":244,"URL":"https://twitch.tv/videos/456?t=0h4m06s"}]}]}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjoiMjk1MjA0Mzk3OjI1Mzpib29rbWFyazoxMDZiOGQ1Y"}}`},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"0",
			0,
			101,
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetStreamMarkers(&StreamMarkersParams{
			UserID: testCase.userID,
			First:  testCase.first,
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

		if len(resp.Data.StreamMarkers[0].Videos[0].Markers) != testCase.markerCount {
			t.Errorf("expected \"%d\" stream markers, got \"%d\"", testCase.markerCount, len(resp.Data.StreamMarkers))
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

	_, err := c.GetStreamMarkers(&StreamMarkersParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestCreateStreamMarker(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		userID      string
		description string
		markerCount int
		respBody    string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"123",
			"a notable moment",
			1,
			`{"data":[{"id":"123","created_at":"2018-08-20T20:10:03Z","description":"hello, this is a marker!","position_seconds":244}]}`},
		{
			http.StatusForbidden,
			&Options{ClientID: "my-client-id"},
			"124",
			"another notable moment",
			0,
			`{"error":"Forbidden","status":403,"message":"Not authorized to create a stream marker for channel test."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreateStreamMarker(&CreateStreamMarkerParams{
			UserID:      testCase.userID,
			Description: testCase.description,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Forbidden Request Responses
		if resp.StatusCode == http.StatusForbidden {
			firstErrStr := "Not authorized to create a stream marker for channel test."
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.CreateStreamMarkers) != testCase.markerCount {
			t.Errorf("expected \"%d\" stream markers, got \"%d\"", testCase.markerCount, len(resp.Data.CreateStreamMarkers))
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

	_, err := c.CreateStreamMarker(&CreateStreamMarkerParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
