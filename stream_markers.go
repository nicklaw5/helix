package helix

type Marker struct {
	ID              string `json:"id"`
	CreatedAt       Time   `json:"created_at"`
	Description     string `json:"description"`
	PositionSeconds int    `json:"position_seconds"`
	URL             string `json:"URL"`
}

type VideoMarker struct {
	VideoID string   `json:"video_id"`
	Markers []Marker `json:"markers"`
}

type StreamMarker struct {
	UserID    string        `json:"user_id"`
	UserName  string        `json:"user_name"`
	UserLogin string        `json:"user_login"`
	Videos    []VideoMarker `json:"videos"`
}

type ManyStreamMarkers struct {
	StreamMarkers []StreamMarker `json:"data"`
	Pagination    Pagination     `json:"pagination"`
}

type StreamMarkersResponse struct {
	ResponseCommon
	Data ManyStreamMarkers
}

// StreamMarkersParams requires _either_ UserID or VideoID set
//
// UserID: fetches stream markers of the current livestream of the given user
// (VOD recording must be enabled).
// VideoID: fetches streams markers of the VOD.
type StreamMarkersParams struct {
	UserID  string `query:"user_id"`
	VideoID string `query:"video_id"`

	// Optional
	After  string `query:"after"`
	Before string `query:"before"`
	First  int    `query:"first,20"` // Limit 100
}

// GetStreamMarkers gets stream markers of a VOD or of the current live stream
// of an user being recorded as VOD.
//
// Required Scope: user:read:broadcast
func (c *Client) GetStreamMarkers(params *StreamMarkersParams) (*StreamMarkersResponse, error) {
	resp, err := c.get("/streams/markers", &ManyStreamMarkers{}, params)
	if err != nil {
		return nil, err
	}

	markers := &StreamMarkersResponse{}
	resp.HydrateResponseCommon(&markers.ResponseCommon)
	markers.Data.StreamMarkers = resp.Data.(*ManyStreamMarkers).StreamMarkers
	markers.Data.Pagination = resp.Data.(*ManyStreamMarkers).Pagination

	return markers, nil
}

type CreateStreamMarker struct {
	ID              string `json:"id"`
	CreatedAt       Time   `json:"created_at"`
	Description     string `json:"description"`
	PositionSeconds int    `json:"position_seconds"`
}

type ManyCreateStreamMarkers struct {
	CreateStreamMarkers []CreateStreamMarker `json:"data"`
}

type CreateStreamMarkerResponse struct {
	ResponseCommon
	Data ManyCreateStreamMarkers
}

type CreateStreamMarkerParams struct {
	UserID string `query:"user_id"`

	// Optional
	Description string `query:"description"`
}

// CreateStreamMarker creates a stream marker for a live stream at the current time.
// The user has to be the stream owner or an editor. Stream markers cannot be created
// in some cases, see:
// https://dev.twitch.tv/docs/api/reference/#create-stream-marker
//
// Required Scope: user:edit:broadcast
func (c *Client) CreateStreamMarker(params *CreateStreamMarkerParams) (*CreateStreamMarkerResponse, error) {
	resp, err := c.post("/streams/markers", &ManyCreateStreamMarkers{}, params)
	if err != nil {
		return nil, err
	}

	markers := &CreateStreamMarkerResponse{}
	resp.HydrateResponseCommon(&markers.ResponseCommon)
	markers.Data.CreateStreamMarkers = resp.Data.(*ManyCreateStreamMarkers).CreateStreamMarkers

	return markers, nil
}
