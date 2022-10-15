package helix

import (
	"net/http"
	"testing"
	"time"
)

func TestGetBannedUsers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		UserID        string
		After         string
		Before        string
		respBody      string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"", // missing broadcaster id
			"",
			"",
			"",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"23161357",
			"",
			"",
			"",
			`{"data":[{"expires_at":"","user_id":"54946241","user_name":"chronophylos","user_name":"chronophylos"},{"expires_at":"2022-03-15T02:00:28Z","user_id":"423374343","user_name":"glowillig"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6MX19"}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetBannedUsers(&BannedUsersParams{
			BroadcasterID: testCase.BroadcasterID,
			UserID:        testCase.UserID,
			After:         testCase.After,
			Before:        testCase.Before,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := `Missing required parameter "broadcaster_id"`
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
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

	_, err := c.GetBannedUsers(&BannedUsersParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestBanUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		params         *BanUserParams
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&BanUserParams{BroadcasterID: "1234", ModeratorId: "5678", Body: BanUserRequestBody{
				UserId:   "9876",
				Duration: 300,
				Reason:   "no reason",
			}},
			`{"error":"Bad Request","status": 400,"message":"user is already banned"}`,
			"user is already banned",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&BanUserParams{BroadcasterID: "1234", ModeratorId: "5678", Body: BanUserRequestBody{
				UserId:   "9876",
				Duration: 300,
				Reason:   "no reason",
			}},
			`{"data": [{"broadcaster_id": "1234","moderator_id": "5678","user_id": "9876","created_at": "2021-09-28T19:22:31Z","end_time": "2021-09-28T19:27:31Z"}]}`,
			"",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&BanUserParams{BroadcasterID: "1234", ModeratorId: "5678", Body: BanUserRequestBody{
				UserId: "9876",
				Reason: "no reason",
			}},
			`{"data": [{"broadcaster_id": "1234","moderator_id": "5678","user_id": "9876","created_at": "2021-09-28T18:22:31Z","end_time": null}]}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.BanUser(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}
		}

		if len(resp.Data.Bans) != 0 {
			if resp.Data.Bans[0].BoardcasterId != testCase.params.BroadcasterID {
				t.Errorf("expected broadcaster id to be \"%s\", got \"%s\"", testCase.params.BroadcasterID, resp.Data.Bans[0].BoardcasterId)
			}

			if resp.Data.Bans[0].ModeratorId != testCase.params.ModeratorId {
				t.Errorf("expected moderator id to be \"%s\", got \"%s\"", testCase.params.ModeratorId, resp.Data.Bans[0].ModeratorId)
			}

			if resp.Data.Bans[0].UserId != testCase.params.Body.UserId {
				t.Errorf("expected user id to be \"%s\", got \"%s\"", testCase.params.Body.UserId, resp.Data.Bans[0].UserId)
			}

			if resp.Data.Bans[0].EndTime != "" {
				layout := "2006-01-02T15:04:05Z"
				createdTime, _ := time.Parse(layout, resp.Data.Bans[0].CreatedAt)
				endTime, _ := time.Parse(layout, resp.Data.Bans[0].EndTime)

				expireTime := createdTime.Add(time.Duration(testCase.params.Body.Duration * int(time.Second)))

				if !expireTime.Equal(endTime) {
					t.Errorf("expected endtime to be \"%s\", got \"%s\"", expireTime, resp.Data.Bans[0].EndTime)
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

	_, err := c.BanUser(&BanUserParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUnbanUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		params         *UnbanUserParams
		respBody       string
		expectedErrMsg string
	}{
		{
			204,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&UnbanUserParams{BroadcasterID: "1234", ModeratorID: "5678", UserID: "9876"},
			"",
			"",
		},
		{
			400,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&UnbanUserParams{BroadcasterID: "1234", ModeratorID: "5678", UserID: "5432"},
			`{"error": "Bad Request", "status": 400, "message": "user is not banned"}`,
			"user is not banned",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UnbanUser(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode != http.StatusNoContent {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
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

	_, err := c.BanUser(&BanUserParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
