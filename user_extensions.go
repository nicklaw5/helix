package helix

// UserExtension ...
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
	resp.HydrateResponseCommon(&userExtensions.ResponseCommon)
	userExtensions.Data.UserExtensions = resp.Data.(*ManyUserExtensions).UserExtensions

	return userExtensions, nil
}

// UserActiveExtensionInfo ...
type UserActiveExtensionInfo struct {
	Active  bool   `json:"active"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

// UserActiveExtension ...
type UserActiveExtension struct {
	Component map[string]UserActiveExtensionInfo `json:"component"`
	Overlay   map[string]UserActiveExtensionInfo `json:"overlay"`
	Panel     map[string]UserActiveExtensionInfo `json:"panel"`
}

// UserActiveExtensionSet ...
type UserActiveExtensionSet struct {
	UserActiveExtensions UserActiveExtension `json:"data"`
}

// UserActiveExtensionsResponse ...
type UserActiveExtensionsResponse struct {
	ResponseCommon
	Data UserActiveExtensionSet
}

// UserActiveExtensionsParams ...
type UserActiveExtensionsParams struct {
	UserID string `query:"user_id"` // Optional, limit 1
}

// GetUserActiveExtensions Gets information about active extensions installed by a specified user, identified
// by a user ID or Bearer token.
//
// Optional scope: user:read:broadcast or user:edit:broadcast
func (c *Client) GetUserActiveExtensions(params *UserActiveExtensionsParams) (*UserActiveExtensionsResponse, error) {
	resp, err := c.get("/users/extensions", &UserActiveExtensionSet{}, params)
	if err != nil {
		return nil, err
	}

	userActiveExtensions := &UserActiveExtensionsResponse{}
	resp.HydrateResponseCommon(&userActiveExtensions.ResponseCommon)
	userActiveExtensions.Data.UserActiveExtensions = resp.Data.(*UserActiveExtensionSet).UserActiveExtensions

	return userActiveExtensions, nil
}
