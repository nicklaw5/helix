package helix

type Video struct {
	ID            string              `json:"id"`
	UserID        string              `json:"user_id"`
	UserLogin     string              `json:"user_login"`
	UserName      string              `json:"user_name"`
	StreamID      string              `json:"stream_id"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	CreatedAt     string              `json:"created_at"`
	PublishedAt   string              `json:"published_at"`
	URL           string              `json:"url"`
	ThumbnailURL  string              `json:"thumbnail_url"`
	Viewable      string              `json:"viewable"`
	ViewCount     int                 `json:"view_count"`
	Language      string              `json:"language"`
	Type          string              `json:"type"`
	Duration      string              `json:"duration"`
	MutedSegments []VideoMutedSegment `json:"muted_segments"`
}

type VideoMutedSegment struct {
	Duration int `json:"duration"`
	Offest   int `json:"offset"`
}

type ManyVideos struct {
	Videos     []Video    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type VideosParams struct {
	IDs    []string `query:"id"`      // Limit 100
	UserID string   `query:"user_id"` // Limit 1
	GameID string   `query:"game_id"` // Limit 1

	// Optional
	After    string `query:"after"`
	Before   string `query:"before"`
	First    int    `query:"first,20"`   // Limit 100
	Language string `query:"language"`   // Limit 1
	Period   string `query:"period,all"` // "all" (default), "day", "month", and "week"
	Sort     string `query:"sort,time"`  // "time" (default), "trending", and "views"
	Type     string `query:"type,all"`   // "all" (default), "upload", "archive", and "highlight"
}

type DeleteVideosParams struct {
	IDs []string `query:"id"` // Limit 5
}

type VideosResponse struct {
	ResponseCommon
	Data ManyVideos
}

type DeleteVideosResponse struct {
	ResponseCommon
}

// GetVideos gets video information by video ID (one or more), user ID (one only),
// or game ID (one only).
func (c *Client) GetVideos(params *VideosParams) (*VideosResponse, error) {
	resp, err := c.get("/videos", &ManyVideos{}, params)
	if err != nil {
		return nil, err
	}

	videos := &VideosResponse{}
	resp.HydrateResponseCommon(&videos.ResponseCommon)
	videos.Data.Videos = resp.Data.(*ManyVideos).Videos
	videos.Data.Pagination = resp.Data.(*ManyVideos).Pagination

	return videos, nil
}

// DeleteVideos delete one or more videos (max 5)
// Required scope: channel:manage:videos
func (c *Client) DeleteVideos(params *DeleteVideosParams) (*DeleteVideosResponse, error) {
	resp, err := c.delete("/videos", &DeleteVideosResponse{}, params)
	if err != nil {
		return nil, err
	}

	videos := &DeleteVideosResponse{}
	resp.HydrateResponseCommon(&videos.ResponseCommon)

	return videos, nil
}
