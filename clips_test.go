package helix

import (
	"net/http"
	"strconv"
	"testing"
)

func TestGetClips(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		slug       string
		respBody   string
	}{
		{
			http.StatusOK,
			"EncouragingPluckySlothSSSsss",
			`{"data":[{"id":"EncouragingPluckySlothSSSsss","url":"https://clips.twitch.tv/EncouragingPluckySlothSSSsss","embed_url":"https://clips.twitch.tv/embed?clip=EncouragingPluckySlothSSSsss","broadcaster_id":"26490481","creator_id":"143839181","video_id":"222004532","game_id":"490377","language":"en","title":"summit and fat tim discover how to use maps","view_count":81808,"created_at":"2018-01-25T04:04:15Z","thumbnail_url":"https://clips-media-assets.twitch.tv/182509178-preview-480x272.jpg"}]}`,
		},
		{
			http.StatusNotFound,
			"bad-slug",
			`{"error":"Not Found","status":404,"message":"clip not found"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetClips(&ClipsParams{
			IDs: []string{testCase.slug},
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if testCase.statusCode == http.StatusNotFound {
			if resp.Error != "Not Found" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Not Found", resp.Error)
			}

			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be %d, got %d", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != "clip not found" {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", "clip not found", resp.ErrorMessage)
			}

			continue
		}

		if resp.Data.Clips[0].ID != testCase.slug {
			t.Errorf("expected clip id to be \"%s\", got \"%s\"", testCase.slug, resp.Data.Clips[0].ID)
		}
	}
}

func TestCreateClip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode      int
		broadcasterID   string
		respBody        string
		headerLimit     int
		headerRemaining int
	}{
		{
			http.StatusAccepted,
			"26490481", // summit1g
			`{"data":[{"id":"IronicHedonisticOryxSquadGoals","edit_url":"https://clips.twitch.tv/IronicHedonisticOryxSquadGoals/edit"}]}`,
			600,
			598,
		},
		{
			http.StatusUnauthorized,
			"26490481", // summit1g
			`{"error":"Unauthorized","status":401,"message":"Missing clips:edit scope"}`, // missing required scope
			600,
			597,
		},
	}

	for _, testCase := range testCases {
		mockRespHeaders := map[string]string{
			"Ratelimit-Helixclipscreation-Limit":     strconv.Itoa(testCase.headerLimit),
			"Ratelimit-Helixclipscreation-Remaining": strconv.Itoa(testCase.headerRemaining),
		}

		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody, mockRespHeaders))

		resp, err := c.CreateClip(testCase.broadcasterID)
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

		if resp.ClipsCreationRateLimit.Limit < 1 {
			t.Errorf("expected clip create rate limit limit not to be \"0\", got \"%d\"", resp.ClipsCreationRateLimit.Limit)
		}

		if resp.ClipsCreationRateLimit.Remaining < 1 {
			t.Errorf("expected clip create rate limit remaining not to be \"0\", got \"%d\"", resp.ClipsCreationRateLimit.Remaining)
		}
	}
}
