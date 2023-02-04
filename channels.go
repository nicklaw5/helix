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
	Tags             []string `json:"tags"`
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

type GetChannelFollowsParams struct {
	BroadcasterID string `query:"broadcaster_id"` // required
	UserID        string `query:"user_id"`
	First         int    `query:"first"` // max 100
	After         string `query:"after"`
}

// SearchChannelsResponse is the response from SearchChannels
type GetChannelFollowersResponse struct {
	ResponseCommon
	Data ManyChannelFollows
}

// ManySearchChannels is the response data from SearchChannels
type ManyChannelFollows struct {
	Channels   []ChannelFollow `json:"data"`
	Pagination Pagination      `json:"pagination"`
	Total      int             `json:"total"`
}

// Channel describes a follow of a channel
type ChannelFollow struct {
	UserID    string `json:"user_id"`
	Username  string `json:"user_name"`
	UserLogin string `json:"user_login"`
	Followed  Time   `json:"followed_at"`
}

type GetFollowedChannelParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"` // required
	First         int    `query:"first"`   // max 100
	After         string `query:"after"`
}

// SearchChannelsResponse is the response from SearchChannels
type GetFollowedChannelResponse struct {
	ResponseCommon
	Data ManyFollowedChannels
}

// ManySearchChannels is the response data from SearchChannels
type ManyFollowedChannels struct {
	FollowedChannels []FollowedChannel `json:"data"`
	Pagination       Pagination        `json:"pagination"`
	Total            int64             `json:"total"`
}

// Channel describes a followed channel
type FollowedChannel struct {
	BroadcasterID   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	BroadcaserLogin string `json:"broadcaster_login"`
	Followed        Time   `json:"followed_at"`
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

type GetChannelInformationParams struct {
	// Deprecated: BroadcasterID will be removed in a future version. Use BroadcasterIDs instead.
	BroadcasterID  string   `query:"broadcaster_id"`
	BroadcasterIDs []string `query:"broadcaster_id"` // Limit 100
}

type EditChannelInformationParams struct {
	BroadcasterID       string   `query:"broadcaster_id" json:"-"`
	GameID              string   `json:"game_id"`
	BroadcasterLanguage string   `json:"broadcaster_language"`
	Title               string   `json:"title"`
	Delay               int      `json:"delay,omitempty"`
	Tags                []string `json:"tags"`
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
	BroadcasterID       string   `json:"broadcaster_id"`
	BroadcasterName     string   `json:"broadcaster_name"`
	BroadcasterLanguage string   `json:"broadcaster_language"`
	GameID              string   `json:"game_id"`
	GameName            string   `json:"game_name"`
	Title               string   `json:"title"`
	Delay               int      `json:"delay"`
	Tags                []string `json:"tags"`
}

func (c *Client) GetChannelInformation(params *GetChannelInformationParams) (*GetChannelInformationResponse, error) {
	resp, err := c.get("/channels", &ManyChannelInformation{}, params)
	if err != nil {
		return nil, err
	}

	channels := &GetChannelInformationResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Channels = resp.Data.(*ManyChannelInformation).Channels

	return channels, nil
}

func (c *Client) EditChannelInformation(params *EditChannelInformationParams) (*EditChannelInformationResponse, error) {
	resp, err := c.patchAsJSON("/channels", &EditChannelInformationResponse{}, params)
	if err != nil {
		return nil, err
	}

	channels := &EditChannelInformationResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)

	return channels, nil
}

// GetChannelFollows Gets a list of users that follow the specified broadcaster.
// You can also use this endpoint to see whether a specific user follows the broadcaster..
// requires moderator:read:followers
func (c *Client) GetChannelFollows(params *GetChannelFollowsParams) (*GetChannelFollowersResponse, error) {
	resp, err := c.get("/channels/followers", &ManyChannelFollows{}, params)
	if err != nil {
		return nil, err
	}

	channelFollows := &GetChannelFollowersResponse{}
	resp.HydrateResponseCommon(&channelFollows.ResponseCommon)
	channelFollows.Data.Total = resp.Data.(*ManyChannelFollows).Total
	channelFollows.Data.Channels = resp.Data.(*ManyChannelFollows).Channels
	channelFollows.Data.Pagination = resp.Data.(*ManyChannelFollows).Pagination

	return channelFollows, nil
}

// GetFollowedChannels Gets a list of broadcasters that the specified user follows.
// You can also use this endpoint to see whether a user follows a specific broadcaster.
// requires user:read:follows
func (c *Client) GetFollowedChannels(params *GetFollowedChannelParams) (*GetFollowedChannelResponse, error) {
	resp, err := c.get("/channels/followed", &ManyFollowedChannels{}, params)
	if err != nil {
		return nil, err
	}

	followedChannels := &GetFollowedChannelResponse{}
	resp.HydrateResponseCommon(&followedChannels.ResponseCommon)
	followedChannels.Data.Total = resp.Data.(*ManyFollowedChannels).Total
	followedChannels.Data.FollowedChannels = resp.Data.(*ManyFollowedChannels).FollowedChannels
	followedChannels.Data.Pagination = resp.Data.(*ManyFollowedChannels).Pagination

	return followedChannels, nil
}
