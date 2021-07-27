package helix

import (
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
	}

	_, err := c.GetHypeTrainEvents(&HypeTrainEventsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
