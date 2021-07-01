package helix

// Prediction ... same struct as Poll
type Prediction struct {
	ID                   string     `json:"id"`
	BroadcasterUserID    string     `json:"broadcaster_id"`
	BroadcasterUserLogin string     `json:"broadcaster_login"`
	BroadcasterUserName  string     `json:"broadcaster_name"`
	Title                string     `json:"title"`
	WinningOutcomeID     string     `json:"winning_outcome_id"`
	Outcomes             []Outcomes `json:"outcomes"`
	PredictionWindow     int        `json:"prediction_window"`
	Status               string     `json:"status"`
	CreatedAt            Time       `json:"created_at"`
	EndedAt              Time       `json:"ended_at"`
	LockedAt             Time       `json:"locked_at"`
}

type Outcomes struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Users         int            `json:"users"`
	ChannelPoints int            `json:"channel_points"`
	TopPredictors []TopPredictor `json:"top_predictors"`
	Color         string         `json:"color"`
}

type TopPredictor struct {
	UserID            string `json:"user_id"`
	UserName          string `json:"user_name"`
	UserLogin         string `json:"user_login"`
	ChannelPointsUsed int    `json:"channel_points_used"`
	ChannelPointsWon  int    `json:"channel_points_won"`
}

type ManyPredictions struct {
	Predictions []Prediction `json:"data"`
	Pagination  Pagination   `json:"pagination"`
}

type PredictionsResponse struct {
	ResponseCommon
	Data ManyPredictions
}

type PredictionsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	ID            string `query:"id"`
	After         string `query:"after"`
	First         string `query:"first"`
}

type GetPredictionsResponse struct {
	ResponseCommon
	Data ManyPredictions
}

// Required scope: channel:read:predictions
func (c *Client) GetPredictions(params *PredictionsParams) (*PredictionsResponse, error) {
	resp, err := c.get("/predictions", &ManyPredictions{}, params)
	if err != nil {
		return nil, err
	}

	predictions := &PredictionsResponse{}
	resp.HydrateResponseCommon(&predictions.ResponseCommon)
	predictions.Data.Predictions = resp.Data.(*ManyPredictions).Predictions
	predictions.Data.Pagination = resp.Data.(*ManyPredictions).Pagination

	return predictions, nil
}

type CreatePredictionParams struct {
	BroadcasterID    string                  `json:"broadcaster_id"`
	Title            string                  `json:"title"`             // Maximum: 45 characters.
	Outcomes         []PredictionChoiceParam `json:"outcomes"`          // 2 choices mandatory
	PredictionWindow int                     `json:"prediction_window"` // Minimum: 1. Maximum: 1800.
}

type PredictionChoiceParam struct {
	Title string `json:"title"` // Maximum: 25 characters.
}

// Required scope: channel:manage:predictions
func (c *Client) CreatePrediction(params *CreatePredictionParams) (*PredictionsResponse, error) {
	resp, err := c.postAsJSON("/predictions", &ManyPredictions{}, params)
	if err != nil {
		return nil, err
	}

	predictions := &PredictionsResponse{}
	resp.HydrateResponseCommon(&predictions.ResponseCommon)
	predictions.Data.Predictions = resp.Data.(*ManyPredictions).Predictions
	predictions.Data.Pagination = resp.Data.(*ManyPredictions).Pagination

	return predictions, nil
}

type EndPredictionParams struct {
	BroadcasterID    string `json:"broadcaster_id"`
	ID               string `json:"id"`
	Status           string `json:"status"`
	WinningOutcomeID string `json:"winning_outcome_id"`
}

// Required scope: channel:manage:predictions
func (c *Client) EndPrediction(params *EndPredictionParams) (*PredictionsResponse, error) {
	resp, err := c.patchAsJSON("/predictions", &ManyPredictions{}, params)
	if err != nil {
		return nil, err
	}

	predictions := &PredictionsResponse{}
	resp.HydrateResponseCommon(&predictions.ResponseCommon)
	predictions.Data.Predictions = resp.Data.(*ManyPredictions).Predictions
	predictions.Data.Pagination = resp.Data.(*ManyPredictions).Pagination

	return predictions, nil
}
