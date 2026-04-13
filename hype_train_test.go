package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetHypeTrainEvents(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		PollsParams *HypeTrainEventsParams
		respBody    string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&HypeTrainEventsParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&HypeTrainEventsParams{BroadcasterID: "121445595"},
			`{"data":[{"id":"1vsv0LoqZEl1YVRUKehGesQJQFY","event_type":"hypetrain.progression","event_timestamp":"2021-07-27T05:23:06Z","version":"1.0","event_data":{"broadcaster_id":"32043018","cooldown_end_time":"2021-07-27T06:27:11Z","expires_at":"2021-07-27T05:27:11Z","goal":1600,"id":"207ab35e-26f3-4b3a-b4c5-0c9025079906","last_contribution":{"total":100,"type":"BITS","user":"36099963"},"level":1,"started_at":"2021-07-27T05:22:11Z","top_contributions":[{"total":100,"type":"BITS","user":"86228897"},{"total":500,"type":"SUBS","user":"36099963"}],"total":1300}},{"id":"1vsuvsQkf1dFPbOXcxvlS3NiAum","event_type":"hypetrain.progression","event_timestamp":"2021-07-27T05:22:30Z","version":"1.0","event_data":{"broadcaster_id":"32043018","cooldown_end_time":"2021-07-27T06:27:11Z","expires_at":"2021-07-27T05:27:11Z","goal":1600,"id":"207ab35e-26f3-4b3a-b4c5-0c9025079906","last_contribution":{"total":100,"type":"BITS","user":"44526173"},"level":1,"started_at":"2021-07-27T05:22:11Z","top_contributions":[{"total":100,"type":"BITS","user":"86228897"},{"total":500,"type":"SUBS","user":"36099963"}],"total":1200}},{"id":"1vsutOWymap8tfJMTuihbwzBQHC","event_type":"hypetrain.progression","event_timestamp":"2021-07-27T05:22:11Z","version":"1.0","event_data":{"broadcaster_id":"32043018","cooldown_end_time":"2021-07-27T06:27:11Z","expires_at":"2021-07-27T05:27:11Z","goal":1600,"id":"207ab35e-26f3-4b3a-b4c5-0c9025079906","last_contribution":{"total":500,"type":"SUBS","user":"158357643"},"level":1,"started_at":"2021-07-27T05:22:11Z","top_contributions":[{"total":100,"type":"BITS","user":"86228897"},{"total":500,"type":"SUBS","user":"36099963"}],"total":1100}},{"id":"1vsutUYGOAFjhD3U8d0NbG4BbxA","event_type":"hypetrain.progression","event_timestamp":"2021-07-27T05:22:11Z","version":"1.0","event_data":{"broadcaster_id":"32043018","cooldown_end_time":"2021-07-27T06:27:11Z","expires_at":"2021-07-27T05:27:11Z","goal":1600,"id":"207ab35e-26f3-4b3a-b4c5-0c9025079906","last_contribution":{"total":500,"type":"SUBS","user":"158357643"},"level":1,"started_at":"2021-07-27T05:22:11Z","top_contributions":[{"total":500,"type":"SUBS","user":"36099963"}],"total":1100}},{"id":"1vsutTHkvLoRBbRpIPNuu6SGZzd","event_type":"hypetrain.progression","event_timestamp":"2021-07-27T05:22:11Z","version":"1.0","event_data":{"broadcaster_id":"32043018","cooldown_end_time":"2021-07-27T06:27:11Z","expires_at":"2021-07-27T05:27:11Z","goal":1600,"id":"207ab35e-26f3-4b3a-b4c5-0c9025079906","last_contribution":{"total":500,"type":"SUBS","user":"158357643"},"level":1,"started_at":"2021-07-27T05:22:11Z","top_contributions":[{"total":500,"type":"SUBS","user":"36099963"}],"total":1100}}],"pagination":null}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetHypeTrainEvents(testCase.PollsParams)
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

		if len(resp.Data.Events) != 5 {
			t.Errorf("expected hype train events len to be 5, got %d", len(resp.Data.Events))
		}

		if len(resp.Data.Events[0].Event.TopContributions) != 2 {
			t.Errorf("expected hype train event top contributors len to be 2, got %d", len(resp.Data.Events[0].Event.TopContributions))
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

	_, err := c.GetHypeTrainEvents(&HypeTrainEventsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetHypeTrainStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *HypeTrainStatusParams
		respBody   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&HypeTrainStatusParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&HypeTrainStatusParams{BroadcasterID: "1337"},
			`{"data":[{"current":{"id":"1b0AsbInCHZW2SQFQkCzqN07Ib2","broadcaster_user_id":"1337","broadcaster_user_login":"cool_user","broadcaster_user_name":"Cool_User","level":2,"total":700,"progress":200,"goal":1000,"top_contributions":[{"user_id":"123","user_login":"pogchamp","user_name":"PogChamp","type":"bits","total":50},{"user_id":"456","user_login":"kappa","user_name":"Kappa","type":"subscription","total":45}],"shared_train_participants":[{"broadcaster_user_id":"456","broadcaster_user_login":"pogchamp","broadcaster_user_name":"PogChamp"},{"broadcaster_user_id":"321","broadcaster_user_login":"pogchamp","broadcaster_user_name":"PogChamp"}],"started_at":"2020-07-15T17:16:03.17106713Z","expires_at":"2020-07-15T17:16:11.17106713Z","type":"golden_kappa"},"all_time_high":{"level":6,"total":2850,"achieved_at":"2020-04-24T20:12:21.003802269Z"},"shared_all_time_high":{"level":16,"total":23850,"achieved_at":"2020-04-27T20:12:21.003802269Z"}}]}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&HypeTrainStatusParams{BroadcasterID: "1338"},
			`{"data":[]}`, // No active hype train
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetHypeTrainStatus(testCase.params)
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

		if testCase.params.BroadcasterID == "1337" {
			// Test case with active hype train
			if len(resp.Data.Statuses) != 1 {
				t.Errorf("expected hype train statuses len to be 1, got %d", len(resp.Data.Statuses))
			}

			status := resp.Data.Statuses[0]
			if status.Current == nil {
				t.Error("expected current hype train status to not be nil")
			} else {
				if status.Current.BroadcasterUserID != "1337" {
					t.Errorf("expected broadcaster_user_id to be '1337', got '%s'", status.Current.BroadcasterUserID)
				}

				if status.Current.Level != 2 {
					t.Errorf("expected level to be 2, got %d", status.Current.Level)
				}

				if len(status.Current.TopContributions) != 2 {
					t.Errorf("expected top contributions len to be 2, got %d", len(status.Current.TopContributions))
				}

				if len(status.Current.SharedTrainParticipants) != 2 {
					t.Errorf("expected shared train participants len to be 2, got %d", len(status.Current.SharedTrainParticipants))
				}
			}

			if status.AllTimeHigh == nil {
				t.Error("expected all_time_high to not be nil")
			} else {
				if status.AllTimeHigh.Level != 6 {
					t.Errorf("expected all_time_high level to be 6, got %d", status.AllTimeHigh.Level)
				}
			}

			if status.SharedAllTimeHigh == nil {
				t.Error("expected shared_all_time_high to not be nil")
			} else {
				if status.SharedAllTimeHigh.Level != 16 {
					t.Errorf("expected shared_all_time_high level to be 16, got %d", status.SharedAllTimeHigh.Level)
				}
			}
		} else if testCase.params.BroadcasterID == "1338" {
			// Test case with no active hype train
			if len(resp.Data.Statuses) != 0 {
				t.Errorf("expected hype train statuses len to be 0, got %d", len(resp.Data.Statuses))
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

	_, err := c.GetHypeTrainStatus(&HypeTrainStatusParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
