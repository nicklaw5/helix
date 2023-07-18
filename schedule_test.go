package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetScheduleParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			&GetScheduleParams{
				BroadcasterID: "141981764",
			},
			`{
				"data": {
				  "segments": [
					{
					  "id": "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=",
					  "start_time": "2021-07-01T18:00:00Z",
					  "end_time": "2021-07-01T19:00:00Z",
					  "title": "TwitchDev Monthly Update // July 1, 2021",
					  "canceled_until": null,
					  "category": {
						  "id": "509670",
						  "name": "Science & Technology"
					  },
					  "is_recurring": false
					}
				  ],
				  "broadcaster_id": "141981764",
				  "broadcaster_name": "TwitchDev",
				  "broadcaster_login": "twitchdev",
				  "vacation": null
				},
				"pagination": {}
			}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetSchedule(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusOK {
			if resp.Data.Schedule.BroadcasterID != testCase.params.BroadcasterID {
				t.Errorf("Expected broadcasterID = %s, got %s", testCase.params.BroadcasterID, resp.Data.Schedule.BroadcasterID)
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

	_, err := c.GetSchedule(&GetScheduleParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateSchedule(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *UpdateScheduleParams
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			&UpdateScheduleParams{
				BroadcasterID:     "141981764",
				IsVacationEnabled: true,
				Timezone:          "America/New_York",
			},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, ``, nil))

		resp, err := c.UpdateSchedule(testCase.params)
		if err != nil {
			t.Error(err)
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

	_, err := c.UpdateSchedule(&UpdateScheduleParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestCreateScheduleSegment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *CreateScheduleSegmentParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			&CreateScheduleSegmentParams{
				BroadcasterID: "141981764",
				Timezone:      "America/New_York",
				IsRecurring:   false,
				Duration:      "60",
				CategoryID:    "509670",
				Title:         "TwitchDev Monthly Update // July 1, 2021",
			},
			`{
				"data": {
				  "segments": [
					{
					  "id": "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=",
					  "start_time": "2021-07-01T18:00:00Z",
					  "end_time": "2021-07-01T19:00:00Z",
					  "title": "TwitchDev Monthly Update // July 1, 2021",
					  "canceled_until": null,
					  "category": {
						  "id": "509670",
						  "name": "Science & Technology"
					  },
					  "is_recurring": false
					}
				  ],
				  "broadcaster_id": "141981764",
				  "broadcaster_name": "TwitchDev",
				  "broadcaster_login": "twitchdev",
				  "vacation": null
				}
			}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreateScheduleSegment(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusOK {
			if resp.Data.Schedule.BroadcasterID != testCase.params.BroadcasterID {
				t.Errorf("Expected broadcasterID = %s, got %s", testCase.params.BroadcasterID, resp.Data.Schedule.BroadcasterID)
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

	_, err := c.CreateScheduleSegment(&CreateScheduleSegmentParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateScheduleSegment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *UpdateScheduleSegmentParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			&UpdateScheduleSegmentParams{
				ID:            "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=",
				BroadcasterID: "141981764",
				Duration:      "120",
			},
			`{
				"data": {
				  "segments": [
					{
					  "id": "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=",
					  "start_time": "2021-07-01T18:00:00Z",
					  "end_time": "2021-07-01T20:00:00Z",
					  "title": "TwitchDev Monthly Update // July 1, 2021",
					  "canceled_until": null,
					  "category": {
						  "id": "509670",
						  "name": "Science & Technology"
					  },
					  "is_recurring": false
					}
				  ],
				  "broadcaster_id": "141981764",
				  "broadcaster_name": "TwitchDev",
				  "broadcaster_login": "twitchdev",
				  "vacation": null
				}
			}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateScheduleSegment(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusOK {
			if resp.Data.Schedule.BroadcasterID != testCase.params.BroadcasterID {
				t.Errorf("Expected broadcasterID = %s, got %s", testCase.params.BroadcasterID, resp.Data.Schedule.BroadcasterID)
			}
			if len(resp.Data.Schedule.Segments) > 0 {
				if resp.Data.Schedule.Segments[0].ID != testCase.params.ID {
					t.Errorf("Expected ID = %s, got %s", testCase.params.ID, resp.Data.Schedule.Segments[0].ID)
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

	_, err := c.UpdateScheduleSegment(&UpdateScheduleSegmentParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestDeleteScheduleSegment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *DeleteScheduleSegmentParams
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			&DeleteScheduleSegmentParams{
				BroadcasterID: "141981764",
				ID:            "eyJzZWdtZW50SUQiOiI4Y2EwN2E2NC0xYTZkLTRjYWItYWE5Ni0xNjIyYzNjYWUzZDkiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyMX0=",
			},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, ``, nil))

		resp, err := c.DeleteScheduleSegment(testCase.params)
		if err != nil {
			t.Error(err)
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

	_, err := c.DeleteScheduleSegment(&DeleteScheduleSegmentParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
