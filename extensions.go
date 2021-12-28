package helix

import "fmt"

type ExtensionTransaction struct {
	ID               string `json:"id"`
	Timestamp        Time   `json:"timestamp"`
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterLogin string `json:"broadcaster_login"`
	BroadcasterName  string `json:"broadcaster_name"`
	UserID           string `json:"user_id"`
	UserLogin        string `json:"user_login"`
	UserName         string `json:"user_name"`
	ProductType      string `json:"product_type"`
	ProductData      struct {
		Domain     string `json:"domain"`
		Broadcast  bool   `json:"broadcast"`
		Expiration string `json:"expiration"`
		SKU        string `json:"sku"`
		Cost       struct {
			Amount int    `json:"amount"`
			Type   string `json:"type"`
		} `json:"cost"`
		DisplayName   string `json:"displayName"`
		InDevelopment bool   `json:"inDevelopment"`
	} `json:"product_data"`
}

type ManyExtensionTransactions struct {
	ExtensionTransactions []ExtensionTransaction `json:"data"`
	Pagination            Pagination             `json:"pagination"`
}

type ExtensionTransactionsResponse struct {
	ResponseCommon
	Data ManyExtensionTransactions
}

type ExtensionTransactionsParams struct {
	ExtensionID string   `query:"extension_id"` // Required
	ID          []string `query:"id"`           // Optional, Limit 100
	After       string   `query:"after"`        // Optional
	First       int      `query:"first,20"`     // Optional, Limit 100
}

type ExtensionSendChatMessageParams struct {
	BroadcasterID    string `query:"broadcaster_id" json:"-"`
	Text             string `json:"text"` // Limit 280
	ExtensionVersion string `json:"extension_version"`
	ExtensionID      string `json:"extension_id"`
}

type ExtensionSendChatMessageResponse struct {
	ResponseCommon
}

type ExtensionLiveChannel struct {
	BroadcasterID   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	GameName        string `json:"game_name"`
	GameID          string `json:"game_id"`
	Title           string `json:"title"`
}

type ManyExtensionLiveChannels struct {
	LiveChannels []ExtensionLiveChannel `json:"data"`
	Pagination   string                 `json:"pagination"`
}

type ExtensionLiveChannelsParams struct {
	ExtensionID string `query:"extension_id"` // Required
	After       string `query:"after"`        // Optional
	First       int    `query:"first,20"`     // Optional, Limit 100
}

type ExtensionLiveChannelsResponse struct {
	ResponseCommon
	Data ManyExtensionLiveChannels
}

// GetExtensionTransactions allows extension back end servers to fetch a list of transactions that
// have occurred for their extension across all of Twitch. A transaction is a record of a user
// exchanging Bits for an in-Extension digital good.
//
// See https://dev.twitch.tv/docs/api/reference/#get-extension-transactions
func (c *Client) GetExtensionTransactions(params *ExtensionTransactionsParams) (*ExtensionTransactionsResponse, error) {
	resp, err := c.get("/extensions/transactions", &ManyExtensionTransactions{}, params)
	if err != nil {
		return nil, err
	}

	extTxnResp := &ExtensionTransactionsResponse{}
	resp.HydrateResponseCommon(&extTxnResp.ResponseCommon)
	extTxnResp.Data.ExtensionTransactions = resp.Data.(*ManyExtensionTransactions).ExtensionTransactions
	extTxnResp.Data.Pagination = resp.Data.(*ManyExtensionTransactions).Pagination
	return extTxnResp, nil
}

// SendExtensionChatMessage  Sends a specified chat message to a specified channel.
// The message will appear in the channel’s chat as a normal message,
// The author of the message is the Extension name.
//
// see https://dev.twitch.tv/docs/api/reference#send-extension-chat-message
func (c *Client) SendExtensionChatMessage(params *ExtensionSendChatMessageParams) (*ExtensionSendChatMessageResponse, error) {

	if len(params.Text) > 280 {
		return nil, fmt.Errorf("error: chat message length exceeds 280 characters")
	}

	if params.BroadcasterID == "" {
		return nil, fmt.Errorf("error: broadcaster ID must be specified")
	}

	resp, err := c.postAsJSON("/extensions/chat", &ExtensionSendChatMessageResponse{}, params)
	if err != nil {
		return nil, err
	}

	sndExtMsgResp := &ExtensionSendChatMessageResponse{}
	resp.HydrateResponseCommon(&sndExtMsgResp.ResponseCommon)

	return sndExtMsgResp, nil
}

func (c *Client) GetExtensionLiveChannels(params *ExtensionLiveChannelsParams) (*ExtensionLiveChannelsResponse, error) {

	if params.ExtensionID == "" {
		return nil, fmt.Errorf("error: extension ID must be specified")
	}

	resp, err := c.get("/extensions/live", &ManyExtensionLiveChannels{}, params)
	if err != nil {
		return nil, err
	}

	liveChannels := &ExtensionLiveChannelsResponse{}
	resp.HydrateResponseCommon(&liveChannels.ResponseCommon)
	liveChannels.Data.LiveChannels = resp.Data.(*ManyExtensionLiveChannels).LiveChannels
	liveChannels.Data.Pagination = resp.Data.(*ManyExtensionLiveChannels).Pagination
	return liveChannels, nil
}
