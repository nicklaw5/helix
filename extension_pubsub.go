package helix

// PublishType The Pub/Sub broadcast type
type ExtensionPubSubPublishType string

type pubSubNotification struct {
	Message     string                       `json:"message"`
	Targets     []ExtensionPubSubPublishType `json:"targets"`
	ContentType string                       `json:"content_type"`
}

// Types of Pub/Sub Permissions or targets
const (
	ExtensionPubSubGenericPublish   ExtensionPubSubPublishType = "*"
	ExtensionPubSubBroadcastPublish ExtensionPubSubPublishType = "broadcast"
	ExtensionPubSubGlobalPublish    ExtensionPubSubPublishType = "global"
)

func (c *Client) createExtensionPubSubWhisper(opaqueId string) ExtensionPubSubPublishType {
	return ExtensionPubSubPublishType("whisper-" + opaqueId)
}

// FormWhisperSendPubSubPermissions create the pubsub permissions
// for publishing a whisper message type
func (c *Client) FormWhisperSendPubSubPermissions(opaqueId string) *PubSubPermissions {
	return &PubSubPermissions{
		Send: []ExtensionPubSubPublishType{c.createExtensionPubSubWhisper(opaqueId)},
	}
}

// FormBroadcastSendPubSubPermissions create the pubsub permissions
// for publishing a broadcast message type
func (c *Client) FormBroadcastSendPubSubPermissions() *PubSubPermissions {
	return &PubSubPermissions{
		Send: []ExtensionPubSubPublishType{ExtensionPubSubBroadcastPublish},
	}
}

// FormGlobalSendPubSubPermissions create the pubsub permissions
// for publishing a global targeted message
func FormGlobalSendPubSubPermissions() *PubSubPermissions {
	return &PubSubPermissions{
		Send: []ExtensionPubSubPublishType{ExtensionPubSubGlobalPublish},
	}
}

// FormGenericPubSubPermissions create the pubsub permissions
// for publishing to message for any target type
func FormGenericPubSubPermissions() *PubSubPermissions {
	return &PubSubPermissions{
		Send: []ExtensionPubSubPublishType{ExtensionPubSubGenericPublish},
	}
}

type SendExtensionPubSubMessageParams struct {
	BroadcasterID     string `json:"broadcaster_id"`
	Message           string `json:"message"`
	Target            string `json:"target"`
	IsGlobalBroadcast bool   `json:"is_global_broadcast"`
}

type SendExtensionPubSubMessageResponse struct {
	ResponseCommon
}

func (c *Client) SendExtensionPubSubMessage(params *SendExtensionPubSubMessageParams) (*SendExtensionPubSubMessageResponse, error) {
	resp, err := c.postAsJSON("/extensions/pubsub", &SendExtensionPubSubMessageResponse{}, params)
	if err != nil {
		return nil, err
	}

	sndExtPubSubMsgRsp := &SendExtensionPubSubMessageResponse{}
	resp.HydrateResponseCommon(&sndExtPubSubMsgRsp.ResponseCommon)

	return sndExtPubSubMsgRsp, nil
}
