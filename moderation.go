package helix

type Ban struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	ExpiresAt string `json:"expires_at"`
}

type ManyBans struct {
	Bans       []Ban      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type BannedUsersResponse struct {
	ResponseCommon
	Data ManyBans
}

type BannedUsersParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"`
	After         string `query:"after"`
	Before        string `query:"before"`
}

func (c *Client) GetBannedUsers(params *BannedUsersParams) (*BannedUsersResponse, error) {
	resp, err := c.get("/moderation/banned", &ManyBans{}, params)
	if err != nil {
		return nil, err
	}

	bans := &BannedUsersResponse{}
	bans.StatusCode = resp.StatusCode
	bans.Header = resp.Header
	bans.Error = resp.Error
	bans.ErrorStatus = resp.ErrorStatus
	bans.ErrorMessage = resp.ErrorMessage
	bans.Data.Pagination = resp.Data.(*ManyBans).Pagination

	return bans, nil
}
