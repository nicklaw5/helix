package helix

type Team struct {
	ID                 string     `json:"id"`
	Users              []TeamUser `json:"users"`
	BackgroundImageURL string     `json:"background_image_url"`
	Banner             string     `json:"banner"`
	CreatedAt          Time       `json:"created_at"`
	UpdatedAt          Time       `json:"updated_at"`
	Info               string     `json:"info"`
	ThumbnailURL       string     `json:"thumbnail_url"`
	TeamName           string     `json:"team_name"`
	TeamDisplayName    string     `json:"team_display_name"`
}

type ChannelTeam struct {
	ID                 string `json:"id"`
	BroadcasterID      string `json:"broadcaster_id"`
	BroadcasterLogin   string `json:"broadcaster_login"`
	BroadcasterName    string `json:"broadcaster_name"`
	BackgroundImageURL string `json:"background_image_url"`
	Banner             string `json:"banner"`
	CreatedAt          Time   `json:"created_at"`
	UpdatedAt          Time   `json:"updated_at"`
	Info               string `json:"info"`
	ThumbnailURL       string `json:"thumbnail_url"`
	TeamName           string `json:"team_name"`
	TeamDisplayName    string `json:"team_display_name"`
}

type TeamUser struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ManyTeams struct {
	Teams []Team `json:"data"`
}

type ManyChannelTeams struct {
	ChannelTeams []ChannelTeam `json:"data"`
}

type TeamsResponse struct {
	ResponseCommon
	Data ManyTeams
}

type TeamsParams struct {
	ID   string `query:"id"`
	Name string `query:"name"`
}

type ChannelTeamsResponse struct {
	ResponseCommon
	Data ManyChannelTeams
}

type ChannelTeamsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

func (c *Client) GetTeams(params *TeamsParams) (*TeamsResponse, error) {
	resp, err := c.get("/teams", &ManyTeams{}, params)
	if err != nil {
		return nil, err
	}

	teams := &TeamsResponse{}
	resp.HydrateResponseCommon(&teams.ResponseCommon)
	teams.Data.Teams = resp.Data.(*ManyTeams).Teams

	return teams, nil
}

func (c *Client) GetChannelTeams(params *ChannelTeamsParams) (*ChannelTeamsResponse, error) {
	resp, err := c.get("/teams/channel", &ManyChannelTeams{}, params)
	if err != nil {
		return nil, err
	}

	channel_teams := &ChannelTeamsResponse{}
	resp.HydrateResponseCommon((&channel_teams.ResponseCommon))
	channel_teams.Data.ChannelTeams = resp.Data.(*ManyChannelTeams).ChannelTeams

	return channel_teams, nil
}
