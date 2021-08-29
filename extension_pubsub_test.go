package helix

import (
	"net/http"
	"strings"
	"testing"
)

func TestExtensionFormPubSubPerms(t *testing.T) {
	t.Parallel()

	c := newMockClient(&Options{}, newMockHandler(http.StatusOK, "", nil))

	perms := c.FormGlobalSendPubSubPermissions()
	if len(perms.Send) == 0 {
		t.Errorf(
			"invalid pubsub permission type, expected 1 got %d",
			len(perms.Send),
		)
	}
	if len(perms.Send) > 1 && perms.Send[0] != ExtensionPubSubGlobalPublish {
		t.Errorf(
			"invalid pubsub permission type, expected %s got %s",
			ExtensionPubSubGlobalPublish.String(),
			perms.Send[0],
		)
	}

	perms = c.FormBroadcastSendPubSubPermissions()
	if len(perms.Send) == 0 {
		t.Errorf(
			"invalid pubsub permission type, expected 1 got %d",
			len(perms.Send),
		)
	}
	if len(perms.Send) > 1 && perms.Send[0] != ExtensionPubSubBroadcastPublish {
		t.Errorf(
			"invalid pubsub permission type, expected %s got %s",
			ExtensionPubSubBroadcastPublish.String(),
			perms.Send[0],
		)
	}

	userId := "100249558"
	perms = c.FormWhisperSendPubSubPermissions(userId)
	if len(perms.Send) == 0 {
		t.Errorf(
			"invalid pubsub permission type, expected 1 got %d",
			len(perms.Send),
		)
	}
	if len(perms.Send) > 1 && strings.Contains(perms.Send[0].String(), ("whisper-"+userId)) {
		t.Errorf(
			"invalid whisper pubsub permission type, does not contain whisper-%s",
			userId,
		)
	}

	perms = c.FormGenericPubSubPermissions()
	if len(perms.Send) == 0 {
		t.Errorf(
			"invalid pubsub permission type, expected 1 got %d",
			len(perms.Send),
		)
	}
	if len(perms.Send) > 1 && perms.Send[0] != ExtensionPubSubGenericPublish {
		t.Errorf(
			"invalid pubsub permission type, expected %s got %s",
			ExtensionPubSubBroadcastPublish.String(),
			perms.Send[0],
		)
	}
}

func TestExtensionSendPubSubMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *SendExtensionPubSubMessageParams
		respBody      string
		validationErr string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&SendExtensionPubSubMessageParams{},
			`{"error":"Unauthorized","status":401,"message":"JWT token is missing"}`,
			"",
		},
		{
			http.StatusNoContent,
			&Options{
				ClientID: "my-client-id",
				ExtensionOpts: ExtensionOptions{
					Secret:      "my-ext-secret",
					OwnerUserID: "ext-owner-id",
				},
			},
			&SendExtensionPubSubMessageParams{
				BroadcasterID: "100249558",
				Message:       "{}",
				Target:        []ExtensionPubSubPublishType{ExtensionPubSubBroadcastPublish},
			},
			"",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendExtensionPubSubMessage(testCase.params)
		if err != nil {
			if err.Error() == testCase.validationErr {
				continue
			}

			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusUnauthorized, resp.ErrorStatus)
			}

			expectedErrMsg := "JWT token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}
