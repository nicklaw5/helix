package helix

import (
	"net/http"
	"strconv"
	"testing"
)

func TestGetStreams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		First      int
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			`{"data":[{"id":"27833742640","user_id":"19571641","game_id":"33214","community_ids":[],"type":"live","title":"I have lost my voice D: | twitter.com/Ninja","viewer_count":72124,"started_at":"2018-03-06T15:07:45Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_ninja-{width}x{height}.jpg"},{"id":"27834185424","user_id":"17337557","game_id":"33214","community_ids":[],"type":"live","title":"Turbo Treehouses || @DrDisRespect","viewer_count":29687,"started_at":"2018-03-06T16:05:00Z","language":"en","thumbnail_url":"https://static-cdn.jtvnw.net/previews-ttv/live_user_drdisrespectlive-{width}x{height}.jpg"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			101,
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetStreams(&StreamsParams{
			First: testCase.First,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "The parameter \"first\" was malformed: the value must be less than or equal to 100"
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Streams) != testCase.First {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.First, len(resp.Data.Streams))
		}
	}
}

func TestGetStreamsMetadata(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode          int
		options             *Options
		First               int
		respBody            string
		expectBroadcastHero []string
		expectOpponentHero  []string
		headerLimit         string
		headerRemaining     string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			1,
			`{"data":[{"user_id":"43356746","game_id":"138585","overwatch":null,"hearthstone":{"broadcaster":{"hero":{"type":"Alternate hero","class":"Mage","name":"Medivh"}},"opponent":{"hero":{"type":"Classic hero","class":"Warlock","name":"Guldan"}}}}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6MX19"}}`,
			[]string{"Alternate hero", "Mage", "Medivh"},  // type, class, name
			[]string{"Classic hero", "Warlock", "Guldan"}, // type, class, name
			"15000",
			"14119",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			1,
			`{"data":[{"user_id":"132395117","game_id":"488552","overwatch":{"broadcaster":{"hero":{"role":"Support","name":"Lucio","ability":"Sonic Amplifier"}}},"hearthstone":null}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6MX19"}}`,
			[]string{"Support", "Lucio", "Sonic Amplifier"}, // role, name, ability
			[]string{}, // N/A
			"15000",
			"14119",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			101, // exceeds 100 limit
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
			[]string{}, // N/A
			[]string{}, // N/A
			"0",
			"0",
		},
	}

	for _, testCase := range testCases {
		mockRespHeaders := map[string]string{
			"Ratelimit-Helixstreamsmetadata-Limit":     testCase.headerLimit,
			"Ratelimit-Helixstreamsmetadata-Remaining": testCase.headerRemaining,
		}

		mockHandler := newMockHandler(testCase.statusCode, testCase.respBody, mockRespHeaders)
		c := newMockClient(testCase.options, mockHandler)

		resp, err := c.GetStreamsMetadata(&StreamsMetadataParams{
			First: testCase.First,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test metadata headers response
		headerLimit, _ := strconv.Atoi(testCase.headerLimit)
		if resp.GetStreamsMetadataRateLimit() != headerLimit {
			t.Errorf("expected metadata limit header to be \"%d\", got \"%d\"", headerLimit, resp.GetStreamsMetadataRateLimit())
		}
		headerRemaining, _ := strconv.Atoi(testCase.headerRemaining)
		if resp.GetStreamsMetadataRateLimitRemaining() != headerRemaining {
			t.Errorf("expected metadata remaining header to be \"%d\", got \"%d\"", headerRemaining, resp.GetStreamsMetadataRateLimitRemaining())
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "The parameter \"first\" was malformed: the value must be less than or equal to 100"
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if len(resp.Data.Streams) != testCase.First {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.First, len(resp.Data.Streams))
		}

		// Test Overwatch response
		overwatch := resp.Data.Streams[0].Overwatch
		if overwatch.Broadcaster.Hero.Name != "" {
			// test overwatch hero
			if overwatch.Broadcaster.Hero.Role != testCase.expectBroadcastHero[0] {
				t.Errorf("expected broadcast hero role to be \"%s\", got \"%s\"", overwatch.Broadcaster.Hero.Role, testCase.expectBroadcastHero[0])
			}
			if overwatch.Broadcaster.Hero.Name != testCase.expectBroadcastHero[1] {
				t.Errorf("expected broadcast hero name to be \"%s\", got \"%s\"", overwatch.Broadcaster.Hero.Name, testCase.expectBroadcastHero[1])
			}
			if overwatch.Broadcaster.Hero.Ability != testCase.expectBroadcastHero[2] {
				t.Errorf("expected broadcast hero ability to be \"%s\", got \"%s\"", overwatch.Broadcaster.Hero.Ability, testCase.expectBroadcastHero[2])
			}
		}

		// Test Hearthstone response
		hearthstone := resp.Data.Streams[0].Hearthstone
		if hearthstone.Broadcaster.Hero.Name != "" {
			// test broadcaster hero
			if hearthstone.Broadcaster.Hero.Type != testCase.expectBroadcastHero[0] {
				t.Errorf("expected broadcast hero type to be \"%s\", got \"%s\"", hearthstone.Broadcaster.Hero.Type, testCase.expectBroadcastHero[0])
			}
			if hearthstone.Broadcaster.Hero.Class != testCase.expectBroadcastHero[1] {
				t.Errorf("expected broadcast hero class to be \"%s\", got \"%s\"", hearthstone.Broadcaster.Hero.Class, testCase.expectBroadcastHero[1])
			}
			if hearthstone.Broadcaster.Hero.Name != testCase.expectBroadcastHero[2] {
				t.Errorf("expected broadcast hero name to be \"%s\", got \"%s\"", hearthstone.Broadcaster.Hero.Name, testCase.expectBroadcastHero[2])
			}

			// test opponent hero
			if hearthstone.Opponent.Hero.Type != testCase.expectOpponentHero[0] {
				t.Errorf("expected broadcast hero type to be \"%s\", got \"%s\"", hearthstone.Opponent.Hero.Type, testCase.expectOpponentHero[0])
			}
			if hearthstone.Opponent.Hero.Class != testCase.expectOpponentHero[1] {
				t.Errorf("expected broadcast hero class to be \"%s\", got \"%s\"", hearthstone.Opponent.Hero.Class, testCase.expectOpponentHero[1])
			}
			if hearthstone.Opponent.Hero.Name != testCase.expectOpponentHero[2] {
				t.Errorf("expected broadcast hero name to be \"%s\", got \"%s\"", hearthstone.Opponent.Hero.Name, testCase.expectOpponentHero[2])
			}
		}
	}
}
