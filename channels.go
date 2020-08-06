package helix

// SearchChannelsParams is parameters for SearchChannels
type SearchChannelsParams struct {
	Channel  string `query:"query"`
	After    string `query:"after"`
	First    int    `query:"first,20"` // Limit 100
	LiveOnly bool   `query:"live_only"`
}

// ManySearchChannels is the response data from SearchChannels
type ManySearchChannels struct {
	Channels   []Channel  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Channel describes a channel from SearchChannel
type Channel struct {
	ID           string   `json:"id"`
	GameID       string   `json:"game_id"`
	DisplayName  string   `json:"display_name"`
	Language     string   `json:"broadcaster_language"`
	Title        string   `json:"title"`
	ThumbnailURL string   `json:"thumbnail_url"`
	IsLive       bool     `json:"is_live"`
	StartedAt    Time     `json:"started_at"`
	TagIDs       []string `json:"tag_ids"`
}

// SearchChannelsResponse is the response from SearchChannels
type SearchChannelsResponse struct {
	ResponseCommon
	Data ManySearchChannels
}

// SearchChannels searches for Twitch channels based on the given search
// parameters. Unlike GetStreams, this can also return offline channels.
func (c *Client) SearchChannels(params *SearchChannelsParams) (*SearchChannelsResponse, error) {
	resp, err := c.get("/search/channels", &ManySearchChannels{}, params)
	if err != nil {
		return nil, err
	}

	channels := &SearchChannelsResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Channels = resp.Data.(*ManySearchChannels).Channels
	channels.Data.Pagination = resp.Data.(*ManySearchChannels).Pagination

	return channels, nil
}
