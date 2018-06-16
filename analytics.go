package helix

// GameAnalytic ...
type GameAnalytic struct {
	GameID string `json:"game_id"`
	URL    string `json:"url"`
}

// ManyGameAnalytics ...
type ManyGameAnalytics struct {
	GameAnalytics []GameAnalytic `json:"data"`
}

// GameAnalyticsResponse ...
type GameAnalyticsResponse struct {
	ResponseCommon
	Data ManyGameAnalytics
}

type gameAnalyticsParams struct {
	GameID string `query:"game_id"`
}

// GetGameAnalytics returns a URL to the downloadable CSV file
// containing analytics data for the specified game. Valid for 1 minute.
func (c *Client) GetGameAnalytics(gameID string) (*GameAnalyticsResponse, error) {
	params := &gameAnalyticsParams{
		GameID: gameID,
	}

	resp, err := c.get("/analytics/games", &ManyGameAnalytics{}, params)
	if err != nil {
		return nil, err
	}

	users := &GameAnalyticsResponse{}
	users.StatusCode = resp.StatusCode
	users.Header = resp.Header
	users.Error = resp.Error
	users.ErrorStatus = resp.ErrorStatus
	users.ErrorMessage = resp.ErrorMessage
	users.Data.GameAnalytics = resp.Data.(*ManyGameAnalytics).GameAnalytics

	return users, nil
}
