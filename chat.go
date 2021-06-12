package helix

// GetChatBadgeParams ...
type GetChatBadgeParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

// GetChatBadgeResponse ...
type GetChatBadgeResponse struct {
	ResponseCommon
	Data ManyChatBadge
}

// ManyChatBadge ...
type ManyChatBadge struct {
	Badges []ChatBadge `json:"data"`
}

// ChatBadge ...
type ChatBadge struct {
	SetID    string         `json:"set_id"`
	Versions []BadgeVersion `json:"versions"`
}

// BadgeVersion ...
type BadgeVersion struct {
	ID         string `json:"id"`
	ImageUrl1x string `json:"image_url_1x"`
	ImageUrl2x string `json:"image_url_2x"`
	ImageUrl4x string `json:"image_url_4x"`
}

// GetChannelChatBadges ...
func (c *Client) GetChannelChatBadges(params *GetChatBadgeParams) (*GetChatBadgeResponse, error) {
	resp, err := c.get("/chat/badges", &ManyChatBadge{}, params)
	if err != nil {
		return nil, err
	}

	channels := &GetChatBadgeResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Badges = resp.Data.(*ManyChatBadge).Badges

	return channels, nil
}

// GetGlobalChatBadges ...
func (c *Client) GetGlobalChatBadges() (*GetChatBadgeResponse, error) {
	resp, err := c.get("/chat/badges/global", &ManyChatBadge{}, nil)
	if err != nil {
		return nil, err
	}

	channels := &GetChatBadgeResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Badges = resp.Data.(*ManyChatBadge).Badges

	return channels, nil
}

// GetChannelEmotesParams ...
type GetChannelEmotesParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

// GetEmoteSetsParams ...
type GetEmoteSetsParams struct {
	EmoteSetIDs []string `query:"emote_set_id"` // Minimum: 1. Maximum: 25.
}

// GetChannelEmotesResponse ...
type GetChannelEmotesResponse struct {
	ResponseCommon
	Data ManyEmotes
}

// ManyEmotes ...
type ManyEmotes struct {
	Emotes []Emote `json:"data"`
}

// Emote ...
type Emote struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Images EmoteImage `json:"images"`
	Tier   string     `json:"tier"`
	Type   string     `json:"emote_type"`
	Set_ID string     `json:"emote_set_id"`
}

// EmoteImage ...
type EmoteImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

// GetChannelChatEmotes ...
func (c *Client) GetChannelEmotes(params *GetChannelEmotesParams) (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes", &ManyEmotes{}, params)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}

// GetGlobalEmotes ...
func (c *Client) GetGlobalEmotes() (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes/global", &ManyEmotes{}, nil)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}

// GetEmoteSets
func (c *Client) GetEmoteSets(params *GetEmoteSetsParams) (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes/set", &ManyEmotes{}, params)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}
