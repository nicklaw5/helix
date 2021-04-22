package helix

import "time"

// Stream ...
type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	TagIDs       []string  `json:"tag_ids"`
	IsMature     bool      `json:"is_mature"`
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
	After      string   `query:"after"`
	Before     string   `query:"before"`
	First      int      `query:"first,20"`   // Limit 100
	GameIDs    []string `query:"game_id"`    // Limit 100
	Language   []string `query:"language"`   // Limit 100
	Type       string   `query:"type,all"`   // "all" (default), "live" and "vodcast"
	UserIDs    []string `query:"user_id"`    // limit 100
	UserLogins []string `query:"user_login"` // limit 100
}

// GetStreams returns a list of live channels based on the search parameters.
// To query offline channels, use SearchChannels.
func (c *Client) GetStreams(params *StreamsParams) (*StreamsResponse, error) {
	resp, err := c.get("/streams", &ManyStreams{}, params)
	if err != nil {
		return nil, err
	}

	streams := &StreamsResponse{}
	resp.HydrateResponseCommon(&streams.ResponseCommon)
	streams.Data.Streams = resp.Data.(*ManyStreams).Streams
	streams.Data.Pagination = resp.Data.(*ManyStreams).Pagination

	return streams, nil
}

// FollowedStreamsParams ...
type FollowedStreamsParams struct {
	After      string   `query:"after"`
	Before     string   `query:"before"`
	First      int      `query:"first,20"`   // Limit 100
	UserID     string   `query:"user_id"`
}

// GetFollowedStream : Gets information about active streams belonging to channels
// that the authenticated user follows. Streams are returned sorted by number of
// current viewers, in descending order. Across multiple pages of results, there
// may be duplicate or missing streams, as viewers join and leave streams.
//
// Required scope: user:read:follows
func (c *Client) GetFollowedStream(params *FollowedStreamsParams) (*StreamsResponse, error) {
	resp, err := c.get("/streams/followed", &ManyStreams{}, params)
	if err != nil {
		return nil, err
	}

	streams := &StreamsResponse{}
	resp.HydrateResponseCommon(&streams.ResponseCommon)
	streams.Data.Streams = resp.Data.(*ManyStreams).Streams
	streams.Data.Pagination = resp.Data.(*ManyStreams).Pagination

	return streams, nil
}


