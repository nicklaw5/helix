package helix

type SendUserWhisperParams struct {
	FromUserID string `query:"from_user_id"`
	ToUserID   string `query:"to_user_id"`
	Message    string `json:"message"`
}

type SendUserWhisperResponse struct {
	ResponseCommon
}

// SendUserWhisper
// requires user access token with user:manage:whispers scope.
// The user sending the whisper must have a verified phone number.
// The API may silently drop whispers that it suspects of violating Twitch policies 204 still returned.
// You may whisper to a maximum of 40 unique recipients per day. Within the per day limit.
// you may whisper a maximum of 3 whispers per second and a maximum of 100 whispers per minute.
// messages character limit:
//   - 500 chars to new recipient
//   - 10,000 if recurring recipient,
//   - > 10,000 chars are truancated.
func (c *Client) SendUserWhisper(params *SendUserWhisperParams) (*SendUserWhisperResponse, error) {
	resp, err := c.postAsJSON("/whispers", nil, params)
	if err != nil {
		return nil, err
	}

	whisperResp := &SendUserWhisperResponse{}
	resp.HydrateResponseCommon(&whisperResp.ResponseCommon)

	return whisperResp, nil
}
