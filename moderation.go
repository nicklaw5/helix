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
func (c *Client) GetBannedUsers(params *BannedUsersParams) (*BannedUsersResponse, error) {
	resp, err := c.get("/moderation/banned", &ManyBans{}, params)
	if err != nil {
		return nil, err
	}

	bans := &BannedUsersResponse{}
	resp.HydrateResponseCommon(&bans.ResponseCommon)
	bans.Data.Bans = resp.Data.(*ManyBans).Bans
	bans.Data.Pagination = resp.Data.(*ManyBans).Pagination

	return bans, nil
}

type BanUserParams struct {
	BroadcasterID string             `json:"broadcaster_id"`
	ModeratorId   string             `json:"moderator_id"`
	Body          BanUserRequestBody `json:"data"`
}

type BanUserRequestBody struct {
	Duration int    `json:"duration"` // optional
	Reason   string `json:"reason"`   // required
	UserId   string `json:"user_id"`  // required
}

type BanUserResponse struct {
	ResponseCommon
	Data ManyBanUser
}

type ManyBanUser struct {
	Bans []BanUser `json:"data"`
}

type BanUser struct {
	BoardcasterId string `json:"broadcaster_id"`
	CreatedAt     string `json:"created_at"`
	EndTime       string `json:"end_time"`
	ModeratorId   string `json:"moderator_id"`
	UserId        string `json:"user_id"`
}

// BanUser Bans a user from participating in a broadcaster’s chat room, or puts them in a timeout.
// Required scope: moderator:manage:banned_users
func (c *Client) BanUser(params *BanUserParams) (*BanUserResponse, error) {
	resp, err := c.postAsJSON("/moderation/bans", &ManyBanUser{}, params)
	if err != nil {
		return nil, err
	}

	banResp := &BanUserResponse{}
	resp.HydrateResponseCommon(&banResp.ResponseCommon)
	banResp.Data.Bans = resp.Data.(*ManyBanUser).Bans

	return banResp, nil
}

type UnbanUserParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ModeratorID   string `json:"moderator_id"`
	UserID        string `json:"user_id"`
}

type UnbanUserResponse struct {
	ResponseCommon
}

func (c *Client) UnbanUser(params *UnbanUserParams) (*UnbanUserResponse, error) {
	resp, err := c.delete("/moderation/bans", nil, params)
	if err != nil {
		return nil, err
	}

	unbanResp := &UnbanUserResponse{}
	resp.HydrateResponseCommon(&unbanResp.ResponseCommon)
	return unbanResp, nil
}
