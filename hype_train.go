package helix

type HypeTrainContribuition struct {
	Total int64  `json:"total"`
	Type  string `json:"type"`
	User  string `json:"user"`
}

// HypeTrainStatusContribution contains information about a contribution to a Hype Train
type HypeTrainStatusContribution struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Type      string `json:"type"`
	Total     int64  `json:"total"`
}

type SharedTrainParticipant struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

type CurrentHypeTrainStatus struct {
	ID                      string                        `json:"id"`
	BroadcasterUserID       string                        `json:"broadcaster_user_id"`
	BroadcasterUserLogin    string                        `json:"broadcaster_user_login"`
	BroadcasterUserName     string                        `json:"broadcaster_user_name"`
	Level                   int64                         `json:"level"`
	Total                   int64                         `json:"total"`
	Progress                int64                         `json:"progress"`
	Goal                    int64                         `json:"goal"`
	TopContributions        []HypeTrainStatusContribution `json:"top_contributions"`
	SharedTrainParticipants []SharedTrainParticipant      `json:"shared_train_participants"`
	StartedAt               Time                          `json:"started_at"`
	ExpiresAt               Time                          `json:"expires_at"`
	Type                    string                        `json:"type"`
}

type AllTimeHighHypeTrainStatus struct {
	Level      int64 `json:"level"`
	Total      int64 `json:"total"`
	AchievedAt Time  `json:"achieved_at"`
}

type HypeTrainStatus struct {
	Current            *CurrentHypeTrainStatus     `json:"current"`
	AllTimeHigh        *AllTimeHighHypeTrainStatus `json:"all_time_high"`
	SharedAllTimeHigh  *AllTimeHighHypeTrainStatus `json:"shared_all_time_high"`
}

type ManyHypeTrainStatuses struct {
	Statuses []HypeTrainStatus `json:"data"`
}

type HypeTrainStatusResponse struct {
	ResponseCommon
	Data ManyHypeTrainStatuses
}

type HypeTrainStatusParams struct {
	BroadcasterID string `query:"broadcaster_id"`
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

// GetHypeTrainStatus gets the Hype Train status for the specified broadcaster.
// Required scope: channel:read:hype_train
func (c *Client) GetHypeTrainStatus(params *HypeTrainStatusParams) (*HypeTrainStatusResponse, error) {
	resp, err := c.get("/hypetrain/status", &ManyHypeTrainStatuses{}, params)
	if err != nil {
		return nil, err
	}

	status := &HypeTrainStatusResponse{}
	resp.HydrateResponseCommon(&status.ResponseCommon)
	status.Data.Statuses = resp.Data.(*ManyHypeTrainStatuses).Statuses

	return status, nil
}

// GetHypeTrainEvents gets Hype Train events for a broadcaster.
// Deprecated: Use GetHypeTrainStatus instead.
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
