package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetChannelChatChatterss(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode            int
		options               *Options
		GetChatChattersParams *GetChatChattersParams
		respBody              string
		validationErr         string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChatChattersParams{BroadcasterID: "", ModeratorID: "1234"},
			``,
			"error: broadcaster and moderator identifiers must be provided",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetChatChattersParams{BroadcasterID: "121445595", ModeratorID: "1234"},
			`{"data": [{"user_login": "smittysmithers", "user_name": "example", "user_id": "100249558"}]}`,
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChatChattersParams{BroadcasterID: "1231", ModeratorID: "1234"},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChannelChatChatters(testCase.GetChatChattersParams)
		if err != nil {
			if err.Error() == testCase.validationErr {
				continue
			}
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

			continue
		}

		if len(resp.Data.Chatters) != 1 {
			t.Errorf("expected %d chatters got %d", 1, len(resp.Data.Chatters))
		}

		if resp.Data.Chatters[0].UserID != "100249558" {
			t.Errorf("expected %s chatters got %s", "100249558", resp.Data.Chatters[0].UserID)
		}
	}
}

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
		ctx:  context.Background(),
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
		ctx:  context.Background(),
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
		ctx:  context.Background(),
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
		ctx:  context.Background(),
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
		ctx:  context.Background(),
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
		ctx:  context.Background(),
	}

	_, err := c.SendChatAnnouncement(&SendChatAnnouncementParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetChatSettings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *GetChatSettingsParams
		respBody      string
		validationErr string
		expected      *ChatSettings
	}{
		{ // Early-out error thrown by us
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&GetChatSettingsParams{BroadcasterID: "", ModeratorID: "82008718"},
			``,
			"error: broadcaster id must be specified",
			nil,
		},
		{ // Error thrown by Twitch
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetChatSettingsParams{BroadcasterID: "2248463222222222222222222222"},
			`{"error":"Bad Request","status":400,"message":"The parameter \"broadcaster_id\" was malformed: value must be a numeric"}`,
			"",
			nil,
		},
		{ // Request made with a valid moderator ID
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&GetChatSettingsParams{BroadcasterID: "22484632", ModeratorID: "11148817"},
			`{"data":[{"broadcaster_id":"22484632","emote_mode":false,"follower_mode":true,"follower_mode_duration":10080,"moderator_id":"11148817","non_moderator_chat_delay":true,"non_moderator_chat_delay_duration":1,"slow_mode":false,"slow_mode_wait_time":null,"subscriber_mode":false,"unique_chat_mode":false}]}`,
			``,
			&ChatSettings{
				BroadcasterID: "22484632",

				EmoteMode: false,

				FollowerMode:         true,
				FollowerModeDuration: 10080,

				SlowMode:         false,
				SlowModeWaitTime: 0,

				SubscriberMode: false,

				UniqueChatMode: false,

				ModeratorID: "11148817",

				NonModeratorChatDelay:         true,
				NonModeratorChatDelayDuration: 1,
			},
		},
		{ // Request made with no moderator ID
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&GetChatSettingsParams{BroadcasterID: "22484632", ModeratorID: ""},
			`{"data":[{"broadcaster_id":"22484632","emote_mode":false,"follower_mode":true,"follower_mode_duration":10080,"slow_mode":false,"slow_mode_wait_time":null,"subscriber_mode":false,"unique_chat_mode":false}]}`,
			``,
			&ChatSettings{
				BroadcasterID: "22484632",

				EmoteMode: false,

				FollowerMode:         true,
				FollowerModeDuration: 10080,

				SlowMode:         false,
				SlowModeWaitTime: 0,

				SubscriberMode: false,

				UniqueChatMode: false,

				ModeratorID: "",

				NonModeratorChatDelay:         false,
				NonModeratorChatDelayDuration: 0,
			},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetChatSettings(testCase.params)
		if err != nil {
			if err.Error() == testCase.validationErr {
				continue
			}
			t.Errorf("Unmatched error, expected '%v', got '%v'", testCase.validationErr, err)
			continue
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

			continue
		}

		if len(resp.Data.Settings) != 1 {
			t.Errorf("expected %d channel settings, got %d", 1, len(resp.Data.Settings))
		}

		if resp.Data.Settings[0].BroadcasterID != "22484632" {
			t.Errorf("expected broadcaster_id to be %s, got %s", "22484632", resp.Data.Settings[0].BroadcasterID)
		}

		if testCase.expected != nil {
			expected := testCase.expected
			actual := resp.Data.Settings[0]

			if *expected != actual {
				t.Errorf("expected %v channel settings, got %v", *expected, actual)
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

	_, err := c.GetChatSettings(&GetChatSettingsParams{BroadcasterID: "123"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	const expectedHTTPError = "Failed to execute API request: Oops, that's bad :("

	if err.Error() != expectedHTTPError {
		t.Errorf("expected error does match return error, got '%s'", err.Error())
	}
}

func TestUpdateChatSettings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *UpdateChatSettingsParams
		respBody      string
		validationErr string
		expected      *ChatSettings
	}{
		{ // Early-out error thrown by us
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&UpdateChatSettingsParams{BroadcasterID: "", ModeratorID: "82008718"},
			``,
			"error: broadcaster id must be specified",
			nil,
		},
		{ // Early-out error thrown by us
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&UpdateChatSettingsParams{BroadcasterID: "11148817", ModeratorID: ""},
			``,
			"error: moderator id must be specified",
			nil,
		},
		{ // Request that updates nothing still return all data
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&UpdateChatSettingsParams{BroadcasterID: "22484632", ModeratorID: "11148817"},
			`{"data":[{"broadcaster_id":"22484632","emote_mode":false,"follower_mode":true,"follower_mode_duration":10080,"moderator_id":"11148817","non_moderator_chat_delay":true,"non_moderator_chat_delay_duration":1,"slow_mode":false,"slow_mode_wait_time":null,"subscriber_mode":false,"unique_chat_mode":false}]}`,
			``,
			&ChatSettings{
				BroadcasterID: "22484632",
				ModeratorID:   "11148817",

				EmoteMode: false,

				FollowerMode:         true,
				FollowerModeDuration: 10080,

				SlowMode:         false,
				SlowModeWaitTime: 0,

				SubscriberMode: false,

				UniqueChatMode: false,

				NonModeratorChatDelay:         true,
				NonModeratorChatDelayDuration: 1,
			},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateChatSettings(testCase.params)
		if err != nil {
			if err.Error() == testCase.validationErr {
				continue
			}
			t.Errorf("Unmatched error, expected '%v', got '%v'", testCase.validationErr, err)
			continue
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

			continue
		}

		if len(resp.Data.Settings) != 1 {
			t.Errorf("expected %d channel settings, got %d", 1, len(resp.Data.Settings))
		}

		if resp.Data.Settings[0].BroadcasterID != "22484632" {
			t.Errorf("expected broadcaster_id to be %s, got %s", "22484632", resp.Data.Settings[0].BroadcasterID)
		}

		if testCase.expected != nil {
			expected := testCase.expected
			actual := resp.Data.Settings[0]

			if *expected != actual {
				t.Errorf("expected %v channel settings, got %v", *expected, actual)
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

	_, err := c.UpdateChatSettings(&UpdateChatSettingsParams{BroadcasterID: "123", ModeratorID: "123"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	const expectedHTTPError = "Failed to execute API request: Oops, that's bad :("

	if err.Error() != expectedHTTPError {
		t.Errorf("expected error does match return error, got '%s'", err.Error())
	}
}

func TestGetUserChatColor(t *testing.T) {
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
			"22484632",
			`{"data": [{"user_id": "11111","user_name": "SpeedySpeedster1","user_login": "speedyspeedster1","color": "#9146FF"},{"user_id": "44444","user_name": "SpeedySpeedster2","user_login": "speedyspeedster2","color": ""}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"22484632",
			`{"error":"Bad Request","status":400,"message":"The ID in the user_id query parameter is not valid."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUserChatColor(&GetUserChatColorParams{
			UserID: testCase.UserID,
		})
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

			expectedErrMsg := "The ID in the user_id query parameter is not valid."
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}

func TestUpdateUserChatColor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		UserID     string
		Color      string
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"22484632",
			"blue",
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"22484632",
			"bad_color",
			`{"error":"Bad Request","status":400,"message":"The named color in the color query parameter is not valid."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateUserChatColor(&UpdateUserChatColorParams{
			UserID: testCase.UserID,
			Color:  testCase.Color,
		})
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

			expectedErrMsg := "The named color in the color query parameter is not valid."
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}

func TestSendChatMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *SendChatMessageParams
		respBody   string
		err        string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&SendChatMessageParams{
				BroadcasterID: "1234",
				SenderID:      "5678",
				Message:       "Hello, world! twitchdevHype",
			},
			`{"data":[{"message_id": "abc-123-def","is_sent": true}]}`,
			``,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&SendChatMessageParams{
				BroadcasterID: "",
				SenderID:      "5678",
				Message:       "Hello, world! twitchdevHype",
			},
			``,
			`error: broadcaster id must be specified`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&SendChatMessageParams{
				BroadcasterID: "1234",
				SenderID:      "",
				Message:       "Hello, world! twitchdevHype",
			},
			``,
			`error: sender id must be specified`,
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&SendChatMessageParams{
				BroadcasterID: "1234",
				SenderID:      "5678",
				Message:       "Hello, world! twitchdevHype",
			},
			`{"error":"Unauthorized","status":401,"message":"Missing user:write:chat scope"}`, // missing required scope
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendChatMessage(testCase.params)
		if err != nil {
			if err.Error() == testCase.err {
				continue
			}
			t.Errorf("Unmatched error, expected '%v', got '%v'", testCase.err, err)
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be %d, got %d", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != "Missing user:write:chat scope" {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", "Missing user:write:chat scope", resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Messages) < 1 {
			t.Errorf("Expected the number of messages to be a positive number")
		}

		if len(resp.Data.Messages[0].MessageID) == 0 {
			t.Errorf("Expected message_id not to be empty")
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

	_, err := c.SendChatMessage(&SendChatMessageParams{BroadcasterID: "123", SenderID: "456", Message: "Hello, world! twitchdevHype"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	const expectedHTTPError = "Failed to execute API request: Oops, that's bad :("

	if err.Error() != expectedHTTPError {
		t.Errorf("expected error does match return error, got '%s'", err.Error())
	}
}
