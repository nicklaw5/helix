package helix

import (
	"time"
)

// StreamMarkerAPIDetails ...
type StreamMarkerAPIDetails struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Description     string    `json:"description"`
	PositionSeconds int       `json:"position_seconds"`
	URL             string    `json:"URL"`
}

// StreamMarkerAPIVideo ...
type StreamMarkerAPIVideo struct {
	VideoID string                   `json:"video_id"`
	Markers []StreamMarkerAPIDetails `json:"markers"`
}

// StreamMarkersAPIResponseData ...
type StreamMarkersAPIResponseData struct {
	UserID   string                 `json:"user_id"`
	UserName string                 `json:"user_name"`
	Videos   []StreamMarkerAPIVideo `json:"videos"`
}

// StreamMarkersAPIResponse ...
type StreamMarkersAPIResponse struct {
	Data       []StreamMarkersAPIResponseData `json:"data"`
	Pagination Pagination                     `json:"pagination"`
}

// StreamMarkersResponseData ...
type StreamMarkersResponseData struct {
	UserID   string                   `json:"user_id"`
	UserName string                   `json:"user_name"`
	VideoID  string                   `json:"video_id"`
	Markers  []StreamMarkerAPIDetails `json:"markers"`
}

// StreamMarkersResponse ...
type StreamMarkersResponse struct {
	ResponseCommon

	Data       StreamMarkersResponseData
	Pagination Pagination `json:"pagination"`
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
	apiResponse := &StreamMarkersAPIResponse{}
	resp, err := c.get("/streams/markers", apiResponse, params)
	if err != nil {
		return nil, err
	}

	responseData := StreamMarkersResponseData{}
	if len(apiResponse.Data) > 0 {
		responseData.UserID = apiResponse.Data[0].UserID
		responseData.UserName = apiResponse.Data[0].UserName
		responseData.VideoID = apiResponse.Data[0].Videos[0].VideoID
		responseData.Markers = apiResponse.Data[0].Videos[0].Markers
	}

	streamMarkers := &StreamMarkersResponse{
		Data: responseData,
		ResponseCommon: ResponseCommon{
			StatusCode:   resp.StatusCode,
			Header:       resp.Header,
			Error:        resp.Error,
			ErrorStatus:  resp.ErrorStatus,
			ErrorMessage: resp.ErrorMessage,
		},
	}

	return streamMarkers, nil
}
