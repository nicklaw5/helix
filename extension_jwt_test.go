package helix

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestValidateJwtParameters(t *testing.T) {
	t.Parallel()
	c := newMockClient(&Options{}, newMockHandler(http.StatusOK, "", nil))

	_, err := c.ExtensionCreateClaims("", c.FormBroadcastSendPubSubPermissions(), 0)
	if err == nil {
		t.Errorf("expected to get an error got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "extension owner id is empty") {
		t.Errorf("expected error extension owner id is empty, got err:%s", err)
	}

	c = newMockClient(&Options{
		ExtensionOpts: ExtensionOptions{OwnerUserID: "100249558"},
	}, newMockHandler(http.StatusOK, "", nil))

	_, err = c.ExtensionCreateClaims("", c.FormBroadcastSendPubSubPermissions(), 0)
	if err == nil {
		t.Errorf("expected to get an error got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "extension secret is empty") {
		t.Errorf("expected error extension secret is empty, got err:%s", err)
	}
}

func TestCreateClaims(t *testing.T) {
	t.Parallel()

	userId := "100249558"
	c := newMockClient(&Options{
		ExtensionOpts: ExtensionOptions{
			OwnerUserID: userId,
			Secret:      "12012311231",
		},
	}, newMockHandler(http.StatusOK, "", nil))

	claims, err := c.ExtensionCreateClaims("", c.FormBroadcastSendPubSubPermissions(), 0)
	if err != nil {
		t.Errorf("unexpected error generating claims %s", err)
	}
	if claims.UserID != userId {
		t.Errorf("claims userId doesn't match got %s expected %s", claims.UserID, userId)
	}
	if claims.ExpiresAt < time.Now().Add(4*time.Minute).UnixNano() && claims.ExpiresAt > time.Now().Add(-2*time.Minute).UnixNano() {
		t.Errorf("claims expiry less than 3 minutes")
	}

	expiration := time.Now().Add(10*time.Minute).UnixNano() / int64(time.Millisecond)
	claims, err = c.ExtensionCreateClaims("100249558", c.FormBroadcastSendPubSubPermissions(), expiration)
	if err != nil {
		t.Errorf("unexpected error generating claims %s", err)
	}

	overTime := time.Now().Add(15 * time.Minute).Unix()
	if claims.ExpiresAt < overTime {
		t.Errorf("claims expiry does not confine to 10 minutes expiry")
	}
}

func TestSignClaimsToJWT(t *testing.T) {
	t.Parallel()

	userId := "100249558"
	c := newMockClient(&Options{
		ExtensionOpts: ExtensionOptions{
			OwnerUserID: userId,
			Secret:      "TYkWIXLIKljq0e4u9id6KvqOxa80uSKKPTreIT12ERk=",
		},
	}, newMockHandler(http.StatusOK, "", nil))

	claims, err := c.ExtensionCreateClaims("100249558", c.FormBroadcastSendPubSubPermissions(), 0)
	if err != nil {
		t.Errorf("unexpected error generating claims %s", err)
	}
	jwt, err := c.ExtensionJWTSign(claims)
	if err != nil {
		t.Errorf("failed to sign claims %s", err)
	}
	if jwt == "" {
		t.Errorf("JWT token is empty")
	}
}

func TestVerifyJWT(t *testing.T) {
	t.Parallel()

	userId := "100249558"
	c := newMockClient(&Options{
		ExtensionOpts: ExtensionOptions{
			OwnerUserID: userId,
			Secret:      "TYkWIXLIKljq0e4u9id6KvqOxa80uSKKPTreIT12ERk=",
		},
	}, newMockHandler(http.StatusOK, "", nil))

	broadcasterId := "1337"
	claims, err := c.ExtensionCreateClaims(broadcasterId, c.FormBroadcastSendPubSubPermissions(), 0)
	if err != nil {
		t.Errorf("unexpected error generating claims %s", err)
	}
	jwt, err := c.ExtensionJWTSign(claims)
	if err != nil {
		t.Errorf("failed to sign claims %s", err)
	}
	if jwt == "" {
		t.Errorf("JWT token is empty")
	}

	claims, err = c.ExtensionJWTVerify("")
	if err != nil && !strings.Contains(err.Error(), "JWT token string missing") {
		t.Errorf("unexpected error verifying JWT err:%s", err)
	}

	claims, err = c.ExtensionJWTVerify(jwt)
	if err != nil && !strings.Contains(err.Error(), "JWT token string missing") {
		t.Errorf("unexpected error verifying JWT err:%s", err)
	}
	if claims.ChannelID != broadcasterId {
		t.Errorf("found unexpected broadcaster in claims got:%s expected:%s", claims.ChannelID, broadcasterId)
	}
	if claims.UserID != userId {
		t.Errorf("found unexpected userId in claims got:%s expected:%s", claims.UserID, userId)
	}
}
