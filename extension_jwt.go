package helix

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// RoleType The user role type
type RoleType string

// Types of user roles used within the JWT Claims
const (
	BroadcasterRole RoleType = "broadcaster"
	ExternalRole    RoleType = "external"
	ModeratorRole   RoleType = "moderator"
	ViewerRole      RoleType = "viewer"

	// toAllChannels this 'channelID' is used for sending global pubsub messages
	toAllChannels = "all"
)

// PubSubPermissions publish permissions used within
// JWT claims
type PubSubPermissions struct {
	Send   []ExtensionPubSubPublishType `json:"send,omitempty"`
	Listen []ExtensionPubSubPublishType `json:"listen,omitempty"`
}

// TwitchJWTClaims contains information
// containing twitch specific JWT information.
type TwitchJWTClaims struct {
	OpaqueUserID string             `json:"opaque_user_id,omitempty"`
	UserID       string             `json:"user_id"`
	ChannelID    string             `json:"channel_id,omitempty"`
	Role         RoleType           `json:"role"`
	Unlinked     bool               `json:"is_unlinked,omitempty"`
	Permissions  *PubSubPermissions `json:"pubsub_perms"`
	jwt.StandardClaims
}

type ExtensionCreateClaimsParams struct {
	// ChannelID if this value is empty it will default to 'all'
	ChannelID string
	// PubSub is the pubsub permission to attach to the claim
	PubSub *PubSubPermissions
	// Expiration is the epoch of jwt expiration, default 3 minutes from time.Now
	Expiration int64
}

// CreateClaims will construct a claims suitable for generating a JWT token,
// containing necessary information required by the Twitch Helix Extension API endpoints.
func (c *Client) ExtensionCreateClaims(
	params *ExtensionCreateClaimsParams,
) (
	*TwitchJWTClaims,
	error,
) {
	err := c.validateExtensionOpts()
	if err != nil {
		return nil, err
	}

	// default expiration to 3 minutes
	if params.Expiration == 0 {
		params.Expiration = time.Now().Add(time.Minute*3).UnixNano() / int64(time.Millisecond)
	}

	// default channelID to 'all'
	if params.ChannelID == "" {
		params.ChannelID = toAllChannels
	}

	claims := &TwitchJWTClaims{
		UserID:      c.opts.ExtensionOpts.OwnerUserID,
		ChannelID:   params.ChannelID,
		Role:        ExternalRole,
		Permissions: params.PubSub,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: params.Expiration,
		},
	}

	return claims, nil
}

// ExtensionJWTSign Sign the a JWT Claim to produce a base64 token.
func (c *Client) ExtensionJWTSign(claims *TwitchJWTClaims) (tokenString string, err error) {

	err = c.validateExtensionOpts()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := base64.StdEncoding.DecodeString(c.opts.ExtensionOpts.Secret)
	if err != nil {
		return
	}

	tokenString, err = token.SignedString(key)
	if err != nil {
		return
	}

	return
}

// ExtensionJWTVerify validates a extension client side twitch base64 token and converts it
// into a twitch claim type, containing relevant information.
func (c *Client) ExtensionJWTVerify(token string) (claims *TwitchJWTClaims, err error) {
	if token == "" {
		err = fmt.Errorf("JWT token string missing")
		return
	}

	err = c.validateExtensionOpts()
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.ParseWithClaims(token, &TwitchJWTClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", tkn.Header["alg"])
		}

		key, err := base64.StdEncoding.DecodeString(c.opts.ExtensionOpts.Secret)

		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return
	}

	claims, ok := parsedToken.Claims.(*TwitchJWTClaims)
	if !ok || !parsedToken.Valid {
		err = fmt.Errorf("could not parse JWT")
		return
	}

	return
}

func (c *Client) validateExtensionOpts() error {
	if c.opts.ExtensionOpts.OwnerUserID == "" {
		return fmt.Errorf("extension owner id is empty")
	}

	if c.opts.ExtensionOpts.Secret == "" {
		return fmt.Errorf("extension secret is empty")
	}

	return nil
}
