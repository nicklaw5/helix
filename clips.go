package helix

type Clip struct {
	ID              string  `json:"id"`
	URL             string  `json:"url"`
	EmbedURL        string  `json:"embed_url"`
	BroadcasterID   string  `json:"broadcaster_id"`
	BroadcasterName string  `json:"broadcaster_name"`
	CreatorID       string  `json:"creator_id"`
	CreatorName     string  `json:"creator_name"`
	Duration        float64 `json:"duration"`
	VideoID         string  `json:"video_id"`
	GameID          string  `json:"game_id"`
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	ViewCount       int     `json:"view_count"`
	CreatedAt       string  `json:"created_at"`
	ThumbnailURL    string  `json:"thumbnail_url"`
	VodOffset       int     `json:"vod_offset"`
}

type ManyClips struct {
	Clips      []Clip     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type ClipsResponse struct {
	ResponseCommon
	Data ManyClips
}

type ClipsParams struct {
	// One of the below
	BroadcasterID string   `query:"broadcaster_id"`
	GameID        string   `query:"game_id"`
	IDs           []string `query:"id"` // Limit 100

	// Optional
	First     int    `query:"first,20"` // Maximum 100
	After     string `query:"after"`
	Before    string `query:"before"`
	StartedAt Time   `query:"started_at"`
	EndedAt   Time   `query:"ended_at"`
}

// GetClips returns information about a specified clip.
func (c *Client) GetClips(params *ClipsParams) (*ClipsResponse, error) {
	resp, err := c.get("/clips", &ManyClips{}, params)
	if err != nil {
		return nil, err
	}

	clips := &ClipsResponse{}
	resp.HydrateResponseCommon(&clips.ResponseCommon)
	clips.Data.Clips = resp.Data.(*ManyClips).Clips
	clips.Data.Pagination = resp.Data.(*ManyClips).Pagination

	return clips, nil
}

type ClipEditURL struct {
	ID      string `json:"id"`
	EditURL string `json:"edit_url"`
}

type ManyClipEditURLs struct {
	ClipEditURLs []ClipEditURL `json:"data"`
}

type CreateClipResponse struct {
	ResponseCommon
	Data ManyClipEditURLs
}

// GetClipsCreationRateLimit returns the "Ratelimit-Helixclipscreation-Limit"
// header as an int.
func (ccr *CreateClipResponse) GetClipsCreationRateLimit() int {
	return ccr.convertHeaderToInt(ccr.Header.Get("Ratelimit-Helixclipscreation-Limit"))
}

// GetClipsCreationRateLimitRemaining returns the "Ratelimit-Helixclipscreation-Remaining"
// header as an int.
func (ccr *CreateClipResponse) GetClipsCreationRateLimitRemaining() int {
	return ccr.convertHeaderToInt(ccr.Header.Get("Ratelimit-Helixclipscreation-Remaining"))
}

type CreateClipParams struct {
	BroadcasterID string `query:"broadcaster_id"`

	// Optional
	HasDelay bool `query:"has_delay,false"`
}

// CreateClip creates a clip programmatically. This returns both an ID and
// an edit URL for the new clip. Clip creation takes time. We recommend that
// you query Get Clip, with the clip ID that is returned here. If Get Clip
// returns a valid clip, your clip creation was successful. If, after 15 seconds,
// you still have not gotten back a valid clip from Get Clip, assume that the
// clip was not created and retry Create Clip.
//
// Required scope: clips:edit
func (c *Client) CreateClip(params *CreateClipParams) (*CreateClipResponse, error) {
	resp, err := c.post("/clips", &ManyClipEditURLs{}, params)
	if err != nil {
		return nil, err
	}

	clips := &CreateClipResponse{}
	resp.HydrateResponseCommon(&clips.ResponseCommon)
	clips.Data.ClipEditURLs = resp.Data.(*ManyClipEditURLs).ClipEditURLs

	return clips, nil
}
