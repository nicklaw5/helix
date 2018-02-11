package helix

import (
	"net/http"
	"testing"
)

func TestGetGames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		IDs        []string
		Names      []string
		respBody   string
		expectGame string
	}{
		{
			http.StatusOK,
			[]string{"27471"},
			[]string{},
			`{"data":[{"id":"27471","name":"Minecraft","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/Minecraft-{width}x{height}.jpg"}]}`,
			"Minecraft",
		},
		{
			http.StatusOK,
			[]string{},
			[]string{"Sea of Thieves"},
			`{"data":[{"id":"490377","name":"Sea of Thieves","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/Sea%20of%20Thieves-{width}x{height}.jpg"}]}`,
			"Sea of Thieves",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody))

		resp, err := c.GetGames(&GamesParams{
			IDs:   testCase.IDs,
			Names: testCase.Names,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.Data.Games[0].Name != testCase.expectGame {
			t.Errorf("expected game name to be %s, got %s", testCase.expectGame, resp.Data.Games[0].Name)
		}
	}
}
