package helix

import (
	"context"
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
			`{"data":[{"broadcaster_language":"en","display_name":"Ninja","game_id":"33214","game_name":"Loserfruit","id":"27833742640","is_live":false,"tag_ids":[],"thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg","title":"I have lost my voice D: | twitter.com/Ninja","started_at":"2018-03-06T15:07:45Z"},{"broadcaster_language":"en","display_name":"DrDisrespect","game_id":"33214","game_name":"Loserfruit","id":"27834185424","is_live":false,"tag_ids":[],"thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg","title":"Turbo Treehouses || @DrDisRespect","started_at":"2018-03-06T16:05:00Z"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
			[]Channel{
				{
					ID:               "27833742640",
					GameID:           "33214",
					GameName:         "Loserfruit",
					BroadcasterLogin: "ninja",
					DisplayName:      "Ninja",
					Language:         "en",
					Title:            "I have lost my voice D: | twitter.com/Ninja",
					ThumbnailURL:     "https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg",
					IsLive:           false,
					TagIDs:           []string{},
				},
				{
					ID:               "27834185424",
					GameID:           "33214",
					GameName:         "Loserfruit",
					BroadcasterLogin: "drdisrespect",
					DisplayName:      "DrDisrespect",
					Language:         "en",
					Title:            "Turbo Treehouses || @DrDisRespect",
					ThumbnailURL:     "https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg",
					IsLive:           false,
					TagIDs:           []string{},
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

	_, err := c.SearchChannels(&SearchChannelsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetChannelInformation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		BroadcasterID  string
		BroadcasterIDs []string
		respBody       string
		parsed         []ChannelInformation
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"44445592",
			[]string{},
			`{"data":[{"broadcaster_id":"44445592","broadcaster_login":"pokimane","broadcaster_name":"pokimane","broadcaster_language":"en","game_id":"509658","game_name":"Just Chatting","title":"See you Wednesday 8am for Among Us ^_^", "delay": 2}]}`,
			[]ChannelInformation{
				{
					BroadcasterID:       "44445592",
					BroadcasterName:     "pokimane",
					BroadcasterLanguage: "en",
					GameID:              "509658",
					GameName:            "Just Chatting",
					Title:               "See you Wednesday 8am for Among Us ^_^",
					Delay:               2,
				},
			},
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"9999999999999999999",
			[]string{},
			`{"error":"Bad Request","status":400,"message":"Invalid broadcasterID"}`,
			[]ChannelInformation{},
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"",
			[]string{"99631238", "11148817"},
			`{"data":[{"broadcaster_id":"11148817","broadcaster_name":"pajlada","broadcaster_language":"en","game_id":"509670","game_name":"Science \u0026 Technology","title":"Programming.S04E71.1440p.MP4-XD.NVENC","delay":0},{"broadcaster_id":"99631238","broadcaster_name":"zneix","broadcaster_language":"en","game_id":"509670","game_name":"Science \u0026 Technology","title":"sunday coding, ranting and gaming (serious sam 4 later)","delay":0}]}`,
			[]ChannelInformation{
				{
					BroadcasterID:       "11148817",
					BroadcasterName:     "pajlada",
					BroadcasterLanguage: "en",
					GameID:              "509670",
					GameName:            "Science & Technology",
					Title:               "Programming.S04E71.1440p.MP4-XD.NVENC",
					Delay:               0,
				},
				{
					BroadcasterID:       "99631238",
					BroadcasterName:     "zneix",
					BroadcasterLanguage: "en",
					GameID:              "509670",
					GameName:            "Science & Technology",
					Title:               "sunday coding, ranting and gaming (serious sam 4 later)",
					Delay:               0,
				},
			},
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"",
			[]string{"99631238", "forsen"},
			`{"error": "Bad Request","status": 400,"message": "broadcaster_id \"forsen\" must be numeric"}`,
			[]ChannelInformation{},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		var resp *GetChannelInformationResponse
		var err error
		if testCase.BroadcasterID != "" {
			resp, err = c.GetChannelInformation(&GetChannelInformationParams{
				BroadcasterID: testCase.BroadcasterID,
			})
		} else {
			resp, err = c.GetChannelInformation(&GetChannelInformationParams{
				BroadcasterIDs: testCase.BroadcasterIDs,
			})
		}
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			broadcasterIDErrStr := "Invalid broadcasterID"
			if resp.ErrorMessage != broadcasterIDErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", broadcasterIDErrStr, resp.ErrorMessage)
				continue
			}
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		for i, channel := range resp.Data.Channels {
			if channel.BroadcasterID != testCase.parsed[i].BroadcasterID {
				t.Errorf("Expected struct field ChannelID = %s, was %s", testCase.parsed[i].BroadcasterID, channel.BroadcasterID)
			}

			if channel.BroadcasterName != testCase.parsed[i].BroadcasterName {
				t.Errorf("Expected struct field BroadcasterName = %s, was %s", testCase.parsed[i].BroadcasterName, channel.BroadcasterName)
			}

			if channel.BroadcasterLanguage != testCase.parsed[i].BroadcasterLanguage {
				t.Errorf("Expected struct field BroadcasterLanguage = %s, was %s", testCase.parsed[i].BroadcasterLanguage, channel.BroadcasterLanguage)
			}

			if channel.GameID != testCase.parsed[i].GameID {
				t.Errorf("Expected struct field GameID = %s, was %s", testCase.parsed[i].GameID, channel.GameID)
			}

			if channel.GameName != testCase.parsed[i].GameName {
				t.Errorf("Expected struct field GameName = %s, was %s", testCase.parsed[i].GameName, channel.GameName)
			}

			if channel.Title != testCase.parsed[i].Title {
				t.Errorf("Expected struct field Title = %s, was %s", testCase.parsed[i].Title, channel.Title)
			}
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

	_, err := c.GetChannelInformation(&GetChannelInformationParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestEditChannelInformation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *EditChannelInformationParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&EditChannelInformationParams{
				BroadcasterID:       "123",
				GameID:              "456",
				BroadcasterLanguage: "en",
				Title:               "Test title",
			},
			`{"data":[{"game_id":"498566","broadcaster_language":"en","title":"Test Twitch API"}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&EditChannelInformationParams{
				BroadcasterID:       "789",
				GameID:              "456",
				BroadcasterLanguage: "en",
				Title:               "Test title",
				Delay:               3,
			},
			`{"error":"Bad Request","status":400,"message":"the broadcaster is not partnered, failed to set delay"}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&EditChannelInformationParams{
				BroadcasterID:       "789",
				GameID:              "-1",
				BroadcasterLanguage: "en",
				Title:               "Test title",
			},
			`{"error":"Bad Request","status":400,"message":"invalid game_id"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.EditChannelInformation(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			broadcasterIDErrStr := "the broadcaster is not partnered, failed to set delay"

			if testCase.params.GameID == "-1" {
				broadcasterIDErrStr = "invalid game_id"
			}

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

	_, err := c.EditChannelInformation(&EditChannelInformationParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestChannelFollows(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetChannelFollowsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetChannelFollowsParams{
				BroadcasterID: "123",
			},
			`{ "total": 8, "data": [{ "user_id": "11111", "user_name": "UserDisplayName", "user_login": "userloginname", "followed_at": "2022-05-24T22:22:08Z" }], "pagination": { "cursor": "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6NX19" } }`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChannelFollowsParams{
				BroadcasterID: "",
			},
			`{"error":"Bad Request","status":400,"message":"the broadcaster id was not provided"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChannelFollows(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			broadcasterIDErrStr := "the broadcaster id was not provided"

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

	_, err := c.GetChannelFollows(&GetChannelFollowsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestFollowedChannels(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetFollowedChannelParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetFollowedChannelParams{
				UserID: "123",
			},
			`{ "total": 8, "data": [{ "broadcaster_id": "11111", "broadcaster_login": "userloginname", "broadcaster_name": "UserDisplayName", "followed_at": "2022-05-24T22:22:08Z" }], "pagination": { "cursor": "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6NX19" } }`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetFollowedChannelParams{
				UserID: "",
			},
			`{"error":"Bad Request","status":400,"message":"the user id was not provided"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetFollowedChannels(testCase.params)
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

	_, err := c.GetFollowedChannels(&GetFollowedChannelParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
