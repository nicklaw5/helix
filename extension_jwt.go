package helix

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// RoleType The user role type
type roleType string

// Types of user roles used within the JWT Claims
// TODO expose these when helix supports them
const (
	BroadcasterRole roleType = "broadcaster"
	ExternalRole    roleType = "external"
	ModeratorRole   roleType = "moderator"
	ViewerRole      roleType = "viewer"

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
	Role         roleType           `json:"role"`
	Unlinked     bool               `json:"is_unlinked,omitempty"`
	Permissions  *PubSubPermissions `json:"pubsub_perms"`
	jwt.StandardClaims
}

// CreateClaims will construct a claims suitable for generating a JWT token,
// containing necessary information required by the Twitch API.
// @param BroadcasterID if this value is empty it will default to 'all'
// @param pubsub the pubsub permission to attach to the claim
// @param expiration the epoch of jwt expiration, default 3 minutes from time.Now
func (c *Client) ExtensionCreateClaims(
	broadcasterID string,
	pubsub *PubSubPermissions,
	expiration int64,
) (
	*TwitchJWTClaims,
	error,
) {
	err := c.validateExtensionOpts()
	if err != nil {
		return nil, err
	}

	// default expiration to 3 minutes
	if expiration == 0 {
		expiration = time.Now().Add(time.Minute*3).UnixNano() / int64(time.Millisecond)
	}

	if broadcasterID == "" {
		broadcasterID = toAllChannels
	}

	claims := &TwitchJWTClaims{
		UserID:      c.opts.ExtensionOpts.OwnerUserID,
		ChannelID:   broadcasterID,
		Role:        ExternalRole,
		Permissions: pubsub,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
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
