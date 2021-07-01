package helix

type Poll struct {
	ID                         string       `json:"id"`
	BroadcasterID              string       `json:"broadcaster_id"`
	BroadcasterLogin           string       `json:"broadcaster_login"`
	BroadcasterName            string       `json:"broadcaster_name"`
	Title                      string       `json:"title"`
	Choices                    []PollChoice `json:"choices"`
	BitsVotingEnabled          bool         `json:"bits_voting_enabled"`
	BitsPerVote                int          `json:"bits_per_vote"`
	ChannelPointsVotingEnabled bool         `json:"channel_points_voting_enabled"`
	ChannelPointsPerVote       int          `json:"channel_points_per_vote"`
	Status                     string       `json:"status"`
	Duration                   int          `json:"duration"`
	StartedAt                  Time         `json:"started_at"`
	EndedAt                    Time         `json:"ended_at"`
}

type PollChoice struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	BitsVotes          int    `json:"bits_votes"`
	ChannelPointsVotes int    `json:"channel_points_votes"`
	Votes              int    `json:"votes"`
}

type ManyPolls struct {
	Polls      []Poll     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type PollsResponse struct {
	ResponseCommon
	Data ManyPolls
}

type PollsParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	ID            string `query:"id"`
	After         string `query:"after"`
	First         string `query:"first"`
}

type GetPollsResponse struct {
	ResponseCommon
	Data ManyPolls
}

// Required scope: channel:read:polls
func (c *Client) GetPolls(params *PollsParams) (*PollsResponse, error) {
	resp, err := c.get("/polls", &ManyPolls{}, params)
	if err != nil {
		return nil, err
	}

	polls := &PollsResponse{}
	resp.HydrateResponseCommon(&polls.ResponseCommon)
	polls.Data.Polls = resp.Data.(*ManyPolls).Polls
	polls.Data.Pagination = resp.Data.(*ManyPolls).Pagination

	return polls, nil
}

type CreatePollParams struct {
	BroadcasterID              string            `json:"broadcaster_id"`
	Title                      string            `json:"title"`                         // Maximum: 60 characters.
	Choices                    []PollChoiceParam `json:"choices"`                       // Minimum: 2 choices. Maximum: 5 choices.
	Duration                   int               `json:"duration"`                      // Minimum: 15. Maximum: 1800.
	BitsVotingEnabled          bool              `json:"bits_voting_enabled"`           // Default: false
	BitsPerVote                int               `json:"bits_per_vote"`                 // Minimum: 0. Maximum: 10000.
	ChannelPointsVotingEnabled bool              `json:"channel_points_voting_enabled"` // Default: false
	ChannelPointsPerVote       int               `json:"channel_points_per_vote"`       // Minimum: 0. Maximum: 1000000.
}

type PollChoiceParam struct {
	Title string `json:"title"` // Maximum: 25 characters.
}

// Required scope: channel:manage:polls
func (c *Client) CreatePoll(params *CreatePollParams) (*PollsResponse, error) {
	resp, err := c.postAsJSON("/polls", &ManyPolls{}, params)
	if err != nil {
		return nil, err
	}

	polls := &PollsResponse{}
	resp.HydrateResponseCommon(&polls.ResponseCommon)
	polls.Data.Polls = resp.Data.(*ManyPolls).Polls
	polls.Data.Pagination = resp.Data.(*ManyPolls).Pagination

	return polls, nil
}

type EndPollParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ID            string `json:"id"`
	Status        string `json:"status"`
}

// Required scope: channel:manage:polls
func (c *Client) EndPoll(params *EndPollParams) (*PollsResponse, error) {
	resp, err := c.patchAsJSON("/polls", &ManyPolls{}, params)
	if err != nil {
		return nil, err
	}

	polls := &PollsResponse{}
	resp.HydrateResponseCommon(&polls.ResponseCommon)
	polls.Data.Polls = resp.Data.(*ManyPolls).Polls
	polls.Data.Pagination = resp.Data.(*ManyPolls).Pagination

	return polls, nil
}
