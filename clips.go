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

// ClipsResponse ...
type ClipsResponse struct {
	ResponseCommon
	Data []Clip `json:"data"`
}

// GetClip ...
func (c *Client) GetClip(clipID string) (*ClipsResponse, error) {
	resp := &ClipsResponse{}
	err := c.Get("/clips?id="+clipID, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
