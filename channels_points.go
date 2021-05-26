package helix

// ChannelCustomRewardsParams ...
type ChannelCustomRewardsParams struct {
	BroadcasterID                     string `query:"broadcaster_id"`
	Title                             string `json:"title"`
	Cost                              int    `json:"cost"`
	Prompt                            string `json:"prompt"`
	IsEnabled                         bool   `json:"is_enabled"`
	BackgroundColor                   string `json:"background_color"`
	IsUserInputRequired               bool   `json:"is_user_input_required"`
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled"`
	MaxPerStream                      int    `json:"max_per_stream"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled"`
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream"`
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled"`
	GlobalCooldownSeconds             int    `json:"global_cooldown_seconds"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
}

// ChannelCustomRewards ...
type ManyChannelCustomRewards struct {
	ChannelCustomRewards []ChannelCustomReward `json:"data"`
}

// ChannelCustomReward ...
type ChannelCustomReward struct {
	BroadcasterID                     string                      `json:"broadcaster_id"`
	BroadcasterLogin                  string                      `json:"broadcaster_login"`
	BroadcasterName                   string                      `json:"broadcaster_name"`
	ID                                string                      `json:"id"`
	Title                             string                      `json:"title"`
	Prompt                            string                      `json:"prompt"`
	Cost                              int                         `json:"cost"`
	Image                             RewardImage                 `json:"image"`
	DefaultImage                      RewardImage                 `json:"default_image"`
	IsEnabled                         bool                        `json:"is_enabled"`
	IsUserInputRequired               bool                        `json:"is_user_input_required"`
	MaxPerStreamSetting               MaxPerStreamSettings        `json:"max_per_stream_setting"`
	MaxPerUserPerStreamSetting        MaxPerUserPerStreamSettings `json:"max_per_user_per_stream_setting"`
	GlobalCooldownSetting             GlobalCooldownSettings      `json:"global_cooldown_setting"`
	IsPaused                          bool                        `json:"is_paused"`
	IsInStock                         bool                        `json:"is_in_stock"`
	ShouldRedemptionsSkipRequestQueue bool                        `json:"should_redemptions_skip_request_queue"`
	RedemptionsRedeemedCurrentStream  int                         `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 string                      `json:"cooldown_expires_at"`
}

// RewardImage ...
type RewardImage struct {
	Url1 string `json:"url_1x"`
	Url2 string `json:"url_2x"`
	Url4 string `json:"url_4x"`
}

// MaxPerUserPerStreamSettings ...
type MaxPerUserPerStreamSettings struct {
	IsEnabled           bool `json:"is_enabled"`
	MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
}

// MaxPerStreamSettings ...
type MaxPerStreamSettings struct {
	IsEnabled    bool `json:"is_enabled"`
	MaxPerStream int  `json:"max_per_stream"`
}

// GlobalCooldownSettings ...
type GlobalCooldownSettings struct {
	IsEnabled             bool `json:"is_enabled"`
	GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
}

// ChannelCustomRewardResponse ...
type ChannelCustomRewardResponse struct {
	ResponseCommon
	Data ManyChannelCustomRewards
}

// CreateCustomReward : Creates a Custom Reward on a channel.
// Required scope: channel:manage:redemptions
func (c *Client) CreateCustomReward(params *ChannelCustomRewardsParams) (*ChannelCustomRewardResponse, error) {
	resp, err := c.postAsJSON("/channel_points/custom_rewards", &ManyChannelCustomRewards{}, params)
	if err != nil {
		return nil, err
	}

	reward := &ChannelCustomRewardResponse{}
	resp.HydrateResponseCommon(&reward.ResponseCommon)
	reward.Data.ChannelCustomRewards = resp.Data.(*ManyChannelCustomRewards).ChannelCustomRewards

	return reward, nil
}
