package helix

// ExpiresAt must be parsed manually since an empty string means perma ban
type Ban struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	ExpiresAt Time   `json:"expires_at"`
}

type ManyBans struct {
	Bans       []Ban      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type BannedUsersResponse struct {
	ResponseCommon
	Data ManyBans
}

// BroadcasterID must match the auth tokens user_id
type BannedUsersParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"`
	After         string `query:"after"`
	Before        string `query:"before"`
}

// GetBannedUsers returns all banned and timed-out users in a channel.
//
// Required scope: moderation:read
func (c *Client) GetBannedUsers(params *BannedUsersParams, opts ...Options) (*BannedUsersResponse, error) {
	resp, err := c.get("/moderation/banned", &ManyBans{}, params, opts...)
	if err != nil {
		return nil, err
	}

	bans := &BannedUsersResponse{}
	resp.HydrateResponseCommon(&bans.ResponseCommon)
	bans.Data.Bans = resp.Data.(*ManyBans).Bans
	bans.Data.Pagination = resp.Data.(*ManyBans).Pagination

	return bans, nil
}
