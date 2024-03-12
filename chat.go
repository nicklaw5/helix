package helix

import "errors"

type GetChatChattersParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	After         string `query:"after"`
	First         string `query:"first"`
}

type ChatChatter struct {
	UserLogin string `json:"user_login"`
	UserID    string `json:"user_id"`
	Username  string `json:"user_name"`
}

type ManyChatChatters struct {
	Chatters   []ChatChatter `json:"data"`
	Pagination Pagination    `json:"pagination"`
	Total      int           `json:"total"`
}

type GetChatChattersResponse struct {
	ResponseCommon
	Data ManyChatChatters
}

// Required scope: moderator:read:chatters
func (c *Client) GetChannelChatChatters(params *GetChatChattersParams) (*GetChatChattersResponse, error) {
	if params.BroadcasterID == "" || params.ModeratorID == "" {
		return nil, errors.New("error: broadcaster and moderator identifiers must be provided")
	}
	resp, err := c.get("/chat/chatters", &ManyChatChatters{}, params)
	if err != nil {
		return nil, err
	}

	chatters := &GetChatChattersResponse{}
	resp.HydrateResponseCommon(&chatters.ResponseCommon)
	chatters.Data.Chatters = resp.Data.(*ManyChatChatters).Chatters
	chatters.Data.Total = resp.Data.(*ManyChatChatters).Total
	chatters.Data.Pagination = resp.Data.(*ManyChatChatters).Pagination

	return chatters, nil
}

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

type GetChatSettingsParams struct {
	// Required, the ID of the broadcaster whose chat settings you want to get
	BroadcasterID string `query:"broadcaster_id"`

	// Optional, can be specified if you want the `non_moderator_chat_delay` and `non_moderator_chat_delay_duration` fields in the response. The ID should be a user that has moderation privileges in the broadcaster's chat.
	// The ID must match the specified User Access Token & the User Access Token must have the `moderator:read:chat_settings` scope
	ModeratorID string `query:"moderator_id,omitempty"`
}

type ChatSettings struct {
	BroadcasterID string `json:"broadcaster_id"`

	EmoteMode bool `json:"emote_mode"`

	FollowerMode bool `json:"follower_mode"`
	// Follower mode duration in minutes
	FollowerModeDuration int `json:"follower_mode_duration"`

	SlowMode bool `json:"slow_mode"`
	// Slow mode wait time in seconds
	SlowModeWaitTime int `json:"slow_mode_wait_time"`

	SubscriberMode bool `json:"subscriber_mode"`

	UniqueChatMode bool `json:"unique_chat_mode"`

	// Only included if the user access token includes the `moderator:read:chat_settings` scope
	ModeratorID string `json:"moderator_id"`

	// Boolean value denoting whether the "Non moderator chat delay" setting is enabled.
	// Only included if the request specifies a user access token that includes the moderator:read:chat_settings scope and the user in the moderator_id query parameter is one of the broadcaster’s moderators.
	NonModeratorChatDelay bool `json:"non_moderator_chat_delay"`
	// The amount of time, in seconds, that messages are delayed before appearing in chat.
	// Only included if the request specifies a user access token that includes the moderator:read:chat_settings scope and the user in the moderator_id query parameter is one of the broadcaster’s moderators.
	NonModeratorChatDelayDuration int `json:"non_moderator_chat_delay_duration"`
}

type ManyChatSettings struct {
	Settings []ChatSettings `json:"data"`
}

type GetChatSettingsResponse struct {
	ResponseCommon
	Data ManyChatSettings
}

// GetChatSettings gets the chat settings for the broadcaster's chat room.
// Optional scope: moderator:read:chat_settings
func (c *Client) GetChatSettings(params *GetChatSettingsParams) (*GetChatSettingsResponse, error) {
	if params.BroadcasterID == "" {
		return nil, errors.New("error: broadcaster id must be specified")
	}
	resp, err := c.get("/chat/settings", &ManyChatSettings{}, params)
	if err != nil {
		return nil, err
	}

	settings := &GetChatSettingsResponse{}
	resp.HydrateResponseCommon(&settings.ResponseCommon)
	settings.Data.Settings = resp.Data.(*ManyChatSettings).Settings

	return settings, nil
}

