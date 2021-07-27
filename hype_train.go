package helix

type HypeTrainContribuition struct {
	Total int64  `json:"total"`
	Type  string `json:"type"`
	User  string `json:"user"`
}

type HypeTrainEvent struct {
	ID             string             `json:"id"`
	EventType      string             `json:"event_type"`
	EventTimestamp Time               `json:"event_timestamp"`
	Version        string             `json:"version"`
	Event          HypeTrainEventData `json:"event_data"`
}

type HypeTrainEventData struct {
	ID               string                   `json:"id"`
	BroadcasterID    string                   `json:"broadcaster_id"`
	CooldownEndTime  Time                     `json:"cooldown_end_time"`
	ExpiresAt        Time                     `json:"expires_at"`
	Goal             int64                    `json:"goal"`
	LastContribution HypeTrainContribuition   `json:"last_contribution"`
	Level            int64                    `json:"level"`
	StartedAt        Time                     `json:"started_at"`
	TopContributions []HypeTrainContribuition `json:"top_contributions"`
	Total            int64                    `json:"total"`
}

type ManyHypeTrainEvents struct {
	Events     []HypeTrainEvent `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

type HypeTrainEventsResponse struct {
	ResponseCommon
	Data ManyHypeTrainEvents
}

type HypeTrainEventsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	After         string `query:"after"`
	First         int    `query:"first,20"` // Limit 100
	ID            string `query:"id"`
}

// Required scope: channel:read:hype_train
func (c *Client) GetHypeTrainEvents(params *HypeTrainEventsParams) (*HypeTrainEventsResponse, error) {
	resp, err := c.get("/hypetrain/events", &ManyHypeTrainEvents{}, params)
	if err != nil {
		return nil, err
	}

	events := &HypeTrainEventsResponse{}
	resp.HydrateResponseCommon(&events.ResponseCommon)
	events.Data.Events = resp.Data.(*ManyHypeTrainEvents).Events
	events.Data.Pagination = resp.Data.(*ManyHypeTrainEvents).Pagination

	return events, nil
}
