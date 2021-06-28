package helix

import "time"

type UserBitTotal struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Rank      int    `json:"rank"`
	Score     int    `json:"score"`
}

type ManyUserBitTotals struct {
	Total         int            `json:"total"`
	DateRange     DateRange      `json:"date_range"`
	UserBitTotals []UserBitTotal `json:"data"`
}

type BitsLeaderboardResponse struct {
	ResponseCommon
	Data ManyUserBitTotals
}

type BitsLeaderboardParams struct {
	Count     int       `query:"count,10"`   // Maximum 100
	Period    string    `query:"period,all"` // "all" (default), "day", "week", "month" and "year"
	StartedAt time.Time `query:"started_at"`
	UserID    string    `query:"user_id"`
}

// GetBitsLeaderboard gets a ranked list of Bits leaderboard
// information for an authorized broadcaster.
//
// Required Scope: bits:read
func (c *Client) GetBitsLeaderboard(params *BitsLeaderboardParams) (*BitsLeaderboardResponse, error) {
	resp, err := c.get("/bits/leaderboard", &ManyUserBitTotals{}, params)
	if err != nil {
		return nil, err
	}

	bits := &BitsLeaderboardResponse{}
	resp.HydrateResponseCommon(&bits.ResponseCommon)
	bits.Data.Total = resp.Data.(*ManyUserBitTotals).Total
	bits.Data.DateRange = resp.Data.(*ManyUserBitTotals).DateRange
	bits.Data.UserBitTotals = resp.Data.(*ManyUserBitTotals).UserBitTotals

	return bits, nil
}

type CheermotesParams struct {
	BroadcasterID string `query:"broadcaster_id"` // optional
}

type TierImages struct {
	Image1   string `json:"1"`
	Image1_5 string `json:"1.5"`
	Image2   string `json:"2"`
	Image3   string `json:"3"`
	Image4   string `json:"4"`
}

type TierImageTypes struct {
	Animated TierImages `json:"animated"`
	Static   TierImages `json:"static"`
}

type CheermoteTierImages struct {
	Dark  TierImageTypes `json:"dark"`
	Light TierImageTypes `json:"light"`
}

type CheermoteTiers struct {
	MinBits        uint                `json:"min_bits"`
	ID             string              `json:"id"`
	Color          string              `json:"color"`
	Images         CheermoteTierImages `json:"images"`
	CanCheer       bool                `json:"can_cheer"`
	ShowInBitsCard bool                `json:"show_in_bits_card"`
}

type Cheermotes struct {
	Prefix       string           `json:"prefix"`
	Tiers        []CheermoteTiers `json:"tiers"`
	Type         string           `json:"type"` // global_first_party, global_third_party, channel_custom, display_only, sponsored
	Order        uint             `json:"order"`
	LastUpdated  Time             `json:"last_updated"`
	IsCharitable bool             `json:"is_charitable"`
}

type ManyCheermotes struct {
	Cheermotes []Cheermotes `json:"data"`
}

type CheermotesResponse struct {
	ResponseCommon
	Data ManyCheermotes
}

func (c *Client) GetCheermotes(params *CheermotesParams) (*CheermotesResponse, error) {
	resp, err := c.get("/bits/cheermotes", &ManyCheermotes{}, params)
	if err != nil {
		return nil, err
	}

	cheermoteResp := &CheermotesResponse{}
	resp.HydrateResponseCommon(&cheermoteResp.ResponseCommon)
	cheermoteResp.Data.Cheermotes = resp.Data.(*ManyCheermotes).Cheermotes

	return cheermoteResp, nil
}
