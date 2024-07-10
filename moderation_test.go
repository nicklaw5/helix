package helix

import (
	"context"
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
		UserID        []string
		After         string
		Before        string
		respBody      string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"", // missing broadcaster id
			[]string{},
			"",
			"",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"23161357",
			[]string{},
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
		ctx:  context.Background(),
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

			if !resp.Data.Bans[0].EndTime.IsZero() {
				expireTime := resp.Data.Bans[0].CreatedAt.Add(time.Duration(testCase.params.Body.Duration * int(time.Second)))

				if !expireTime.Equal(resp.Data.Bans[0].EndTime.Time) {
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
		ctx:  context.Background(),
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
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
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

	_, err := c.UnbanUser(&UnbanUserParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetBlockedTerms(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *BlockedTermsParams
		respBody   string
		errorMsg   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&BlockedTermsParams{BroadcasterID: "1234", ModeratorID: "5678", First: 2},
			`{
				"data": [
				  {
					"broadcaster_id": "1234",
					"moderator_id": "5678",
					"id": "520e4d4e-0cda-49c7-821e-e5ef4f88c2f2",
					"text": "A phrase I’m not fond of",
					"created_at": "2021-09-29T19:45:37Z",
					"updated_at": "2021-09-29T19:45:37Z",
					"expires_at": null
				  },
				  {
					"broadcaster_id": "1234",
					"moderator_id": "5678",
					"id": "520e4d4e-0cda-49c7-821e-e5ef4f88c2f2",
					"text": "A phrase I’m not fond of",
					"created_at": "2021-09-29T19:45:37Z",
					"updated_at": "2021-09-29T19:45:37Z",
					"expires_at": null
				  }
				],
				"pagination": {
				  "cursor": "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6I..."
				}
			}`,
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&BlockedTermsParams{BroadcasterID: "", ModeratorID: "5678", First: 2},
			``,
			"broadcaster id and moderator id must be provided",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetBlockedTerms(testCase.params)
		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Terms) != testCase.params.First {
			t.Errorf("expected terms len to be %d, got %d", testCase.params.First, len(resp.Data.Terms))
		}

		if len(resp.Data.Terms) != 0 {
			if resp.Data.Terms[0].BroadcasterID != testCase.params.BroadcasterID {
				t.Errorf("expected broadcaster id to be %s, got %s", testCase.params.BroadcasterID, resp.Data.Terms[0].BroadcasterID)
			}

			if resp.Data.Terms[0].ModeratorID != testCase.params.ModeratorID {
				t.Errorf("expected moderator id to be %s, got %s", testCase.params.ModeratorID, resp.Data.Terms[0].ModeratorID)
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

	_, err := c.GetBlockedTerms(&BlockedTermsParams{BroadcasterID: "1234", ModeratorID: "1234"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestAddBlockedTerm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *AddBlockedTermParams
		respBody   string
		errorMsg   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&AddBlockedTermParams{ModeratorID: "5678", Text: "A phrase I’m not fond of"},
			``,
			"broadcaster id and moderator id must be provided",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&AddBlockedTermParams{BroadcasterID: "1234", ModeratorID: "5678", Text: "a"},
			``,
			"the term len must be between 2 and 500",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&AddBlockedTermParams{BroadcasterID: "1234", ModeratorID: "5678", Text: "crac*"},
			`{
				"data": [
				  {
					"broadcaster_id": "1234",
					"moderator_id": "5678",
					"id": "520e4d4e-0cda-49c7-821e-e5ef4f88c2f2",
					"text": "crac*",
					"created_at": "2021-09-29T19:45:37Z",
					"updated_at": "2021-09-29T19:45:37Z",
					"expires_at": null
				  }
				]
			  }`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.AddBlockedTerm(testCase.params)
		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Terms) != 1 {
			t.Errorf("expected terms len to be %d, got %d", 1, len(resp.Data.Terms))
		}

		if resp.Data.Terms[0].Text != testCase.params.Text {
			t.Errorf("expected blocked word to be %s, got %s", testCase.params.Text, resp.Data.Terms[0].Text)
		}

		if !resp.Data.Terms[0].ExpiresAt.Time.IsZero() {
			t.Errorf("expected expiration time to be %s, got %s", "nil", resp.Data.Terms[0].ExpiresAt)
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

	_, err := c.AddBlockedTerm(&AddBlockedTermParams{BroadcasterID: "1234", ModeratorID: "1234", Text: "test"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestRemoveBlockedTerm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *RemoveBlockedTermParams
		respBody   string
		errorMsg   string
	}{
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&RemoveBlockedTermParams{BroadcasterID: "1234", ModeratorID: "5678", ID: "c9fc79b8-0f63-4ef7-9d38-efd811e74ac2"},
			``,
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&RemoveBlockedTermParams{ModeratorID: "5678", ID: "c9fc79b8-0f63-4ef7-9d38-efd811e74ac2"},
			``,
			"broadcaster id and moderator id must be provided",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&RemoveBlockedTermParams{BroadcasterID: "1234", ModeratorID: "5678"},
			``,
			"id must be provided",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RemoveBlockedTerm(testCase.params)
		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
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

	_, err := c.RemoveBlockedTerm(&RemoveBlockedTermParams{BroadcasterID: "1234", ModeratorID: "1234", ID: "test"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestDeleteChatMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		ModeratorID   string
		MessageID     string
		errorMsg      string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"",
			"test-moderator-id",
			"test-message-id",
			"broadcaster id and moderator id must be provided",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"test-broadcaster-id",
			"",
			"test-message-id",
			"broadcaster id and moderator id must be provided",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"test-broadcaster-id",
			"test-moderator-id",
			"",
			"message id must be provided",
		},
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id"},
			"test-broadcaster-id",
			"test-moderator-id",
			"test-message-id",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, "", nil))

		resp, err := c.DeleteChatMessage(&DeleteChatMessageParams{
			BroadcasterID: testCase.BroadcasterID,
			ModeratorID:   testCase.ModeratorID,
			MessageID:     testCase.MessageID,
		})

		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
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

	_, err := c.RemoveBlockedTerm(&RemoveBlockedTermParams{BroadcasterID: "1234", ModeratorID: "1234", ID: "test"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestDeleteAllChatMessages(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		ModeratorID   string
		errorMsg      string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"",
			"test-moderator-id",
			"broadcaster id and moderator id must be provided",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"test-broadcaster-id",
			"",
			"broadcaster id and moderator id must be provided",
		},
		{
			http.StatusNoContent,
			&Options{ClientID: "my-client-id"},
			"test-broadcaster-id",
			"test-moderator-id",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, "", nil))

		resp, err := c.DeleteAllChatMessages(&DeleteAllChatMessagesParams{
			BroadcasterID: testCase.BroadcasterID,
			ModeratorID:   testCase.ModeratorID,
		})

		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
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

	_, err := c.RemoveBlockedTerm(&RemoveBlockedTermParams{BroadcasterID: "1234", ModeratorID: "1234", ID: "test"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetModerators(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetModeratorsParams
		respBody   string
		parsed     []Moderator
		errorMsg   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&GetModeratorsParams{BroadcasterID: "424596340", First: 2},
			`{
				"data": [
					{
						"user_id": "424596340",
						"user_login": "quotrok",
						"user_name": "quotrok"
					}
				],
				"pagination": {
				  "cursor": "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6I..."
				}
			}`,
			[]Moderator{
				{
					UserID:    "424596340",
					UserLogin: "quotrok",
					UserName:  "quotrok",
				},
			},
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&GetModeratorsParams{BroadcasterID: ""},
			``,
			[]Moderator{},
			"broadcaster id must be provided",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetModerators(testCase.params)

		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		for i, moderator := range resp.Data.Moderators {
			if moderator.UserID != testCase.parsed[i].UserID {
				t.Errorf("Expected struct field UserID = %s, was %s", testCase.parsed[i].UserID, moderator.UserID)
			}

			if moderator.UserLogin != testCase.parsed[i].UserLogin {
				t.Errorf("Expected struct field BroadcasterName = %s, was %s", testCase.parsed[i].UserLogin, moderator.UserLogin)
			}

			if moderator.UserName != testCase.parsed[i].UserName {
				t.Errorf("Expected struct field BroadcasterLanguage = %s, was %s", testCase.parsed[i].UserName, moderator.UserName)
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

	_, err := c.GetModerators(&GetModeratorsParams{BroadcasterID: "424596340", First: 2})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestAddChannelModerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *AddChannelModeratorParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&AddChannelModeratorParams{
				BroadcasterID: "123",
				UserID:        "456",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&AddChannelModeratorParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.AddChannelModerator(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
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

	_, err := c.AddChannelModerator(&AddChannelModeratorParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestRemoveChannelModerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *RemoveChannelModeratorParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&RemoveChannelModeratorParams{
				BroadcasterID: "123",
				UserID:        "456",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&RemoveChannelModeratorParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.RemoveChannelModerator(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
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

	_, err := c.RemoveChannelModerator(&RemoveChannelModeratorParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetModeratedChannels(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetModeratedChannelsParams
		respBody   string
		parsed     *ManyModeratedChannels
		errorMsg   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderatedchannels-access-token"},
			&GetModeratedChannelsParams{UserID: "154315414", First: 2},
			`{
				"data": [
					{
						"broadcaster_id": "183094685",
						"broadcaster_login": "spaceashes",
						"broadcaster_name": "spaceashes"
					},
					{
						"broadcaster_id": "113944563",
						"broadcaster_login": "reapex_1",
						"broadcaster_name": "Reapex_1"
					}
				],
				"pagination": {
					"cursor": "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6ImV5SjBjQ0k2SW5WelpYSTZNVFUwTXpFMU5ERTBPbTF2WkdWeVlYUmxjeUlzSW5Seklqb2lZMmhoYm01bGJEb3hNVE01TkRRMU5qTWlMQ0pwY0NJNkluVnpaWEk2TVRVME16RTFOREUwT20xdlpHVnlZWFJsY3lJc0ltbHpJam9pTVRjeE5EVXhNelF4T0RFNE9UTXlPREV4TnlKOSJ9fQ"
				}
			}`,
			&ManyModeratedChannels{
				ModeratedChannels: []ModeratedChannel{
					{
						BroadcasterID:    "183094685",
						BroadcasterLogin: "spaceashes",
						BroadcasterName:  "spaceashes",
					},
					{
						BroadcasterID:    "113944563",
						BroadcasterLogin: "reapex_1",
						BroadcasterName:  "Reapex_1",
					},
				},
				Pagination: Pagination{
					Cursor: "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6ImV5SjBjQ0k2SW5WelpYSTZNVFUwTXpFMU5ERTBPbTF2WkdWeVlYUmxjeUlzSW5Seklqb2lZMmhoYm01bGJEb3hNVE01TkRRMU5qTWlMQ0pwY0NJNkluVnpaWEk2TVRVME16RTFOREUwT20xdlpHVnlZWFJsY3lJc0ltbHpJam9pTVRjeE5EVXhNelF4T0RFNE9UTXlPREV4TnlKOSJ9fQ",
				},
			},
			"",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderatedchannels-access-token"},
			&GetModeratedChannelsParams{UserID: "154315414", After: "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6ImV5SjBjQ0k2SW5WelpYSTZNVFUwTXpFMU5ERTBPbTF2WkdWeVlYUmxjeUlzSW5Seklqb2lZMmhoYm01bGJEb3hNVE01TkRRMU5qTWlMQ0pwY0NJNkluVnpaWEk2TVRVME16RTFOREUwT20xdlpHVnlZWFJsY3lJc0ltbHpJam9pTVRjeE5EVXhNelF4T0RFNE9UTXlPREV4TnlKOSJ9fQ"},
			`{
				"data": [
					{
						"broadcaster_id": "106590483",
						"broadcaster_login": "vaiastol",
						"broadcaster_name": "vaiastol"
					}
				],
				"pagination": {}
			}`,
			&ManyModeratedChannels{
				ModeratedChannels: []ModeratedChannel{
					{
						BroadcasterID:    "106590483",
						BroadcasterLogin: "vaiastol",
						BroadcasterName:  "vaiastol",
					},
				},
			},
			"",
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id", UserAccessToken: "invalid-access-token"},
			&GetModeratedChannelsParams{UserID: "154315414"},
			`{"error":"Unauthorized","status":401,"message":"Invalid OAuth token"}`,
			&ManyModeratedChannels{},
			"",
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderatedchannels-access-token"},
			&GetModeratedChannelsParams{UserID: "123456789"},
			`{"error":"Unauthorized","status":401,"message":"The ID in user_id must match the user ID found in the request's OAuth token."}`,
			&ManyModeratedChannels{},
			"",
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id", UserAccessToken: "missingscope-access-token"},
			&GetModeratedChannelsParams{UserID: "154315414"},
			`{"error":"Unauthorized","status":401,"message":"Missing scope: user:read:moderated_channels"}`,
			&ManyModeratedChannels{},
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderatedchannels-access-token"},
			&GetModeratedChannelsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"user_id\""}`,
			&ManyModeratedChannels{},
			"user id is required",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetModeratedChannels(testCase.params)

		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		for i, channel := range resp.Data.ModeratedChannels {
			if channel.BroadcasterID != testCase.parsed.ModeratedChannels[i].BroadcasterID {
				t.Errorf("Expected ModeratedChannel field BroadcasterID = %s, was %s", testCase.parsed.ModeratedChannels[i].BroadcasterID, channel.BroadcasterID)
			}

			if channel.BroadcasterLogin != testCase.parsed.ModeratedChannels[i].BroadcasterLogin {
				t.Errorf("Expected ModeratedChannel field BroadcasterLogin = %s, was %s", testCase.parsed.ModeratedChannels[i].BroadcasterLogin, channel.BroadcasterLogin)
			}

			if channel.BroadcasterName != testCase.parsed.ModeratedChannels[i].BroadcasterName {
				t.Errorf("Expected ModeratedChannel field BroadcasterName = %s, was %s", testCase.parsed.ModeratedChannels[i].BroadcasterName, channel.BroadcasterName)
			}
		}

		if resp.Data.Pagination.Cursor != testCase.parsed.Pagination.Cursor {
			t.Errorf("Expected Pagination field Cursor = %s, was %s", testCase.parsed.Pagination.Cursor, resp.Data.Pagination.Cursor)
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

	_, err := c.GetModeratedChannels(&GetModeratedChannelsParams{UserID: "154315414"})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestSendModeratorWarnMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *SendModeratorWarnChatMessageParams
		respBody   string
		errorMsg   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "moderator-access-token"},
			&SendModeratorWarnChatMessageParams{
				BroadcasterID: "1234",
				ModeratorID:   "5678",
				UserID:        "9876",
				Reason:        "Test warning message",
			},
			`{"data": [{"broadcaster_id": "1234", "moderator_id": "5678", "user_id": "9876", "reason": "Test warning message"}]}`,
			"",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "invalid-access-token"},
			&SendModeratorWarnChatMessageParams{
				BroadcasterID: "1234",
				ModeratorID:   "5678",
				Reason:        "Test warning message",
			},
			"",
			"error: user id must be specified",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "invalid-access-token"},
			&SendModeratorWarnChatMessageParams{
				UserID:      "1234",
				ModeratorID: "5678",
				Reason:      "Test warning message",
			},
			"",
			"error: broadcaster id must be specified",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "invalid-access-token"},
			&SendModeratorWarnChatMessageParams{
				UserID:        "1234",
				BroadcasterID: "12345",
				Reason:        "Test warning message",
			},
			"",
			"error: moderator id must be specified",
		},
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id", UserAccessToken: "invalid-access-token"},
			&SendModeratorWarnChatMessageParams{
				BroadcasterID: "1234",
				ModeratorID:   "5678",
				UserID:        "9876",
				Reason:        "Test warning message",
			},
			`{"error":"Unauthorized","status":401,"message":"Invalid OAuth token"}`,
			"",
		},
		// Add more test cases as needed
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendModeratorWarnMessage(testCase.params)
		if err != nil {
			if err.Error() != testCase.errorMsg {
				t.Errorf("expected error message to be %s, got %s", testCase.errorMsg, err.Error())
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode > http.StatusOK {
			continue
		}

		if len(resp.Data.Warnings) == 0 {
			continue
		}
		warning := resp.Data.Warnings[0]

		if warning.BroadcasterID != testCase.params.BroadcasterID {
			t.Errorf("expected broadcaster id to be %s, got %s", testCase.params.BroadcasterID, warning.BroadcasterID)
		}

		if warning.ModeratorID != testCase.params.ModeratorID {
			t.Errorf("expected moderator id to be %s, got %s", testCase.params.ModeratorID, warning.ModeratorID)
		}

		if warning.UserID != testCase.params.UserID {
			t.Errorf("expected user id to be %s, got %s", testCase.params.UserID, warning.UserID)
		}

		if warning.Reason != testCase.params.Reason {
			t.Errorf("expected reason to be %s, got %s", testCase.params.Reason, warning.Reason)
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

	_, err := c.SendModeratorWarnMessage(&SendModeratorWarnChatMessageParams{
		BroadcasterID: "1234",
		ModeratorID:   "5678",
		UserID:        "9876",
		Reason:        "Test warning message",
	})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