type UpdateChatSettingsParams struct {
	// Required, the ID of the broadcaster whose chat settings you want to update
	BroadcasterID string `query:"broadcaster_id"`

	// Required, the ID of a user that has moderator privileges in the BroadcasterID's channel.
	// The ID must match the specified User Access Token
	ModeratorID string `query:"moderator_id"`

	// Optional, set to true if only emotes are allowed
	// If unset (i.e. nil), no change to this setting will be made
	EmoteMode *bool `json:"emote_mode,omitempty"`

	// Optional, set to true if only followers may chat
	// If unset (i.e. nil), no change to this setting will be made
	FollowerMode *bool `json:"follower_mode,omitempty"`

	// Optional, time in minutes a user must have been following to chat.
	// If unset (i.e. nil), no change to this setting will be made
	// If set, FollowerMode must be set to true.
	// Possible values are 0 (no time restriction) through 129600 (3 months)
	FollowerModeDuration *int `json:"follower_mode_duration,omitempty"`

	// Optional, set to true if there's a delay before chat messages appear for non-moderators
	// If unset (i.e. nil), no change to this setting will be made
	NonModeratorChatDelay *bool `json:"non_moderator_chat_delay,omitempty"`

	// Optional, time in seconds before messages appear for non-moderators
	// If unset (i.e. nil), no change to this setting will be made
	// If set, FollowerMode must be set to true.
	// Possible values are 2, 4, or 6
	NonModeratorChatDelayDuration *int `json:"non_moderator_chat_delay_duration,omitempty"`

	// Optional, set to true if chatters must wait some extra time between sending more messages
	// If unset (i.e. nil), no change to this setting will be made
	SlowMode *bool `json:"slow_mode,omitempty"`

	// Optional, time in seconds chatters must wait between sending messages
	// If unset (i.e. nil), no change to this setting will be made
	// If set, SlowMode must be set to true.
	// Possible values are 3 through 120 seconds
	SlowModeWaitTime *int `json:"slow_mode_wait_time,omitempty"`

	// Optional, set to true if only subscribers may chat
	// If unset (i.e. nil), no change to this setting will be made
	SubscriberMode *bool `json:"subscriber_mode,omitempty"`

	// Optional, set to true if users may only post "unique messages" in chat
	// If unset (i.e. nil), no change to this setting will be made
	UniqueChatMode *bool `json:"unique_chat_mode,omitempty"`
}

type UpdateChatSettingsResponse struct {
	ResponseCommon
	Data ManyChatSettings
}

// UpdateChatSettings updates the broadcaster's chat settings.
// Required scope: moderator:manage:chat_settings
func (c *Client) UpdateChatSettings(params *UpdateChatSettingsParams) (*UpdateChatSettingsResponse, error) {
	if params.BroadcasterID == "" {
		return nil, errors.New("error: broadcaster id must be specified")
	}
	if params.ModeratorID == "" {
		return nil, errors.New("error: moderator id must be specified")
	}
	resp, err := c.patchAsJSON("/chat/settings", &ManyChatSettings{}, params)
	if err != nil {
		return nil, err
	}

	settings := &UpdateChatSettingsResponse{}
	resp.HydrateResponseCommon(&settings.ResponseCommon)
	settings.Data.Settings = resp.Data.(*ManyChatSettings).Settings

	return settings, nil
}

// UserChatColorResponse is the response from GetUserChatColor
type UserChatColorResponse struct {
	ResponseCommon
	Data GetUserChatColorResponse
}

// GetUserChatColorParams are the parameters for GetUserChatColor
type GetUserChatColorParams struct {
	UserID string `json:"user_id"`
}

// GetUserChatColorResponse is the response data in UserChatColorResponse
type GetUserChatColorResponse struct {
	Data []GetUserChatColorUser `json:"data"`
}

// GetUserChatColorUser describes the user and their color
type GetUserChatColorUser struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Color     string `json:"color"`
}

