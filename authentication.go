package helix

import (
	"strings"
)

// GetAuthorizationURL ...
func (c *Client) GetAuthorizationURL(state string, forceVerify bool) string {
	url := AuthBaseURL + "/authorize?response_type=code"
	url += "&client_id=" + c.clientID
	url += "&redirect_uri=" + c.redirectURI

	if state != "" {
		url += "&state=" + state
	}

	if forceVerify {
		url += "&force_verify=true"
	}

	if len(c.scopes) > 0 {
		url += "&scope=" + strings.Join(c.scopes, "%20")
	}

	return url
}
