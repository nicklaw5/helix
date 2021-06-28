package helix

type ExtensionAnalytic struct {
	ExtensionID string    `json:"extension_id"`
	URL         string    `json:"URL"`
	Type        string    `json:"type"`
	DateRange   DateRange `json:"date_range"`
}

type ManyExtensionAnalytics struct {
	ExtensionAnalytics []ExtensionAnalytic `json:"data"`
	Pagination         Pagination          `json:"pagination"`
}

type ExtensionAnalyticsResponse struct {
	ResponseCommon
	Data ManyExtensionAnalytics
}

type ExtensionAnalyticsParams struct {
	ExtensionID string `query:"extension_id"`
	First       int    `query:"first,20"`
	After       string `query:"after"`
	StartedAt   Time   `query:"started_at"`
	EndedAt     Time   `query:"ended_at"`
	Type        string `query:"type"`
}

// GetExtensionAnalytics returns a URL to the downloadable CSV file
// containing analytics data. Valid for 5 minutes.
func (c *Client) GetExtensionAnalytics(params *ExtensionAnalyticsParams) (*ExtensionAnalyticsResponse, error) {
	resp, err := c.get("/analytics/extensions", &ManyExtensionAnalytics{}, params)
	if err != nil {
		return nil, err
	}

	users := &ExtensionAnalyticsResponse{}
	resp.HydrateResponseCommon(&users.ResponseCommon)
	users.Data.ExtensionAnalytics = resp.Data.(*ManyExtensionAnalytics).ExtensionAnalytics
	users.Data.Pagination = resp.Data.(*ManyExtensionAnalytics).Pagination
	return users, nil
}

type GameAnalytic struct {
	GameID    string    `json:"game_id"`
	URL       string    `json:"URL"`
	Type      string    `json:"type"`
	DateRange DateRange `json:"date_range"`
}

type ManyGameAnalytics struct {
	GameAnalytics []GameAnalytic `json:"data"`
	Pagination    Pagination     `json:"pagination"`
}

type GameAnalyticsResponse struct {
	ResponseCommon
	Data ManyGameAnalytics
}

type GameAnalyticsParams struct {
	GameID    string `query:"game_id"`
	First     int    `query:"first,20"`
	After     string `query:"after"`
	StartedAt Time   `query:"started_at"`
	EndedAt   Time   `query:"ended_at"`
	Type      string `query:"type"`
}

// GetGameAnalytics returns a URL to the downloadable CSV file
// containing analytics data for the specified game. Valid for 5 minutes.
func (c *Client) GetGameAnalytics(params *GameAnalyticsParams) (*GameAnalyticsResponse, error) {

	resp, err := c.get("/analytics/games", &ManyGameAnalytics{}, params)
	if err != nil {
		return nil, err
	}

	users := &GameAnalyticsResponse{}
	resp.HydrateResponseCommon(&users.ResponseCommon)
	users.Data.GameAnalytics = resp.Data.(*ManyGameAnalytics).GameAnalytics
	users.Data.Pagination = resp.Data.(*ManyGameAnalytics).Pagination

	return users, nil
}
