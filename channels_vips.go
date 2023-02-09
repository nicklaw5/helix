package helix

type GetChannelVipsParams struct {
	UserID        string `query:"user_id"`
	BroadcasterID string `query:"broadcaster_id"` // required
	First         int    `query:"first"`
	After         string `query:"after"`
}

type ManyChannelVips struct {
	ChannelsVips []ChannelVips `json:"data"`
	Pagination   Pagination    `json:"pagination"`
}

type ChannelVips struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type ChannelVipsResponse struct {
	ResponseCommon
	Data ManyChannelVips
}

type AddChannelVipParams struct {
	UserID        string `query:"user_id"`        // required
	BroadcasterID string `query:"broadcaster_id"` // required
}

type AddChannelVipResponse struct {
	ResponseCommon
}

type RemoveChannelVipParams struct {
	UserID        string `query:"user_id"`        // required
	BroadcasterID string `query:"broadcaster_id"` // required
}

type RemoveChannelVipResponse struct {
	ResponseCommon
}

// GetChannelVips Gets a list of the broadcaster’s VIPs.
// Required scope: channel:read:vips
func (c *Client) GetChannelVips(params *GetChannelVipsParams) (*ChannelVipsResponse, error) {
	resp, err := c.get("/channels/vips", &ManyChannelVips{}, params)
	if err != nil {
		return nil, err
	}

	vips := &ChannelVipsResponse{}
	resp.HydrateResponseCommon(&vips.ResponseCommon)
	vips.Data.ChannelsVips = resp.Data.(*ManyChannelVips).ChannelsVips
	vips.Data.Pagination = resp.Data.(*ManyChannelVips).Pagination

	return vips, nil
}

// AddChannelVip Adds the specified user as a VIP in the broadcaster’s channel.
// Required scope: channel:manage:vips
// Rate Limits: The broadcaster may add a maximum of 10 VIPs within a 10-second window.
func (c *Client) AddChannelVip(params *AddChannelVipParams) (*AddChannelVipResponse, error) {
	resp, err := c.post("/channels/vips", nil, params)
	if err != nil {
		return nil, err
	}

	vips := &AddChannelVipResponse{}
	resp.HydrateResponseCommon(&vips.ResponseCommon)

	return vips, nil
}

// RemoveChannelVip : Removes the specified user as a VIP in the broadcaster’s channel.
// Required scope: channel:manage:vips
// Rate Limits: The broadcaster may remove a maximum of 10 VIPs within a 10-second window.
func (c *Client) RemoveChannelVip(params *RemoveChannelVipParams) (*RemoveChannelVipResponse, error) {
	resp, err := c.delete("/channels/vips", nil, params)
	if err != nil {
		return nil, err
	}

	vips := &RemoveChannelVipResponse{}
	resp.HydrateResponseCommon(&vips.ResponseCommon)

	return vips, nil
}
