package helix

import (
	"fmt"
)

// User (Use Helix Twitch API)
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

// UsersResponse ...
type UsersResponse struct {
	ResponseCommon
	Data []User `json:"data"`
}

// UsersRequest ...
type UsersRequest struct {
	IDs    []string
	Logins []string
}

// GetUsers ...
func (c *Client) GetUsers(req *UsersRequest) (*UsersResponse, error) {
	var query string

	if req.IDs != nil {
		query = fmt.Sprintf("id=%s", concatString(req.IDs, "&id="))
	}
	if req.Logins != nil {
		if query != "" {
			query += "&"
		}
		query = fmt.Sprintf("%slogin=%s", query, concatString(req.Logins, "&login="))
	}

	resp := &UsersResponse{}
	err := c.Get("/users?"+query, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
