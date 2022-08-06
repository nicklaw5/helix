package helix

import (
	"net/http"
	"testing"
)

func TestGetChannelChatBadges(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode         int
		options            *Options
		GetChatBadgeParams *GetChatBadgeParams
		respBody           string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChatBadgeParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetChatBadgeParams{BroadcasterID: "121445595"},
			`{"data": [{"set_id": "bits","versions": [{"id": "1","image_url_1x": "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/1","image_url_2x": "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/2","image_url_4x": "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/3"}]},{"set_id": "subscriber","versions": [{"id": "0","image_url_1x": "https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/1","image_url_2x": "https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/2","image_url_4x": "https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/3"}]}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChannelChatBadges(testCase.GetChatBadgeParams)
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

	_, err := c.GetChannelChatBadges(&GetChatBadgeParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetGlobalChatBadges(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			`{"error":"Unauthorized","status":401,"message":"OAuth token is missing"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-user-access-token"},
			`{"data": [{"set_id": "vip","versions": [{"id": "1","image_url_1x": "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/1","image_url_2x": "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/2","image_url_4x": "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/3"}]}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetGlobalChatBadges()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be %s, got %s", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "OAuth token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
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

	_, err := c.GetGlobalChatBadges()
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetChannelEmotes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode             int
		options                *Options
		GetChannelEmotesParams *GetChannelEmotesParams
		respBody               string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChannelEmotesParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetChannelEmotesParams{BroadcasterID: "121445595"},
			`{"data":[{"id":"300678378","name":"scoorfPoko","images":{"url_1x":"https://static-cdn.jtvnw.net/emoticons/v1/300678378/1.0","url_2x":"https://static-cdn.jtvnw.net/emoticons/v1/300678378/2.0","url_4x":"https://static-cdn.jtvnw.net/emoticons/v1/300678378/3.0"},"tier":"1000","emote_type":"subscriptions","emote_set_id":"1347400"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChannelEmotes(testCase.GetChannelEmotesParams)
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

	_, err := c.GetChannelEmotes(&GetChannelEmotesParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetGlobalEmotes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			`{"error":"Unauthorized","status":401,"message":"OAuth token is missing"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-user-access-token"},
			`{"data":[{"id":"1220086","name":"TwitchRPG","images":{"url_1x":"https://static-cdn.jtvnw.net/emoticons/v1/1220086/1.0","url_2x":"https://static-cdn.jtvnw.net/emoticons/v1/1220086/2.0","url_4x":"https://static-cdn.jtvnw.net/emoticons/v1/1220086/3.0"}},{"id":"196892","name":"TwitchUnity","images":{"url_1x":"https://static-cdn.jtvnw.net/emoticons/v1/196892/1.0","url_2x":"https://static-cdn.jtvnw.net/emoticons/v1/196892/2.0","url_4x":"https://static-cdn.jtvnw.net/emoticons/v1/196892/3.0"}}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetGlobalEmotes()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be %s, got %s", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "OAuth token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
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

	_, err := c.GetGlobalEmotes()
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetEmoteSets(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode         int
		options            *Options
		GetEmoteSetsParams *GetEmoteSetsParams
		expectedEmotes     []EmoteWithOwner
		respBody           string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetEmoteSetsParams{EmoteSetIDs: nil},
			nil,
			`{"error":"Bad Request","status":400,"message":"The parameter \"emote_set_id\" was malformed: the value must be greater than or equal to 1"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetEmoteSetsParams{EmoteSetIDs: []string{"300678379"}},
			[]EmoteWithOwner{
				{
					Emote: Emote{
						ID:   "301147694",
						Name: "sixone3SixDab",
						Images: EmoteImage{
							Url1x: "https://static-cdn.jtvnw.net/emoticons/v1/301147694/1.0",
							Url2x: "https://static-cdn.jtvnw.net/emoticons/v1/301147694/2.0",
							Url4x: "https://static-cdn.jtvnw.net/emoticons/v1/301147694/3.0",
						},
						EmoteType:  "subscriptions",
						EmoteSetId: "300678379",
					},
					OwnerID: "44931651",
				},
			},
			`{"data":[{"id":"301147694","name":"sixone3SixDab","images":{"url_1x":"https://static-cdn.jtvnw.net/emoticons/v1/301147694/1.0","url_2x":"https://static-cdn.jtvnw.net/emoticons/v1/301147694/2.0","url_4x":"https://static-cdn.jtvnw.net/emoticons/v1/301147694/3.0"},"emote_type":"subscriptions","emote_set_id":"300678379","owner_id":"44931651"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetEmoteSets(testCase.GetEmoteSetsParams)
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

			expectedErrMsg := "The parameter \"emote_set_id\" was malformed: the value must be greater than or equal to 1"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Emotes) != len(testCase.expectedEmotes) {
			t.Errorf("returned emotes were different lengths")
		} else {
			for i, expectedEmote := range testCase.expectedEmotes {
				actualEmote := resp.Data.Emotes[i]
				if expectedEmote != actualEmote {
					t.Errorf("mismatching emotes %#v != %#v", expectedEmote, actualEmote)
				}
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
	}

	_, err := c.GetEmoteSets(&GetEmoteSetsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestSendChatAnnouncement(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode                 int
		options                    *Options
		SendChatAnnouncementParams *SendChatAnnouncementParams
		respBody                   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendChatAnnouncementParams{BroadcasterID: "100249558", ModeratorID: "100249558", Message: "hello world", Color: "blue"},
			`{"error":"Bad Request","status":400,"message":"The parameter \"Color\" was malformed: the value must be a valid color"}`,
		},
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendChatAnnouncementParams{BroadcasterID: "100249558", ModeratorID: "100249558", Message: "hello twitch chat", Color: "blue"},
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendChatAnnouncement(testCase.SendChatAnnouncementParams)
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

			expectedErrMsg := "The parameter \"Color\" was malformed: the value must be a valid color"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
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

	_, err := c.SendChatAnnouncement(&SendChatAnnouncementParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
