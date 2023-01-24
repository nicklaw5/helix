package helix

type SendShoutoutParams struct {
	FromBroadcasterID string `query:"from_broadcaster_id"` // required
	ToBroadcasterID   string `query:"to_broadcaster_id"`   // required
	ModeratorID       string `query:"moderator_id"`        // required
}

type SendShoutoutResponse struct {
	ResponseCommon
}

// SendShoutout sends a Shoutout to the specified broadcaster.
// Required scope: moderator:manage:shoutouts
// The broadcaster may send a Shoutout once every 2 minutes.
// They may send the same broadcaster a Shoutout once every 60 minutes.
func (c *Client) SendShoutout(params *SendShoutoutParams) (*SendShoutoutResponse, error) {
	resp, err := c.post("/chat/shoutouts", nil, params)
	if err != nil {
		return nil, err
	}

	shoutoutResp := &SendShoutoutResponse{}
	resp.HydrateResponseCommon(&shoutoutResp.ResponseCommon)

	return shoutoutResp, nil
}
