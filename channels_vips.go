package helix

type ChannelsVipsParams struct {
	UserID        string `query:"user_id"`
	BroadcasterID string `query:"broadcaster_id"` // required
	First         int    `query:"first"`
	After         string `query:"after"`
}

type ManyChannelsVips struct {
	ChannelsVips []ChannelsVips `json:"data"`
	Pagination   Pagination     `json:"pagination"`
}

type ChannelsVips struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type ChannelsVipsResponse struct {
	ResponseCommon
	Data ManyChannelsVips
}

type AddChannelsVipsParams struct {
	UserID        string `query:"user_id"`        // required
	BroadcasterID string `query:"broadcaster_id"` // required
}

type AddChannelsVipsResponse struct {
	ResponseCommon
}

type RemoveChannelsVipsParams struct {
	UserID        string `query:"user_id"`        // required
	BroadcasterID string `query:"broadcaster_id"` // required
}

type RemoveChannelsVipsResponse struct {
	ResponseCommon
}

// GetChannelVips Gets a list of the broadcaster’s VIPs.
// Required scope: channel:read:vips
func (c *Client) GetChannelVips(params *ChannelsVipsParams) (*ChannelsVipsResponse, error) {
	resp, err := c.get("/channels/vips", &ManyChannelsVips{}, params)
	if err != nil {
		return nil, err
	}

	vips := &ChannelsVipsResponse{}
	resp.HydrateResponseCommon(&vips.ResponseCommon)
	vips.Data.ChannelsVips = resp.Data.(*ManyChannelsVips).ChannelsVips
	vips.Data.Pagination = resp.Data.(*ManyChannelsVips).Pagination

	return vips, nil
}

// AddChannelVips Adds the specified user as a VIP in the broadcaster’s channel.
// Required scope: channel:manage:vips
// Rate Limits: The broadcaster may add a maximum of 10 VIPs within a 10-second window.
func (c *Client) AddChannelVips(params *AddChannelsVipsParams) (*AddChannelsVipsResponse, error) {
	resp, err := c.post("/channels/vips", nil, params)
	if err != nil {
		return nil, err
	}

	vips := &AddChannelsVipsResponse{}
	resp.HydrateResponseCommon(&vips.ResponseCommon)

	return vips, nil
}

// RemoveChannelVips : Removes the specified user as a VIP in the broadcaster’s channel.
// Required scope: channel:manage:vips
// Rate Limits: The broadcaster may remove a maximum of 10 VIPs within a 10-second window.
func (c *Client) RemoveChannelVips(params *RemoveChannelsVipsParams) (*RemoveChannelsVipsResponse, error) {
	resp, err := c.delete("/channels/vips", nil, params)
	if err != nil {
		return nil, err
	}

	vips := &RemoveChannelsVipsResponse{}
	resp.HydrateResponseCommon(&vips.ResponseCommon)

	return vips, nil
}
