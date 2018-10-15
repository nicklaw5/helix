package helix

// ExtensionAnalytic ...
type ExtensionAnalytic struct {
	ExtensionID string    `json:"extension_id"`
	URL         string    `json:"URL"`
	Type        string    `json:"type"`
	dateRange   DateRange `json:"date_range"`
}

// DateRange ...
type DateRange struct {
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

// ManyExtensionAnalytics ...
type ManyExtensionAnalytics struct {
	ExtensionAnalytics []ExtensionAnalytic `json:"data"`
}

// ExtensionAnalyticsResponse ...
type ExtensionAnalyticsResponse struct {
	ResponseCommon
	Data ManyExtensionAnalytics
}

type ExtensionAnalyticsParams struct {
	ExtensionID string `query:"extension_id"`
	First       int    `query:"first"`
	After       string `query:"after"`
	StartedAt   string `query:"started_at"`
	EndedAt     string `query:"ended_at"`
	Type        string `query:"type"`
}

// GetExtensionAnalytics returns a URL to the downloadable CSV file
// containing analytics data . Valid for 5 minutes.
func (c *Client) GetExtensionAnalytics(params *ExtensionAnalyticsParams) (*ExtensionAnalyticsResponse, error) {
	resp, err := c.get("/analytics/extensions", &ManyExtensionAnalytics{}, params)
	if err != nil {
		return nil, err
	}

	users := &ExtensionAnalyticsResponse{}
	users.StatusCode = resp.StatusCode
	users.Header = resp.Header
	users.Error = resp.Error
	users.ErrorStatus = resp.ErrorStatus
	users.ErrorMessage = resp.ErrorMessage
	users.Data.ExtensionAnalytics = resp.Data.(*ManyExtensionAnalytics).ExtensionAnalytics

	return users, nil
}
