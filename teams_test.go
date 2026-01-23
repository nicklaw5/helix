package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetTeams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		GetTeamsParams *GetTeamsParams
		respBody       string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetTeamsParams{},
			`{"error":"Bad Request","status":400,"message":"Must provide either team_id or team_name"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetTeamsParams{Name: "weightedblanket"},
			`{"data":[{"users":[{"user_id":"278217731","user_login":"mastermndio","user_name":"MasterMndio"},{"user_id":"41284990","user_login":"jenninexus","user_name":"JenniNexus"}],"background_image_url":null,"banner":null,"created_at":"2020-02-20T19:45:45Z","updated_at":"2021-03-15T10:20:30Z","info":"A cozy community of streamers.","thumbnail_url":"https://static-cdn.jtvnw.net/team-profile-images/weightedblanket-profile_image-12345.png","team_name":"weightedblanket","team_display_name":"WeightedBlanket","id":"1234567"}]}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetTeamsParams{ID: "1234567"},
			`{"data":[{"users":[{"user_id":"278217731","user_login":"mastermndio","user_name":"MasterMndio"}],"background_image_url":"https://example.com/bg.png","banner":"https://example.com/banner.png","created_at":"2020-02-20T19:45:45Z","updated_at":"2021-03-15T10:20:30Z","info":"Team info here","thumbnail_url":"https://static-cdn.jtvnw.net/team-profile-images/test-profile.png","team_name":"testteam","team_display_name":"Test Team","id":"1234567"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetTeams(testCase.GetTeamsParams)
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

			expectedErrMsg := "Must provide either team_id or team_name"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Teams) == 0 {
			t.Error("expected teams data but got none")
		}

		if len(resp.Data.Teams[0].Users) == 0 {
			t.Error("expected team users but got none")
		}

		if resp.Data.Teams[0].ID == "" {
			t.Error("expected team ID to be set")
		}

		if resp.Data.Teams[0].TeamName == "" {
			t.Error("expected team name to be set")
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

	_, err := c.GetTeams(&GetTeamsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
