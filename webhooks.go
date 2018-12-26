package helix

// WebhookSubscription ...
type WebhookSubscription struct {
	Topic     string `json:"topic"`
	Callback  string `json:"callback"`
	ExpiresAt Time   `json:"expires_at"`
}

// ManyWebhookSubscriptions ...
type ManyWebhookSubscriptions struct {
	Total                int                   `json:"total"`
	WebhookSubscriptions []WebhookSubscription `json:"data"`
	Pagination           Pagination            `json:"pagination"`
}

// WebhookSubscriptionsResponse ...
type WebhookSubscriptionsResponse struct {
	ResponseCommon
	Data ManyWebhookSubscriptions
}

// WebhookSubscriptionsParams ...
type WebhookSubscriptionsParams struct {
	After string `query:"after"`
	First int    `query:"first,20"` // Limit 100
}

// GetWebhookSubscriptions gets webhook subscriptions, in order of expiration.
// Requires an app access token.
func (c *Client) GetWebhookSubscriptions(params *WebhookSubscriptionsParams) (*WebhookSubscriptionsResponse, error) {
	resp, err := c.get("/webhooks/subscriptions", &ManyWebhookSubscriptions{}, params)
	if err != nil {
		return nil, err
	}

	webhooks := &WebhookSubscriptionsResponse{}
	webhooks.StatusCode = resp.StatusCode
	webhooks.Header = resp.Header
	webhooks.Error = resp.Error
	webhooks.ErrorStatus = resp.ErrorStatus
	webhooks.ErrorMessage = resp.ErrorMessage
	webhooks.Data.Total = resp.Data.(*ManyWebhookSubscriptions).Total
	webhooks.Data.WebhookSubscriptions = resp.Data.(*ManyWebhookSubscriptions).WebhookSubscriptions
	webhooks.Data.Pagination = resp.Data.(*ManyWebhookSubscriptions).Pagination

	return webhooks, nil
}

// WebhookSubscriptionResponse ...
type WebhookSubscriptionResponse struct {
	ResponseCommon
}

// WebhookSubscriptionPayload ...
type WebhookSubscriptionPayload struct {
	Mode         string `json:"hub.mode"`
	Topic        string `json:"hub.topic"`
	Callback     string `json:"hub.callback"`
	LeaseSeconds int    `json:"hub.lease_seconds,omitempty"`
	Secret       string `json:"secret,omitempty"`
}

// PostWebhookSubscription ...
func (c *Client) PostWebhookSubscription(payload *WebhookSubscriptionPayload) (*WebhookSubscriptionResponse, error) {
	resp, err := c.post("/webhooks/hub", nil, payload)
	if err != nil {
		return nil, err
	}

	webhook := &WebhookSubscriptionResponse{}
	webhook.StatusCode = resp.StatusCode
	webhook.Header = resp.Header
	webhook.Error = resp.Error
	webhook.ErrorStatus = resp.ErrorStatus
	webhook.ErrorMessage = resp.ErrorMessage

	return webhook, nil
}