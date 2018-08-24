# Webhook Documentation

## Get Webhook Subscriptions

This is an example of how to get webhook subscriptions.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    AppAccessToken: "your-app-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetWebhookSubscriptions(&helix.WebhookSubscriptionsParams{
    First: 10,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
