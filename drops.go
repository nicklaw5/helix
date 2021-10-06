package helix

type GetDropEntitlementsParams struct {
	ID                string `query:"id"`
	UserID            string `query:"user_id"`
	GameID            string `query:"game_id"`
	FulfillmentStatus string `query:"fulfillment_status"` // Valid values "CLAIMED", "FULFILLED"
	After             string `query:"after"`
	First             int    `query:"first,20"` // Limit 1000
}

type UpdateDropsEntitlementsParams struct {
	EntitlementIDs    []string `json:"entitlement_ids"`    // Limit 100
	FulfillmentStatus string   `json:"fulfillment_status"` // Valid values "CLAIMED", "FULFILLED"
}

type Entitlement struct {
	ID                string `json:"id"`
	BenefitID         string `json:"benefit_id"`
	Timestamp         Time   `json:"timestamp"`
	UserID            string `json:"user_id"`
	GameID            string `json:"game_id"`
	FulfillmentStatus string `json:"fulfillment_status"` // Valid values "CLAIMED", "FULFILLED"
	UpdatedAt         Time   `json:"updated_at"`
}
type UpdatedEntitlementSet struct {
	Status string   `json:"status"` // Valid values "SUCCESS", "INVALID_ID", "NOT_FOUND", "UNAUTHORIZED", "UPDATE_FAILED"
	IDs    []string `json:"ids"`
}

type ManyEntitlements struct {
	Entitlements []Entitlement `json:"data"`
}

type ManyUpdatedEntitlementSet struct {
	EntitlementSets []UpdatedEntitlementSet `json:"data"`
}

type ManyEntitlementsWithPagination struct {
	ManyEntitlements
	Pagination `json:"pagination"`
}

type GetDropsEntitlementsResponse struct {
	ResponseCommon
	Data ManyEntitlementsWithPagination
}

type UpdateDropsEntitlementsResponse struct {
	ResponseCommon
	Data ManyUpdatedEntitlementSet
}

// GetDropsEntitlements returns a list of entitlements, which have been awarded to users by your organization.
// Filtering by UserID returns all of the entitlements related to that specific user.
// Filtering by GameID returns all of the entitlements related to that game.
// Filtering by GameID and UserID returns all of the entitlements related to that game and that user.
// Filtering by FulfillmentStatus returns all of the entitlements with the specified fulfillment status.
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

// UpdateDropsEntitlements updates the fulfillment status of a set of entitlements, owned by the authenticated user or
// your organization. It returns a list of the entitlement ids requested, grouped by a status code used to indicate
// partial success.
// "SUCCESS" means the entitlement was successfully updated, "INVALID_ID" means invalid format for the entitlement,
// "NOT_FOUND" means the entitlement was not found, "UNAUTHORIZED" means entitlement is not owned by the organization or
// the user when called with a user OAuth token and "UPDATE_FAILED" indicates a possible transient error and the
// operation should be retried again later.
// Entitlements are digital items that users are entitled to use. Twitch entitlements are granted based on viewership
// engagement with a content creator, based on the game developers' campaign.
func (c *Client) UpdateDropsEntitlements(params *UpdateDropsEntitlementsParams) (*UpdateDropsEntitlementsResponse, error) {
	resp, err := c.patchAsJSON("/entitlements/drops", &ManyUpdatedEntitlementSet{}, params)
	if err != nil {
		return nil, err
	}

	entitlementSets := &UpdateDropsEntitlementsResponse{}
	resp.HydrateResponseCommon(&entitlementSets.ResponseCommon)
	entitlementSets.Data.EntitlementSets = resp.Data.(*ManyUpdatedEntitlementSet).EntitlementSets

	return entitlementSets, nil
}
