package helix

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
