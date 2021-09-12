# Analytics Documentation

## Get Game Analytics

This is an example of how to get the downloadable CSV file containing analytics data:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

gameID := "493057"

resp, err := client.GetGameAnalytics(gameID)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Extensions Analytics

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

params := helix.ExtensionAnalyticsParams{
    ExtensionID: "abcd",
    Type:        "overview_v1",
}

resp, err := client.GetExtensionAnalytics(&params)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
