package helix

type CharityCampaignAmount struct {
	Value         int64  `json:"amount"`
	DecimalPlaces int64  `json:"decimal_places"`
	Currency      string `json:"currency"`
}

type CharityCampaignData struct {
	ID               string                `json:"id"`
	BroadcasterID    string                `json:"broadcaster_id"`
	BroadcasterName  string                `json:"broadcaster_name"`
	BroadcasterLogin string                `json:"broadcaster_login"`
	Name             string                `json:"charity_name"`
	Description      string                `json:"charity_description"`
	LogoUrl          string                `json:"charity_logo"`
	WebsiteUrl       string                `json:"charity_website"`
	TargetAmount     CharityCampaignAmount `json:"target_amount"`
	CurrentAmount    CharityCampaignAmount `json:"current_amount"`
}

type ManyCharityCampaigns struct {
	Campaigns  []CharityCampaignData `json:"data"`
	Pagination Pagination            `json:"pagination"`
}

type CharityCampaignsResponse struct {
	ResponseCommon
	Data ManyCharityCampaigns
}

type CharityCampaignsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	After         string `query:"after"`
	First         int    `query:"first,20"` // Limit 100
}

// Required scope: channel:read:charity
func (c *Client) GetCharityCampaigns(params *CharityCampaignsParams) (*CharityCampaignsResponse, error) {
	resp, err := c.get("/charity/campaigns", &ManyCharityCampaigns{}, params)
	if err != nil {
		return nil, err
	}

	events := &CharityCampaignsResponse{}
	resp.HydrateResponseCommon(&events.ResponseCommon)
	events.Data.Campaigns = resp.Data.(*ManyCharityCampaigns).Campaigns
	events.Data.Pagination = resp.Data.(*ManyCharityCampaigns).Pagination

	return events, nil
}
