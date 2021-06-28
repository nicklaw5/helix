package helix

type UserExtension struct {
	CanActivate bool     `json:"can_activate"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        []string `json:"type"`
	Version     string   `json:"version"`
}

type ManyUserExtensions struct {
	UserExtensions []UserExtension `json:"data"`
}

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

type UserActiveExtensionInfo struct {
	Active  bool   `json:"active"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

type UserActiveExtension struct {
	Component map[string]UserActiveExtensionInfo `json:"component"`
	Overlay   map[string]UserActiveExtensionInfo `json:"overlay"`
	Panel     map[string]UserActiveExtensionInfo `json:"panel"`
}

type UserActiveExtensionSet struct {
	UserActiveExtensions UserActiveExtension `json:"data"`
}

type UserActiveExtensionsResponse struct {
	ResponseCommon
	Data UserActiveExtensionSet
}

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

type UpdateUserExtensionsPayload struct {
	Component map[string]UserActiveExtensionInfo `json:"component,omitempty"`
	Overlay   map[string]UserActiveExtensionInfo `json:"overlay,omitempty"`
	Panel     map[string]UserActiveExtensionInfo `json:"panel,omitempty"`
}

type wrappedUpdateUserExtensionsPayload struct {
	UpdateUserExtensionsPayload `json:"data"`
}

// UpdateUserExtensions Updates the activation state, extension ID, and/or version number of installed extensions for a specified user, identified by a Bearer token.
// If you try to activate a given extension under multiple extension types, the last write wins (and there is no guarantee of write order).
//
// Required scope: user:edit:broadcast
func (c *Client) UpdateUserExtensions(payload *UpdateUserExtensionsPayload) (*UserActiveExtensionsResponse, error) {
	normalizedPayload := &wrappedUpdateUserExtensionsPayload{UpdateUserExtensionsPayload: *payload}
	resp, err := c.putAsJSON("/users/extensions", &UserActiveExtensionSet{}, normalizedPayload)
	if err != nil {
		return nil, err
	}

	userActiveExtensions := &UserActiveExtensionsResponse{}
	userActiveExtensions.StatusCode = resp.StatusCode
	userActiveExtensions.Header = resp.Header
	userActiveExtensions.Error = resp.Error
	userActiveExtensions.ErrorStatus = resp.ErrorStatus
	userActiveExtensions.ErrorMessage = resp.ErrorMessage
	userActiveExtensions.Data.UserActiveExtensions = resp.Data.(*UserActiveExtensionSet).UserActiveExtensions

	return userActiveExtensions, nil
}
