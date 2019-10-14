# Subscriptions Documentation

## Get Broadcaster Subscriptions

This is an example of how to get the broadcaster subscriptions.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetSubscriptions(&helix.SubscriptionsParams{
    BroadcasterID:  "29776980",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
