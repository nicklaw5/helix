package helix

type Goal struct {
	ID               string `json:"id"`
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterName  string `json:"broadcaster_name"`
	BroadcasterLogin string `json:"broadcaster_login"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	CurrentAmount    int    `json:"current_amount"`
	TargetAmount     int    `json:"target_amount"`
	CreatedAt        Time   `json:"created_at"`
}

type ManyGoals struct {
	Goals      []Goal     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type CreatorGoalsResponse struct {
	ResponseCommon
	Data ManyGoals
}

type GetCreatorGoalsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

// Required scope: channel:read:goals
func (c *Client) GetCreatorGoals(payload *GetCreatorGoalsParams) (*CreatorGoalsResponse, error) {
	resp, err := c.get("/goals", &ManyGoals{}, payload)
	if err != nil {
		return nil, err
	}

	goals := &CreatorGoalsResponse{}
	resp.HydrateResponseCommon(&goals.ResponseCommon)
	goals.Data.Goals = resp.Data.(*ManyGoals).Goals
	goals.Data.Pagination = resp.Data.(*ManyGoals).Pagination

	return goals, nil
}
