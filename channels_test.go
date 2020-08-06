package helix

import (
	"net/http"
	"testing"
)

func TestSearchChannels(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		First      int
		respBody   string
		parsed     []Channel
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			`{"data":[{"broadcaster_language":"en","display_name":"Ninja","game_id":"33214","id":"27833742640","is_live":false,"tag_ids":[],"thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg","title":"I have lost my voice D: | twitter.com/Ninja","started_at":"2018-03-06T15:07:45Z"},{"broadcaster_language":"en","display_name":"DrDisrespect","game_id":"33214","id":"27834185424","is_live":false,"tag_ids":[],"thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg","title":"Turbo Treehouses || @DrDisRespect","started_at":"2018-03-06T16:05:00Z"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
			[]Channel{
				Channel{
					ID:           "27833742640",
					GameID:       "33214",
					DisplayName:  "Ninja",
					Language:     "en",
					Title:        "I have lost my voice D: | twitter.com/Ninja",
					ThumbnailURL: "https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg",
					IsLive:       false,
					TagIDs:       []string{},
				},
				Channel{
					ID:           "27834185424",
					GameID:       "33214",
					DisplayName:  "DrDisrespect",
					Language:     "en",
					Title:        "Turbo Treehouses || @DrDisRespect",
					ThumbnailURL: "https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg",
					IsLive:       false,
					TagIDs:       []string{},
				},
			},
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			101,
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
			[]Channel{},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SearchChannels(&SearchChannelsParams{
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

		if len(resp.Data.Channels) != testCase.First {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.First, len(resp.Data.Channels))
		}

		for i, channel := range resp.Data.Channels {
			if channel.ID != testCase.parsed[i].ID {
				t.Errorf("Expected struct field ID = %s, was %s", testCase.parsed[i].ID, channel.ID)
			}
			if channel.GameID != testCase.parsed[i].GameID {
				t.Errorf("Expected struct field GameID = %s, was %s", testCase.parsed[i].GameID, channel.GameID)
			}
			if channel.DisplayName != testCase.parsed[i].DisplayName {
				t.Errorf("Expected struct field DisplayName = %s, was %s", testCase.parsed[i].DisplayName, channel.DisplayName)
			}
			if channel.Language != testCase.parsed[i].Language {
				t.Errorf("Expected struct field Language = %s, was %s", testCase.parsed[i].Language, channel.Language)
			}
			if channel.Title != testCase.parsed[i].Title {
				t.Errorf("Expected struct field Title = %s, was %s", testCase.parsed[i].Title, channel.Title)
			}
			if channel.ThumbnailURL != testCase.parsed[i].ThumbnailURL {
				t.Errorf("Expected struct field ThumbnailURL = %s, was %s", testCase.parsed[i].ThumbnailURL, channel.ThumbnailURL)
			}
			if channel.IsLive != testCase.parsed[i].IsLive {
				t.Errorf("Expected struct field IsLive = %t, was %t", testCase.parsed[i].IsLive, channel.IsLive)
			}
			if len(channel.TagIDs) != len(testCase.parsed[i].TagIDs) {
				t.Errorf("Expected struct field TagIDs length = %d, was %d", len(testCase.parsed[i].TagIDs), len(channel.TagIDs))
			}
		}
	}
}
