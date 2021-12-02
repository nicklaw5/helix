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
	ID               string   `json:"id"`
	GameID           string   `json:"game_id"`
	GameName         string   `json:"game_name"`
	BroadcasterLogin string   `json:"broadcaster_login"`
	DisplayName      string   `json:"display_name"`
	Language         string   `json:"broadcaster_language"`
	Title            string   `json:"title"`
	ThumbnailURL     string   `json:"thumbnail_url"`
	IsLive           bool     `json:"is_live"`
	StartedAt        Time     `json:"started_at"`
	TagIDs           []string `json:"tag_ids"`
}

// SearchChannelsResponse is the response from SearchChannels
type SearchChannelsResponse struct {
	ResponseCommon
	Data ManySearchChannels
}

// SearchChannels searches for Twitch channels based on the given search
// parameters. Unlike GetStreams, this can also return offline channels.
func (c *Client) SearchChannels(params *SearchChannelsParams, opts ...Options) (*SearchChannelsResponse, error) {
	resp, err := c.get("/search/channels", &ManySearchChannels{}, params, opts...)
	if err != nil {
		return nil, err
	}

	channels := &SearchChannelsResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Channels = resp.Data.(*ManySearchChannels).Channels
	channels.Data.Pagination = resp.Data.(*ManySearchChannels).Pagination

	return channels, nil
}

type GetChannelInformationParams struct {
	// Deprecated: BroadcasterID will be removed in a future version. Use BroadcasterIDs instead.
	BroadcasterID  string   `query:"broadcaster_id"`
	BroadcasterIDs []string `query:"broadcaster_id"` // Limit 100
}

type EditChannelInformationParams struct {
	BroadcasterID       string `query:"broadcaster_id" json:"-"`
	GameID              string `json:"game_id"`
	BroadcasterLanguage string `json:"broadcaster_language"`
	Title               string `json:"title"`
	Delay               int    `json:"delay,omitempty"`
}

type GetChannelInformationResponse struct {
	ResponseCommon
	Data ManyChannelInformation
}

type EditChannelInformationResponse struct {
	ResponseCommon
}

type ManyChannelInformation struct {
	Channels []ChannelInformation `json:"data"`
}

type ChannelInformation struct {
	BroadcasterID       string `json:"broadcaster_id"`
	BroadcasterName     string `json:"broadcaster_name"`
	BroadcasterLanguage string `json:"broadcaster_language"`
	GameID              string `json:"game_id"`
	GameName            string `json:"game_name"`
	Title               string `json:"title"`
	Delay               int    `json:"delay"`
}

func (c *Client) GetChannelInformation(params *GetChannelInformationParams, opts ...Options) (*GetChannelInformationResponse, error) {
	resp, err := c.get("/channels", &ManyChannelInformation{}, params, opts...)
	if err != nil {
		return nil, err
	}

	channels := &GetChannelInformationResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Channels = resp.Data.(*ManyChannelInformation).Channels

	return channels, nil
}

func (c *Client) EditChannelInformation(params *EditChannelInformationParams, opts ...Options) (*EditChannelInformationResponse, error) {
	resp, err := c.patchAsJSON("/channels", &EditChannelInformationResponse{}, params, opts...)
	if err != nil {
		return nil, err
	}

	channels := &EditChannelInformationResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)

	return channels, nil
}