// GetUserChatColor fetches the color used for the user’s name in chat.
func (c *Client) GetUserChatColor(params *GetUserChatColorParams) (*UserChatColorResponse, error) {
	resp, err := c.get("/chat/color", &GetUserChatColorResponse{}, params)
	if err != nil {
		return nil, err
	}

	userColor := &UserChatColorResponse{}
	resp.HydrateResponseCommon(&userColor.ResponseCommon)

	return userColor, nil
}

// UpdateUserChatColorResponse is the response for UpdateUserChatColor
type UpdateUserChatColorResponse struct {
	ResponseCommon
}

// UpdateUserChatColorParams are the parameters for UpdateUserChatColor
type UpdateUserChatColorParams struct {
	UserID string `query:"user_id"`
	Color  string `query:"color"`
}

// UpdateUserChatcolor updates the color used for the user’s name in chat.
//
// Required scope: user:manage:chat_color
//
// Prime and Turbo users can specify a Hex color code, everyone can use the default colors:
//   - blue
//   - blue_violet
//   - cadet_blue
//   - chocolate
//   - coral
//   - dodger_blue
//   - firebrick
//   - golden_rod
//   - green
//   - hot_pink
//   - orange_red
//   - red
//   - sea_green
//   - spring_green
//   - yellow_green
func (c *Client) UpdateUserChatColor(params *UpdateUserChatColorParams) (*UpdateUserChatColorResponse, error) {
	resp, err := c.put("/chat/color", nil, params)
	if err != nil {
		return nil, err
	}

	update := &UpdateUserChatColorResponse{}
	resp.HydrateResponseCommon(&update.ResponseCommon)

	return update, nil
}

type SendChatMessageParams struct {
	// The ID of the broadcaster whose chat room the message will be sent to
	BroadcasterID string `json:"broadcaster_id"`

	// The ID of the user sending the message. This ID must match the user ID in the user access token
	SenderID string `json:"sender_id"`

	// The message to send. The message is limited to a maximum of 500 characters.
	// Chat messages can also include emoticons.
	// To include emoticons, use the name of the emote.
	// The names are case sensitive.
	// Don’t include colons around the name (e.g., :bleedPurple:).
	// If Twitch recognizes the name, Twitch converts the name to the emote before writing the chat message to the chat room
	Message string `json:"message"`

	// The ID of the chat message being replied to
	ReplyParentMessageID string `json:"reply_parent_message_id,omitempty"`
}

type ChatMessageResponse struct {
	ResponseCommon

	Data ManyChatMessages
}

type ManyChatMessages struct {
	Messages []ChatMessage `json:"data"`
}

type ChatMessage struct {
	// The message id for the message that was sent
	MessageID string `json:"message_id"`

	// If the message passed all checks and was sent
	IsSent bool `json:"is_sent"`

	// The reason the message was dropped, if any
	DropReasons ManyDropReasons `json:"drop_reason"`
}

type ManyDropReasons struct {
	Data DropReason
}

type DropReason struct {
	// Code for why the message was dropped
	Code string `json:"code"`

	// Message for why the message was dropped
	Message string `json:"message"`
}

// Requires an app access token or user access token that includes the user:write:chat scope.
// If app access token used, then additionally requires user:bot scope from chatting user,
// and either channel:bot scope from broadcaster or moderator status
func (c *Client) SendChatMessage(params *SendChatMessageParams) (*ChatMessageResponse, error) {
	if params.BroadcasterID == "" {
		return nil, errors.New("error: broadcaster id must be specified")
	}
	if params.SenderID == "" {
		return nil, errors.New("error: sender id must be specified")
	}

	resp, err := c.postAsJSON("/chat/messages", &ManyChatMessages{}, params)
	if err != nil {
		return nil, err
	}

	chatMessages := &ChatMessageResponse{}
	resp.HydrateResponseCommon(&chatMessages.ResponseCommon)
	chatMessages.Data.Messages = resp.Data.(*ManyChatMessages).Messages

	return chatMessages, nil
}
