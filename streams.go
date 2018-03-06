package helix

import "time"

// Stream ...
type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	GameID       string    `json:"game_id"`
	CommunityIDs []string  `json:"community_ids"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

// ManyStreams ...
type ManyStreams struct {
	Streams    []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// StreamsResponse ...
type StreamsResponse struct {
	ResponseCommon
	Data ManyStreams
}

// StreamsParams ...
type StreamsParams struct {
	After        string   `query:"after"`
	Before       string   `query:"before"`
	CommunityIDs []string `query:"community_id"`
	First        int      `query:"first,20"`   // Limit 100
	GameIDs      []string `query:"game_id"`    // Limit 100
	Language     []string `query:"language"`   // Limit 100
	Type         string   `query:"type,all"`   // "all" (default), "live" and "vodcast"
	UserIDs      []string `query:"user_id"`    // limit 100
	UserLogins   []string `query:"user_login"` // limit 100
}

// GetStreams ...
func (c *Client) GetStreams(params *StreamsParams) (*StreamsResponse, error) {
	resp, err := c.get("/streams", &ManyStreams{}, params)
	if err != nil {
		return nil, err
	}

	streams := &StreamsResponse{}
	streams.StatusCode = resp.StatusCode
	streams.Error = resp.Error
	streams.ErrorStatus = resp.ErrorStatus
	streams.ErrorMessage = resp.ErrorMessage
	streams.RateLimit.Limit = resp.RateLimit.Limit
	streams.RateLimit.Remaining = resp.RateLimit.Remaining
	streams.RateLimit.Reset = resp.RateLimit.Reset
	streams.Data.Streams = resp.Data.(*ManyStreams).Streams
	streams.Data.Pagination = resp.Data.(*ManyStreams).Pagination

	return streams, nil
}
