package helix

type ChannelCustomRewardsParams struct {
	BroadcasterID                     string `query:"broadcaster_id"`
	Title                             string `json:"title"`
	Cost                              int    `json:"cost"`
	Prompt                            string `json:"prompt"`
	IsEnabled                         bool   `json:"is_enabled"`
	BackgroundColor                   string `json:"background_color,omitempty"`
	IsUserInputRequired               bool   `json:"is_user_input_required"`
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled"`
	MaxPerStream                      int    `json:"max_per_stream"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled"`
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream"`
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled"`
	GlobalCooldownSeconds             int    `json:"global_cooldown_seconds"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
}

type UpdateChannelCustomRewardsParams struct {
	ID                                string `query:"id"`
	BroadcasterID                     string `query:"broadcaster_id"`
	Title                             string `json:"title"`
	Cost                              int    `json:"cost"`
	Prompt                            string `json:"prompt"`
	IsEnabled                         bool   `json:"is_enabled"`
	BackgroundColor                   string `json:"background_color,omitempty"`
	IsUserInputRequired               bool   `json:"is_user_input_required"`
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled"`
	MaxPerStream                      int    `json:"max_per_stream"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled"`
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream"`
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled"`
	GlobalCooldownSeconds             int    `json:"global_cooldown_seconds"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
}

type DeleteCustomRewardsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	ID            string `query:"id"`
}

type GetCustomRewardsParams struct {
	BroadcasterID         string `query:"broadcaster_id"`
	ID                    string `query:"id"`
	OnlyManageableRewards bool   `query:"only_manageable_rewards"`
}

type ManyChannelCustomRewards struct {
	ChannelCustomRewards []ChannelCustomReward `json:"data"`
}

type ChannelCustomReward struct {
	BroadcasterID                     string                      `json:"broadcaster_id"`
	BroadcasterLogin                  string                      `json:"broadcaster_login"`
	BroadcasterName                   string                      `json:"broadcaster_name"`
	ID                                string                      `json:"id"`
	Title                             string                      `json:"title"`
	Prompt                            string                      `json:"prompt"`
	Cost                              int                         `json:"cost"`
	Image                             RewardImage                 `json:"image"`
	BackgroundColor                   string                      `json:"background_color"`
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

type RewardImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

type MaxPerUserPerStreamSettings struct {
	IsEnabled           bool `json:"is_enabled"`
	MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
}

type MaxPerStreamSettings struct {
	IsEnabled    bool `json:"is_enabled"`
	MaxPerStream int  `json:"max_per_stream"`
}

type GlobalCooldownSettings struct {
	IsEnabled             bool `json:"is_enabled"`
	GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
}

type ChannelCustomRewardResponse struct {
	ResponseCommon
	Data ManyChannelCustomRewards
}

// Response for removing a custom reward
type DeleteCustomRewardsResponse struct {
	ResponseCommon
}

type UpdateChannelCustomRewardsRedemptionStatusParams struct {
	ID            string `query:"id"`
	BroadcasterID string `query:"broadcaster_id"`
	RewardID      string `query:"reward_id"`
	Status        string `json:"status"`
}

type ChannelCustomRewardsRedemptionResponse struct {
	ResponseCommon
	Data ManyChannelCustomRewardsRedemptions
}

type ManyChannelCustomRewardsRedemptions struct {
	Redemptions []ChannelCustomRewardsRedemption `json:"data"`
}

type ChannelCustomRewardsRedemption struct {
	ID               string              `json:"id"`
	BroadcasterID    string              `json:"broadcaster_id"`
	BroadcasterLogin string              `json:"broadcaster_login"`
	BroadcasterName  string              `json:"broadcaster_name"`
	UserID           string              `json:"user_id"`
	UserName         string              `json:"user_name"`
	UserLogin        string              `json:"user_login"`
	UserInput        string              `json:"user_input"`
	Status           string              `json:"status"`
	RedeemedAt       Time                `json:"redeemed_at"`
	Reward           ChannelCustomReward `json:"reward"`
}

type GetCustomRewardsRedemptionsParams struct {
	BroadcasterID string `query:"broadcaster_id"` // required
	RewardID      string `query:"reward_id"`      // required
	Status        string `query:"status"`         // required if ID is null
	ID            string `query:"id"`             // max 50
	Sort          string `query:"sort"`
	First         int    `query:"first"` // max 50
	After         string `query:"after"`
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

// UpdateCustomReward : Update a Custom Reward on a channel.
// Required scope: channel:manage:redemptions
func (c *Client) UpdateCustomReward(params *UpdateChannelCustomRewardsParams) (*ChannelCustomRewardResponse, error) {
	resp, err := c.patchAsJSON("/channel_points/custom_rewards", &ManyChannelCustomRewards{}, params)
	if err != nil {
		return nil, err
	}

	reward := &ChannelCustomRewardResponse{}
	resp.HydrateResponseCommon(&reward.ResponseCommon)
	reward.Data.ChannelCustomRewards = resp.Data.(*ManyChannelCustomRewards).ChannelCustomRewards

	return reward, nil
}

// DeleteCustomRewards : Deletes a Custom Rewards on a channel
// Required scope: channel:manage:redemptions
func (c *Client) DeleteCustomRewards(params *DeleteCustomRewardsParams) (*DeleteCustomRewardsResponse, error) {
	resp, err := c.delete("/channel_points/custom_rewards", nil, params)
	if err != nil {
		return nil, err
	}

	reward := &DeleteCustomRewardsResponse{}
	resp.HydrateResponseCommon(&reward.ResponseCommon)

	return reward, nil
}

// GetCustomRewards : Get Custom Rewards on a channel
// Required scope: channel:read:redemptions
func (c *Client) GetCustomRewards(params *GetCustomRewardsParams) (*ChannelCustomRewardResponse, error) {
	resp, err := c.get("/channel_points/custom_rewards", &ManyChannelCustomRewards{}, params)
	if err != nil {
		return nil, err
	}

	rewards := &ChannelCustomRewardResponse{}
	resp.HydrateResponseCommon(&rewards.ResponseCommon)
	rewards.Data.ChannelCustomRewards = resp.Data.(*ManyChannelCustomRewards).ChannelCustomRewards

	return rewards, nil
}

// GetCustomRewardsRedemptions : Gets Custom Reward Redemption statuses on a channel.
// Required scope: channel:manage:redemptions
func (c *Client) GetCustomRewardsRedemptions(params *GetCustomRewardsRedemptionsParams) (*ChannelCustomRewardsRedemptionResponse, error) {
	resp, err := c.get("/channel_points/custom_rewards/redemptions", &ManyChannelCustomRewardsRedemptions{}, params)
	if err != nil {
		return nil, err
	}

	redemptions := &ChannelCustomRewardsRedemptionResponse{}
	resp.HydrateResponseCommon(&redemptions.ResponseCommon)
	redemptions.Data.Redemptions = resp.Data.(*ManyChannelCustomRewardsRedemptions).Redemptions

	return redemptions, nil
}

// UpdateChannelCustomRewardsRedemptionStatus : Update a Custom Reward Redemption status on a channel.
// Required scope: channel:manage:redemptions
func (c *Client) UpdateChannelCustomRewardsRedemptionStatus(params *UpdateChannelCustomRewardsRedemptionStatusParams) (*ChannelCustomRewardsRedemptionResponse, error) {
	resp, err := c.patchAsJSON("/channel_points/custom_rewards/redemptions", &ManyChannelCustomRewardsRedemptions{}, params)
	if err != nil {
		return nil, err
	}

	redemptions := &ChannelCustomRewardsRedemptionResponse{}
	resp.HydrateResponseCommon(&redemptions.ResponseCommon)
	redemptions.Data.Redemptions = resp.Data.(*ManyChannelCustomRewardsRedemptions).Redemptions

	return redemptions, nil
}
