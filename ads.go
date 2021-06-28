package helix

type AdLengthEnum int

const (
	_ AdLengthEnum = iota * 30
	AdLen30
	AdLen60
	AdLen90
	AdLen120
	AdLen150
	AdLen180
)

type StartCommercialParams struct {
	BroadcasterID string       `query:"broadcaster_id"`
	Length        AdLengthEnum `query:"length"`
}

type AdDetails struct {
	Length     AdLengthEnum `json:"length"`
	Message    string       `json:"message"`
	RetryAfter int          `json:"retry_after"`
}

type ManyAdDetails struct {
	AdDetails []AdDetails `json:"data"`
}

type StartCommercialResponse struct {
	ResponseCommon
	Data ManyAdDetails
}

// StartCommercial starts a commercial on a specified channel
// OAuth Token required
// Requires channel:edit:commercial scope
func (c *Client) StartCommercial(params *StartCommercialParams) (*StartCommercialResponse, error) {
	resp, err := c.post("/channels/commercial", &ManyAdDetails{}, params)
	if err != nil {
		return nil, err
	}

	commercials := &StartCommercialResponse{}
	resp.HydrateResponseCommon(&commercials.ResponseCommon)
	commercials.Data.AdDetails = resp.Data.(*ManyAdDetails).AdDetails

	return commercials, nil
}
