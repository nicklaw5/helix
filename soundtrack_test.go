package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetSoundTrackCurrentTrack(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		BroadcasterID string
		respBody      string
		parsed        []TrackItem
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"145328278",
			`{
			  "data": [
				{
				  "track": {
					"artists": [
					  {
						"id": "B07S7JG3TK",
						"name": "Enoth",
						"creator_channel_id": "39051113"
					  }
					],
					"id": "B08D6QFS38",
					"isrc": "CCXXXYYNNNNN",
					"duration": 153,
					"title": "Please stay",
					"album": {
					  "id": "B08D6PMKYL",
					  "name": "Summer 2020",
					  "image_url": "https://m.media-amazon.com/images/I/51zs1JZY8tL.jpg"
					}
				  },
				  "source": {
					"id": "B08HCW84SF",
					"content_type": "PLAYLIST",
					"title": "Beats To Stream To",
					"image_url": "https://m.media-amazon.com/images/I/419WuvMXzEL.jpg",
					"soundtrack_url": "https://soundtrack.twitch.tv/playlist?playlistID=B08HCW84SF",
					"spotify_url": "https://open.spotify.com/playlist/1LOP14236oTUscowY3NvYN"
				  }
				}
			  ]
			}`,
			[]TrackItem{
				{
					Track: SoundtrackTrack{
						Artists: []SoundtrackTrackArtist{
							{ID: "B07S7JG3TK", Name: "Enoth", CreatorChannelID: "39051113"},
						},
						ID:       "B08D6QFS38",
						ISRC:     "CCXXXYYNNNNN",
						Duration: 153,
						Title:    "Please stay",
						Album: SoundtrackTrackAlbum{
							ID:       "B08D6PMKYL",
							Name:     "Summer 2020",
							ImageURL: "https://m.media-amazon.com/images/I/51zs1JZY8tL.jpg",
						},
					},
					Source: SoundtrackSource{
						ID:            "B08HCW84SF",
						ContentType:   "PLAYLIST",
						Title:         "Beats To Stream To",
						ImageURL:      "https://m.media-amazon.com/images/I/419WuvMXzEL.jpg",
						SoundtrackURL: "https://soundtrack.twitch.tv/playlist?playlistID=B08HCW84SF",
						SpotifyURL:    "https://open.spotify.com/playlist/1LOP14236oTUscowY3NvYN",
					}},
			},
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			"",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
			nil,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetSoundTrackCurrentTrack(&SoundtrackCurrentTrackParams{
			BroadcasterID: testCase.BroadcasterID,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		for i, track := range resp.Data.Tracks {
			if track.Track.ID != testCase.parsed[i].Track.ID {
				t.Errorf("Expected struct field ID = %s, was %s", testCase.parsed[i].Track.ID, track.Track.ID)
			}
			if track.Track.Title != testCase.parsed[i].Track.Title {
				t.Errorf("Expected struct field Title = %s, was %s", testCase.parsed[i].Track.Title, track.Track.Title)
			}
			if track.Track.ISRC != testCase.parsed[i].Track.ISRC {
				t.Errorf("Expected struct field ISRC = %s, was %s", testCase.parsed[i].Track.ISRC, track.Track.ISRC)
			}
			if track.Track.Duration != testCase.parsed[i].Track.Duration {
				t.Errorf("Expected struct field Duration = %d, was %d", testCase.parsed[i].Track.Duration, track.Track.Duration)
			}
			if track.Track.Album.ID != testCase.parsed[i].Track.Album.ID {
				t.Errorf("Expected struct field Album.ID = %s, was %s", testCase.parsed[i].Track.Album.ID, track.Track.Album.ID)
			}
			if track.Track.Album.Name != testCase.parsed[i].Track.Album.Name {
				t.Errorf("Expected struct field Album.Name = %s, was %s", testCase.parsed[i].Track.Album.Name, track.Track.Album.Name)
			}
			if track.Track.Album.ImageURL != testCase.parsed[i].Track.Album.ImageURL {
				t.Errorf("Expected struct field Album.ImageURL = %s, was %s", testCase.parsed[i].Track.Album.ImageURL, track.Track.Album.ImageURL)
			}
			if track.Source.ID != testCase.parsed[i].Source.ID {
				t.Errorf("Expected struct field Source.ID = %s, was %s", testCase.parsed[i].Source.ID, track.Source.ID)
			}
			if track.Source.ContentType != testCase.parsed[i].Source.ContentType {
				t.Errorf("Expected struct field Source.ContentType = %s, was %s", testCase.parsed[i].Source.ContentType, track.Source.ContentType)
			}
			if track.Source.Title != testCase.parsed[i].Source.Title {
				t.Errorf("Expected struct field Source.Title = %s, was %s", testCase.parsed[i].Source.Title, track.Source.Title)
			}
			if track.Source.ImageURL != testCase.parsed[i].Source.ImageURL {
				t.Errorf("Expected struct field Source.ImageURL = %s, was %s", testCase.parsed[i].Source.ImageURL, track.Source.ImageURL)
			}
			if track.Source.SoundtrackURL != testCase.parsed[i].Source.SoundtrackURL {
				t.Errorf("Expected struct field Source.SoundtrackURL = %s, was %s", testCase.parsed[i].Source.SoundtrackURL, track.Source.SoundtrackURL)
			}
			if track.Source.SpotifyURL != testCase.parsed[i].Source.SpotifyURL {
				t.Errorf("Expected struct field Source.SpotifyURL = %s, was %s", testCase.parsed[i].Source.SpotifyURL, track.Source.SpotifyURL)
			}

		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
		ctx:  context.Background(),
	}

	_, err := c.GetSoundTrackCurrentTrack(&SoundtrackCurrentTrackParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
