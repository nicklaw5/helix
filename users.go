package helix

// User ...
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

// ManyUsers ...
type ManyUsers struct {
	Users []User `json:"data"`
}

// UsersResponse ...
type UsersResponse struct {
	ResponseCommon
	Data ManyUsers
}

// UsersParams ...
type UsersParams struct {
	IDs    []string `query:"id"`    // Limit 100
	Logins []string `query:"login"` // Limit 100
}

// GetUsers ...
func (c *Client) GetUsers(params *UsersParams) (*UsersResponse, error) {
	resp, err := c.get("/users", &ManyUsers{}, params)
	if err != nil {
		return nil, err
	}

	users := &UsersResponse{}
	users.StatusCode = resp.StatusCode
	users.Error = resp.Error
	users.ErrorStatus = resp.ErrorStatus
	users.ErrorMessage = resp.ErrorMessage
	users.RatelimitLimit = resp.RatelimitLimit
	users.RatelimitRemaining = resp.RatelimitRemaining
	users.RatelimitReset = resp.RatelimitReset
	users.Data.Users = resp.Data.(*ManyUsers).Users

	return users, nil
}
