package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetClips(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *ClipsParams
		respBody   string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "invalid-client-id"}, // invalid client id
			&ClipsParams{IDs: []string{"EncouragingPluckySlothSSSsss"}},
			`{"error":"Unauthorized","status":401,"message":"Must provide a valid Client-ID or OAuth token"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ClipsParams{IDs: []string{"bad-id"}}, // invalid clip id
			`{"data":[],"pagination":{}}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ClipsParams{BroadcasterID: "bad-broadcaster-id"}, // invalid broadcaster id
			`{"data":[],"pagination":{}}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ClipsParams{GameID: "bad-game-id"}, // invalid game id
			`{"data":[],"pagination":{}}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ClipsParams{IDs: []string{"EncouragingPluckySlothSSSsss"}}, // valid clip id
			`{"data":[{"id":"EncouragingPluckySlothSSSsss","url":"https://clips.twitch.tv/EncouragingPluckySlothSSSsss","embed_url":"https://clips.twitch.tv/embed?clip=EncouragingPluckySlothSSSsss","broadcaster_id":"26490481","broadcaster_name":"summit1g","creator_id":"143839181","creator_name":"nB00ts","video_id":"","game_id":"490377","language":"en","title":"summit and fat tim discover how to use maps","view_count":91876,"created_at":"2018-01-25T04:04:15Z","thumbnail_url":"https://clips-media-assets2.twitch.tv/182509178-preview-480x272.jpg","duration":22.3}],"pagination":{}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetClips(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if testCase.statusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be %d, got %d", testCase.statusCode, resp.ErrorStatus)
			}

			errMsg := "Must provide a valid Client-ID or OAuth token"
			if resp.ErrorMessage != errMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", errMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Clips) != 0 && resp.Data.Clips[0].ID != testCase.params.IDs[0] {
			t.Errorf("expected clip id to be \"%s\", got \"%s\"", testCase.params.IDs[0], resp.Data.Clips[0].ID)
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

	_, err := c.GetClips(&ClipsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestCreateClip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode      int
		options         *Options
		params          *CreateClipParams
		respBody        string
		headerLimit     string
		headerRemaining string
	}{
		{
			http.StatusAccepted,
			&Options{ClientID: "my-client-id"},
			&CreateClipParams{BroadcasterID: "26490481"}, // summit1g
			`{"data":[{"id":"IronicHedonisticOryxSquadGoals","edit_url":"https://clips.twitch.tv/IronicHedonisticOryxSquadGoals/edit"}]}`,
			"600",
			"598",
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&CreateClipParams{BroadcasterID: "26490481"},                                 // summit1g
			`{"error":"Unauthorized","status":401,"message":"Missing clips:edit scope"}`, // missing required scope
			"600",
			"597",
		},
	}

	for _, testCase := range testCases {
		mockRespHeaders := map[string]string{
			"Ratelimit-Helixclipscreation-Limit":     testCase.headerLimit,
			"Ratelimit-Helixclipscreation-Remaining": testCase.headerRemaining,
		}

		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, mockRespHeaders))

		resp, err := c.CreateClip(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if testCase.statusCode != http.StatusAccepted {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be %d, got %d", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != "Missing clips:edit scope" {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", "Missing clips:edit scope", resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.ClipEditURLs[0].ID) < 1 {
			t.Errorf("expected clip id not to be empty, got \"%s\"", resp.Data.ClipEditURLs[0].ID)
		}

		if len(resp.Data.ClipEditURLs[0].EditURL) < 1 {
			t.Errorf("expected clip edit url not to be empty, got \"%s\"", resp.Data.ClipEditURLs[0].EditURL)
		}

		if resp.GetClipsCreationRateLimit() < 1 {
			t.Errorf("expected clip create rate limit limit not to be \"0\", got \"%d\"", resp.GetClipsCreationRateLimit())
		}

		if resp.GetClipsCreationRateLimitRemaining() < 1 {
			t.Errorf("expected clip create rate limit remaining not to be \"0\", got \"%d\"", resp.GetClipsCreationRateLimitRemaining())
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

	_, err := c.CreateClip(&CreateClipParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
