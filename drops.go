package helix

type GetDropEntitlementsParams struct {
	ID     string `query:"id"`
	UserID string `query:"user_id"`
	GameID string `query:"game_id"`
	After  string `query:"after"`
	First  int    `query:"first,20"` // Limit 1000
}

type Entitlement struct {
	ID        string `json:"id"`
	BenefitID string `json:"benefit_id"`
	Timestamp Time   `json:"timestamp"`
	UserID    string `json:"user_id"`
	GameID    string `json:"game_id"`
}

type ManyEntitlements struct {
	Entitlements []Entitlement `json:"data"`
}

type ManyEntitlementsWithPagination struct {
	ManyEntitlements
	Pagination `json:"pagination"`
}

type GetDropsEntitlementsResponse struct {
	ResponseCommon
	Data ManyEntitlementsWithPagination
}

// GetDropsEntitlements returns a list of entitlements, which have been awarded to users by your organization.
// Filtering by UserID returns all of the entitlements related to that specific user.
// Filtering by GameID returns all of the entitlements related to that game.
// Filtering by GameID and UserID returns all of the entitlements related to that game and that user.
// Entitlements are digital items that users are entitled to use. Twitch entitlements are granted based on viewership
// engagement with a content creator, based on the game developers' campaign.
func (c *Client) GetDropsEntitlements(params *GetDropEntitlementsParams) (*GetDropsEntitlementsResponse, error) {
	resp, err := c.get("/entitlements/drops", &ManyEntitlementsWithPagination{}, params)
	if err != nil {
		return nil, err
	}

	entitlements := &GetDropsEntitlementsResponse{}
	resp.HydrateResponseCommon(&entitlements.ResponseCommon)
	entitlements.Data.Entitlements = resp.Data.(*ManyEntitlementsWithPagination).Entitlements
	entitlements.Data.Pagination = resp.Data.(*ManyEntitlementsWithPagination).Pagination

	return entitlements, nil
}
