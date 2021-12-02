package helix

type GetChatBadgeParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetChatBadgeResponse struct {
	ResponseCommon
	Data ManyChatBadge
}

type ManyChatBadge struct {
	Badges []ChatBadge `json:"data"`
}

type ChatBadge struct {
	SetID    string         `json:"set_id"`
	Versions []BadgeVersion `json:"versions"`
}

type BadgeVersion struct {
	ID         string `json:"id"`
	ImageUrl1x string `json:"image_url_1x"`
	ImageUrl2x string `json:"image_url_2x"`
	ImageUrl4x string `json:"image_url_4x"`
}

func (c *Client) GetChannelChatBadges(params *GetChatBadgeParams, opts ...Options) (*GetChatBadgeResponse, error) {
	resp, err := c.get("/chat/badges", &ManyChatBadge{}, params, opts...)
	if err != nil {
		return nil, err
	}

	channels := &GetChatBadgeResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Badges = resp.Data.(*ManyChatBadge).Badges

	return channels, nil
}

func (c *Client) GetGlobalChatBadges(opts ...Options) (*GetChatBadgeResponse, error) {
	resp, err := c.get("/chat/badges/global", &ManyChatBadge{}, nil, opts...)
	if err != nil {
		return nil, err
	}

	channels := &GetChatBadgeResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Badges = resp.Data.(*ManyChatBadge).Badges

	return channels, nil
}

type GetChannelEmotesParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetEmoteSetsParams struct {
	EmoteSetIDs []string `query:"emote_set_id"` // Minimum: 1. Maximum: 25.
}

type GetChannelEmotesResponse struct {
	ResponseCommon
	Data ManyEmotes
}

type GetEmoteSetsResponse struct {
	ResponseCommon
	Data ManyEmotesWithOwner
}

type ManyEmotes struct {
	Emotes []Emote `json:"data"`
}

type ManyEmotesWithOwner struct {
	Emotes []EmoteWithOwner `json:"data"`
}

type Emote struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Images     EmoteImage `json:"images"`
	Tier       string     `json:"tier"`
	EmoteType  string     `json:"emote_type"`
	EmoteSetId string     `json:"emote_set_id"`
}

type EmoteWithOwner struct {
	Emote
	OwnerID string `json:"owner_id"`
}

type EmoteImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

func (c *Client) GetChannelEmotes(params *GetChannelEmotesParams, opts ...Options) (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes", &ManyEmotes{}, params, opts...)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}

func (c *Client) GetGlobalEmotes(opts ...Options) (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes/global", &ManyEmotes{}, nil, opts...)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}

// GetEmoteSets
func (c *Client) GetEmoteSets(params *GetEmoteSetsParams, opts ...Options) (*GetEmoteSetsResponse, error) {
	resp, err := c.get("/chat/emotes/set", &ManyEmotesWithOwner{}, params, opts...)
	if err != nil {
		return nil, err
	}

	emotes := &GetEmoteSetsResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotesWithOwner).Emotes

	return emotes, nil
}
