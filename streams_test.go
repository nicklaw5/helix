package helix

import (
	"net/http"
	"strconv"
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
			`{"data":[{"id":"27833742640","user_id":"19571641","user_name":"Ninja","game_id":"33214","tag_ids":[],"type":"live","title":"I have lost my voice D: | twitter.com/Ninja","viewer_count":72124,"started_at":"2018-03-06T15:07:45Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg"},{"id":"27834185424","user_id":"17337557","user_name":"DrDisrespect","game_id":"33214","tag_ids":[],"type":"live","title":"Turbo Treehouses || @DrDisRespect","viewer_count":29687,"started_at":"2018-03-06T16:05:00Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
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
}
