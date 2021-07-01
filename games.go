package helix

type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

type ManyGames struct {
	Games []Game `json:"data"`
}

type GamesResponse struct {
	ResponseCommon
	Data ManyGames
}

type GamesParams struct {
	IDs   []string `query:"id"`   // Limit 100
	Names []string `query:"name"` // Limit 100
}

func (c *Client) GetGames(params *GamesParams) (*GamesResponse, error) {
	resp, err := c.get("/games", &ManyGames{}, params)
	if err != nil {
		return nil, err
	}

	games := &GamesResponse{}
	resp.HydrateResponseCommon(&games.ResponseCommon)
	games.Data.Games = resp.Data.(*ManyGames).Games

	return games, nil
}

type ManyGamesWithPagination struct {
	ManyGames
	Pagination Pagination `json:"pagination"`
}

type TopGamesParams struct {
	After  string `query:"after"`
	Before string `query:"before"`
	First  int    `query:"first,20"` // Limit 100
}

type TopGamesResponse struct {
	ResponseCommon
	Data ManyGamesWithPagination
}

func (c *Client) GetTopGames(params *TopGamesParams) (*TopGamesResponse, error) {
	resp, err := c.get("/games/top", &ManyGamesWithPagination{}, params)
	if err != nil {
		return nil, err
	}

	games := &TopGamesResponse{}
	resp.HydrateResponseCommon(&games.ResponseCommon)
	games.Data.Games = resp.Data.(*ManyGamesWithPagination).Games
	games.Data.Pagination = resp.Data.(*ManyGamesWithPagination).Pagination

	return games, nil
}
