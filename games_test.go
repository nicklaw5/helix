package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetGames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		IDs        []string
		Names      []string
		respBody   string
		expectGame string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			[]string{"27471"},
			[]string{},
			`{"data":[{"id":"27471","name":"Minecraft","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/Minecraft-{width}x{height}.jpg"}]}`,
			"Minecraft",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			[]string{},
			[]string{"Sea of Thieves"},
			`{"data":[{"id":"490377","name":"Sea of Thieves","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/Sea%20of%20Thieves-{width}x{height}.jpg"}]}`,
			"Sea of Thieves",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetGames(&GamesParams{
			IDs:   testCase.IDs,
			Names: testCase.Names,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.Data.Games[0].Name != testCase.expectGame {
			t.Errorf("expected game name to be \"%s\", got \"%s\"", testCase.expectGame, resp.Data.Games[0].Name)
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

	_, err := c.GetGames(&GamesParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetTopGames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		options     *Options
		First       int
		AfterCursor string
		respBody    string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			3,
			"",
			`{"data":[{"id":"33214","name":"Fortnite","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg"},{"id":"21779","name":"League of Legends","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/League%20of%20Legends-{width}x{height}.jpg"},{"id":"29595","name":"Dota 2","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/Dota%202-{width}x{height}.jpg"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6M319"}}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			101, // exceeds 100 limit
			"",
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			20,
			"a-non-cursor-string", // invalid cursor
			`{"error":"Bad Request","status":400,"message":"Invalid cursor."}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetTopGames(&TopGamesParams{
			First: testCase.First,
			After: testCase.AfterCursor,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			if testCase.First == 101 {
				firstErrStr := "The parameter \"first\" was malformed: the value must be less than or equal to 100"
				if resp.ErrorMessage != firstErrStr {
					t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
				}

				continue
			}

			errorStr := "Invalid cursor."
			if resp.ErrorMessage != errorStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", errorStr, resp.ErrorMessage)
			}
		}

		// Test Success Response
		if resp.StatusCode == http.StatusOK {
			if len(resp.Data.Games) != testCase.First {
				t.Errorf("expected \"%d\" games, got \"%d\"", testCase.First, len(resp.Data.Games))
			}

			gameOne := "Fortnite"
			if resp.Data.Games[0].Name != gameOne {
				t.Errorf("expected game 1 name to be \"%s\", got \"%s\"", gameOne, resp.Data.Games[0].Name)
			}

			gameTwo := "League of Legends"
			if resp.Data.Games[1].Name != gameTwo {
				t.Errorf("expected game 2 name to be \"%s\", got \"%s\"", gameTwo, resp.Data.Games[0].Name)
			}

			gameThree := "Dota 2"
			if resp.Data.Games[2].Name != gameThree {
				t.Errorf("expected game 3 name to be \"%s\", got \"%s\"", gameThree, resp.Data.Games[0].Name)
			}

			cursor := "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6M319"
			if resp.Data.Pagination.Cursor != cursor {
				t.Errorf("expected cursor to be \"%s\", got \"%s\"", cursor, resp.Data.Pagination.Cursor)
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

	_, err := c.GetTopGames(&TopGamesParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
