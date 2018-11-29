package helix

import (
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
			`{"data":[{"user_id":"123","user_name":"DisplayName","videos":[{"video_id":"456","markers":[{"id":"106b8d6243a4f883d25ad75e6cdffdc4","created_at":"2018-08-20T20:10:03Z","description":"hello,thisisamarker!","position_seconds":244,"URL":"https://twitch.tv/videos/456?t=0h4m06s"}]}]}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjoiMjk1MjA0Mzk3OjI1Mzpib29rbWFyazoxMDZiOGQ1Y"}}`,
		},
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

		if len(resp.Data.Markers) != testCase.markerCount {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.markerCount, len(resp.Data.Markers))
		}
	}
}
