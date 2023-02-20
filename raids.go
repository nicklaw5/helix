package helix

// RaidResponse is the response from StartRaid
type RaidResponse struct {
	ResponseCommon
	Data StartRaidResponse
}

// StartRaidParams are the parameters for StartRaid
type StartRaidParams struct {
	FromBroadcasterID string `query:"from_broadcaster_id"`
	ToBroadcasterID   string `query:"to_broadcaster_id"`
}

// StartRaidResponse is the response data in RaidResponse
type StartRaidResponse struct {
	Data []RaidDetails `json:"data"`
}

// RaidDetails describes details of the ongoing raid
type RaidDetails struct {
	CreatedAt Time `json:"created_at"`
	IsMature  bool `json:"is_mature"`
}

// StartRaid raids another channel by sending the broadcaster’s viewers to the targeted channel.
// When called, the Twitch UX pops up a window at the top of the chat room that identifies the number of viewers in the raid.
// The raid occurs when the broadcaster clicks Raid Now or after the 90-second countdown expires.
// Required scope: channel:manage:raids
// Rate limit: 10 requests within a 10-minute window.
func (c *Client) StartRaid(params *StartRaidParams) (*RaidResponse, error) {
	resp, err := c.post("/raids", &StartRaidResponse{}, params)
	if err != nil {
		return nil, err
	}

	raid := &RaidResponse{}
	resp.HydrateResponseCommon(&raid.ResponseCommon)

	return raid, nil
}

// CancelRaidResponse is the response from StartRaid
type CancelRaidResponse struct {
	ResponseCommon
}

// CancelRaidParams are the parameters for CancelRaid
type CancelRaidParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

// CancelRaid cancels a pending raid.
// Required scope: channel:manage:raids
// Rate limit: 10 requests within a 10-minute window.
func (c *Client) CancelRaid(params *CancelRaidParams) (*CancelRaidResponse, error) {
	resp, err := c.delete("/raids", nil, params)
	if err != nil {
		return nil, err
	}

	canceledRaid := &CancelRaidResponse{}
	resp.HydrateResponseCommon(&canceledRaid.ResponseCommon)

	return canceledRaid, nil
}
