package helix

import (
	"net/http"
	"testing"
	"time"
)

func TestClient_GetBitsLeaderboard(t *testing.T) {
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
			`{"data":[{"user_id":"158010205","user_login":"tundracowboy","user_name":"TundraCowboy","rank":1,"score":12543},{"user_id":"7168163","rank":2,"score":6900}],"date_range":{"started_at":"2018-02-05T08:00:00Z","ended_at":"2018-02-12T08:00:00Z"},"total":2}`,
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

	_, err := c.GetBitsLeaderboard(&BitsLeaderboardParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestClient_GetCheermotes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode        int
		options           *Options
		count             int
		broadcasterID     string
		initialPrefix     string
		initialOrder      uint
		initialTiersCount int
		respBody          string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			"",
			"Cheer",
			1,
			5,
			`{
  "data": [
    {
      "prefix": "Cheer",
      "tiers": [
        {
          "min_bits": 1,
          "id": "1",
          "color": "#979797",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": true
        },
        {
          "min_bits": 100,
          "id": "100",
          "color": "#9c3ee8",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/100/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/100/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/100/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/100/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/100/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/100/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/100/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/100/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/100/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/100/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/100/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/100/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/100/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/100/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/100/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/100/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/100/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/100/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/100/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/100/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": true
        },
        {
          "min_bits": 1000,
          "id": "1000",
          "color": "#1db2a5",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1000/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1000/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": true
        },
        {
          "min_bits": 5000,
          "id": "5000",
          "color": "#0099fe",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/5000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/5000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/5000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/5000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/5000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/5000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/5000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/5000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/5000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/5000/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/5000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/5000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/5000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/5000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/5000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/5000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/5000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/5000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/5000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/5000/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": true
        },
        {
          "min_bits": 10000,
          "id": "10000",
          "color": "#f43021",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/10000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/10000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/10000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/10000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/10000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/10000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/10000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/10000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/10000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/10000/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/10000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/10000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/10000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/10000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/10000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/10000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/10000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/10000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/10000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/10000/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": true
        }
      ],
      "type": "global_first_party",
      "order": 1,
      "last_updated": "2018-05-22T00:06:04Z",
      "is_charitable": false
    },
    {
      "prefix": "DoodleCheer",
      "tiers": [
        {
          "min_bits": 1,
          "id": "1",
          "color": "#979797",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": false
        },
        {
          "min_bits": 100,
          "id": "100",
          "color": "#9c3ee8",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/100/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/100/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/100/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/100/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/100/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/100/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/100/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/100/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/100/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/100/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/100/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/100/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/100/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/100/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/100/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/100/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/100/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/100/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/100/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/100/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": false
        },
        {
          "min_bits": 1000,
          "id": "1000",
          "color": "#1db2a5",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/1000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/1000/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/1000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/1000/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": false
        },
        {
          "min_bits": 5000,
          "id": "5000",
          "color": "#0099fe",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/5000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/5000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/5000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/5000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/5000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/5000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/5000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/5000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/5000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/5000/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/5000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/5000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/5000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/5000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/5000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/5000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/5000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/5000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/5000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/5000/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": false
        },
        {
          "min_bits": 10000,
          "id": "10000",
          "color": "#f43021",
          "images": {
            "dark": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/10000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/10000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/10000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/10000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/animated/10000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/10000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/10000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/10000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/10000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/dark/static/10000/4.png"
              }
            },
            "light": {
              "animated": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/10000/1.gif",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/10000/1.5.gif",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/10000/2.gif",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/10000/3.gif",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/animated/10000/4.gif"
              },
              "static": {
                "1": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/10000/1.png",
                "1.5": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/10000/1.5.png",
                "2": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/10000/2.png",
                "3": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/10000/3.png",
                "4": "https://d3aqoihi2n8ty8.cloudfront.net/actions/doodlecheer/light/static/10000/4.png"
              }
            }
          },
          "can_cheer": true,
          "show_in_bits_card": false
        }
      ],
      "type": "global_third_party",
      "order": 1,
      "last_updated": "2018-05-22T00:06:05Z",
      "is_charitable": false
    }
  ]
}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		params := &CheermotesParams{
			BroadcasterID: testCase.broadcasterID,
		}

		resp, err := c.GetCheermotes(params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// TODO: Test error cases

		// Test success cases
		if len(resp.Data.Cheermotes) != testCase.count {
			t.Errorf("expected number of results to be \"%d\", got \"%d\"", testCase.count, len(resp.Data.Cheermotes))
		}

		cheermotes := resp.Data.Cheermotes[0]
		if cheermotes.Prefix != testCase.initialPrefix {
			t.Errorf("expected prefix of \"%s\", got \"%s\"", testCase.initialPrefix, cheermotes.Prefix)
		}

		if cheermotes.Order != testCase.initialOrder {
			t.Errorf("expected order of \"%d\", got \"%d\"", testCase.initialOrder, cheermotes.Order)
		}

		if len(cheermotes.Tiers) != testCase.initialTiersCount {
			t.Errorf("expected tier count of \"%d\", got \"%d\"", testCase.initialTiersCount, len(cheermotes.Tiers))
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

	_, err := c.GetCheermotes(&CheermotesParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
