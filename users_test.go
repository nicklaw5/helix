package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		IDs         []string
		Logins      []string
		respBody    string
		expectUsers []string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			[]string{},
			[]string{},
			`{"error":"Bad Request","status":400,"message":"Must provide an ID, Login or OAuth Token"}`,
			[]string{},
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			[]string{"26301881"},
			[]string{"summit1g"},
			`{"data":[{"id":"26301881","login":"sodapoppin","display_name":"sodapoppin","type":"","broadcaster_type":"partner","description":"Wtf do i write here? Click my stream, or i scream.","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-profile_image-10049b6200f90c14-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-channel_offline_image-2040c6fcacec48db-1920x1080.jpeg","view_count":190154823,"created_at":"2011-11-22T04:40:56.75883Z"},{"id":"26490481","login":"summit1g","display_name":"summit1g","type":"","broadcaster_type":"partner","description":"I'm a competitive CounterStrike player who likes to play casually now and many other games. You will mostly see me play CS, H1Z1,and single player games at night. There will be many othergames played on this stream in the future as they come out:D","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/200cea12142f2384-profile_image-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/summit1g-channel_offline_image-e2f9a1df9e695ec1-1920x1080.png","view_count":202707885,"created_at":"2011-12-01T06:33:31.487567Z"}]}`,
			[]string{"sodapoppin", "summit1g"},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUsers(&UsersParams{
			IDs:    testCase.IDs,
			Logins: testCase.Logins,
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

			expectedErrMsg := "Must provide an ID, Login or OAuth Token"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if resp.Data.Users[0].Login != testCase.expectUsers[0] { // sodapoppin
			t.Errorf("expected username 1 to be \"%s\", got \"%s\"", testCase.expectUsers[0], resp.Data.Users[0].Login)
		}

		if resp.Data.Users[1].Login != testCase.expectUsers[1] { // summit1g
			t.Errorf("expected username 2 to be \"%s\", got \"%s\"", testCase.expectUsers[0], resp.Data.Users[0].Login)
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

	_, err := c.GetUsers(&UsersParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		respBody    string
		description string
	}{
		{
			http.StatusForbidden,
			&Options{ClientID: "my-client-id"},
			`{"error":"Forbidden","status":403,"message":"Missing user:edit scope"}`, // missing required scope
			"new description",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			`{"data":[{"id":"26301881","login":"sodapoppin","display_name":"sodapoppin","type":"","broadcaster_type":"partner","description":"new description","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-profile_image-10049b6200f90c14-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-channel_offline_image-2040c6fcacec48db-1920x1080.jpeg","view_count":190154823,"created_at":"2011-11-22T04:40:56.75883Z"}]}`,
			"new description",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateUser(&UpdateUserParams{
			Description: testCase.description,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusForbidden {
			if resp.Error != "Forbidden" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusForbidden {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusForbidden, resp.ErrorStatus)
			}

			expectedErrMsg := "Missing user:edit scope"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if resp.Data.Users[0].Description != testCase.description {
			t.Errorf("expected description to be \"%s\", got \"%s\"", testCase.description, resp.Data.Users[0].Description)
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

	_, err := c.UpdateUser(&UpdateUserParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetUsersFollows(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		FromID     string
		First      int
		respBody   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"", // missing from_id
			2,
			`{"error":"Bad Request","status":400,"message":"Must provide either from_id or to_id"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"23161357",
			2,
			`{"total":89,"data":[{"from_id":"23161357","from_name":"lirik","to_id":"23528098","to_name":"avoidingthepuddle","followed_at":"2017-10-01T03:57:21Z"},{"from_id":"23161357","from_name":"lirik","to_id":"127506955","to_name":"playbattlegrounds","followed_at":"2017-08-23T15:04:20Z"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjoiMTUwMzUwMDY2MDYwNzAyNTAwMCJ9"}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUsersFollows(&UsersFollowsParams{
			First:  testCase.First,
			FromID: testCase.FromID,
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

			expectedErrMsg := "Must provide either from_id or to_id"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Follows) != testCase.First {
			t.Errorf("expected result length to be \"%d\", got \"%d\"", testCase.First, len(resp.Data.Follows))
		}

		for _, follow := range resp.Data.Follows {
			if follow.FromID != testCase.FromID {
				t.Errorf("expected from_id to be \"%s\", got \"%s\"", testCase.FromID, follow.FromID)
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

	_, err := c.GetUsersFollows(&UsersFollowsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetUsersBlocked(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		First         int
		respBody      string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"",
			1,
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"23161357",
			1,
			`{"data":[{"user_id":"199340135","user_login":"jlarkyzus","display_name":"JLArkyzus"}],"pagination":{"cursor":"xxx"}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUsersBlocked(&UsersBlockedParams{
			First:         testCase.First,
			BroadcasterID: testCase.BroadcasterID,
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

			expectedErrMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Users) != testCase.First {
			t.Errorf("expected result length to be \"%d\", got \"%d\"", testCase.First, len(resp.Data.Users))
		}

		userID := "199340135"
		if resp.Data.Users[0].UserID != userID {
			t.Errorf("expected user id to be \"%s\", got \"%s\"", userID, resp.Data.Users[0].UserID)
		}

		userLogin := "jlarkyzus"
		if resp.Data.Users[0].UserLogin != userLogin {
			t.Errorf("expected user id to be \"%s\", got \"%s\"", userLogin, resp.Data.Users[0].UserLogin)
		}

		displayName := "JLArkyzus"
		if resp.Data.Users[0].DisplayName != displayName {
			t.Errorf("expected user id to be \"%s\", got \"%s\"", displayName, resp.Data.Users[0].DisplayName)
		}

		cursor := "xxx"
		if resp.Data.Pagination.Cursor != cursor {
			t.Errorf("expected cursor to be \"%s\", got \"%s\"", cursor, resp.Data.Pagination.Cursor)
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

	_, err := c.GetUsersBlocked(&UsersBlockedParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestBlockUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		Params     *BlockUserParams
		respBody   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&BlockUserParams{
				TargetUserID:  "199340135",
				SourceContext: "twitter",
				Reason:        "spam",
			},
			`{"error":"Bad Request","status":400,"message":"The parameter \"source_context\" was malformed: value must be one of \"chat\", \"whisper\", \"\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&BlockUserParams{
				TargetUserID:  "199340135",
				SourceContext: "chat",
				Reason:        "spam",
			},
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.BlockUser(testCase.Params)
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

			expectedErrMsg := "The parameter \"source_context\" was malformed: value must be one of \"chat\", \"whisper\", \"\""
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

	_, err := c.BlockUser(&BlockUserParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUnblockUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		Params     *UnblockUserParams
		respBody   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&UnblockUserParams{
				TargetUserID: "",
			},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"target_user_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&UnblockUserParams{
				TargetUserID: "199340135",
			},
			``,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UnblockUser(testCase.Params)
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

			expectedErrMsg := "Missing required parameter \"target_user_id\""
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

	_, err := c.UnblockUser(&UnblockUserParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
