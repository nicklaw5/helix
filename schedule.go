package helix

// GetScheduleParams are the parameters for GetSchedule
type GetScheduleParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ID            string `json:"id"`
	StartTime     Time   `json:"start_time"`
	UTCOffset     string `json:"utc_offset"`
	First         int    `json:"first"`
	After         string `json:"after"`
}

// GetScheduleResponse is the response data in GetSchedule
type GetScheduleResponse struct {
	ResponseCommon

	Data GetScheduleData
}

type GetScheduleData struct {
	Schedule   ScheduleData          `json:"data"`
	Pagination GetSchedulePagination `json:"pagination"`
}

type ScheduleData struct {
	Segments         []GetScheduleSegment `json:"segments"`
	BroadcasterID    string               `json:"broadcaster_id"`
	BroadcasterName  string               `json:"broadcaster_name"`
	BroadcasterLogin string               `json:"broadcaster_login"`
	Vacation         GetScheduleVacation  `json:"vacation"`
}

type GetScheduleSegment struct {
	ID            string                     `json:"id"`
	StartTime     Time                       `json:"start_time"`
	EndTime       Time                       `json:"end_time"`
	Title         string                     `json:"title"`
	CanceledUntil string                     `json:"canceled_until"`
	Category      GetScheduleSegmentCategory `json:"category"`
	IsRecurring   bool                       `json:"is_recurring"`
}

type GetScheduleSegmentCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetScheduleVacation struct {
	StartTime Time `json:"start_time"`
	EndTime   Time `json:"end_time"`
}

type GetSchedulePagination struct {
	Cursor string `json:"cursor"`
}

// Gets the broadcaster’s streaming schedule.
// You can get the entire schedule or specific segments of the schedule
func (c *Client) GetSchedule(params *GetScheduleParams) (*GetScheduleResponse, error) {
	resp, err := c.get("/schedule", &GetScheduleData{}, params)
	if err != nil {
		return nil, err
	}

	schedule := &GetScheduleResponse{}
	resp.HydrateResponseCommon(&schedule.ResponseCommon)
	schedule.Data.Schedule = resp.Data.(*GetScheduleData).Schedule
	schedule.Data.Pagination = resp.Data.(*GetScheduleData).Pagination

	return schedule, nil
}

type UpdateScheduleParams struct {
	BroadcasterID     string `json:"broadcaster_id"`
	IsVacationEnabled bool   `json:"is_vacation_enabled"`
	VacationStartTime Time   `json:"vacation_start_time"`
	VacationEndTime   Time   `json:"vacation_end_time"`
	Timezone          string `json:"timezone"`
}

type UpdateScheduleResponse struct {
	ResponseCommon
}

// Updates the broadcaster’s schedule settings, such as scheduling a vacation
func (c *Client) UpdateSchedule(params *UpdateScheduleParams) (*UpdateScheduleResponse, error) {
	resp, err := c.get("/schedule/settings", nil, params)
	if err != nil {
		return nil, err
	}

	schedule := &UpdateScheduleResponse{}
	resp.HydrateResponseCommon(&schedule.ResponseCommon)

	return schedule, nil
}

type CreateScheduleSegmentParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	StartTime     Time   `json:"start_time"`
	Timezone      string `json:"timezone"`
	Duration      string `json:"duration"`
	IsRecurring   bool   `json:"is_recurring"`
	CategoryID    string `json:"category_id"`
	Title         string `json:"title"`
}

type CreateScheduleSegmentResponse struct {
	ResponseCommon

	Data CreateScheduleSegmentData
}

type CreateScheduleSegmentData struct {
	Schedule ScheduleData `json:"data"`
}

// Updates the broadcaster’s schedule settings, such as scheduling a vacation
func (c *Client) CreateScheduleSegment(params *CreateScheduleSegmentParams) (*CreateScheduleSegmentResponse, error) {
	resp, err := c.post("/schedule/segment", &CreateScheduleSegmentData{}, params)
	if err != nil {
		return nil, err
	}

	schedule := &CreateScheduleSegmentResponse{}
	resp.HydrateResponseCommon(&schedule.ResponseCommon)
	schedule.Data.Schedule = resp.Data.(*CreateScheduleSegmentData).Schedule

	return schedule, nil
}

type UpdateScheduleSegmentParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ID            string `json:"id"`
	StartTime     Time   `json:"start_time"`
	Duration      string `json:"duration"`
	CategoryID    string `json:"category_id"`
	Title         string `json:"title"`
	IsCanceled    bool   `json:"is_canceled"`
	Timezone      string `json:"timezone"`
}

type UpdateScheduleSegmentResponse struct {
	ResponseCommon

	Data UpdateScheduleSegmentData
}

type UpdateScheduleSegmentData struct {
	Schedule ScheduleData `json:"data"`
}

// Updates the broadcaster’s schedule settings, such as scheduling a vacation
func (c *Client) UpdateScheduleSegment(params *UpdateScheduleSegmentParams) (*UpdateScheduleSegmentResponse, error) {
	resp, err := c.patchAsJSON("/schedule/segment", &UpdateScheduleSegmentData{}, params)
	if err != nil {
		return nil, err
	}

	schedule := &UpdateScheduleSegmentResponse{}
	resp.HydrateResponseCommon(&schedule.ResponseCommon)
	schedule.Data.Schedule = resp.Data.(*UpdateScheduleSegmentData).Schedule

	return schedule, nil
}

type DeleteScheduleSegmentParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ID            string `json:"id"`
}

type DeleteScheduleSegmentResponse struct {
	ResponseCommon
}

// Removes a broadcast segment from the broadcaster’s streaming schedule
func (c *Client) DeleteScheduleSegment(params *DeleteScheduleSegmentParams) (*DeleteScheduleSegmentResponse, error) {
	resp, err := c.delete("/schedule/segment", nil, params)
	if err != nil {
		return nil, err
	}

	schedule := &DeleteScheduleSegmentResponse{}
	resp.HydrateResponseCommon(&schedule.ResponseCommon)

	return schedule, nil
}
