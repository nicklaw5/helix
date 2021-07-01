package helix

type ChannelEditorsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type ManyChannelEditors struct {
	ChannelEditors []ChannelEditor `json:"data"`
}

// ChannelEditor
type ChannelEditor struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	CreatedAt Time   `json:"created_at"`
}

type ChannelEditorsResponse struct {
	ResponseCommon
	Data ManyChannelEditors
}

// GetChannelEditors Get a list of users who have editor permissions for a specific channel
// Required scope: channel:read:editors
func (c *Client) GetChannelEditors(params *ChannelEditorsParams) (*ChannelEditorsResponse, error) {
	resp, err := c.get("/channels/editors", &ManyChannelEditors{}, params)
	if err != nil {
		return nil, err
	}

	editors := &ChannelEditorsResponse{}
	resp.HydrateResponseCommon(&editors.ResponseCommon)
	editors.Data.ChannelEditors = resp.Data.(*ManyChannelEditors).ChannelEditors

	return editors, nil
}
