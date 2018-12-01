package helix

import (
	"net/http"
	"testing"
	"time"
)

func TestGetBitsLeaderboard(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		count      int
		period     string
		startedAt  time.Time
		respBody   string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			1,
			"all",
			time.Time{},
			`{"error":"Bad Request","status":400,"message":"The parameter \"count\" was malformed: the value must be less than or equal to 100"}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			"week",
			time.Time{},
			`{"data":[{"user_id":"158010205","user_name":"TundraCowboy","rank":1,"score":12543},{"user_id":"7168163","rank":2,"score":6900}],"date_range":{"started_at":"2018-02-05T08:00:00Z","ended_at":"2018-02-12T08:00:00Z"},"total":2}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			"week",
			time.Now().Add(-744 * time.Hour),
			`{"data":[{"user_id":"158010205","user_name":"TundraCowboy","rank":1,"score":12543},{"user_id":"7168163","rank":2,"score":6900}],"date_range":{"started_at":"2018-02-05T08:00:00Z","ended_at":"2018-02-12T08:00:00Z"},"total":2}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		params := &BitsLeaderboardParams{
			Count:  testCase.count,
			Period: testCase.period,
		}

		if !testCase.startedAt.IsZero() {
			params.StartedAt = testCase.startedAt
		}

		resp, err := c.GetBitsLeaderboard(params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if testCase.statusCode != http.StatusOK {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			errMsg := "The parameter \"count\" was malformed: the value must be less than or equal to 100"
			if resp.ErrorMessage != errMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", errMsg, resp.ErrorMessage)
			}

			continue
		}

		// Test success cases
		if len(resp.Data.UserBitTotals) != testCase.count {
			t.Errorf("expected number of results to be \"%d\", got \"%d\"", testCase.count, len(resp.Data.UserBitTotals))
		}

		userData := resp.Data.UserBitTotals[0]
		if userData.UserID != "158010205" || userData.Rank != 1 || userData.Score != 12543 {
			t.Error("expected bits user data does not match expected values")
		}

		if resp.Data.DateRange.EndedAt.IsZero() {
			t.Error("expected DateRange.EndedAt to not be zero")
		}

		if resp.Data.DateRange.StartedAt.IsZero() {
			t.Error("expected DateRange.Started to not be zero")
		}

		if resp.Data.Total < 1 {
			t.Error("expected Total to be more than zero")
		}

	}
}
