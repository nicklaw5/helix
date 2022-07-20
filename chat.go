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

type GetChannelEmotesParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetEmoteSetsParams struct {
	EmoteSetIDs []string `query:"emote_set_id"` // Minimum: 1. Maximum: 25.
}

type SendChatAnnouncementParams struct {
	BroadcasterID string `query:"broadcaster_id"` // required
	ModeratorID   string `query:"moderator_id"`   // required
	Message       string `json:"message"`         // upto 500 chars, thereafter str is truncated
	// blue || green || orange || purple are valid, default 'primary' or empty str result in channel accent color.
	Color string `json:"color"`
}

type SendChatAnnouncementResponse struct {
	ResponseCommon
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
func (c *Client) GetEmoteSets(params *GetEmoteSetsParams) (*GetEmoteSetsResponse, error) {
	resp, err := c.get("/chat/emotes/set", &ManyEmotesWithOwner{}, params)
	if err != nil {
		return nil, err
	}

	emotes := &GetEmoteSetsResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotesWithOwner).Emotes

	return emotes, nil
}

// SendChatAnnouncement sends an announcement to the broadcaster’s chat room.
// Required scope: moderator:manage:announcements
func (c *Client) SendChatAnnouncement(params *SendChatAnnouncementParams) (*SendChatAnnouncementResponse, error) {
	resp, err := c.postAsJSON("/chat/announcements", nil, params)
	if err != nil {
		return nil, err
	}

	chatResp := &SendChatAnnouncementResponse{}
	resp.HydrateResponseCommon(&chatResp.ResponseCommon)

	return chatResp, nil
}
