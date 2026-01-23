package helix

type GetTeamsParams struct {
	Name string `query:"name"`
	ID   string `query:"id"`
}

type GetTeamsResponse struct {
	ResponseCommon
	Data ManyTeams
}

type ManyTeams struct {
	Teams []Team `json:"data"`
}

type Team struct {
	Users              []TeamUser `json:"users"`
	BackgroundImageURL string     `json:"background_image_url"`
	Banner             string     `json:"banner"`
	CreatedAt          Time       `json:"created_at"`
	UpdatedAt          Time       `json:"updated_at"`
	Info               string     `json:"info"`
	ThumbnailURL       string     `json:"thumbnail_url"`
	TeamName           string     `json:"team_name"`
	TeamDisplayName    string     `json:"team_display_name"`
	ID                 string     `json:"id"`
}

type TeamUser struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

// GetTeams gets information about the specified Twitch team.
// Specify the team using the name or ID parameter (but not both).
func (c *Client) GetTeams(params *GetTeamsParams) (*GetTeamsResponse, error) {
	resp, err := c.get("/teams", &ManyTeams{}, params)
	if err != nil {
		return nil, err
	}

	teams := &GetTeamsResponse{}
	resp.HydrateResponseCommon(&teams.ResponseCommon)
	teams.Data.Teams = resp.Data.(*ManyTeams).Teams

	return teams, nil
}
