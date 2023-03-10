package helix

type SoundtrackCurrentTrackParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type SoundtrackTrackAlbum struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type SoundtrackTrackArtist struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	CreatorChannelID string `json:"creator_channel_id"`
}

type SoundtrackTrack struct {
	Artists  []SoundtrackTrackArtist `json:"artists"`
	ID       string                  `json:"id"`
	ISRC     string                  `json:"isrc"`
	Duration int                     `json:"duration"`
	Title    string                  `json:"title"`
	Album    SoundtrackTrackAlbum    `json:"album"`
}

type SoundtrackSource struct {
	ID            string `json:"id"`
	ContentType   string `json:"content_type"`
	Title         string `json:"title"`
	ImageURL      string `json:"image_url"`
	SoundtrackURL string `json:"soundtrack_url"`
	SpotifyURL    string `json:"spotify_url"`
}

type TrackItem struct {
	Track  SoundtrackTrack  `json:"track"`
	Source SoundtrackSource `json:"source"`
}

type ManyGetSoundTrackCurrent struct {
	Tracks []TrackItem `json:"data"`
}

type SoundtrackCurrentTrackResponse struct {
	ResponseCommon
	Data ManyGetSoundTrackCurrent
}

func (c *Client) GetSoundTrackCurrentTrack(params *SoundtrackCurrentTrackParams) (*SoundtrackCurrentTrackResponse, error) {
	resp, err := c.get("/soundtrack/current_track", &ManyGetSoundTrackCurrent{}, params)
	if err != nil {
		return nil, err
	}

	tracks := &SoundtrackCurrentTrackResponse{}
	resp.HydrateResponseCommon(&tracks.ResponseCommon)
	tracks.Data.Tracks = resp.Data.(*ManyGetSoundTrackCurrent).Tracks

	return tracks, nil
}
