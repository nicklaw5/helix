package helix

// Game ...
type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

// ManyGames ...
type ManyGames struct {
	Games []Game `json:"data"`
}

// GamesResponse ...
type GamesResponse struct {
	ResponseCommon
	Data ManyGames
}

// GamesParams ...
type GamesParams struct {
	IDs   []string `query:"id"`   // Limit 100
	Names []string `query:"name"` // Limit 100
}

// GetGames ...
func (c *Client) GetGames(params *GamesParams) (*GamesResponse, error) {
	resp, err := c.get("/games", &ManyGames{}, params)
	if err != nil {
		return nil, err
	}

	games := &GamesResponse{}
	games.StatusCode = resp.StatusCode
	games.Error = resp.Error
	games.ErrorStatus = resp.ErrorStatus
	games.ErrorMessage = resp.ErrorMessage
	games.RateLimit.Limit = resp.RateLimit.Limit
	games.RateLimit.Remaining = resp.RateLimit.Remaining
	games.RateLimit.Reset = resp.RateLimit.Reset
	games.Data.Games = resp.Data.(*ManyGames).Games

	return games, nil
}
