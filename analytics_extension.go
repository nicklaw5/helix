package helix

// ExtensionAnalytic ...
type ExtensionAnalytic struct {
	ExtensionID string `json:"extension_id"`
	URL    string `json:"url"`
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

type extensionAnalyticsParams struct {
	ExtensionID string `query:"extension_id"`
}

// GetExtensionAnalytics returns a URL to the downloadable CSV file
// containing analytics data . Valid for 1 minute.
func (c *Client) GetExtensionAnalytics(extensionID string) (*ExtensionAnalyticsResponse, error) {
	params := &extensionAnalyticsParams{
		ExtensionID: extensionID,
	}

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
