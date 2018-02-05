package helix

// Clip ...
type Clip struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	EmbedURL      string `json:"embed_url"`
	BroadcasterID string `json:"broadcaster_id"`
	CreatorID     string `json:"creator_id"`
	VideoID       string `json:"video_id"`
	GameID        string `json:"game_id"`
	Language      string `json:"language"`
	Title         string `json:"title"`
	ViewCount     int    `json:"view_count"`
	CreatedAt     string `json:"created_at"`
	ThumbnailURL  string `json:"thumbnail_url"`
}

// ManyClips ...
type ManyClips struct {
	Clips []Clip `json:"data"`
}

// ClipsResponse ...
type ClipsResponse struct {
	ResponseCommon
	Data ManyClips
}

// ClipsParams ...
type ClipsParams struct {
	IDs []string `query:"id"` // Limit 1
}

// GetClips ...
func (c *Client) GetClips(params *ClipsParams) (*ClipsResponse, error) {
	resp, err := c.get("/clips", &ManyClips{}, params)
	if err != nil {
		return nil, err
	}

	clips := &ClipsResponse{}
	clips.StatusCode = resp.StatusCode
	clips.Error = resp.Error
	clips.ErrorStatus = resp.ErrorStatus
	clips.ErrorMessage = resp.ErrorMessage
	clips.RateLimit.Limit = resp.RateLimit.Limit
	clips.RateLimit.Remaining = resp.RateLimit.Remaining
	clips.RateLimit.Reset = resp.RateLimit.Reset
	clips.Data.Clips = resp.Data.(*ManyClips).Clips

	return clips, nil
}
