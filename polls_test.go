package helix

import (
	"net/http"
	"testing"
)

func TestGetPolls(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		PollsParams *PollsParams
		respBody    string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&PollsParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&PollsParams{BroadcasterID: "121445595"},
			`{"data":[{"id":"ed961efd-8a3f-4cf5-a9d0-e616c590cd2a","broadcaster_id":"55696719","broadcaster_name":"TwitchDev","broadcaster_login":"twitchdev","title":"Heads or Tails?","choices":[{"id": "4c123012-1351-4f33-84b7-43856e7a0f47","title": "Heads","votes": 0,"channel_points_votes": 0,"bits_votes": 0},{"id": "279087e3-54a7-467e-bcd0-c1393fcea4f0","title": "Tails","votes": 0,"channel_points_votes": 0,"bits_votes": 0}],"bits_voting_enabled": false,"bits_per_vote": 0,"channel_points_voting_enabled": false,"channel_points_per_vote": 0,"status": "ACTIVE","duration": 1800,"started_at": "2021-03-19T06:08:33.871278372Z"}],"pagination": {}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetPolls(testCase.PollsParams)
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

	_, err := c.GetPolls(&PollsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestCreatePoll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode       int
		options          *Options
		CreatePollParams *CreatePollParams
		respBody         string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&CreatePollParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&CreatePollParams{
				BroadcasterID: "145328278",
				Title:         "Test",
				Choices: []PollChoiceParam{
					PollChoiceParam{Title: "choix 1"},
					PollChoiceParam{Title: "choix 2"},
				},
				Duration: 30,
			},
			`{"data":[{"id":"fb156390-a4de-4acb-8ca6-52ef125f533b","broadcaster_id":"145328278","broadcaster_name":"Scorfly","broadcaster_login":"scorfly","title":"Test","choices":[{"id":"3fea0835-6059-4bb2-95ab-1318659f0282","title":"choix 1","votes":0,"channel_points_votes":0,"bits_votes":0},{"id":"fc4d6457-f32b-492b-a93e-aaa88c343598","title":"choix 2","votes":0,"channel_points_votes":0,"bits_votes":0}],"bits_voting_enabled":false,"bits_per_vote":0,"channel_points_voting_enabled":false,"channel_points_per_vote":0,"status":"ACTIVE","duration":30,"started_at":"2021-05-06T20:43:24.60506479Z"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreatePoll(testCase.CreatePollParams)
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

	_, err := c.CreatePoll(&CreatePollParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestEndPoll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		EndPollParams *EndPollParams
		respBody      string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&EndPollParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&EndPollParams{
				BroadcasterID: "145328278",
				ID:            "25b14b42-d4d8-4756-86ce-842bf76f82a0",
				Status:        "TERMINATED",
			},
			`{"data":[{"id":"6aee6aae-e536-4eb0-afa5-d64567aec2c6","broadcaster_id":"145328278","broadcaster_name":"Scorfly","broadcaster_login":"scorfly","title":"Test","choices":[{"id":"cdebad56-ea5a-4d7a-8caf-cdf80c71514e","title":"choix 1","votes":0,"channel_points_votes":0,"bits_votes":0},{"id":"e3027452-e3ab-4ee4-bc4a-03d9bac37dcc","title":"choix 2","votes":0,"channel_points_votes":0,"bits_votes":0}],"bits_voting_enabled":false,"bits_per_vote":0,"channel_points_voting_enabled":false,"channel_points_per_vote":0,"status":"TERMINATED","duration":300,"started_at":"2021-05-06T21:15:08.661352925Z","ended_at":"2021-05-06T21:15:26.894542904Z"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.EndPoll(testCase.EndPollParams)
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

	_, err := c.EndPoll(&EndPollParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
