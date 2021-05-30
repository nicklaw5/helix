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
