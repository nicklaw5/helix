package helix

// UserExtensions ...
type UserExtension struct {
	CanActivate bool     `json:"can_activate"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        []string `json:"type"`
	Version     string   `json:"version"`
}

// ManyUserExtensions ...
type ManyUserExtensions struct {
	UserExtensions []UserExtension `json:"data"`
}

// UserExtensionsResponse ...
type UserExtensionsResponse struct {
	ResponseCommon
	Data ManyUserExtensions
}

// GetUserExtensions gets a list of all extensions (both active and inactive) for a specified user,
// identified by a Bearer token
//
// Required scope: user:read:broadcast
func (c *Client) GetUserExtensions() (*UserExtensionsResponse, error) {
	resp, err := c.get("/users/extensions/list", &ManyUserExtensions{}, nil)
	if err != nil {
		return nil, err
	}

	userExtensions := &UserExtensionsResponse{}
	userExtensions.StatusCode = resp.StatusCode
	userExtensions.Header = resp.Header
	userExtensions.Error = resp.Error
	userExtensions.ErrorStatus = resp.ErrorStatus
	userExtensions.ErrorMessage = resp.ErrorMessage
	userExtensions.Data.UserExtensions = resp.Data.(*ManyUserExtensions).UserExtensions

	return userExtensions, nil
}
